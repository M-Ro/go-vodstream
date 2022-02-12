package sql

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"git.thorn.sh/Thorn/go-vodstream/internal/domain"
	"git.thorn.sh/Thorn/go-vodstream/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var testDb *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		testDb, err = sqlx.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return testDb.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	testDb.MustExec(fmt.Sprintf(`CREATE TABLE %s (
		id          BIGSERIAL   PRIMARY KEY,
		username    TEXT        NOT NULL,
		email       TEXT        NOT NULL,
		password    TEXT,
		publish_key TEXT,
		can_publish BOOL,
		can_stream  BOOL,
		created_at  TIMESTAMP,
		updated_at  TIMESTAMP
	);`, UsersTableName))

	seed(testDb)

	// Run tests
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// seed prepares the database with testing data
func seed(db *sqlx.DB) {
	// Ensure table is empty
	db.MustExec(fmt.Sprintf("DELETE FROM %s", UsersTableName))
	db.MustExec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", UsersTableName))

	time := time.Now().Truncate(time.Microsecond)

	users := []storage.User{
		{
			Username:  "testUser1",
			Email:     "testUser1@example.com",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			Username:  "testUser2",
			Email:     "testUser2@example.com",
			CreatedAt: time,
			UpdatedAt: time,
		},
		{
			Username:  "testUser3",
			Email:     "testUser3@example.com",
			CreatedAt: time,
			UpdatedAt: time,
		},
	}

	for _, user := range users {
		db.MustExec(fmt.Sprintf(`
			INSERT INTO %s (username, email, password, publish_key, can_publish, can_stream, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		`, UsersTableName), user.Username, user.Email, user.Password, user.PublishKey, user.CanPublish, user.CanStream,
			user.CreatedAt, user.UpdatedAt)
	}
}

func TestUserStorage_All(t *testing.T) {
	ctx := context.Background()

	curTime := time.Now()

	tests := []struct {
		testName       string
		expectedReturn []storage.User
		expectedError  error
	}{
		{
			testName: "Expect matching rows",
			expectedReturn: []storage.User{
				{
					Id:        1,
					Username:  "testUser1",
					Email:     "testUser1@example.com",
					CreatedAt: curTime,
					UpdatedAt: curTime,
				},
				{
					Id:        2,
					Username:  "testUser2",
					Email:     "testUser2@example.com",
					CreatedAt: curTime,
					UpdatedAt: curTime,
				},
				{
					Id:        3,
					Username:  "testUser3",
					Email:     "testUser3@example.com",
					CreatedAt: curTime,
					UpdatedAt: curTime,
				},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			storage := NewUserStorage(testDb)
			time.Sleep(time.Second * 2)
			seed(testDb)

			// Perform check and compare output
			users, err := storage.All(ctx)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(users, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(users, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestUserStorage_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		testName        string
		paginateOptions domain.PaginateQueryOptions
		expectedReturn  []storage.User
		expectedError   error
	}{
		{
			testName:        "test paginate defaults",
			paginateOptions: domain.NewPaginateOptions(),
			expectedError:   nil,
			expectedReturn: []storage.User{
				{
					Id: 1,
				},
				{
					Id: 2,
				},
				{
					Id: 3,
				},
			},
		},
		{
			testName: "order by descending",
			paginateOptions: domain.PaginateQueryOptions{
				Limit:  25,
				Offset: 0,
				Order: domain.PaginateOrder{
					Field:  "id",
					Method: domain.OrderMethodDesc,
				},
			},
			expectedError: nil,
			expectedReturn: []storage.User{
				{
					Id: 3,
				},
				{
					Id: 2,
				},
				{
					Id: 1,
				},
			},
		},
		{
			testName: "test order field",
			paginateOptions: domain.PaginateQueryOptions{
				Limit:  25,
				Offset: 0,
				Order: domain.PaginateOrder{
					Field:  "username",
					Method: domain.OrderMethodAsc,
				},
			},
			expectedError: nil,
			expectedReturn: []storage.User{
				{
					Id: 1,
				},
				{
					Id: 2,
				},
				{
					Id: 3,
				},
			},
		},
		{
			testName: "test limit",
			paginateOptions: domain.PaginateQueryOptions{
				Limit:  2,
				Offset: 0,
				Order: domain.PaginateOrder{
					Field:  "id",
					Method: domain.OrderMethodAsc,
				},
			},
			expectedError: nil,
			expectedReturn: []storage.User{
				{
					Id: 1,
				},
				{
					Id: 2,
				},
			},
		},
		{
			testName: "test offset",
			paginateOptions: domain.PaginateQueryOptions{
				Limit:  2,
				Offset: 1,
				Order: domain.PaginateOrder{
					Field:  "id",
					Method: domain.OrderMethodAsc,
				},
			},
			expectedError: nil,
			expectedReturn: []storage.User{
				{
					Id: 2,
				},
				{
					Id: 3,
				},
			},
		},
		{
			testName: "test offset larger than row count",
			paginateOptions: domain.PaginateQueryOptions{
				Limit:  25,
				Offset: 250,
				Order: domain.PaginateOrder{
					Field:  "id",
					Method: domain.OrderMethodAsc,
				},
			},
			expectedError:  nil,
			expectedReturn: []storage.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			userStorage := NewUserStorage(testDb)
			seed(testDb)

			// Perform check and compare output
			users, err := userStorage.List(ctx, test.paginateOptions)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			ignoreOpts := cmpopts.IgnoreFields(storage.User{}, "Username", "Email", "Password", "CreatedAt", "UpdatedAt")

			if !cmp.Equal(users, test.expectedReturn, ignoreOpts) {
				t.Fatal(cmp.Diff(users, test.expectedReturn, ignoreOpts))
			}
		})
	}
}

func TestUserStorage_GetByID(t *testing.T) {
	ctx := context.Background()
	curTime := time.Now()

	tests := []struct {
		testName       string
		Id             uint64
		expectedReturn *storage.User
	}{
		{
			testName: "expect correct user",
			Id:       2,
			expectedReturn: &storage.User{
				Id:        2,
				Username:  "testUser2",
				Email:     "testUser2@example.com",
				CreatedAt: curTime,
				UpdatedAt: curTime,
			},
		},
		{
			testName:       "expect nil for invalid user",
			Id:             60,
			expectedReturn: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			storage := NewUserStorage(testDb)
			seed(testDb)

			// Perform check and compare output
			user := storage.GetByID(ctx, test.Id)

			if !cmp.Equal(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestUserStorage_GetByUsername(t *testing.T) {
	ctx := context.Background()
	curTime := time.Now()

	tests := []struct {
		testName       string
		Username       string
		expectedReturn *storage.User
	}{
		{
			testName: "expect correct user",
			Username: "testUser2",
			expectedReturn: &storage.User{
				Id:        2,
				Username:  "testUser2",
				Email:     "testUser2@example.com",
				CreatedAt: curTime,
				UpdatedAt: curTime,
			},
		},
		{
			testName: "expect correct user, test case sensitivity",
			Username: "TESTUSER2",
			expectedReturn: &storage.User{
				Id:        2,
				Username:  "testUser2",
				Email:     "testUser2@example.com",
				CreatedAt: curTime,
				UpdatedAt: curTime,
			},
		},
		{
			testName:       "expect nil for invalid user",
			Username:       "thisUserDoesNotExist",
			expectedReturn: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			storage := NewUserStorage(testDb)
			seed(testDb)

			// Perform check and compare output
			user := storage.GetByUsername(ctx, test.Username)

			if !cmp.Equal(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestUserStorage_GetByEmail(t *testing.T) {
	ctx := context.Background()
	curTime := time.Now()

	tests := []struct {
		testName       string
		Email          string
		expectedReturn *storage.User
	}{
		{
			testName: "expect correct user",
			Email:    "testUser2@example.com",
			expectedReturn: &storage.User{
				Id:        2,
				Username:  "testUser2",
				Email:     "testUser2@example.com",
				CreatedAt: curTime,
				UpdatedAt: curTime,
			},
		},
		{
			testName: "expect correct user, test case sensitivity",
			Email:    "TESTUSER2@eXample.COM",
			expectedReturn: &storage.User{
				Id:        2,
				Username:  "testUser2",
				Email:     "testUser2@example.com",
				CreatedAt: curTime,
				UpdatedAt: curTime,
			},
		},
		{
			testName:       "expect nil for invalid user",
			Email:          "thisUserDoesNotExist@testuser.com",
			expectedReturn: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			storage := NewUserStorage(testDb)
			seed(testDb)

			// Perform check and compare output
			user := storage.GetByEmail(ctx, test.Email)

			if !cmp.Equal(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedReturn, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestUserStorage_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		testName      string
		Id            uint64
		expectedError error
	}{
		{
			testName:      "Expect success deleting valid user",
			Id:            1,
			expectedError: nil,
		},
		{
			testName:      "Expect success deleting invalid user",
			Id:            61,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			userStorage := NewUserStorage(testDb)
			seed(testDb)

			// Perform check and compare output
			storageErr := userStorage.Delete(ctx, test.Id)

			if !cmp.Equal(storageErr, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(storageErr, test.expectedError, cmpopts.EquateErrors()))
			}

			// Ensure row is gone
			row := testDb.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), test.Id)
			var user storage.User
			err := row.StructScan(&user)

			if !cmp.Equal(err, sql2.ErrNoRows, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, sql2.ErrNoRows, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestUserStorage_Insert(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		testName      string
		user          storage.User
		expectedError error
	}{
		{
			testName: "expect success on valid user",
			user: storage.User{
				Username: "insertUser1",
				Email:    "insertUser1@example.com",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			userStorage := NewUserStorage(testDb)
			seed(testDb)

			err := userStorage.Insert(ctx, &test.user)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			// Ensure row is added
			row := testDb.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE username = $1`), test.user.Username)
			var user storage.User
			err = row.StructScan(&user)

			if err != nil {
				t.Fatalf("Inserted user not found %v", err)
			}
		})
	}
}

func TestUserStorage_Update(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		testName      string
		user          storage.User
		expectedError error
	}{
		{
			testName: "expect success on valid update",
			user: storage.User{
				Id:       1,
				Username: "updateTest1",
				Email:    "updateTest1@example.com",
			},
			expectedError: nil,
		},
		{
			testName: "expect fail updating non-existent user",
			user: storage.User{
				Id:       25,
				Username: "updateTest2",
				Email:    "updateTest2@example.com",
			},
			expectedError: ErrNoRowsAffected,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			userStorage := NewUserStorage(testDb)
			seed(testDb)

			err := userStorage.Update(ctx, test.user.Id, &test.user)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if err == nil {
				// Ensure row is updated
				row := testDb.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), test.user.Id)
				var user storage.User
				err = row.StructScan(&user)

				if err != nil {
					t.Fatalf("Updated user not found %v", err)
				}

				if !cmp.Equal(user.Username, test.user.Username) {
					t.Fatal(cmp.Diff(user.Username, test.user.Username))
				}
			}
		})
	}
}

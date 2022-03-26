package user

import (
	"context"
	"errors"
	"github.com/M-Ro/go-vodstream/internal/domain/user"
	"github.com/M-Ro/go-vodstream/internal/paginate"
	"github.com/M-Ro/go-vodstream/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
)

func TestRepository_All(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		expectedError error
		expectedValue []user.User
	}{
		{
			testName: "expect nothing with empty empty storage",
			mockStorage: mockUserStorage{
				ReturnAllUsers: []storage.User{},
				ReturnAllError: nil,
			},
			expectedError: nil,
			expectedValue: []user.User{},
		},
		{
			testName: "expect empty return, with error on storage failure",
			mockStorage: mockUserStorage{
				ReturnAllUsers: []storage.User{},
				ReturnAllError: mockStorageErr,
			},
			expectedError: mockStorageErr,
			expectedValue: []user.User{},
		},
		{
			testName: "expect correct return from given storage",
			mockStorage: mockUserStorage{
				ReturnAllUsers: []storage.User{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				ReturnAllError: nil,
			},
			expectedError: nil,
			expectedValue: []user.User{
				{
					Id: 0,
				},
				{
					Id: 1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			users, err := r.All(ctx)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(users, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(users, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_List(t *testing.T) {
	tests := []struct {
		testName        string
		mockStorage     mockUserStorage
		paginateOptions paginate.QueryOptions
		expectedError   error
		expectedValue   []user.User
	}{
		{
			testName: "expect empty array with empty storage",
			mockStorage: mockUserStorage{
				ReturnListUsers: []storage.User{},
				ReturnListError: nil,
			},
			paginateOptions: paginate.QueryOptions{},
			expectedError:   nil,
			expectedValue:   []user.User{},
		},
		{
			testName: "expect error when storage fails",
			mockStorage: mockUserStorage{
				ReturnListUsers: []storage.User{},
				ReturnListError: ErrUserNotFound,
			},
			paginateOptions: paginate.QueryOptions{},
			expectedError:   ErrUserNotFound,
			expectedValue:   []user.User{},
		},
		{
			testName: "expect correct output",
			mockStorage: mockUserStorage{
				ReturnListError: nil,
				ReturnListUsers: []storage.User{
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
			paginateOptions: paginate.QueryOptions{},
			expectedError:   nil,
			expectedValue: []user.User{
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
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			users, err := r.List(ctx, test.paginateOptions)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(users, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(users, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		userId        uint64
		expectedError error
		expectedValue user.User
	}{
		{
			testName: "expect error with empty storage",
			mockStorage: mockUserStorage{
				ReturnGetByIDUser: nil,
			},
			userId:        1,
			expectedValue: user.User{},
			expectedError: ErrUserNotFound,
		},
		{
			testName: "expect success with valid user",
			mockStorage: mockUserStorage{
				ReturnGetByIDUser: &storage.User{
					Id: 1,
				},
			},
			userId: 1,
			expectedValue: user.User{
				Id: 1,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			user, err := r.GetByID(ctx, test.userId)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		email         string
		expectedError error
		expectedValue user.User
	}{
		{
			testName: "expect error with empty storage",
			mockStorage: mockUserStorage{
				ReturnGetByEmailUser: nil,
			},
			email:         "testUser@testUser.com",
			expectedValue: user.User{},
			expectedError: ErrUserNotFound,
		},
		{
			testName: "expect success with valid user",
			mockStorage: mockUserStorage{
				ReturnGetByEmailUser: &storage.User{
					Email: "testUser@testUser.com",
				},
			},
			email: "testUser@testUser.com",
			expectedValue: user.User{
				Email: "testUser@testUser.com",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			user, err := r.GetByEmail(ctx, test.email)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_GetByUsername(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		username      string
		expectedError error
		expectedValue user.User
	}{
		{
			testName: "expect error with empty storage",
			mockStorage: mockUserStorage{
				ReturnGetByUsernameUser: nil,
			},
			username:      "testUser",
			expectedValue: user.User{},
			expectedError: ErrUserNotFound,
		},
		{
			testName: "expect success with valid user",
			mockStorage: mockUserStorage{
				ReturnGetByUsernameUser: &storage.User{
					Username: "testUser",
				},
			},
			username: "testUser",
			expectedValue: user.User{
				Username: "testUser",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			user, err := r.GetByUsername(ctx, test.username)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_Insert(t *testing.T) {
	tests := []struct {
		testName       string
		mockStorage    mockUserStorage
		user           user.User
		expectedError  error
		expectedResult user.User
	}{
		{
			testName: "expect error on storage failure",
			mockStorage: mockUserStorage{
				ReturnInsertError: mockStorageErr,
				ReturnInsertId:    0,
			},
			user: user.User{
				Username: "testUser",
				Email:    "testUser@blah",
			},
			expectedError: mockStorageErr,
			expectedResult: user.User{
				Username: "testUser",
				Email:    "testUser@blah",
			},
		},
		{
			testName: "expect valid domain user on insert",
			mockStorage: mockUserStorage{
				ReturnInsertError: nil,
				ReturnInsertId:    1,
			},
			user: user.User{
				Id:       0,
				Username: "testUser",
				Email:    "testUser@blah",
			},
			expectedError: nil,
			expectedResult: user.User{
				Id:       1,
				Username: "testUser",
				Email:    "testUser@blah",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			err := r.Insert(ctx, &test.user)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if err != nil {
				return
			}

			if test.user.Id == 0 {
				t.Fatal("UserID should not be 0 after insert")
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		user          user.User
		userId        uint64
		expectedValue user.User
		expectedError error
	}{
		{
			testName: "expect failure for non existing user",
			mockStorage: mockUserStorage{
				ReturnUpdateError: ErrUserNotFound,
			},
			user:          user.User{},
			userId:        32,
			expectedValue: user.User{},
			expectedError: ErrUserNotFound,
		},
		{
			testName: "expect valid domain user when updating existing user",
			mockStorage: mockUserStorage{
				ReturnUpdateDateUpdated: time.Now(),
				ReturnUpdateError:       nil,
			},
			user: user.User{
				Id:        1,
				UpdatedAt: time.Time{},
			},
			userId: 1,
			expectedValue: user.User{
				Id:        1,
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			user, err := r.Update(ctx, test.user.Id, test.user)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}

			if !cmp.Equal(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)) {
				t.Fatal(cmp.Diff(user, test.expectedValue, cmpopts.EquateApproxTime(time.Minute)))
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	tests := []struct {
		testName      string
		mockStorage   mockUserStorage
		userId        uint64
		expectedError error
	}{
		{
			testName: "expect error on storage layer error",
			mockStorage: mockUserStorage{
				ReturnDeleteError: ErrUserNotFound,
			},
			userId:        16,
			expectedError: ErrUserNotFound,
		},
		{
			testName: "expect success on valid delete",
			mockStorage: mockUserStorage{
				ReturnDeleteError: nil,
			},
			userId:        1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			ctx := context.Background()

			r := Repository{
				StorageProvider: test.mockStorage,
			}

			err := r.Delete(ctx, test.userId)

			if !cmp.Equal(err, test.expectedError, cmpopts.EquateErrors()) {
				t.Fatal(cmp.Diff(err, test.expectedError, cmpopts.EquateErrors()))
			}
		})
	}
}

func TestNewRepository(t *testing.T) {
	s := mockUserStorage{}
	expected := Repository{
		StorageProvider: s,
	}

	repository := NewRepository(s)

	if !cmp.Equal(expected, repository) {
		t.Fatal(cmp.Diff(expected, repository))
	}
}

var mockStorageErr = errors.New("database error")

type mockUserStorage struct {
	ReturnAllUsers          []storage.User
	ReturnAllError          error
	ReturnListUsers         []storage.User
	ReturnListError         error
	ReturnGetByIDUser       *storage.User
	ReturnGetByUsernameUser *storage.User
	ReturnGetByEmailUser    *storage.User
	ReturnDeleteError       error
	ReturnInsertError       error
	ReturnInsertId          uint64
	ReturnUpdateError       error
	ReturnUpdateDateUpdated time.Time
}

func (m mockUserStorage) All(_ context.Context) ([]storage.User, error) {
	return m.ReturnAllUsers, m.ReturnAllError
}

func (m mockUserStorage) List(_ context.Context, _ paginate.QueryOptions) ([]storage.User, error) {
	return m.ReturnListUsers, m.ReturnListError
}

func (m mockUserStorage) GetByID(_ context.Context, _ uint64) *storage.User {
	return m.ReturnGetByIDUser
}

func (m mockUserStorage) GetByUsername(_ context.Context, _ string) *storage.User {
	return m.ReturnGetByUsernameUser
}

func (m mockUserStorage) GetByEmail(_ context.Context, _ string) *storage.User {
	return m.ReturnGetByEmailUser
}

func (m mockUserStorage) Delete(_ context.Context, _ uint64) error {
	return m.ReturnDeleteError
}

func (m mockUserStorage) Insert(_ context.Context, user *storage.User) error {
	if m.ReturnInsertError != nil {
		return m.ReturnInsertError
	}

	user.Id = m.ReturnInsertId

	return nil
}

func (m mockUserStorage) Update(_ context.Context, _ uint64, user *storage.User) error {
	user.UpdatedAt = m.ReturnUpdateDateUpdated
	return m.ReturnUpdateError
}

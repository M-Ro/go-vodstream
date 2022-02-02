package handlers

import (
	"bytes"
	"encoding/json"
	"git.thorn.sh/Thorn/go-vodstream/api"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	tests := []struct {
		testName   string
		testMethod string
		handler    AuthHandler

		reqEndpoint string
		reqBody     api.LoginRequest

		respStatus int
		respBody   api.LoginResponse
	}{
		{
			testName:   "Expect success (200) with correct credentials, using username.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/login",
			reqBody: api.LoginRequest{
				UsernameOrEmail: "user1",
				Password:        "P@ssword",
				RememberMe:      false,
			},

			respStatus: 200,
			respBody: api.LoginResponse{
				Success:     true,
				Errors:      []string{},
				UserID:      0,
				AccessToken: "", // Manual check.
			},
		},
		{
			testName:   "Expect success (200) with correct credentials, using email.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/login",
			reqBody: api.LoginRequest{
				UsernameOrEmail: "email@example.com",
				Password:        "P@ssword",
				RememberMe:      false,
			},

			respStatus: 200,
			respBody: api.LoginResponse{
				Success:     true,
				Errors:      []string{},
				UserID:      0,
				AccessToken: "", // Manual check.
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.LoginResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			// Manual check for API Token since it could be anything, we only care that it isn't empty for now.
			// TODO proper comparator for a valid JWT
			if result.AccessToken == "" {
				t.Fatal("Access token is empty.")
			}
			result.AccessToken = ""

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}

func TestAuthHandler_Login_Fail(t *testing.T) {
	tests := []struct {
		testName   string
		testMethod string
		handler    AuthHandler

		reqEndpoint string
		reqBody     api.LoginRequest

		respStatus int
		respBody   api.LoginResponse
	}{
		{
			testName:   "Expect error (401) with no credentials.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/login",
			reqBody: api.LoginRequest{
				UsernameOrEmail: "",
				Password:        "",
				RememberMe:      false,
			},

			respStatus: 401,
			respBody: api.LoginResponse{
				Success:     false,
				Errors:      []string{ErrLoginInvalidCredentials.Error()},
				UserID:      0,
				AccessToken: "",
			},
		},
		{
			testName:   "Expect error (401) with invalid username/password.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/login",
			reqBody: api.LoginRequest{
				UsernameOrEmail: "notARealUser",
				Password:        "notARealPassword",
				RememberMe:      false,
			},

			respStatus: 401,
			respBody: api.LoginResponse{
				Success:     false,
				Errors:      []string{ErrLoginInvalidCredentials.Error()},
				UserID:      0,
				AccessToken: "",
			},
		},
		{
			testName:   "Expect error (401) with different users password.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/login",
			reqBody: api.LoginRequest{
				UsernameOrEmail: "user2",
				Password:        "P@ssword",
				RememberMe:      false,
			},

			respStatus: 401,
			respBody: api.LoginResponse{
				Success:     false,
				Errors:      []string{ErrLoginInvalidCredentials.Error()},
				UserID:      0,
				AccessToken: "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.LoginResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}

func TestAuthHandler_Register_Success(t *testing.T) {
	tests := []struct {
		testName   string
		testMethod string
		handler    AuthHandler

		reqEndpoint string
		reqBody     api.RegisterRequest

		respStatus int
		respBody   api.RegisterResponse
	}{
		{
			testName:   "Expect success (200) with valid details.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "user1",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 200,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{},
				Messages: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.RegisterResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}

func TestAuthHandler_Register_Fail(t *testing.T) {
	tests := []struct {
		testName   string
		testMethod string
		handler    AuthHandler

		reqEndpoint string
		reqBody     api.RegisterRequest

		respStatus int
		respBody   api.RegisterResponse
	}{
		{
			testName:   "Expect error (401) with missing details.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "user1",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationFailed,
				Errors:   []string{ErrRegisterMissingDetails.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) with missing username.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrRegisterMissingDetails.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) with missing email.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "",
				Username:        "testuser",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrRegisterMissingDetails.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) with existing email.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "testuser@exists.com",
				Username:        "testuser",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrRegisterEmailExists.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) with existing username.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "userexists",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrRegisterUsernameExists.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) if password does not match requirements.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "testuser",
				Password:        "pass",
				ConfirmPassword: "pass",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrPasswordRequirements.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) if ConfirmPassword != Password.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "testuser",
				Password:        "P@ssword",
				ConfirmPassword: "_P@ssword",
				AcceptedTerms:   true,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrPasswordMatch.Error()},
				Messages: []string{},
			},
		},
		{
			testName:   "Expect error (401) if AcceptedTerms != true",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.RegisterRequest{
				Email:           "user@example.org",
				Username:        "testuser",
				Password:        "P@ssword",
				ConfirmPassword: "P@ssword",
				AcceptedTerms:   false,
			},

			respStatus: 401,
			respBody: api.RegisterResponse{
				Status:   api.RegistrationComplete,
				Errors:   []string{ErrRegisterNoTermsAccepted.Error()},
				Messages: []string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.RegisterResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}

/*
func TestAuthHandler_ResetPassword_Success(t *testing.T) {
	tests := []struct {
		testName	string
		testMethod	string
		handler		AuthHandler

		reqEndpoint	string
		reqBody		api.ChangePasswordRequest

		respStatus	int
		respBody	api.ChangePasswordResponse
	}{
		{
			testName:   "Expect success (200) with valid details.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.ChangePasswordRequest{

			},

			respStatus: 401,
			respBody: api.ChangePasswordResponse{

			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.ChangePasswordResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}

func TestAuthHandler_ResetPassword_Fail(t *testing.T) {
	tests := []struct {
		testName	string
		testMethod	string
		handler		AuthHandler

		reqEndpoint	string
		reqBody		api.ChangePasswordRequest

		respStatus	int
		respBody	api.ChangePasswordResponse
	}{
		{
			testName:   "Expect error (401) blah.",
			testMethod: http.MethodPost,
			handler:    AuthHandler{},

			reqEndpoint: "/v1/auth/register",
			reqBody: api.ChangePasswordRequest{

			},

			respStatus: 401,
			respBody: api.ChangePasswordResponse{

			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Encode request body as JSON
			b, err := json.Marshal(test.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// Perform HTTP test, fetch result.
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(test.testMethod, test.reqEndpoint, bytes.NewReader(b))

			r := mux.NewRouter()

			handler := AuthHandler{}
			handler.RegisterRoutes(r)

			r.ServeHTTP(recorder, req)
			resp := recorder.Result()

			// Compare response output to expected test output.
			if !cmp.Equal(resp.StatusCode, test.respStatus) {
				t.Fatal(cmp.Diff(resp.StatusCode, test.respStatus))
			}

			result := api.ChangePasswordResponse{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(result, test.respBody) {
				t.Fatal(cmp.Diff(result, test.respBody))
			}
		})
	}
}
*/

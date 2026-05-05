package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kidx45/Project-KK/Backend-Team/api"
	"github.com/kidx45/Project-KK/Backend-Team/db/mockdb"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) (db.User, string) {
	password := utils.RandomLetterString(10) // Raw password
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	user := db.User{
		Username:       utils.RandomUserName(),
		HashedPassword: hashedPassword, // Actual hash
		Email:          utils.RandomEmail(),
		FullName:       utils.RandomFullName(),
		PhoneNumber:    utils.RandomPhoneNumber(),
	}
	return user, password
}

type EqCreateUserParamsMatcher struct {
	user     db.CreateUserParams
	password string
}

func (e *EqCreateUserParamsMatcher) Matches(x interface{}) bool {
	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.CheckPassword(e.password, args.HashedPassword)
	if err != nil {
		return false
	}

	e.user.HashedPassword = args.HashedPassword
	return reflect.DeepEqual(e.user, args)
}

func (e *EqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("is %v", e.user)
}

func EqCreateUserParams(user db.CreateUserParams, password string) gomock.Matcher {
	return &EqCreateUserParamsMatcher{user: user, password: password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)
	testCase := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":    user.Username,
				"email":       user.Email,
				"password":    password,
				"fullName":    user.FullName,
				"phoneNumber": user.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username:    user.Username,
					Email:       user.Email,
					FullName:    user.FullName,
					PhoneNumber: user.PhoneNumber,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				gotUser, err := utils.ConvertBuffertoUserType(recorder.Body)
				require.NoError(t, err)
				require.Equal(t, gotUser.Username, user.Username)
				require.Equal(t, gotUser.Email, user.Email)
				require.Equal(t, gotUser.FullName, user.FullName)
				require.Equal(t, gotUser.PhoneNumber, user.PhoneNumber)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":    user.Username,
				"email":       user.Email,
				"password":    password,
				"fullName":    user.FullName,
				"phoneNumber": user.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":    user.Username,
				"email":       user.Email,
				"password":    password,
				"fullName":    user.FullName,
				"phoneNumber": user.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":    user.Username,
				"email":       "invalid-email",
				"password":    password,
				"fullName":    user.FullName,
				"phoneNumber": user.PhoneNumber,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testCase {
		t.Run(testCase[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase[i].buildStubs(store)
			server := api.NewServer(store)
			recorder := httptest.NewRecorder()
			reqBody, err := json.Marshal(testCase[i].body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			server.Router.ServeHTTP(recorder, req)
			testCase[i].checkResponse(t, recorder)
		})

	}
}
func TestGetUser(t *testing.T) {
	user, _ := randomUser(t)
	testCase := []struct {
		name          string
		username      string
		buildstubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name:     "OK",
		username: user.Username,
		buildstubs: func(store *mockdb.MockStore) {
			store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			gotUser, err := utils.ConvertBuffertoUserType(recorder.Body)

			require.NoError(t, err)
			require.Equal(t, gotUser, user)
		},
	}, {
		name:     "Not Found",
		username: user.Username,
		buildstubs: func(store *mockdb.MockStore) {
			store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusNotFound, recorder.Code)
			requireBodyMatchError(t, recorder.Body, sql.ErrNoRows)
		},
	}}

	for i := range testCase {
		t.Run(testCase[i].name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase[i].buildstubs(store)
			server := api.NewServer(store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/user/%s", user.Username)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			server.Router.ServeHTTP(recorder, req)
			testCase[i].checkResponse(t, recorder)
		})

	}
}

func requireBodyMatchError(t *testing.T, body *bytes.Buffer, err error) {
	data, errRead := io.ReadAll(body)
	require.NoError(t, errRead)

	var gotError map[string]string
	errUnmarshal := json.Unmarshal(data, &gotError)
	require.NoError(t, errUnmarshal)

	require.Contains(t, gotError, "error")
	require.Equal(t, err.Error(), gotError["error"])
}

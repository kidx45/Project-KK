package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kidx45/Project-KK/Backend-Team/api"
	"github.com/kidx45/Project-KK/Backend-Team/db/mockdb"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func RandomUser() db.User {
	return db.User{
		Username:       utils.RandomUserName(),
		HashedPassword: utils.RandomPassword(),
		Email:          utils.RandomEmail(),
		FullName:       utils.RandomFullName(),
		PhoneNumber:    utils.RandomPhoneNumber(),
	}
}
func TestGetUser(t *testing.T) {
	user := RandomUser()
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
		name: "Not Found",
		username: user.Username,
		buildstubs: func(store *mockdb.MockStore) {
			store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(db.User{},sql.ErrNoRows)
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

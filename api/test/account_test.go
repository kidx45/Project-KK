package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kidx45/Project-KK/Backend-Team/db/mockdb"
	db "github.com/kidx45/Project-KK/Backend-Team/db/sqlc"
	"github.com/kidx45/Project-KK/Backend-Team/token"
	"github.com/kidx45/Project-KK/Backend-Team/utils"
	"github.com/stretchr/testify/require"
)

func randomAccount(user string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 10),
		Username: user,
		Balance:  utils.RandomInt(0, 1000),
		Currency: utils.RandomCurrency(),
	}
}

func TestCreateAccount(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name: "ok",
		body: gin.H{
			"username": user.Username,
			"currency": account.Currency,
		},
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuthorization(t, request, tokenMaker, "Bearer", user.Username, time.Minute)
		},
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(account, nil).Times(1)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
		},
	}, {
		name: "no Authorization",
		body: gin.H{
			"username": user.Username,
			"currency": account.Currency,
		},
		setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			addAuthorization(t, request, tokenMaker, "", "", time.Minute)
		},
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(account, nil).Times(0)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	}}

	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			testCase := testCases[i]

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			reqBody, err := json.Marshal(testCase.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/account", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}
			testCase.setupAuth(t, req, server.Token)
			server.Router.ServeHTTP(recorder, req)
			testCase.checkResponse(t, recorder)
		})
	}
}

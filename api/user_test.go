package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Ma-hiru/simplebank/db/mock"
	db "github.com/Ma-hiru/simplebank/db/sqlc"
	"github.com/Ma-hiru/simplebank/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	var hashedPassword, err = util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:     util.RandomOwner(),
		HashPassword: hashedPassword,
		FullName:     util.RandomOwner(),
		Email:        util.RandomEmail(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	var data, err = io.ReadAll(body)
	require.NoError(t, err)

	var gotUser UserResponse
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
}

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	var arg, ok = x.(db.CreateUserParams)
	if !ok {
		return false
	}

	var err = util.CheckPassword(e.password, arg.HashPassword)
	if err != nil {
		return false
	}

	e.arg.HashPassword = arg.HashPassword
	return reflect.DeepEqual(e.arg, arg)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	var user, password = randomUser(t)
	var testCases = []struct {
		name          string
		body          createUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				var arg = db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
				requireBodyMatchUser(t, recoder.Body, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var ctrl = gomock.NewController(t)
			defer ctrl.Finish()
			var store = mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			var data, err = json.Marshal(tc.body)
			require.NoError(t, err)
			var url = "/users"
			var request, err2 = http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err2)

			var server = NewServer(store)
			var recoder = httptest.NewRecorder()
			server.router.ServeHTTP(recoder, request)

			tc.checkResponse(recoder)
		})
	}
}

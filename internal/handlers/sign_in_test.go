package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	market "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/internal/service"
	"gophermarket/internal/service/auth"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SignIn(t *testing.T) {

	const (
		salt     = "ManyWivedJackLaughsAtProbesOfSexQuiz"
		tokenKey = "PaintTheTownRed"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		contentType string
		user        market.User
		wantStatus  int
		callRepo    bool
		errRepo     error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success sign up",
			args: args{
				contentType: "application/json",
				user: market.User{
					Username: "user",
					Password: "hard_password",
				},
				wantStatus: http.StatusOK,
				callRepo:   true,
				errRepo:    nil,
			},
		},
		{
			name: "sign up existing user",
			args: args{
				contentType: "application/json",
				user: market.User{
					Username: "user",
					Password: "hard_password",
				},
				wantStatus: http.StatusConflict,
				callRepo:   true,
				errRepo:    pkg.ErrUserAlreadyExists,
			},
		},
		{
			name: "sign up empty auth data",
			args: args{
				contentType: "application/json",
				wantStatus:  http.StatusBadRequest,
				callRepo:    false,
				errRepo:     nil,
			},
		},
		{
			name: "sign up empty auth login",
			args: args{
				contentType: "application/json",
				user: market.User{
					Username: "user",
				},
				wantStatus: http.StatusBadRequest,
				callRepo:   false,
				errRepo:    nil,
			},
		},
		{
			name: "sign up empty auth password",
			args: args{
				contentType: "application/json",
				user: market.User{
					Password: "hard_password",
				},
				wantStatus: http.StatusBadRequest,
				callRepo:   false,
				errRepo:    nil,
			},
		},
		{
			name: "sign up invalid content-type",
			args: args{
				contentType: "text/json",
				user: market.User{
					Password: "hard_password",
				},
				wantStatus: http.StatusUnsupportedMediaType,
				callRepo:   false,
				errRepo:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			authRepoMock := repository.NewMockAuthorization(ctrl)
			mockRepo := repository.Repository{
				Authorization: authRepoMock,
			}

			if tt.args.callRepo {

				userRepo := market.User{
					Username: tt.args.user.Username,
					Password: tt.args.user.Password,
				}

				if len(tt.args.user.Password) > 0 {
					hash, err := pkg.GeneratePasswordHash(tt.args.user.Password, salt)
					require.NoError(t, err)

					userRepo.Password = hash
				}

				authRepoMock.EXPECT().ID(userRepo).Return(int64(0), tt.args.errRepo)
			}

			authService := auth.NewService(&mockRepo, salt, logpack.NewLogger())
			services := service.Service{
				Auth: authService,
			}
			handlers := NewHandler(&services, tokenKey, logpack.NewLogger())

			authBody, err := json.Marshal(tt.args.user)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(authBody))
			request.Header.Set("Content-Type", tt.args.contentType)

			w := httptest.NewRecorder()
			handlers.SignIn(w, request)

			response := w.Result()
			errClose := response.Body.Close()
			assert.NoError(t, errClose)

			assert.Equal(t, tt.args.wantStatus, response.StatusCode)
		})
	}
}

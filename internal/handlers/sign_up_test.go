package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	market "gophermarket/internal"
	"gophermarket/internal/service"
	"gophermarket/internal/service/auth"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SignUp(t *testing.T) {

	const (
		tokenKey = "PaintTheTownRed"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx         context.Context
		contentType string
		user        market.User
		userService bool
	}

	type want struct {
		status     int
		errService error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success sign up",
			args: args{
				ctx:         context.Background(),
				contentType: "application/json",
				user: market.User{
					Username: "user",
					Password: "hard_password",
				},
				userService: true,
			},
			want: want{
				status:     http.StatusOK,
				errService: nil,
			},
		},
		{
			name: "sign up existing user",
			args: args{
				ctx:         context.Background(),
				contentType: "application/json",
				user: market.User{
					Username: "user",
					Password: "hard_password",
				},
				userService: true,
			},
			want: want{
				status:     http.StatusConflict,
				errService: pkg.ErrUserAlreadyExists,
			},
		},
		{
			name: "sign up empty auth data",
			args: args{
				ctx:         context.Background(),
				contentType: "application/json",
				userService: true,
			},
			want: want{
				status:     http.StatusBadRequest,
				errService: pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up empty auth login",
			args: args{
				ctx:         context.Background(),
				contentType: "application/json",
				user: market.User{
					Username: "user",
				},
				userService: true,
			},
			want: want{
				status:     http.StatusBadRequest,
				errService: pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up empty auth password",
			args: args{
				ctx:         context.Background(),
				contentType: "application/json",
				user: market.User{
					Password: "hard_password",
				},
				userService: true,
			},
			want: want{
				status:     http.StatusBadRequest,
				errService: pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up invalid content-type",
			args: args{
				ctx:         context.Background(),
				contentType: "text/plain",
				user: market.User{
					Username: "user",
					Password: "hard_password",
				},
				userService: false,
			},
			want: want{
				status:     http.StatusUnsupportedMediaType,
				errService: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			authServiceMock := auth.NewMockServiceAuth(ctrl)

			if tt.args.userService {
				authServiceMock.EXPECT().SignUp(tt.args.ctx, tt.args.user).Return(tt.want.errService)
			}

			services := service.Service{
				Auth: authServiceMock,
			}

			handlers := NewHandler(&services, tokenKey, logpack.NewLogger())

			authBody, err := json.Marshal(tt.args.user)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewReader(authBody))
			request.Header.Set("Content-Type", tt.args.contentType)

			w := httptest.NewRecorder()
			handlers.SignUp(w, request)

			response := w.Result()
			errClose := response.Body.Close()
			assert.NoError(t, errClose)

			assert.Equal(t, tt.want.status, response.StatusCode)
		})
	}
}

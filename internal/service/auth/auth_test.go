package auth

import (
	"context"
	"testing"

	market "gophermarket/internal"
	"gophermarket/internal/repository"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	salt = "FiveWineExpertsJokinglyQuizzedSampleChablis"
)

func TestAuth_SignUp(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		userRepo      market.User
		userAuth      market.User
		waitErrSignUp error
		waitErrHash   error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success sign up",
			args: args{
				userRepo: market.User{
					Username: "user",
					Password: "user_password",
				},
				userAuth: market.User{
					Username: "user",
					Password: "user_password",
				},
			},
		},
		{
			name: "sign up without password",
			args: args{
				userRepo: market.User{
					Username: "user",
				},
				userAuth: market.User{
					Username: "user",
				},
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up without auth data",
			args: args{
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up without username",
			args: args{
				userRepo: market.User{
					Password: "user_password",
				},
				userAuth: market.User{
					Password: "user_password",
				},
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash, errHash := pkg.GeneratePasswordHash(tt.args.userRepo.Password, salt)
			require.Equal(t, tt.args.waitErrHash, errHash)
			tt.args.userRepo.Password = hash

			authRepoMock := repository.NewMockAuthorization(ctrl)
			ctx := context.Background()

			if tt.args.waitErrSignUp == nil {
				authRepoMock.EXPECT().Create(ctx, tt.args.userRepo).Return(nil)
			}

			repo := repository.Repository{
				Authorization: authRepoMock,
			}

			authService := NewService(&repo, salt, logpack.NewLogger())

			err := authService.SignUp(ctx, tt.args.userAuth)
			assert.Equal(t, err, tt.args.waitErrSignUp)
		})
	}

}

func TestAuth_SignIn(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		userRepo      market.User
		userAuth      market.User
		waitErrSignUp error
		waitErrHash   error
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success sign up",
			args: args{
				userRepo: market.User{
					Username: "user",
					Password: "user_password",
				},
				userAuth: market.User{
					Username: "user",
					Password: "user_password",
				},
			},
		},
		{
			name: "sign up without password",
			args: args{
				userRepo: market.User{
					Username: "user",
				},
				userAuth: market.User{
					Username: "user",
				},
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up without auth data",
			args: args{
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   pkg.ErrEmptyAuthData,
			},
		},
		{
			name: "sign up without username",
			args: args{
				userRepo: market.User{
					Password: "user_password",
				},
				userAuth: market.User{
					Password: "user_password",
				},
				waitErrSignUp: pkg.ErrEmptyAuthData,
				waitErrHash:   nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash, errHash := pkg.GeneratePasswordHash(tt.args.userRepo.Password, salt)
			require.Equal(t, tt.args.waitErrHash, errHash)
			tt.args.userRepo.Password = hash

			authRepoMock := repository.NewMockAuthorization(ctrl)
			ctx := context.Background()

			if tt.args.waitErrSignUp == nil {
				authRepoMock.EXPECT().ID(ctx, tt.args.userRepo).Return(int64(0), nil)
			}

			repo := repository.Repository{
				Authorization: authRepoMock,
			}

			authService := NewService(&repo, salt, logpack.NewLogger())

			err := authService.SignIn(ctx, tt.args.userAuth)
			assert.Equal(t, err, tt.args.waitErrSignUp)
		})
	}
}

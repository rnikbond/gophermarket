package auth

import (
	"testing"

	mock "gophermarket/mocks/pkg/repository"
	market "gophermarket/pkg"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {

	userMock := market.User{
		Username: "user",
		Password: "test_password",
	}
	hash, _ := GeneratePasswordHash(userMock.Password)
	userMock.Password = hash

	userAuth := market.User{
		Username: "user",
		Password: "test_password",
	}

	ctrl := gomock.NewController(t)

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().SignIn(userMock).Return(market.ErrUserNotFound)
	mockRepo.EXPECT().SignUp(userMock).Return(nil)
	mockRepo.EXPECT().SignIn(userMock).Return(nil)

	authService := NewService(mockRepo)

	token, err := authService.SignIn(userAuth)
	assert.Equal(t, err, market.ErrUserNotFound)
	assert.Equal(t, token, "")

	token, err = authService.SignUp(userAuth)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token, err = authService.SignIn(userAuth)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

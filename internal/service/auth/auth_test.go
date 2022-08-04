package auth

import (
	"testing"
)

/*
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
	mockRepo.EXPECT().SignIn(userMock).Return(pkg.ErrUserNotFound)
	mockRepo.EXPECT().SignUp(userMock).Return(nil)
	mockRepo.EXPECT().SignIn(userMock).Return(nil)

	authService := NewService(mockRepo)

	token, err := authService.SignIn(userAuth)
	assert.Equal(t, err, pkg.ErrUserNotFound)
	assert.Equal(t, token, "")

	token, err = authService.SignUp(userAuth)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token, err = authService.SignIn(userAuth)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
*/

func TestGeneratePasswordHash(t *testing.T) {

	/*
		tests := []struct {
			name     string
			password string
			wantErr  error
		}{
			{
				name:     "Valid password",
				password: "PwdQwerty",
				wantErr:  nil,
			},
			{
				name:    "Empty password",
				wantErr: market.ErrEmptyAuthData,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				hash, err := GeneratePasswordHash(tt.password)

				if tt.wantErr != nil {
					assert.Equal(t, tt.wantErr, err)
				} else {
					require.NoError(t, err)
					assert.NotEmpty(t, hash)
				}
			})
		}

	*/
}

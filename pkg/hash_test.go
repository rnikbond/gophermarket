package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePasswordHash(t *testing.T) {

	salt := "PackMyBoxWithFiveDozenLiquorJugs"

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "valid password",
			password: "PwdQwerty",
			wantErr:  nil,
		},
		{
			name:    "empty password",
			wantErr: ErrEmptyAuthData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			hash, err := GeneratePasswordHash(tt.password, salt)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}

}

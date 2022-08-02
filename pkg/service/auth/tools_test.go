package auth

import (
	"testing"

	market "gophermarket/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeneratePasswordHash(t *testing.T) {
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
}

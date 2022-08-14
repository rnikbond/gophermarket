package loyalty

import (
	"context"
	"testing"

	"gophermarket/internal/repository"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoyalty_Balance(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logpack.NewLogger()

	type args struct {
		ctx      context.Context
		username string
	}

	type want struct {
		balance Balance
		err     error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success check balance",
			args: args{
				ctx:      context.Background(),
				username: "user_test",
			},
			want: want{
				balance: Balance{
					Accrual:   123.23,
					Withdrawn: 0,
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loyaltyRepoMock := repository.NewMockLoyalty(ctrl)
			ctx := context.Background()

			if tt.want.err == nil {
				loyaltyRepoMock.EXPECT().HowMatchAvailable(ctx, tt.args.username).Return(tt.want.balance.Accrual, nil)
				loyaltyRepoMock.EXPECT().HowMatchUsed(ctx, tt.args.username).Return(tt.want.balance.Withdrawn, nil)
			}

			repo := repository.Repository{
				Loyalty: loyaltyRepoMock,
			}

			loyaltyService := NewService(&repo, logger)

			balance, err := loyaltyService.Balance(ctx, tt.args.username)
			require.Equal(t, err, tt.want.err)

			assert.Equal(t, balance, tt.want.balance)
		})
	}
}

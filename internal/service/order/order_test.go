package order

import (
	"context"
	"testing"

	"gophermarket/internal/repository"
	market "gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrder_Create(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logpack.NewLogger()

	type args struct {
		ctx      context.Context
		username string
		order    int64
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success create order",
			args: args{
				ctx:      context.Background(),
				username: "user_test",
				order:    417147,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "invalid order",
			args: args{
				ctx:      context.Background(),
				username: "user_test",
				order:    123123,
			},
			want: want{
				err: market.ErrInvalidOrderNumber,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderRepoMock := repository.NewMockOrder(ctrl)
			ctx := context.Background()

			if tt.want.err == nil {
				orderRepoMock.EXPECT().Create(ctx, tt.args.order, tt.args.username).Return(nil)
			}

			repo := repository.Repository{
				Order: orderRepoMock,
			}

			orderService := NewService(&repo, logger)

			err := orderService.Create(ctx, tt.args.order, tt.args.username)
			assert.Equal(t, err, tt.want.err)
		})
	}
}

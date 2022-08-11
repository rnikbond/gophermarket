package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophermarket/internal/service"
	"gophermarket/internal/service/loyalty"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Balance(t *testing.T) {
	const (
		tokenKey = "PaintTheTownRed"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx         context.Context
		username    string
		userService bool
	}

	type want struct {
		balance    loyalty.Balance
		status     int
		errService error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success get balance",
			args: args{
				ctx:         context.Background(),
				username:    "user_test",
				userService: true,
			},
			want: want{
				balance: loyalty.Balance{
					Accrual:   153.23,
					Withdrawn: 0,
				},
				status:     http.StatusOK,
				errService: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			loyaltyServiceMock := loyalty.NewMockServiceLoyalty(ctrl)
			if tt.args.userService {
				loyaltyServiceMock.EXPECT().Balance(tt.args.ctx, tt.args.username).Return(tt.want.balance, tt.want.errService)
			}

			services := service.Service{
				Loyalty: loyaltyServiceMock,
			}

			handlers := NewHandler(&services, tokenKey, logpack.NewLogger())

			request := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			request.SetBasicAuth(tt.args.username, "")

			w := httptest.NewRecorder()
			handlers.Balance(w, request)

			response := w.Result()

			require.Equal(t, response.StatusCode, tt.want.status)

			data, errRead := io.ReadAll(response.Body)
			errBody := response.Body.Close()
			assert.NoError(t, errBody)

			require.NoError(t, errRead)

			var balance loyalty.Balance
			errJSON := json.Unmarshal(data, &balance)
			require.NoError(t, errJSON)

			assert.Equal(t, balance, tt.want.balance)
		})
	}
}

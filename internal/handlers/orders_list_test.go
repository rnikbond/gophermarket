package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gophermarket/internal/service"
	"gophermarket/internal/service/order"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_OrdersList(t *testing.T) {

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
		orders     []pkg.OrderInfo
		status     int
		errService error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "success get orders",
			args: args{
				ctx:         context.Background(),
				username:    "user_test",
				userService: true,
			},
			want: want{
				orders: []pkg.OrderInfo{
					{
						Order:      "417147",
						Status:     pkg.StatusProcessed,
						Accrual:    50.65,
						UploadedAt: time.Now().Format(time.RFC3339),
					},
					{
						Order:      "951913",
						Status:     pkg.StatusProcessed,
						UploadedAt: time.Now().Format(time.RFC3339),
					},
				},
				status:     http.StatusOK,
				errService: nil,
			},
		},
		{
			name: "success orders - user no orders",
			args: args{
				ctx:         context.Background(),
				username:    "user_test",
				userService: true,
			},
			want: want{
				orders:     []pkg.OrderInfo{},
				status:     http.StatusNoContent,
				errService: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			orderServiceMock := order.NewMockServiceOrder(ctrl)
			if tt.args.userService {
				orderServiceMock.EXPECT().UserOrders(tt.args.ctx, tt.args.username).Return(tt.want.orders, tt.want.errService)
			}

			services := service.Service{
				Order: orderServiceMock,
			}

			handlers := NewHandler(&services, tokenKey, logpack.NewLogger())

			request := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			request.SetBasicAuth(tt.args.username, "")

			w := httptest.NewRecorder()
			handlers.OrdersList(w, request)

			response := w.Result()

			require.Equal(t, response.StatusCode, tt.want.status)

			if response.StatusCode == http.StatusNoContent {
				return
			}

			data, errRead := io.ReadAll(response.Body)
			errBody := response.Body.Close()
			assert.NoError(t, errBody)

			require.NoError(t, errRead)

			var orders []pkg.OrderInfo
			errJSON := json.Unmarshal(data, &orders)
			require.NoError(t, errJSON)

			assert.Equal(t, orders, tt.want.orders)
		})
	}
}

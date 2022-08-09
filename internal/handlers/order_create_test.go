package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gophermarket/internal/service"
	"gophermarket/internal/service/order"
	"gophermarket/pkg"
	"gophermarket/pkg/logpack"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateOrder(t *testing.T) {
	const (
		tokenKey = "PaintTheTownRed"
	)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		ctx         context.Context
		contentType string
		username    string
		orderNum    int64
		body        io.Reader
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
			name: "success create order",
			args: args{
				ctx:         context.Background(),
				contentType: "text/plain",
				username:    "user_test",
				orderNum:    417147,
				body:        bytes.NewReader([]byte("417147")),
				userService: true,
			},
			want: want{
				status:     http.StatusAccepted,
				errService: nil,
			},
		},
		{
			name: "failed create order - invalid order number",
			args: args{
				ctx:         context.Background(),
				contentType: "text/plain",
				username:    "user_test",
				orderNum:    111000111,
				body:        bytes.NewReader([]byte("111000111")),
				userService: true,
			},
			want: want{
				status:     http.StatusUnprocessableEntity,
				errService: pkg.ErrInvalidOrderNumber,
			},
		},
		{
			name: "failed create order - invalid Content-Type",
			args: args{
				ctx:         context.Background(),
				username:    "user_test",
				orderNum:    111000111,
				body:        bytes.NewReader([]byte("111000111")),
				userService: false,
			},
			want: want{
				status: http.StatusUnsupportedMediaType,
			},
		},
		{
			name: "failed create order - no order",
			args: args{
				ctx:         context.Background(),
				contentType: "text/plain",
				username:    "user_test",
				orderNum:    111000111,
				body:        nil,
				userService: false,
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			orderServiceMock := order.NewMockServiceOrder(ctrl)
			if tt.args.userService {
				orderServiceMock.EXPECT().Create(tt.args.ctx, tt.args.orderNum, tt.args.username).Return(tt.want.errService)
			}

			services := service.Service{
				Order: orderServiceMock,
			}

			handler := NewHandler(&services, tokenKey, logpack.NewLogger())

			request := httptest.NewRequest(http.MethodPost, "/api/user/orders", tt.args.body)
			request.SetBasicAuth(tt.args.username, "")
			request.Header.Set("Content-Type", tt.args.contentType)

			w := httptest.NewRecorder()
			handler.CreateOrder(w, request)

			response := w.Result()
			assert.Equal(t, response.StatusCode, tt.want.status)
		})
	}
}

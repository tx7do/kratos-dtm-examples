package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type PaymentService struct {
	shopV1.UnimplementedPaymentServiceServer

	log *log.Helper
}

func NewPaymentService(logger log.Logger) *PaymentService {
	l := log.NewHelper(log.With(logger, "module", "payment/service/shop-service"))
	return &PaymentService{
		log: l,
	}
}

func (s *PaymentService) TryMakePayment(ctx context.Context, req *shopV1.TryMakePaymentRequest) (*shopV1.TryMakePaymentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TryMakePayment not implemented")
}

func (s *PaymentService) ConfirmMakePayment(ctx context.Context, req *shopV1.ConfirmMakePaymentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmMakePayment not implemented")
}

func (s *PaymentService) CancelMakePayment(ctx context.Context, req *shopV1.CancelMakePaymentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelMakePayment not implemented")
}

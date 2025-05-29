package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"
)

type BankService struct {
	bankV1.UnimplementedBankServiceServer

	log *log.Helper
}

func NewBankService(logger log.Logger) *BankService {
	l := log.NewHelper(log.With(logger, "module", "bank/service/bank-service"))
	return &BankService{
		log: l,
	}
}

func (s *BankService) TransIn(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("transfer in %d cents to %d", req.Amount, req.UserId)
	return &bankV1.TransferResponse{}, nil
}

func (s *BankService) TransOut(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("transfer out %d cents from %d", req.Amount, req.UserId)
	return &bankV1.TransferResponse{}, nil
}

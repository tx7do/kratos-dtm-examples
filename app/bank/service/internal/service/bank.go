package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

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

func (s *BankService) GetAccount(ctx context.Context, req *bankV1.GetAccountRequest) (*bankV1.Account, error) {
	return &bankV1.Account{}, nil
}

func (s *BankService) TransIn(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("transfer in %d cents from %s to %s", req.Amount, req.FromAccountId, req.ToAccountId)
	return &bankV1.TransferResponse{}, nil
}

func (s *BankService) TransOut(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("transfer out %d cents from %s to %s", req.Amount, req.FromAccountId, req.ToAccountId)
	return &bankV1.TransferResponse{}, nil
}

// TryDeduct TCC - Try阶段：预扣款，冻结金额
func (s *BankService) TryDeduct(ctx context.Context, req *bankV1.TransactionRequest) (*bankV1.TransactionResponse, error) {
	s.log.Infof("try deduct %d cents from %s", req.Amount, req.AccountId)
	return &bankV1.TransactionResponse{}, nil
}

// ConfirmDeduct TCC - Confirm阶段：确认扣款，实际扣除金额
func (s *BankService) ConfirmDeduct(ctx context.Context, req *bankV1.TransactionRequest) (*emptypb.Empty, error) {
	s.log.Infof("confirm deduct %d cents from %s", req.Amount, req.AccountId)
	return &emptypb.Empty{}, nil
}

// CancelDeduct TCC - Cancel阶段：取消扣款，解冻金额
func (s *BankService) CancelDeduct(ctx context.Context, req *bankV1.TransactionRequest) (*emptypb.Empty, error) {
	s.log.Infof("cancel deduct %d cents from %s", req.Amount, req.AccountId)
	return &emptypb.Empty{}, nil
}

// Deduct 扣款操作（SAGA正向操作）
func (s *BankService) Deduct(ctx context.Context, req *bankV1.TransactionRequest) (*bankV1.TransactionResponse, error) {
	s.log.Infof("deduct %d cents from %s", req.Amount, req.AccountId)
	return &bankV1.TransactionResponse{}, nil
}

// Refund 退款操作（SAGA反向操作）
func (s *BankService) Refund(ctx context.Context, req *bankV1.TransactionRequest) (*bankV1.TransactionResponse, error) {
	s.log.Infof("refund %d cents to %s", req.Amount, req.AccountId)
	return &bankV1.TransactionResponse{}, nil
}

// Transfer 转账操作（SAGA正向操作）
func (s *BankService) Transfer(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("transfer %d cents from %s to %s", req.Amount, req.FromAccountId, req.ToAccountId)
	return &bankV1.TransferResponse{}, nil
}

// ReverseTransfer 反转转账操作（SAGA反向操作）
func (s *BankService) ReverseTransfer(ctx context.Context, req *bankV1.TransferRequest) (*bankV1.TransferResponse, error) {
	s.log.Infof("reverse transfer %d cents from %s to %s", req.Amount, req.ToAccountId, req.FromAccountId)
	return &bankV1.TransferResponse{}, nil
}

package service

import (
	"context"

	"github.com/dtm-labs/client/dtmgrpc"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"kratos-dtm-examples/app/shop/service/internal/data"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"
	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/pkg/service"
)

var (
	dtmServer  = service.MakeDiscoveryAddress(service.DTMService)
	bankServer = service.MakeDiscoveryAddress(service.BankService)
)

type ShopService struct {
	shopV1.UnimplementedShopServiceServer

	log *log.Helper

	bankServiceClient bankV1.BankServiceClient
}

func NewShopService(logger log.Logger, _ *data.Data, bankServiceClient bankV1.BankServiceClient) *ShopService {
	return &ShopService{
		log:               log.NewHelper(log.With(logger, "module", "shop/service/shop-service")),
		bankServiceClient: bankServiceClient,
	}
}

func (s *ShopService) Buy(_ context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	s.log.Infof("buy %d items of product %d for user %d", req.Quantity, req.ProductId, req.UserId)
	return &shopV1.BuyResponse{}, nil
}

func (s *ShopService) TestTP(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	m := dtmgrpc.NewMsgGrpc(dtmServer, gid).
		Add(bankServer+bankV1.BankService_TransOut_FullMethodName, &bankV1.TransferRequest{Amount: 30, FromAccountId: "from", ToAccountId: "to"}).
		Add(bankServer+bankV1.BankService_TransIn_FullMethodName, &bankV1.TransferRequest{Amount: 30, FromAccountId: "from", ToAccountId: "to"})
	m.WaitResult = true
	err := m.Submit()
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, bankV1.ErrorInternalServerError(err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestTCC(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	s.log.Infof("testTCC %s", gid)

	err := dtmgrpc.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmgrpc.TccGrpc) error {
		err := tcc.CallBranch(
			&bankV1.TransactionRequest{Amount: 30},
			bankServer+bankV1.BankService_TryDeduct_FullMethodName,
			bankServer+bankV1.BankService_ConfirmDeduct_FullMethodName,
			bankServer+bankV1.BankService_CancelDeduct_FullMethodName,
			&bankV1.TransactionResponse{},
		)
		if err != nil {
			s.log.Errorf("failed to call branch for deduct: %v", err)
			return bankV1.ErrorInternalServerError(err.Error())
		}
		return nil
	})
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, bankV1.ErrorInternalServerError(err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestSAGA(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	req := &bankV1.TransactionRequest{Amount: 30}

	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(bankServer+bankV1.BankService_Deduct_FullMethodName, bankServer+bankV1.BankService_Refund_FullMethodName, req).
		Add(bankServer+bankV1.BankService_Transfer_FullMethodName, bankServer+bankV1.BankService_ReverseTransfer_FullMethodName, req)

	err := saga.Submit()
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, bankV1.ErrorInternalServerError(err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestXA(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	err := dtmgrpc.XaGlobalTransaction(dtmServer, gid, func(xa *dtmgrpc.XaGrpc) error {
		if err := xa.CallBranch(
			&bankV1.TransactionRequest{Amount: 30},
			bankServer+bankV1.BankService_Deduct_FullMethodName,
			&bankV1.TransactionResponse{},
		); err != nil {
			s.log.Errorf("failed to call branch for transOutXa: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, bankV1.ErrorInternalServerError(err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestWorkFlow(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

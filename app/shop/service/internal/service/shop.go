package service

import (
	"context"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"

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
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestTCC(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	err := dtmcli.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmcli.Tcc) (*resty.Response, error) {
		resp, err := tcc.CallBranch(
			&bankV1.TransactionRequest{Amount: 30},
			bankServer+bankV1.BankService_TryDeduct_FullMethodName,
			bankServer+bankV1.BankService_ConfirmDeduct_FullMethodName,
			bankServer+bankV1.BankService_CancelDeduct_FullMethodName,
		)
		if err != nil {
			s.log.Errorf("failed to call branch for deduct: %v", err)
			return resp, err
		}
		return resp, nil
	})
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *ShopService) TestSAGA(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	req := &bankV1.TransactionRequest{Amount: 30}

	saga := dtmcli.NewSaga(dtmServer, uuid.New().String()).
		Add(bankServer+bankV1.BankService_Deduct_FullMethodName, bankServer+bankV1.BankService_Refund_FullMethodName, req).
		Add(bankServer+bankV1.BankService_Transfer_FullMethodName, bankServer+bankV1.BankService_ReverseTransfer_FullMethodName, req)
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	err := saga.Submit()
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

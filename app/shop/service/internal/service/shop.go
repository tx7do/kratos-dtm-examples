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

func (s *ShopService) TestTransaction(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	gid := dtmgrpc.MustGenGid(dtmServer)

	m := dtmgrpc.NewMsgGrpc(dtmServer, gid).
		Add(bankServer+bankV1.BankService_TransOut_FullMethodName, &bankV1.TransferRequest{Amount: 30, UserId: 1}).
		Add(bankServer+bankV1.BankService_TransIn_FullMethodName, &bankV1.TransferRequest{Amount: 30, UserId: 2})
	m.WaitResult = true
	err := m.Submit()
	if err != nil {
		s.log.Errorf("failed to submit transaction: %v", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

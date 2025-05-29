package service

import (
	"context"

	"github.com/dtm-labs/client/dtmgrpc"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc/resolver/discovery"

	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/types/known/emptypb"

	bankV1 "kratos-dtm-examples/api/gen/go/bank/service/v1"
	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	_ "github.com/dtm-labs/driver-kratos"
)

const (
	dtmServer  = "discovery:///dtmservice"
	bankServer = "discovery:///bank-service"
)

type ShopService struct {
	shopV1.UnimplementedShopServiceServer

	log *log.Helper
}

func NewShopService(logger log.Logger, rr registry.Discovery) *ShopService {
	resolver.Register(discovery.NewBuilder(rr, discovery.WithInsecure(true)))

	return &ShopService{
		log: log.NewHelper(log.With(logger, "module", "shop/service/shop-service")),
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

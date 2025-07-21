package service

import (
	"context"

	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/dtm-labs/client/dtmgrpc/dtmgpb"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/pkg/service"
)

var (
	dtmServer  = service.MakeDiscoveryAddress(service.DTMService)
	shopServer = service.MakeDiscoveryAddress(service.ShopService)
)

type ShopService struct {
	shopV1.UnimplementedShopServiceServer

	log *log.Helper

	dtmClient dtmgpb.DtmClient
}

func NewShopService(
	logger log.Logger,
) *ShopService {
	return &ShopService{
		log: log.NewHelper(log.With(logger, "module", "shop/service/shop-service")),
	}
}

func (s *ShopService) Buy(_ context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	s.log.Infof("buy %d items of product %d for user %d", req.Quantity, req.ProductId, req.UserId)
	return &shopV1.BuyResponse{}, nil
}

func (s *ShopService) TestTP(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(dtmServer)

	requestId = gid // 使用 gid 作为 request_id

	// 创建消息事务
	msg := dtmgrpc.NewMsgGrpc(dtmServer, gid).
		Add(
			shopServer+shopV1.StockService_DeductStock_FullMethodName,
			&shopV1.DeductStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
			},
		).
		Add(
			shopServer+shopV1.OrderService_CreateOrder_FullMethodName,
			&shopV1.CreateOrderRequest{
				UserId:    req.UserId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
				OrderNo:   requestId, // 简化使用 requestId 作为订单号
			},
		)

	msg.WaitResult = true

	// 提交事务
	if err := msg.Submit(); err != nil {
		s.log.Errorf("提交购买事务失败: %v", err)
		return nil, shopV1.ErrorInternalServerError(err.Error())
	}

	s.log.Infof("购买事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestTCC(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(dtmServer)

	requestId = gid // 使用 gid 作为 request_id

	s.log.Infof("开始 TCC 事务，GID: %s", gid)

	var err error

	err = dtmgrpc.TccGlobalTransaction(dtmServer, gid, func(tcc *dtmgrpc.TccGrpc) error {
		// Try 阶段：扣减库存
		//err = tcc.CallBranch(
		//	&shopV1.TryDeductStockRequest{
		//		ProductId: req.ProductId,
		//		Quantity:  req.Quantity,
		//		RequestId: requestId,
		//	},
		//	shopServer+shopV1.StockService_TryDeductStock_FullMethodName,
		//	shopServer+shopV1.StockService_ConfirmDeductStock_FullMethodName,
		//	shopServer+shopV1.StockService_CancelDeductStock_FullMethodName,
		//	&shopV1.StockResponse{},
		//)
		//if err != nil {
		//	s.log.Errorf("扣减库存失败: %v", err)
		//	return shopV1.ErrorInternalServerError("扣减库存失败")
		//}

		// Try 阶段：创建订单
		err = tcc.CallBranch(
			&shopV1.TryCreateOrderRequest{
				UserId:    req.UserId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
				OrderNo:   requestId, // 简化使用 requestId 作为订单号
			},
			shopServer+shopV1.OrderService_TryCreateOrder_FullMethodName,
			shopServer+shopV1.OrderService_ConfirmCreateOrder_FullMethodName,
			shopServer+shopV1.OrderService_CancelCreateOrder_FullMethodName,
			&shopV1.OrderResponse{},
		)
		if err != nil {
			s.log.Errorf("TCC创建订单失败: %v", err)
			return shopV1.ErrorInternalServerError("创建订单失败")
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("TCC 事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError(err.Error())
	}

	s.log.Infof("TCC 事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestSAGA(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(dtmServer)

	requestId = gid // 使用 gid 作为 request_id

	s.log.Infof("开始 SAGA 事务，GID: %s", gid)

	saga := dtmgrpc.NewSagaGrpc(dtmServer, gid).
		Add(
			shopServer+shopV1.StockService_DeductStock_FullMethodName,
			shopServer+shopV1.StockService_RefundStock_FullMethodName,
			&shopV1.DeductStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
			},
		).
		Add(
			shopServer+shopV1.OrderService_CreateOrder_FullMethodName,
			shopServer+shopV1.OrderService_RefundOrder_FullMethodName,
			&shopV1.CreateOrderRequest{
				UserId:    req.UserId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
				OrderNo:   requestId, // 简化使用 requestId 作为订单号
			},
		)

	if err := saga.Submit(); err != nil {
		s.log.Errorf("SAGA 事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError(err.Error())
	}

	s.log.Infof("SAGA 事务提交成功，GID: %s", gid)
	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestXA(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(dtmServer)

	requestId = gid // 使用 gid 作为 request_id

	err := dtmgrpc.XaGlobalTransaction(dtmServer, gid, func(xa *dtmgrpc.XaGrpc) error {
		// 扣减库存
		if err := xa.CallBranch(
			&shopV1.DeductStockRequest{ProductId: req.ProductId, Quantity: req.Quantity},
			shopServer+shopV1.StockService_DeductStock_FullMethodName,
			&emptypb.Empty{},
		); err != nil {
			s.log.Errorf("扣减库存失败: %v", err)
			return err
		}

		// 创建订单
		if err := xa.CallBranch(
			&shopV1.CreateOrderRequest{
				UserId:    req.UserId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
				OrderNo:   requestId, // 简化使用 requestId 作为订单号
			},
			shopServer+shopV1.OrderService_CreateOrder_FullMethodName,
			&emptypb.Empty{},
		); err != nil {
			s.log.Errorf("XA创建订单失败: %v", err)
			return shopV1.ErrorInternalServerError("创建订单失败")
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("XA 事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError(err.Error())
	}

	s.log.Infof("XA 事务提交成功，GID: %s", gid)
	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestWorkFlow(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(dtmServer)

	s.log.Infof("开始工作流事务，GID: %s", gid)

	//workflow.InitGrpc(dtmServer, shopServer, gsvr)

	s.log.Infof("工作流事务提交成功，GID: %s", gid)
	return &shopV1.BuyResponse{Success: true}, nil
}

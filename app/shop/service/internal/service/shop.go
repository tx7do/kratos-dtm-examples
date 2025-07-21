package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"

	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/dtm-labs/client/dtmgrpc/dtmgimp"
	"github.com/dtm-labs/client/workflow"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/pkg/service"
)

const (
	WorkflowShopServiceOrderSAGA  = "test_workflow_shop_order_saga"
	WorkflowShopServiceOrderTCC   = "test_workflow_shop_order_tcc"
	WorkflowShopServiceOrderXA    = "test_workflow_shop_order_xa"
	WorkflowShopServiceOrderMixed = "test_workflow_shop_order_mixed"
)

type ShopService struct {
	shopV1.UnimplementedShopServiceServer

	log *log.Helper

	stockService   *StockService
	orderService   *OrderService
	paymentService *PaymentService
}

func NewShopService(
	logger log.Logger,
	stockService *StockService,
	orderService *OrderService,
	paymentService *PaymentService,
) *ShopService {
	svc := &ShopService{
		log:            log.NewHelper(log.With(logger, "module", "shop/service/shop-service")),
		stockService:   stockService,
		orderService:   orderService,
		paymentService: paymentService,
	}

	svc.init()

	return svc
}

func (s *ShopService) init() {
	var err error

	// SAGA工作流注册
	err = workflow.Register(WorkflowShopServiceOrderSAGA, func(wf *workflow.Workflow, data []byte) error {

		var codec = encoding.GetCodec("proto")

		var req shopV1.BuyRequest
		if len(data) > 0 {
			if err = codec.Unmarshal(data, &req); err != nil {
				s.log.Errorf("工作流数据反序列化失败: %v", err)
				return shopV1.ErrorInternalServerError("工作流数据反序列化失败")
			}
		}

		// 扣减库存步骤
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.stockService.RefundStock(wf.Context, &shopV1.RefundStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
			}); err != nil {
				s.log.Errorf("工作流回滚扣减库存失败: %v", err)
				return shopV1.ErrorInternalServerError("工作流回滚扣减库存失败")
			}

			return nil
		})
		if _, err = s.stockService.DeductStock(wf.Context, &shopV1.DeductStockRequest{
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
			RequestId: wf.Gid,
		}); err != nil {
			s.log.Errorf("工作流扣减库存失败: %v", err)
			return shopV1.ErrorInternalServerError("工作流扣减库存失败")
		}

		// 创建订单步骤
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if _, err = s.orderService.RefundOrder(wf.Context, &shopV1.RefundOrderRequest{
				OrderNo: wf.Gid,
			}); err != nil {
				s.log.Errorf("工作流回滚创建订单失败: %v", err)
				return shopV1.ErrorInternalServerError("工作流回滚创建订单失败")
			}
			return nil
		})
		if _, err = s.orderService.CreateOrder(wf.Context, &shopV1.CreateOrderRequest{
			UserId:    req.UserId,
			ProductId: req.ProductId,
			Quantity:  req.Quantity,
			RequestId: wf.Gid,
			OrderNo:   wf.Gid,
		}); err != nil {
			s.log.Errorf("工作流创建订单失败: %v", err)
			return shopV1.ErrorInternalServerError("工作流创建订单失败")
		}

		return nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowShopServiceOrderSAGA, err)
		return
	}

	// TCC工作流注册
	err = workflow.Register(WorkflowShopServiceOrderTCC, func(wf *workflow.Workflow, data []byte) error {
		// TODO TCC工作流注册
		return nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowShopServiceOrderTCC, err)
		return
	}

	// XA工作流注册
	err = workflow.Register(WorkflowShopServiceOrderXA, func(wf *workflow.Workflow, data []byte) error {
		// TODO XA工作流注册
		return nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowShopServiceOrderXA, err)
		return
	}

	// 混合工作流注册
	err = workflow.Register(WorkflowShopServiceOrderMixed, func(wf *workflow.Workflow, data []byte) error {
		// TODO 混合工作流注册
		return nil
	})
	if err != nil {
		s.log.Errorf("工作流[%s] 注册失败: %v", WorkflowShopServiceOrderMixed, err)
		return
	}
}

func (s *ShopService) Buy(_ context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	s.log.Infof("buy %d items of product %d for user %d", req.Quantity, req.ProductId, req.UserId)
	return &shopV1.BuyResponse{}, nil
}

func (s *ShopService) TestTP(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	requestId = gid // 使用 gid 作为 request_id

	// 创建消息事务
	msg := dtmgrpc.NewMsgGrpc(service.DtmServerAddress, gid).
		Add(
			service.ShopServerAddress+shopV1.StockService_DeductStock_FullMethodName,
			&shopV1.DeductStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
			},
		).
		Add(
			service.ShopServerAddress+shopV1.OrderService_CreateOrder_FullMethodName,
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
		s.log.Errorf("提交 二阶消息 事务失败: %v", err)
		return nil, shopV1.ErrorInternalServerError(err.Error())
	}

	s.log.Infof("二阶消息 事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestTCC(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	var requestId string

	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	requestId = gid // 使用 gid 作为 request_id

	s.log.Infof("开始 TCC 事务，GID: %s", gid)

	var err error

	err = dtmgrpc.TccGlobalTransaction(service.DtmServerAddress, gid, func(tcc *dtmgrpc.TccGrpc) error {
		// Try 阶段：扣减库存
		err = tcc.CallBranch(
			&shopV1.TryDeductStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
			},
			service.ShopServerAddress+shopV1.StockService_TryDeductStock_FullMethodName,
			service.ShopServerAddress+shopV1.StockService_ConfirmDeductStock_FullMethodName,
			service.ShopServerAddress+shopV1.StockService_CancelDeductStock_FullMethodName,
			&shopV1.StockResponse{},
		)
		if err != nil {
			s.log.Errorf("扣减库存失败: %v", err)
			return shopV1.ErrorInternalServerError("扣减库存失败")
		}

		// Try 阶段：创建订单
		err = tcc.CallBranch(
			&shopV1.TryCreateOrderRequest{
				UserId:    req.UserId,
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
				OrderNo:   requestId, // 简化使用 requestId 作为订单号
			},
			service.ShopServerAddress+shopV1.OrderService_TryCreateOrder_FullMethodName,
			service.ShopServerAddress+shopV1.OrderService_ConfirmCreateOrder_FullMethodName,
			service.ShopServerAddress+shopV1.OrderService_CancelCreateOrder_FullMethodName,
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
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	requestId = gid // 使用 gid 作为 request_id

	s.log.Infof("开始 SAGA 事务，GID: %s", gid)

	saga := dtmgrpc.NewSagaGrpc(service.DtmServerAddress, gid).
		// 扣减库存
		Add(
			service.ShopServerAddress+shopV1.StockService_DeductStock_FullMethodName,
			service.ShopServerAddress+shopV1.StockService_RefundStock_FullMethodName,
			&shopV1.DeductStockRequest{
				ProductId: req.ProductId,
				Quantity:  req.Quantity,
				RequestId: requestId,
			},
		).
		// 创建订单
		Add(
			service.ShopServerAddress+shopV1.OrderService_CreateOrder_FullMethodName,
			service.ShopServerAddress+shopV1.OrderService_RefundOrder_FullMethodName,
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
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	requestId = gid // 使用 gid 作为 request_id

	err := dtmgrpc.XaGlobalTransaction(service.DtmServerAddress, gid, func(xa *dtmgrpc.XaGrpc) error {
		// 扣减库存
		if err := xa.CallBranch(
			&shopV1.DeductStockRequest{ProductId: req.ProductId, Quantity: req.Quantity},
			service.ShopServerAddress+shopV1.StockService_DeductStock_FullMethodName,
			&shopV1.StockResponse{},
		); err != nil {
			s.log.Errorf("XA扣减库存失败: %v", err)
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
			service.ShopServerAddress+shopV1.OrderService_CreateOrder_FullMethodName,
			&shopV1.OrderResponse{},
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

func (s *ShopService) TestWorkFlowSAGA(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	s.log.Infof("开始SAGA工作流事务，GID: %s", gid)

	// 提交工作流
	if err := workflow.Execute(WorkflowShopServiceOrderSAGA, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("SAGA工作流事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError("SAGA工作流事务提交失败")
	}

	s.log.Infof("SAGA工作流事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestWorkFlowTCC(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	s.log.Infof("开始TCC工作流事务，GID: %s", gid)

	// 提交工作流
	if err := workflow.Execute(WorkflowShopServiceOrderTCC, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("TCC工作流事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError("TCC工作流事务提交失败")
	}

	s.log.Infof("TCC工作流 事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestWorkFlowXA(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	s.log.Infof("开始XA工作流事务，GID: %s", gid)

	// 提交工作流
	if err := workflow.Execute(WorkflowShopServiceOrderXA, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("XA工作流事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError("XA工作流事务提交失败")
	}

	s.log.Infof("XA工作流 事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

func (s *ShopService) TestWorkFlowMixed(ctx context.Context, req *shopV1.BuyRequest) (*shopV1.BuyResponse, error) {
	// 生成全局唯一事务 ID (GID)
	gid := dtmgrpc.MustGenGid(service.DtmServerAddress)

	s.log.Infof("开始混合工作流事务，GID: %s", gid)

	// 提交工作流
	if err := workflow.Execute(WorkflowShopServiceOrderMixed, gid, dtmgimp.MustProtoMarshal(req)); err != nil {
		s.log.Errorf("混合工作流事务提交失败: %v", err)
		return nil, shopV1.ErrorInternalServerError("混合工作流事务提交失败")
	}

	s.log.Infof("混合事务提交成功，GID: %s", gid)

	return &shopV1.BuyResponse{Success: true}, nil
}

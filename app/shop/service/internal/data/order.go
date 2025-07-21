package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"

	"kratos-dtm-examples/pkg/dtmgorm"
)

type OrderRepo struct {
	log  *log.Helper
	data *Data
}

func NewOrderRepo(logger log.Logger, db *Data) *OrderRepo {
	return &OrderRepo{
		data: db,
		log:  log.NewHelper(log.With(logger, "module", "shop/service/data/order")),
	}
}

// CreateOrder 创建订单
func (r *OrderRepo) CreateOrder(dto *shopV1.Order) error {
	return r.CreateOrderWithTx(r.data.db, dto)
}

func (r *OrderRepo) CreateOrderWithTx(tx *gorm.DB, dto *shopV1.Order) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.CreateTime == nil {
		dto.CreateTime = timestamppb.New(time.Now())
	}

	var model models.Order
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return tx.Create(&model).Error
}

// UpdateOrder 更新订单
func (r *OrderRepo) UpdateOrder(dto *shopV1.Order) error {
	return r.UpdateOrderWithTx(r.data.db, dto)
}

func (r *OrderRepo) UpdateOrderWithTx(tx *gorm.DB, dto *shopV1.Order) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.UpdateTime == nil {
		dto.UpdateTime = timestamppb.New(time.Now())
	}

	var model models.Order
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return tx.Save(&model).Error
}

// GetOrderByID 根据 ID 获取订单
func (r *OrderRepo) GetOrderByID(id uint) (*shopV1.Order, error) {
	var model models.Order
	if err := r.data.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var dto shopV1.Order
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *OrderRepo) GetOrderByRequestID(requestId string) (*shopV1.Order, error) {
	var model models.Order
	if err := r.data.db.
		Where("request_id = ?", requestId).
		First(&model).Error; err != nil {
		return nil, err
	}

	var dto shopV1.Order
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *OrderRepo) GetOrderByOrderNo(orderNo string) (*shopV1.Order, error) {
	return r.GetOrderByOrderNoWithTx(r.data.db, orderNo)
}

func (r *OrderRepo) GetOrderByOrderNoWithTx(tx *gorm.DB, orderNo string) (*shopV1.Order, error) {
	var model models.Order
	if err := tx.Where("order_no = ?", orderNo).First(&model).Error; err != nil {
		return nil, err
	}

	var dto shopV1.Order
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

// DeleteOrderByID 删除订单
func (r *OrderRepo) DeleteOrderByID(id uint32) error {
	return r.data.db.Delete(&models.Order{}, id).Error
}

func (r *OrderRepo) DeleteOrderByRequestID(requestId string) error {
	return r.DeleteOrderByRequestIDWithTx(r.data.db, requestId)
}

func (r *OrderRepo) DeleteOrderByRequestIDWithTx(tx *gorm.DB, requestId string) error {
	return tx.Where("request_id = ?", requestId).Delete(&models.Order{}).Error
}

func (r *OrderRepo) DeleteOrderByOrderNo(orderNo string) error {
	return r.DeleteOrderByOrderNoWithTx(r.data.db, orderNo)
}

func (r *OrderRepo) DeleteOrderByOrderNoWithTx(tx *gorm.DB, orderNo string) error {
	return tx.Where("order_no = ?", orderNo).Delete(&models.Order{}).Error
}

func (r *OrderRepo) OrderExists(userId, productId uint32) (bool, error) {
	var count int64
	err := r.data.db.Model(&models.Order{}).
		Where("user_id = ? AND product_id = ?", userId, productId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderRepo) OrderExistsByRequestID(requestId string) (bool, error) {
	return r.OrderExistsByRequestIDWithTx(r.data.db, requestId)
}

func (r *OrderRepo) OrderExistsByRequestIDWithTx(tx *gorm.DB, requestId string) (bool, error) {
	var count int64
	err := tx.Model(&models.Order{}).
		Where("request_id = ?", requestId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderRepo) OrderExistsByOrderNo(orderNo string) (bool, error) {
	var count int64
	err := r.data.db.Model(&models.Order{}).
		Where("order_no = ?", orderNo).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *OrderRepo) TryCreateOrder(ctx context.Context, req *shopV1.TryCreateOrderRequest) (*shopV1.OrderResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 检查订单是否已存在
		var exists bool
		exists, err = r.OrderExistsByRequestIDWithTx(tx, req.RequestId)
		if err != nil {
			r.log.Errorf("Failed to check if order exists: %v", err)
			return shopV1.ErrorInternalServerError("failed to check order existence")
		}
		if exists {
			r.log.Infof("Order already exists for request ID %s", req.RequestId)
			return shopV1.ErrorInternalServerError("order already exists for request ID %s", req.RequestId)
		}

		// 生成订单预创建记录
		order := &shopV1.Order{
			UserId:    req.UserId,
			ProductId: req.ProductId,
			RequestId: req.RequestId,
			OrderNo:   req.OrderNo,
			Quantity:  req.Quantity,
			Status:    shopV1.OrderStatus_PENDING, // 设置订单状态为 PENDING
		}
		err = r.CreateOrderWithTx(tx, order)
		if err != nil {
			r.log.Errorf("Failed to create pending order: %v", err)
			return shopV1.ErrorInternalServerError("failed to create pending order: %v", err)
		}

		return nil
	})

	r.log.Infof("Pending order created successfully for user %d", req.UserId)

	if err != nil {
		return &shopV1.OrderResponse{Success: false, Message: err.Error()}, nil
	}

	return &shopV1.OrderResponse{
		Success: true,
		Message: "Order created successfully, please confirm or cancel",
	}, nil
}

func (r *OrderRepo) ConfirmCreateOrder(ctx context.Context, req *shopV1.ConfirmCreateOrderRequest) (*shopV1.OrderResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询订单是否存在
		var order *shopV1.Order
		if order, err = r.GetOrderByOrderNoWithTx(tx, req.OrderNo); err != nil {
			r.log.Errorf("Failed to get order: %v", err)
			return shopV1.ErrorInternalServerError("failed to get order: %v", err)
		}
		if order == nil {
			return shopV1.ErrorNotFound("order not found")
		}

		// 检查订单状态是否为 PENDING
		if order.Status != shopV1.OrderStatus_PENDING {
			return shopV1.ErrorPreconditionFailed("order is not in PENDING state")
		}

		// 更新订单状态为 CONFIRMED
		order.Status = shopV1.OrderStatus_CONFIRMED
		if err = r.UpdateOrderWithTx(tx, order); err != nil {
			r.log.Errorf("Failed to update order status: %v", err)
			return shopV1.ErrorInternalServerError("failed to update order status: %v", err)
		}

		return nil
	})
	if err != nil {
		return &shopV1.OrderResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Order ID %s confirmed successfully", req.OrderNo)

	return &shopV1.OrderResponse{
		Success: true,
		Message: "Order confirmed successfully",
	}, nil
}

func (r *OrderRepo) CancelCreateOrder(ctx context.Context, req *shopV1.CancelCreateOrderRequest) (*shopV1.OrderResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询订单是否存在
		var order *shopV1.Order
		if order, err = r.GetOrderByOrderNoWithTx(tx, req.OrderNo); err != nil {
			r.log.Errorf("Failed to get order: %v", err)
			return shopV1.ErrorInternalServerError("failed to get order: %v", err)
		}
		if order == nil {
			return shopV1.ErrorNotFound("order not found")
		}

		// 检查订单状态是否为 PENDING
		if order.Status != shopV1.OrderStatus_PENDING {
			return shopV1.ErrorPreconditionFailed("order is not in PENDING state")
		}

		// 删除订单
		if err = r.DeleteOrderByOrderNoWithTx(tx, req.OrderNo); err != nil {
			r.log.Errorf("Failed to delete order: %v", err)
			return shopV1.ErrorInternalServerError("failed to delete order: %v", err)
		}

		return nil
	})
	if err != nil {
		return &shopV1.OrderResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Order ID %s canceled successfully", req.OrderNo)

	return &shopV1.OrderResponse{
		Success: true,
		Message: "Order canceled successfully",
	}, nil
}

func (r *OrderRepo) RefundOrder(ctx context.Context, req *shopV1.RefundOrderRequest) (*shopV1.OrderResponse, error) {
	var err error

	err = dtmgorm.BarrierGorm(ctx, r.data.db, func(tx *gorm.DB) error {
		// 查询订单是否存在
		var order *shopV1.Order
		if order, err = r.GetOrderByOrderNoWithTx(tx, req.OrderNo); err != nil {
			r.log.Errorf("Failed to get order: %v", err)
			return shopV1.ErrorInternalServerError("failed to get order: %v", err)
		}
		if order == nil {
			return shopV1.ErrorNotFound("order not found")
		}

		// 检查订单状态是否允许退款
		if order.Status != shopV1.OrderStatus_CONFIRMED {
			return shopV1.ErrorPreconditionFailed("order is not in a state that allows refund")
		}

		// 更新订单状态为 REFUNDED
		order.Status = shopV1.OrderStatus_REFUNDED
		if err = r.UpdateOrderWithTx(tx, order); err != nil {
			r.log.Errorf("Failed to update order status: %v", err)
			return shopV1.ErrorInternalServerError("failed to update order status: %v", err)
		}

		return nil
	})
	if err != nil {
		return &shopV1.OrderResponse{Success: false, Message: err.Error()}, nil
	}

	r.log.Infof("Order ID %s refunded successfully", req.OrderNo)

	return &shopV1.OrderResponse{
		Success: true,
		Message: "Order refunded successfully",
	}, nil
}

package data

import (
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

// CreateOrder 创建订单
func (r *OrderRepo) CreateOrder(dto *shopV1.Order) error {
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

	return r.db.Create(&model).Error
}

// UpdateOrder 更新订单
func (r *OrderRepo) UpdateOrder(dto *shopV1.Order) error {
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

	return r.db.Save(&model).Error
}

// GetOrderByID 根据 ID 获取订单
func (r *OrderRepo) GetOrderByID(id uint) (*shopV1.Order, error) {
	var model models.Order
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var dto shopV1.Order
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

// DeleteOrderByID 删除订单
func (r *OrderRepo) DeleteOrderByID(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

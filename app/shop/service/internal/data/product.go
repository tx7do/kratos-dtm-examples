package data

import (
	"time"

	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

// CreateProduct 创建商品
func (r *ProductRepo) CreateProduct(dto *shopV1.Product) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.CreateTime == nil {
		dto.CreateTime = timestamppb.New(time.Now())
	}

	var model models.Product
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return r.db.Create(&model).Error
}

// UpdateProduct 更新商品
func (r *ProductRepo) UpdateProduct(dto *shopV1.Product) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.UpdateTime == nil {
		dto.UpdateTime = timestamppb.New(time.Now())
	}

	var model models.Product
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return r.db.Save(&model).Error
}

// GetProductByID 根据 ID 获取商品
func (r *ProductRepo) GetProductByID(id uint) (*shopV1.Product, error) {
	var model models.Product
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var dto shopV1.Product
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

// DeleteProductByID 删除商品
func (r *ProductRepo) DeleteProductByID(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

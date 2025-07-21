package data

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"

	shopV1 "kratos-dtm-examples/api/gen/go/shop/service/v1"
)

type UserRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(logger log.Logger, db *Data) *UserRepo {
	return &UserRepo{
		data: db,
		log:  log.NewHelper(log.With(logger, "module", "shop/service/data/user")),
	}
}

// CreateUser 创建用户
func (r *UserRepo) CreateUser(dto *shopV1.User) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.CreateTime == nil {
		dto.CreateTime = timestamppb.New(time.Now())
	}

	var model models.User
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return r.data.db.Create(&model).Error
}

// UpdateUser 更新用户
func (r *UserRepo) UpdateUser(dto *shopV1.User) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.UpdateTime == nil {
		dto.UpdateTime = timestamppb.New(time.Now())
	}

	var model models.User
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return r.data.db.Save(&model).Error
}

// GetUserByID 根据 ID 获取用户
func (r *UserRepo) GetUserByID(id uint32) (*shopV1.User, error) {
	var model models.User
	if err := r.data.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var dto shopV1.User
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

// DeleteUserByID 删除用户
func (r *UserRepo) DeleteUserByID(id uint32) error {
	return r.data.db.Delete(&models.User{}, id).Error
}

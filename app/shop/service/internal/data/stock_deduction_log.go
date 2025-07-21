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

type StockDeductionLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewStockDeductionLogRepo(logger log.Logger, db *Data) *StockDeductionLogRepo {
	return &StockDeductionLogRepo{
		data: db,
		log:  log.NewHelper(log.With(logger, "module", "shop/service/data/stock_deduction_log")),
	}
}

func (r *StockDeductionLogRepo) CreateLog(dto *shopV1.StockDeductionLog) error {
	if dto == nil {
		return gorm.ErrInvalidData
	}

	if dto.CreateTime == nil {
		dto.CreateTime = timestamppb.New(time.Now())
	}

	var model models.StockDeductionLog
	if err := copier.Copy(&model, dto); err != nil {
		return err
	}

	return r.data.db.Create(&model).Error
}

func (r *StockDeductionLogRepo) GetLogByID(id uint32) (*shopV1.StockDeductionLog, error) {
	var model models.StockDeductionLog
	if err := r.data.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	var dto shopV1.StockDeductionLog
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *StockDeductionLogRepo) GetLogByRequestID(requestID string) (*shopV1.StockDeductionLog, error) {
	var model models.StockDeductionLog
	if err := r.data.db.
		Where("request_id = ?", requestID).
		First(&model).Error; err != nil {
		return nil, err
	}

	var dto shopV1.StockDeductionLog
	if err := copier.Copy(&dto, &model); err != nil {
		return nil, err
	}

	return &dto, nil
}

func (r *StockDeductionLogRepo) ExistLogByRequestID(requestID string) (bool, error) {
	var count int64
	if err := r.data.db.Model(&models.StockDeductionLog{}).
		Where("request_id = ?", requestID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

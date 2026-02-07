package dtmgorm

import (
	"context"
	"database/sql"

	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"gorm.io/gorm"
)

// BarrierGorm 封装Barrier与GORM的集成
func BarrierGorm(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	// 从上下文中获取Barrier
	barrier, err := dtmgrpc.BarrierFromGrpc(ctx)
	if err != nil {
		return err
	}

	aDb, err := db.DB()
	if err != nil {
		return err
	}

	// 使用Barrier调用，并将sql.Tx转换为GORM会话
	return barrier.CallWithDB(aDb, func(sqlTx *sql.Tx) error {

		gormTx := db.
			Session(&gorm.Session{NewDB: true}).
			WithContext(ctx)
		gormTx.Statement.ConnPool = sqlTx

		return fn(gormTx)
	})
}

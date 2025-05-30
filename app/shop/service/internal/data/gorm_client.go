package data

import (
	"github.com/go-kratos/kratos/v2/log"
	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"

	"kratos-dtm-examples/app/shop/service/internal/data/models"
)

// NewGormClient 创建数据库客户端
func NewGormClient(cfg *conf.Bootstrap, logger log.Logger) *gorm.DB {
	l := log.NewHelper(log.With(logger, "module", "gorm/data/shop-service"))

	var driver gorm.Dialector
	switch cfg.Data.Database.Driver {
	default:
		fallthrough
	case "mysql":
		driver = mysql.Open(cfg.Data.Database.Source)
		break
	case "postgres":
		driver = postgres.Open(cfg.Data.Database.Source)
		break
	case "clickhouse":
		driver = clickhouse.Open(cfg.Data.Database.Source)
		break
	case "sqlite":
		driver = sqlite.Open(cfg.Data.Database.Source)
		break
	case "sqlserver":
		driver = sqlserver.Open(cfg.Data.Database.Source)
		break
	}

	client, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		l.Fatalf("failed opening connection to db: %v", err)
		return nil
	}

	// 运行数据库迁移工具
	if cfg.Data.Database.Migrate {
		if err := client.AutoMigrate(
			models.GetMigrates()...,
		); err != nil {
			l.Fatalf("failed creating schema resources: %v", err)
			return nil
		}
	}
	return client
}

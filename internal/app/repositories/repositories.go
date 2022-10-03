package repositories

import (
	"strings"

	"github.com/flytrap/gin-base/internal/app/config"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var RepositorySet = wire.NewSet()

func AutoMigrate(db *gorm.DB) error {
	if dbType := config.C.Gorm.DBType; strings.ToLower(dbType) == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	}

	return db.AutoMigrate() // end
}

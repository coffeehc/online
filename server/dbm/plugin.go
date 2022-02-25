package dbm

import (
	"github.com/jinzhu/gorm"
	"online/common/utils"
	"online/common/utils/pagination"
	"online/server/api/dbutil"
	"online/server/web/gen/restapi/operations"
	"time"
)

type YakitPlugin struct {
	gorm.Model

	Type              string    `json:"type"`
	ScriptName        string    `json:"script_name"`
	Authors           string    `json:"authors"`
	Content           string    `json:"content"`
	PublishedAt       time.Time `json:"published_at"`
	Tags              string    `json:"tags"`
	OwnerUserUniqueId string    // 插件拥有者，只有拥有着才有权限修改插件
	IsOfficial        bool
	DefaultOpen       bool
	Params            string // 插件关联的参数，同 Yakit 中的定义，需要被 str.Quote 包裹
	DownloadTotal     int64  `json:"downloadTotal"`
}

func CreateOrUpdateYakitPlugin(db *gorm.DB, hash string, i interface{}) error {
	db = db.Model(&YakitPlugin{})

	if db := db.Where("hash = ?", hash).Assign(i).FirstOrCreate(&YakitPlugin{}); db.Error != nil {
		return utils.Errorf("create/update YakitPlugin failed: %s", db.Error)
	}

	return nil
}

func GetYakitPlugin(db *gorm.DB, id int64) (*YakitPlugin, error) {
	var req YakitPlugin
	if db := db.Model(&YakitPlugin{}).Where("id = ?", id).First(&req); db.Error != nil {
		return nil, utils.Errorf("get YakitPlugin failed: %s", db.Error)
	}

	return &req, nil
}

func DeleteYakitPluginByID(db *gorm.DB, id int64) error {
	if db := db.Model(&YakitPlugin{}).Where(
		"id = ?", id,
	).Unscoped().Delete(&YakitPlugin{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func QueryYakitPlugin(db *gorm.DB, params operations.GetYakitPluginParams) (*pagination.Paginator, []*YakitPlugin, error) {
	db = db.Model(&YakitPlugin{})

	db = dbutil.QueryOrderP(db, params.OrderBy, params.Order)

	var ret []*YakitPlugin
	paging, db := dbutil.PagingP(db, params.Page, params.Limit, &ret)
	if db.Error != nil {
		return nil, nil, utils.Errorf("paging failed: %s", db.Error)
	}

	return paging, ret, nil
}

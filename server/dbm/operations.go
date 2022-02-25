package dbm

import (
	"github.com/jinzhu/gorm"
	"online/common/utils"
	"online/common/utils/pagination"
	"online/server/api/dbutil"
	"online/server/web/gen/restapi/operations"
)

type OperationType string

const (
	// 基础操作
	OperationType_Star     OperationType = "star"
	OperationType_Fork     OperationType = "fork"
	OperationType_Download OperationType = "download"

	// 增删改查
	OperationType_Submit   OperationType = "submit"
	OperationType_Delete   OperationType = "delete"
	OperationType_Modified OperationType = "modified"
)

type Operation struct {
	gorm.Model

	OperationPluginId   int           `json:"operationPluginId"`
	TriggerUserUniqueId string        `json:"triggerUserUniqueId"`
	Type                OperationType `json:"operationType"`

	Extra string `json:"extra"` // 额外信息
}

func CreateOrUpdateOperation(db *gorm.DB, hash string, i interface{}) error {
	db = db.Model(&Operation{})

	if db := db.Where("hash = ?", hash).Assign(i).FirstOrCreate(&Operation{}); db.Error != nil {
		return utils.Errorf("create/update Operation failed: %s", db.Error)
	}

	return nil
}

func GetOperation(db *gorm.DB, id int64) (*Operation, error) {
	var req Operation
	if db := db.Model(&Operation{}).Where("id = ?", id).First(&req); db.Error != nil {
		return nil, utils.Errorf("get Operation failed: %s", db.Error)
	}

	return &req, nil
}

func DeleteOperationByID(db *gorm.DB, id int64) error {
	if db := db.Model(&Operation{}).Where(
		"id = ?", id,
	).Unscoped().Delete(&Operation{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func QueryOperation(db *gorm.DB, params operations.GetOperationParams) (*pagination.Paginator, []*Operation, error) {
	db = db.Model(&Operation{})

	db = dbutil.QueryOrderP(db, params.OrderBy, params.Order)

	var ret []*Operation
	paging, db := dbutil.PagingP(db, params.Page, params.Limit, &ret)
	if db.Error != nil {
		return nil, nil, utils.Errorf("paging failed: %s", db.Error)
	}

	return paging, ret, nil
}

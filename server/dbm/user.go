package dbm

import (
	"github.com/jinzhu/gorm"
	"online/common/utils"
	"online/common/utils/pagination"
	"online/server/api/dbutil"
	"online/server/web/gen/restapi/operations"
)

type User struct {
	gorm.Model

	UserUniqueId string `json:"userUniqueId" gorm:"unique_index"`
	UserVerbose  string `json:"userVerbose"`
	FromPlatform string `json:"userFromPlatform"`
	Email        string `json:"email"`
	Tags         string
}

func CreateOrUpdateUser(db *gorm.DB, hash string, i interface{}) error {
	db = db.Model(&User{})

	if db := db.Where("hash = ?", hash).Assign(i).FirstOrCreate(&User{}); db.Error != nil {
		return utils.Errorf("create/update User failed: %s", db.Error)
	}

	return nil
}

func GetUser(db *gorm.DB, id int64) (*User, error) {
	var req User
	if db := db.Model(&User{}).Where("id = ?", id).First(&req); db.Error != nil {
		return nil, utils.Errorf("get User failed: %s", db.Error)
	}

	return &req, nil
}

func DeleteUserByID(db *gorm.DB, id int64) error {
	if db := db.Model(&User{}).Where(
		"id = ?", id,
	).Unscoped().Delete(&User{}); db.Error != nil {
		return db.Error
	}
	return nil
}

func QueryUser(db *gorm.DB, params operations.GetUserParams) (*pagination.Paginator, []*User, error) {
	db = db.Model(&User{})

	db = dbutil.QueryOrderP(db, params.OrderBy, params.Order)

	var ret []*User
	paging, db := dbutil.PagingP(db, params.Page, params.Limit, &ret)
	if db.Error != nil {
		return nil, nil, utils.Errorf("paging failed: %s", db.Error)
	}

	return paging, ret, nil
}

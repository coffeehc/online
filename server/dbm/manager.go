package dbm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"online/common/log"
	"online/common/utils"
	"regexp"
)

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var uniqueExistRegex = regexp.MustCompile(`pq\:\sduplicate\s+key\s+value\sviolates\s+unique\s+constraint\s\"uix\_`)

type Manager struct {
	DB *gorm.DB
}

func (m *Manager) Save(i interface{}, msg string, items ...interface{}) {
	var (
		message = msg
	)
	if len(items) > 0 {
		message = fmt.Sprintf(msg, items...)
	}
	if db := m.DB.Save(i); db.Error != nil {
		log.Errorf("save %v failed: %s", message, db.Error)
		return
	}
}

var (
	managerSingleton *Manager
)

func NewDBManager(params string) (*Manager, error) {
	return NewDBManagerWithMigrate(params, true)
}

func NewDBManagerWithMigrate(params string, migrate bool) (*Manager, error) {
	if managerSingleton != nil {
		return managerSingleton, nil
	}

	log.Info("start to connection postgres")
	db, err := gorm.Open("postgres", params)
	if err != nil {
		return nil, errors.Errorf("open postgres db failed: %s", err)
	}
	log.Info("build basic database manager instance")
	m := &Manager{}

	m.DB = db

	if utils.InDebugMode() {
		m.DB = m.DB.Debug()
	}

	if migrate {
		m.AutoMigrate()
	}

	managerSingleton = m
	return m, nil

}

func (m *Manager) AutoMigrate() {
	log.Info("start to migrate models")

	m.DB.AutoMigrate(&User{}, &Operation{}, &YakitPlugin{})
	log.Info("migrate all models finished")
}

func (m *Manager) Init() error {

	return nil
}

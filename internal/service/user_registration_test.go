package service

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"huynhmanh.com/gin/internal/model"
)

type MockProfileRepo struct {
	mock.Mock
}

func (m *MockProfileRepo) CreateTx(tx *gorm.DB,profile *model.UserProFile) error{
	args := m.Called(tx,profile)
	return args.Error(0)
}

type MockAuditRepo struct{
	mock.Mock
}

func (m *MockAuditRepo) CreateTx(tx *gorm.DB, log *model.AuditLog) error {
    args := m.Called(tx, log)
    return args.Error(0)
}

func setupTestDB() *gorm.DB{
	dsn := "root:huynhmanh221199@tcp(127.0.0.1:3006)/test_db?charset=utf8mb4"
	db,err := gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil {
		panic("failed to setup test database: "+ err.Error())
	}

	db.Migrator().DropTable(&model.AuditLog{})
	db.Migrator().DropTable(&model.UserProFile{})
	db.Migrator().DropTable(&model.User{})

	db.AutoMigrate(
		&model.User{},
        &model.UserProFile{},
        &model.AuditLog{},
	)
	return db
}
func cleanupTestDB(db *gorm.DB) {
    sqlDB, _ := db.DB()
    sqlDB.Close()
}

func TestRegisterWithProfile_Success(t *testing.T){
	db := setupTestDB()
	defer cleanupTestDB(db)

}
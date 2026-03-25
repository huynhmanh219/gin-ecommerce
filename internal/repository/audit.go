package repository

import (
    "gorm.io/gorm"
    "huynhmanh.com/gin/internal/model"
)

type AuditRepository interface {
    Create(log *model.AuditLog) (*model.AuditLog, error)
    GetByUserID(userID uint) ([]model.AuditLog, error)
    
    CreateTx(tx *gorm.DB, log *model.AuditLog) (*model.AuditLog, error)
}

type MySQLAuditRepository struct {
    db *gorm.DB
}

func NewMySQLAuditRepository(db *gorm.DB) AuditRepository {
    return &MySQLAuditRepository{db: db}
}

func (r *MySQLAuditRepository) CreateTx(tx *gorm.DB, log *model.AuditLog) (*model.AuditLog, error) {
    if err := tx.Create(log).Error; err != nil {
        return nil, err
    }
    return log, nil
}

func (r *MySQLAuditRepository) Create(log *model.AuditLog) (*model.AuditLog, error) {
    if err := r.db.Create(log).Error; err != nil {
        return nil, err
    }
    return log, nil
}

func (r *MySQLAuditRepository) GetByUserID(userID uint) ([]model.AuditLog, error) {
    var logs []model.AuditLog
    if err := r.db.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
        return nil, err
    }
    return logs, nil
}
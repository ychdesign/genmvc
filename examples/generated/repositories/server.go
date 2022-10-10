package repositories

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ychdesign/genmvc/examples/generated/po"
	"gorm.io/gorm"
)

type ServerRepository interface {
	Create(ctx context.Context, db *gorm.DB, serverPO *po.Server) error
	Delete(ctx context.Context, db *gorm.DB, id uint64) error
	Update(ctx context.Context, db *gorm.DB, serverPO *po.Server) error
	FindByID(ctx context.Context, db *gorm.DB, id uint64) (*po.Server, error)
	FindByPage(ctx context.Context, db *gorm.DB, pageNum int, pageSize int) ([]*po.Server, int, error)
}

type serverRepository struct {
}

func NewServerRepository() ServerRepository {
	return &serverRepository{}
}

func (r serverRepository) Create(ctx context.Context, db *gorm.DB, serverPO *po.Server) error {
	db = db.Table(po.ServerTableName).WithContext(ctx)
	if err := db.Create(serverPO).Error; err != nil {
		return errors.WithMessage(err, "create Server error")
	}
	return nil
}

func (r serverRepository) Delete(ctx context.Context, db *gorm.DB, id uint64) error {
	db = db.Table(po.ServerTableName).WithContext(ctx)
	serverPO := new(po.Server)
	if err := db.Delete(serverPO, "id = ?", id).Error; err != nil {
		return errors.WithMessage(err, "delete Server by id error")
	}
	return nil
}

func (r serverRepository) Update(ctx context.Context, db *gorm.DB, serverPO *po.Server) error {
	db = db.Table(po.ServerTableName).WithContext(ctx)
	if err := db.Save(serverPO).Error; err != nil {
		return errors.WithMessage(err, "update Server error")
	}
	return nil
}

func (r serverRepository) FindByID(ctx context.Context, db *gorm.DB, id uint64) (*po.Server, error) {
	db = db.Table(po.ServerTableName).WithContext(ctx)
	serverPO := new(po.Server)
	if err := db.Where("id = ?", id).First(serverPO).Error; err != nil {
		return nil, errors.WithMessage(err, "find Server error")
	}
	return serverPO, nil
}

func (r serverRepository) FindByPage(ctx context.Context, db *gorm.DB, pageNum int, pageSize int) ([]*po.Server, int, error) {
	db = db.Table(po.ServerTableName).WithContext(ctx)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errors.WithMessage(err, "count Server error")
	}

	serverPOs := make([]*po.Server, 0)
	if err := db.Offset(pageSize * (pageNum - 1)).Limit(pageSize).Find(&serverPOs).Error; err != nil {
		return nil, 0, errors.WithMessage(err, "find Server by page error")
	}
	return serverPOs, int(count), nil
}

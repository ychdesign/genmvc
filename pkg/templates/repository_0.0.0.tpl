package repositories

import (
	"context"
	"{{.ModulePath}}/po"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type {{.RepositoryIfaceName}} interface {
	Create(ctx context.Context, db *gorm.DB, {{.InstanceName}} *po.{{.Name}}) error
	Delete(ctx context.Context, db *gorm.DB, id uint64) error
	Update(ctx context.Context, db *gorm.DB, {{.InstanceName}} *po.{{.Name}}) error
	FindByID(ctx context.Context, db *gorm.DB, id uint64) (*po.{{.Name}}, error)
	FindByPage(ctx context.Context, db *gorm.DB, pageNum int, pageSize int) ([]*po.{{.Name}}, int, error)
}

type {{.RepositoryIfaceInstanceName}} struct {
}

func New{{.RepositoryIfaceName}}() {{.RepositoryIfaceName}} {
	return &{{.RepositoryIfaceInstanceName}}{}
}

func (r {{.RepositoryIfaceInstanceName}}) Create(ctx context.Context, db *gorm.DB, {{.InstanceName}} *po.{{.Name}}) error {
	db = db{{if ne "" .RepositoryTableName }}.Table(po.{{.RepositoryTableName}}){{end}}.WithContext(ctx)
	if err := db.Create({{.InstanceName}}).Error; err != nil {
		return errors.WithMessage(err, "create {{.Name}} error")
	}
	return nil
}

func (r {{.RepositoryIfaceInstanceName}}) Delete(ctx context.Context, db *gorm.DB, id uint64) error {
	db = db{{if ne "" .RepositoryTableName }}.Table(po.{{.RepositoryTableName}}){{end}}.WithContext(ctx)
	{{.InstanceName}} := new(po.{{.Name}})
	if err := db.Delete({{.InstanceName}}, "id = ?", id).Error; err != nil {
		return errors.WithMessage(err, "delete {{.Name}} by id error")
	}
	return nil
}

func (r {{.RepositoryIfaceInstanceName}}) Update(ctx context.Context, db *gorm.DB, {{.InstanceName}} *po.{{.Name}}) error {
	db = db{{if ne "" .RepositoryTableName }}.Table(po.{{.RepositoryTableName}}){{end}}.WithContext(ctx)
	if err := db.Save({{.InstanceName}}).Error; err != nil {
		return errors.WithMessage(err, "update {{.Name}} error")
	}
	return nil
}

func (r {{.RepositoryIfaceInstanceName}}) FindByID(ctx context.Context, db *gorm.DB, id uint64) (*po.{{.Name}}, error) {
	db = db{{if ne "" .RepositoryTableName }}.Table(po.{{.RepositoryTableName}}){{end}}.WithContext(ctx)
	{{.InstanceName}} := new(po.{{.Name}})
	if err := db.Where("id = ?", id).First({{.InstanceName}}).Error; err != nil {
		return nil, errors.WithMessage(err, "find {{.Name}} error")
	}
	return {{.InstanceName}}, nil
}

func (r {{.RepositoryIfaceInstanceName}}) FindByPage(ctx context.Context, db *gorm.DB, pageNum int, pageSize int) ([]*po.{{.Name}}, int, error) {
	db = db{{if ne "" .RepositoryTableName }}.Table(po.{{.RepositoryTableName}}){{end}}.WithContext(ctx)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errors.WithMessage(err, "count {{.Name}} error")
	}
	
	{{.InstanceSliceName}} := make([]*po.{{.Name}}, 0)
	if err := db.Offset(pageSize * (pageNum - 1)).Limit(pageSize).Find(&{{.InstanceSliceName}}).Error; err != nil {
		return nil, 0, errors.WithMessage(err, "find {{.Name}} by page error")
	}
	return {{.InstanceSliceName}}, int(count), nil
}

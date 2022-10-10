package services

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"{{.ModulePath}}/po"
	"{{.ModulePath}}/bo"
	"{{.ModulePath}}/repositories"
)

type {{.ServiceIfaceName}} interface {
	Create(ctx context.Context, {{.InstanceName}} *bo.{{.Name}}) error
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, {{.InstanceName}} *bo.{{.Name}}) error
	FindByID(ctx context.Context, id uint64) (*bo.{{.Name}}, error)
	FindByPage(ctx context.Context, pageNum int, pageSize int) ([]*bo.{{.Name}}, int, error)
}

type {{.ServiceIfaceInstanceName}} struct {
	db                   *gorm.DB
	{{.RepositoryIfaceInstanceName}} repositories.{{.RepositoryIfaceName}}
}

func New{{.ServiceIfaceName}}(
	db *gorm.DB,
 	{{.RepositoryIfaceInstanceName}} repositories.{{.RepositoryIfaceName}},
 ) {{.ServiceIfaceName}} {
	return &{{.ServiceIfaceInstanceName}}{
		db:                   db,
		{{.RepositoryIfaceInstanceName}}: {{.RepositoryIfaceInstanceName}},
	}
}

func (svc {{.ServiceIfaceInstanceName}}) Create(ctx context.Context, {{.InstanceName}} *bo.{{.Name}}) error {
	if err := svc.{{.RepositoryIfaceInstanceName}}.Create(ctx, svc.db, convertTo{{.Name}}PO({{.InstanceName}})); err != nil {
		return errors.WithMessage(err, "call {{.RepositoryIfaceInstanceName}}.Create error")
	}
	return nil
}

func (svc {{.ServiceIfaceInstanceName}}) Delete(ctx context.Context, id uint64) error {
	if err := svc.{{.RepositoryIfaceInstanceName}}.Delete(ctx, svc.db, id); err != nil {
		return errors.WithMessage(err, "call {{.RepositoryIfaceInstanceName}}.Delete error")
	}
	return nil
}

func (svc {{.ServiceIfaceInstanceName}}) Update(ctx context.Context, {{.InstanceName}} *bo.{{.Name}}) error {
	if err := svc.{{.RepositoryIfaceInstanceName}}.Update(ctx, svc.db, convertTo{{.Name}}PO({{.InstanceName}})); err != nil {
		return errors.WithMessage(err, "call {{.RepositoryIfaceInstanceName}}.Update error")
	}
	return nil
}

func (svc {{.ServiceIfaceInstanceName}}) FindByID(ctx context.Context, id uint64) (*bo.{{.Name}}, error) {
	{{.InstanceName}}, err := svc.{{.RepositoryIfaceInstanceName}}.FindByID(ctx, svc.db, id)
	if err != nil {
		return nil, errors.WithMessage(err, "call {{.RepositoryIfaceInstanceName}}.FindByID error")
	}
	return convertTo{{.Name}}BO({{.InstanceName}}), nil
}

func (svc {{.ServiceIfaceInstanceName}}) FindByPage(ctx context.Context, pageNum int, pageSize int) ([]*bo.{{.Name}}, int, error) {
	pos, count, err := svc.{{.RepositoryIfaceInstanceName}}.FindByPage(ctx, svc.db, pageNum, pageSize)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "call {{.RepositoryIfaceInstanceName}}.FindByPage error")
	}
	
	return convertTo{{.Name}}BOs(pos), count, nil
}

// TODO
func convertTo{{.Name}}PO(objPO *bo.{{.Name}}) *po.{{.Name}}{
	panic("Implements Me")
}

// TODO
func convertTo{{.Name}}BO(objPO *po.{{.Name}}) *bo.{{.Name}}{
	panic("Implements Me")
}	

func convertTo{{.Name}}BOs(pos []*po.{{.Name}}) []*bo.{{.Name}}{
	bos := []*bo.{{.Name}}{}
	for _,po := range pos{
		bos = append(bos,convertTo{{.Name}}BO(po))
	}
	return bos
}

func convertTo{{.Name}}POs(bos []*bo.{{.Name}}) []*po.{{.Name}}{
	pos := []*po.{{.Name}}{}
	for _,bo := range bos{
		pos = append(pos,convertTo{{.Name}}PO(bo))
	}
	return pos
}
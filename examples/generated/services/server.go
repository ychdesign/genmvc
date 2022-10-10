package services

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"github.com/ychdesign/genmvc/examples/generated/po"
	"github.com/ychdesign/genmvc/examples/generated/bo"
	"github.com/ychdesign/genmvc/examples/generated/repositories"
)

type ServerService interface {
	Create(ctx context.Context, serverPO *bo.Server) error
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, serverPO *bo.Server) error
	FindByID(ctx context.Context, id uint64) (*bo.Server, error)
	FindByPage(ctx context.Context, pageNum int, pageSize int) ([]*bo.Server, int, error)
}

type serverService struct {
	db                   *gorm.DB
	serverRepository repositories.ServerRepository
}

func NewServerService(
	db *gorm.DB,
 	serverRepository repositories.ServerRepository,
 ) ServerService {
	return &serverService{
		db:                   db,
		serverRepository: serverRepository,
	}
}

func (svc serverService) Create(ctx context.Context, serverPO *bo.Server) error {
	if err := svc.serverRepository.Create(ctx, svc.db, convertToServerPO(serverPO)); err != nil {
		return errors.WithMessage(err, "call serverRepository.Create error")
	}
	return nil
}

func (svc serverService) Delete(ctx context.Context, id uint64) error {
	if err := svc.serverRepository.Delete(ctx, svc.db, id); err != nil {
		return errors.WithMessage(err, "call serverRepository.Delete error")
	}
	return nil
}

func (svc serverService) Update(ctx context.Context, serverPO *bo.Server) error {
	if err := svc.serverRepository.Update(ctx, svc.db, convertToServerPO(serverPO)); err != nil {
		return errors.WithMessage(err, "call serverRepository.Update error")
	}
	return nil
}

func (svc serverService) FindByID(ctx context.Context, id uint64) (*bo.Server, error) {
	serverPO, err := svc.serverRepository.FindByID(ctx, svc.db, id)
	if err != nil {
		return nil, errors.WithMessage(err, "call serverRepository.FindByID error")
	}
	return convertToServerBO(serverPO), nil
}

func (svc serverService) FindByPage(ctx context.Context, pageNum int, pageSize int) ([]*bo.Server, int, error) {
	pos, count, err := svc.serverRepository.FindByPage(ctx, svc.db, pageNum, pageSize)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "call serverRepository.FindByPage error")
	}
	
	return convertToServerBOs(pos), count, nil
}

// TODO
func convertToServerPO(objPO *bo.Server) *po.Server{
	panic("Implements Me")
}

// TODO
func convertToServerBO(objPO *po.Server) *bo.Server{
	panic("Implements Me")
}	

func convertToServerBOs(pos []*po.Server) []*bo.Server{
	bos := []*bo.Server{}
	for _,po := range pos{
		bos = append(bos,convertToServerBO(po))
	}
	return bos
}

func convertToServerPOs(bos []*bo.Server) []*po.Server{
	pos := []*po.Server{}
	for _,bo := range bos{
		pos = append(pos,convertToServerPO(bo))
	}
	return pos
}
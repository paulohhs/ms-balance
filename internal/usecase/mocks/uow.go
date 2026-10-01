package mocks

import (
	"context"

	"github.com/paulohhs/ms-balance/pkg/uow"
	"github.com/stretchr/testify/mock"
)

type UowMock struct {
	mock.Mock
}

func (m *UowMock) Register(name string, fc uow.RepositoryFactory) {
	m.Called(name, fc)
}

func (m *UowMock) GetRepository(ctx context.Context, name string) (interface{}, error) {
	args := m.Called(ctx, name)
	return args.Get(0), args.Error(1)
}

func (m *UowMock) Do(ctx context.Context, fn func(uow *uow.Uow) error) error {
	args := m.Called(ctx, fn)
	if err := args.Error(0); err != nil {
		return err
	}
	return fn(nil)
}

func (m *UowMock) CommitOrRollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func (m *UowMock) UnRegister(name string) {
	m.Called(name)
}

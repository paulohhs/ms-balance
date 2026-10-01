package create_client

import (
	"testing"

	"github.com/paulohhs/ms-wallet/internal/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateClientUseCase_Execute(t *testing.T) {
	clientMock := &mocks.ClientGatewayMock{}
	clientMock.On("Save", mock.Anything).Return(nil)
	uc := NewCreateClientUseCase(clientMock)

	output, err := uc.Execute(CreateClientInputDTO{
		Name:  "John Doe",
		Email: "j@j.com",
	})

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "j@j.com", output.Email)
	clientMock.AssertExpectations(t)
	clientMock.AssertNumberOfCalls(t, "Save", 1)
}

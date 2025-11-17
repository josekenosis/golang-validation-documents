package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/neoway/golang-validation-documents/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Create(document *domain.Document) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockDocumentRepository) FindByID(id uuid.UUID) (*domain.Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockDocumentRepository) FindByNumber(number string) (*domain.Document, error) {
	args := m.Called(number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockDocumentRepository) FindAll(filters map[string]interface{}, orderBy string, limit, offset int) ([]domain.Document, int64, error) {
	args := m.Called(filters, orderBy, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.Document), args.Get(1).(int64), args.Error(2)
}

func (m *MockDocumentRepository) Update(document *domain.Document) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockDocumentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestDocumentService_Create(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	tests := []struct {
		name    string
		number  string
		setup   func()
		wantErr bool
	}{
		{
			name:   "valid CPF",
			number: "11144477735",
			setup: func() {
				mockRepo.On("FindByNumber", "11144477735").Return(nil, errors.New("not found"))
				mockRepo.On("Create", mock.AnythingOfType("*domain.Document")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "invalid CPF",
			number: "11144477736",
			setup: func() {
			},
			wantErr: true,
		},
		{
			name:   "existing document",
			number: "11144477735",
			setup: func() {
				doc := &domain.Document{Number: "11144477735"}
				mockRepo.On("FindByNumber", "11144477735").Return(doc, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setup()

			_, err := service.Create(tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentService_Update(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	id := uuid.New()
	blocked := true

	tests := []struct {
		name    string
		id      uuid.UUID
		blocked *bool
		setup   func()
		wantErr bool
	}{
		{
			name:    "successful update",
			id:      id,
			blocked: &blocked,
			setup: func() {
				doc := &domain.Document{ID: id}
				mockRepo.On("FindByID", id).Return(doc, nil)
				mockRepo.On("Update", mock.AnythingOfType("*domain.Document")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "document not found",
			id:      id,
			blocked: &blocked,
			setup: func() {
				mockRepo.On("FindByID", id).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setup()

			_, err := service.Update(tt.id, tt.blocked)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}


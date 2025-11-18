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

func TestDocumentService_Create_AdditionalCases(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	tests := []struct {
		name    string
		number  string
		setup   func()
		wantErr bool
	}{
		{
			name:   "empty string",
			number: "",
			setup: func() {
			},
			wantErr: true,
		},
		{
			name:   "invalid format",
			number: "123",
			setup: func() {
			},
			wantErr: true,
		},
		{
			name:   "valid CNPJ",
			number: "11222333000181",
			setup: func() {
				mockRepo.On("FindByNumber", "11222333000181").Return(nil, errors.New("not found"))
				mockRepo.On("Create", mock.AnythingOfType("*domain.Document")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "CPF with formatting",
			number: "111.444.777-35",
			setup: func() {
				mockRepo.On("FindByNumber", "11144477735").Return(nil, errors.New("not found"))
				mockRepo.On("Create", mock.AnythingOfType("*domain.Document")).Return(nil)
			},
			wantErr: false,
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

func TestDocumentService_GetByID(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	id := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func()
		wantErr bool
	}{
		{
			name: "successful get",
			id:   id,
			setup: func() {
				doc := &domain.Document{ID: id, Number: "11144477735"}
				mockRepo.On("FindByID", id).Return(doc, nil)
			},
			wantErr: false,
		},
		{
			name: "document not found",
			id:   id,
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

			_, err := service.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentService_GetByNumber(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	tests := []struct {
		name    string
		number  string
		setup   func()
		wantErr bool
	}{
		{
			name:   "successful get",
			number: "11144477735",
			setup: func() {
				doc := &domain.Document{Number: "11144477735"}
				mockRepo.On("FindByNumber", "11144477735").Return(doc, nil)
			},
			wantErr: false,
		},
		{
			name:   "with formatting",
			number: "111.444.777-35",
			setup: func() {
				doc := &domain.Document{Number: "11144477735"}
				mockRepo.On("FindByNumber", "11144477735").Return(doc, nil)
			},
			wantErr: false,
		},
		{
			name:   "document not found",
			number: "11144477735",
			setup: func() {
				mockRepo.On("FindByNumber", "11144477735").Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setup()

			_, err := service.GetByNumber(tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByNumber() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDocumentService_List(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	tests := []struct {
		name      string
		filters   map[string]interface{}
		orderBy   string
		page      int
		pageSize  int
		setup     func()
		wantErr   bool
		wantCount int64
	}{
		{
			name:     "successful list",
			filters:  map[string]interface{}{},
			orderBy:  "created_at DESC",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{{Number: "11144477735"}}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 10, 0).Return(docs, int64(1), nil)
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "with number filter",
			filters: map[string]interface{}{
				"number": "11144477735",
			},
			orderBy:  "number ASC",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{{Number: "11144477735"}}
				mockRepo.On("FindAll", map[string]interface{}{"number": "11144477735"}, "number ASC", 10, 0).Return(docs, int64(1), nil)
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "with type filter",
			filters: map[string]interface{}{
				"type": "CPF",
			},
			orderBy:  "type DESC",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{{Type: "CPF"}}
				mockRepo.On("FindAll", map[string]interface{}{"type": "CPF"}, "type DESC", 10, 0).Return(docs, int64(1), nil)
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name: "with blocked filter",
			filters: map[string]interface{}{
				"blocked": true,
			},
			orderBy:  "blocked ASC",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{{Blocked: true}}
				mockRepo.On("FindAll", map[string]interface{}{"blocked": true}, "blocked ASC", 10, 0).Return(docs, int64(1), nil)
			},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:     "page 2",
			filters:  map[string]interface{}{},
			orderBy:  "",
			page:     2,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 10, 10).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "page size too large",
			filters:  map[string]interface{}{},
			orderBy:  "",
			page:     1,
			pageSize: 200,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 100, 0).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "invalid page",
			filters:  map[string]interface{}{},
			orderBy:  "",
			page:     0,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 10, 0).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "invalid page size",
			filters:  map[string]interface{}{},
			orderBy:  "",
			page:     1,
			pageSize: 0,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 10, 0).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "invalid order by",
			filters:  map[string]interface{}{},
			orderBy:  "invalid_field",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "created_at DESC", 10, 0).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:     "order by with regex pattern",
			filters:  map[string]interface{}{},
			orderBy:  "number ASC",
			page:     1,
			pageSize: 10,
			setup: func() {
				docs := []domain.Document{}
				mockRepo.On("FindAll", map[string]interface{}{}, "number ASC", 10, 0).Return(docs, int64(0), nil)
			},
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil
			tt.setup()

			_, count, err := service.List(tt.filters, tt.orderBy, tt.page, tt.pageSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if count != tt.wantCount {
				t.Errorf("List() count = %v, want %v", count, tt.wantCount)
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
	unblocked := false

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
			name:    "update to unblocked",
			id:      id,
			blocked: &unblocked,
			setup: func() {
				doc := &domain.Document{ID: id, Blocked: true}
				mockRepo.On("FindByID", id).Return(doc, nil)
				mockRepo.On("Update", mock.AnythingOfType("*domain.Document")).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "update with nil blocked",
			id:      id,
			blocked: nil,
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

func TestDocumentService_Delete(t *testing.T) {
	mockRepo := new(MockDocumentRepository)
	service := NewDocumentService(mockRepo)

	id := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		setup   func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   id,
			setup: func() {
				doc := &domain.Document{ID: id}
				mockRepo.On("FindByID", id).Return(doc, nil)
				mockRepo.On("Delete", id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "document not found",
			id:   id,
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

			err := service.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}


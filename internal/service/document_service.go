package service

import (
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/neoway/golang-validation-documents/internal/domain"
	"github.com/neoway/golang-validation-documents/internal/repository"
	"github.com/neoway/golang-validation-documents/internal/utils"
)

type DocumentService interface {
	Create(number string) (*domain.Document, error)
	GetByID(id uuid.UUID) (*domain.Document, error)
	GetByNumber(number string) (*domain.Document, error)
	List(filters map[string]interface{}, orderBy string, page, pageSize int) ([]domain.Document, int64, error)
	Update(id uuid.UUID, blocked *bool) (*domain.Document, error)
	Delete(id uuid.UUID) error
}

type documentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) DocumentService {
	return &documentService{repo: repo}
}

func (s *documentService) Create(number string) (*domain.Document, error) {
	cleaned := utils.CleanDocument(number)
	if cleaned == "" {
		return nil, errors.New("document number is required")
	}

	if !utils.ValidateDocument(cleaned) {
		return nil, errors.New("invalid document number")
	}

	existing, _ := s.repo.FindByNumber(cleaned)
	if existing != nil {
		return nil, errors.New("document already exists")
	}

	docType := utils.GetDocumentType(cleaned)
	if docType == "" {
		return nil, errors.New("invalid document format")
	}

	document := &domain.Document{
		Number: cleaned,
		Type:   docType,
		Blocked: false,
	}

	if err := s.repo.Create(document); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *documentService) GetByID(id uuid.UUID) (*domain.Document, error) {
	return s.repo.FindByID(id)
}

func (s *documentService) GetByNumber(number string) (*domain.Document, error) {
	cleaned := utils.CleanDocument(number)
	return s.repo.FindByNumber(cleaned)
}

func (s *documentService) List(filters map[string]interface{}, orderBy string, page, pageSize int) ([]domain.Document, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	processedFilters := make(map[string]interface{})
	if number, ok := filters["number"].(string); ok && number != "" {
		processedFilters["number"] = utils.CleanDocument(number)
	}
	if docType, ok := filters["type"].(string); ok && docType != "" {
		processedFilters["type"] = strings.ToUpper(docType)
	}
	if blocked, ok := filters["blocked"].(bool); ok {
		processedFilters["blocked"] = blocked
	}

	validOrderBy := validateOrderBy(orderBy)
	offset := (page - 1) * pageSize

	return s.repo.FindAll(processedFilters, validOrderBy, pageSize, offset)
}

func (s *documentService) Update(id uuid.UUID, blocked *bool) (*domain.Document, error) {
	document, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("document not found")
	}

	if blocked != nil {
		document.Blocked = *blocked
	}

	if err := s.repo.Update(document); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *documentService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("document not found")
	}

	return s.repo.Delete(id)
}

func validateOrderBy(orderBy string) string {
	validFields := map[string]bool{
		"created_at":     true,
		"updated_at":     true,
		"number":         true,
		"type":           true,
		"blocked":        true,
		"created_at ASC": true,
		"created_at DESC": true,
		"updated_at ASC": true,
		"updated_at DESC": true,
		"number ASC":     true,
		"number DESC":    true,
		"type ASC":       true,
		"type DESC":      true,
	}

	if validFields[orderBy] {
		return orderBy
	}

	re := regexp.MustCompile(`^(\w+)\s+(ASC|DESC)$`)
	matches := re.FindStringSubmatch(orderBy)
	if len(matches) == 3 && validFields[matches[1]] {
		return orderBy
	}

	return "created_at DESC"
}

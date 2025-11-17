package repository

import (
	"github.com/google/uuid"
	"github.com/neoway/golang-validation-documents/internal/domain"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	Create(document *domain.Document) error
	FindByID(id uuid.UUID) (*domain.Document, error)
	FindByNumber(number string) (*domain.Document, error)
	FindAll(filters map[string]interface{}, orderBy string, limit, offset int) ([]domain.Document, int64, error)
	Update(document *domain.Document) error
	Delete(id uuid.UUID) error
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) DocumentRepository {
	return &documentRepository{db: db}
}

func (r *documentRepository) Create(document *domain.Document) error {
	return r.db.Create(document).Error
}

func (r *documentRepository) FindByID(id uuid.UUID) (*domain.Document, error) {
	var document domain.Document
	err := r.db.Where("id = ?", id).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) FindByNumber(number string) (*domain.Document, error) {
	var document domain.Document
	err := r.db.Where("number = ?", number).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (r *documentRepository) FindAll(filters map[string]interface{}, orderBy string, limit, offset int) ([]domain.Document, int64, error) {
	var documents []domain.Document
	var total int64

	query := r.db.Model(&domain.Document{})

	if number, ok := filters["number"].(string); ok && number != "" {
		query = query.Where("number LIKE ?", "%"+number+"%")
	}

	if docType, ok := filters["type"].(string); ok && docType != "" {
		query = query.Where("type = ?", docType)
	}

	if blocked, ok := filters["blocked"].(bool); ok {
		query = query.Where("blocked = ?", blocked)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&documents).Error
	return documents, total, err
}

func (r *documentRepository) Update(document *domain.Document) error {
	return r.db.Save(document).Error
}

func (r *documentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Document{}, id).Error
}


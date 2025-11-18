package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/neoway/golang-validation-documents/internal/service"
	"github.com/neoway/golang-validation-documents/internal/utils"
)

type DocumentHandler struct {
	service service.DocumentService
}

func NewDocumentHandler(service service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

type CreateDocumentRequest struct {
	Number string `json:"number"`
}

type UpdateDocumentRequest struct {
	Blocked *bool `json:"blocked"`
}

func (h *DocumentHandler) Create(c *fiber.Ctx) error {
	var req CreateDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	document, err := h.service.Create(req.Number)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(document)
}

func (h *DocumentHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	document, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document not found",
		})
	}

	return c.JSON(document)
}

func (h *DocumentHandler) GetByNumber(c *fiber.Ctx) error {
	number := c.Query("number")
	if number == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Number parameter is required",
		})
	}

	cleaned := utils.CleanDocument(number)
	if !utils.ValidateDocument(cleaned) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document number",
		})
	}

	document, err := h.service.GetByNumber(cleaned)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Document not found",
		})
	}

	return c.JSON(document)
}

func (h *DocumentHandler) List(c *fiber.Ctx) error {
	filters := make(map[string]interface{})

	if number := c.Query("number"); number != "" {
		filters["number"] = number
	}

	if docType := c.Query("type"); docType != "" {
		filters["type"] = docType
	}

	if blocked := c.Query("blocked"); blocked != "" {
		blockedBool, err := strconv.ParseBool(blocked)
		if err == nil {
			filters["blocked"] = blockedBool
		}
	}

	orderBy := c.Query("order_by", "created_at DESC")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	documents, total, err := h.service.List(filters, orderBy, page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list documents",
		})
	}

	return c.JSON(fiber.Map{
		"data":       documents,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (int(total) + pageSize - 1) / pageSize,
	})
}

func (h *DocumentHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	var req UpdateDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	document, err := h.service.Update(id, req.Blocked)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(document)
}

func (h *DocumentHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid document ID",
		})
	}

	if err := h.service.Delete(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}


package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ecommerce-catalog-api/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const maxImageSize = 5 << 20

type ImageUploadHandler struct {
	uploadDir string
}

func NewImageUploadHandler(uploadDir string) *ImageUploadHandler {
	return &ImageUploadHandler{uploadDir: uploadDir}
}

func (h *ImageUploadHandler) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Image file is required", err.Error())
	}
	if fileHeader.Size <= 0 || fileHeader.Size > maxImageSize {
		return response.Error(c, fiber.StatusBadRequest, "Image must be between 1 byte and 5 MB", nil)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Unable to read uploaded image", err.Error())
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, maxImageSize+1))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Unable to read uploaded image", err.Error())
	}
	if len(content) == 0 || len(content) > maxImageSize {
		return response.Error(c, fiber.StatusBadRequest, "Image must be between 1 byte and 5 MB", nil)
	}
	contentType := detectImageContentType(content)
	extension, ok := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}[contentType]
	if !ok {
		return response.Error(c, fiber.StatusBadRequest, "Supported image formats are JPEG, PNG, GIF, and WEBP", nil)
	}

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Unable to prepare image storage", err.Error())
	}
	filename := uuid.NewString() + extension
	destination := filepath.Join(h.uploadDir, filename)
	if err := os.WriteFile(destination, content, 0o644); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Unable to save image", err.Error())
	}

	publicBaseURL := strings.TrimRight(os.Getenv("PUBLIC_API_URL"), "/")
	if publicBaseURL == "" {
		publicBaseURL = strings.TrimRight(c.BaseURL(), "/")
	}
	publicURL := fmt.Sprintf("%s/uploads/%s", publicBaseURL, filename)
	return response.Success(c, fiber.StatusCreated, "Image uploaded successfully", fiber.Map{
		"publicUrl": publicURL,
	})
}

func detectImageContentType(content []byte) string {
	if len(content) >= 12 &&
		string(content[:4]) == "RIFF" &&
		string(content[8:12]) == "WEBP" {
		return "image/webp"
	}
	return strings.ToLower(strings.Split(http.DetectContentType(content), ";")[0])
}

package handler

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUploadImageSavesSupportedFileAndReturnsPublicURL(t *testing.T) {
	uploadDir := t.TempDir()
	app := fiber.New()
	app.Post("/api/v1/uploads", NewImageUploadHandler(uploadDir).Upload)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "product.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("\x89PNG\r\n\x1a\n" + strings.Repeat("\x00", 16))); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("POST", "/api/v1/uploads", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusCreated)
	}

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || filepath.Ext(files[0].Name()) != ".png" {
		t.Fatalf("uploaded files = %v, want one PNG file", files)
	}

	var responseBody bytes.Buffer
	if _, err := responseBody.ReadFrom(response.Body); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(responseBody.String(), "/uploads/"+files[0].Name()) {
		t.Fatalf("response does not contain public upload URL: %s", responseBody.String())
	}
}

func TestUploadImageRejectsNonImage(t *testing.T) {
	app := fiber.New()
	app.Post("/api/v1/uploads", NewImageUploadHandler(t.TempDir()).Upload)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "not-image.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("not an image")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest("POST", "/api/v1/uploads", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.StatusCode, fiber.StatusBadRequest)
	}
}

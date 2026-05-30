package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

const maxFileSize = 10 << 20 // 10 MB

// SaveUpload saves a multipart file to uploadDir and returns the filename.
func SaveUpload(uploadDir string, file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("file type %s not allowed", ext)
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, nil
}

// DeleteUpload removes an uploaded file.
func DeleteUpload(uploadDir, filename string) error {
	if filename == "" {
		return nil
	}
	return os.Remove(filepath.Join(uploadDir, filename))
}

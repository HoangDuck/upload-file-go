package uploadfiles

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	localStoragePath string
	secretKey        = []byte("hrdept-secret-key-2025") // Đặt trong env nếu cần bảo mật hơn
)

// InitLocalStorage initializes local storage configuration
func InitLocalStorage(uploadPath string) {
	localStoragePath = uploadPath
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(localStoragePath, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}
}

// UploadFileLocal uploads a file to local storage (similar to MinIO)
func UploadFileLocal(objectName string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	// Generate unique filename using timestamp and original filename
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	filename := timestamp + "_" + objectName
	filePath := filepath.Join(localStoragePath, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("error copying file: %v", err)
	}

	return filename, nil
}

const (
	MaxFileSize = 10 * 1024 * 1024 // 20MB in bytes
)

func ValidateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of 20MB")
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".pdf":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".docx": true,
	}

	if !allowedExts[ext] {
		return fmt.Errorf("unsupported file type. Allowed types: PDF, PNG, JPEG, JPG, DOCX")
	}

	// Open the file to check MIME type
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer src.Close()

	// Read first 512 bytes to determine MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Check MIME type
	mimeType := http.DetectContentType(buffer)
	allowedMimeTypes := map[string]bool{
		"application/pdf": true, // PDF
		"image/png":       true, // PNG
		"image/jpeg":      true, // JPEG/JPG
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // DOCX
	}

	if !allowedMimeTypes[mimeType] {
		return fmt.Errorf("invalid file content type: %s", mimeType)
	}

	return nil
}

// ValidateFileLocal validates the uploaded file (same as MinIO)
func ValidateFileLocal(file *multipart.FileHeader) error {
	// Reuse the same validation logic as MinIO
	return ValidateFile(file)
}

// ServeFileLocal serves a file from local storage with signature validation
func ServeFileLocal(filename string, expiresStr string, sig string) (string, error) {
	// Parse time
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil || time.Now().Unix() > expires {
		return "", fmt.Errorf("link expired or invalid")
	}

	// Verify signature
	data := fmt.Sprintf("%s:%d", filename, expires)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(data))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return "", fmt.Errorf("invalid signature")
	}

	// Check if file exists
	filePath := filepath.Join(localStoragePath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	}

	return filePath, nil
}

// GeneratePublicURL generates a public URL for local storage
func GeneratePublicURL(filename string, host string, scheme string) string {
	return fmt.Sprintf("%s://%s/api/v1/%s", scheme, host, filename)
}

// GenerateSignedURL generates a signed URL for local storage
func GenerateSignedURL(filename string, host string, scheme string, expiryMinutes int) string {
	expiry := time.Now().Add(time.Duration(expiryMinutes) * time.Minute).Unix()
	data := fmt.Sprintf("%s:%d", filename, expiry)

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(data))
	signature := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s://%s/api/v1/upload/private/%s?expires=%d&sig=%s", scheme, host, filename, expiry, signature)
}

// DeleteFileLocal deletes a file from local storage
func DeleteFileLocal(filename string) error {
	filePath := filepath.Join(localStoragePath, filename)
	return os.Remove(filePath)
}

// FileExistsLocal checks if a file exists in local storage
func FileExistsLocal(filename string) bool {
	filePath := filepath.Join(localStoragePath, filename)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFileSizeLocal gets the size of a file in local storage
func GetFileSizeLocal(filename string) (int64, error) {
	filePath := filepath.Join(localStoragePath, filename)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// ListFilesLocal lists all files in local storage
func ListFilesLocal() ([]string, error) {
	var files []string
	err := filepath.Walk(localStoragePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Return relative path from storage directory
			relPath, err := filepath.Rel(localStoragePath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}

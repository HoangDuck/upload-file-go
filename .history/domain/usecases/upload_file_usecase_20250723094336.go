package usecases

import (
	"fmt"
	"mime/multipart"
	uploadfiles "sound_qr_services/infrastructure/upload_files"
	"strings"
	"time"
)

type UploadFileUsecase interface {
	UploadFile(file *multipart.FileHeader) (string, error)
	GetPresignedURL(objectName string) (string, error)
}
type uploadFileUsecase struct {
}

// GetPresignedURL implements UploadFileUsecase.
func (u *uploadFileUsecase) GetPresignedURL(objectName string) (string, error) {
	return uploadfiles.GetPresignedURLLocal(objectName, time.Minute*15)
}

// UploadFile implements UploadFileUsecase.
func (u *uploadFileUsecase) UploadFile(file *multipart.FileHeader) (string, error) {
	fileName := u.GetFinalFileName(file)

	return uploadfiles.UploadFileLocal(fileName, file)
}

// GetFinalFileName returns the final sanitized filename with UUID prefix
func (u *uploadFileUsecase) GetFinalFileName(file *multipart.FileHeader) string {
	id, _ := generates.NewUUID()
	if id == "" {
		id = "default"
	}

	// Replace Vietnamese accented characters with non-accented equivalents
	vietnameseChars := map[rune]rune{
		'à': 'a', 'á': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
		'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
		'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
		'è': 'e', 'é': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
		'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
		'ì': 'i', 'í': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
		'ò': 'o', 'ó': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
		'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
		'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
		'ù': 'u', 'ú': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
		'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
		'ỳ': 'y', 'ý': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
		'đ': 'd',
	}

	// First replace Vietnamese characters
	filename := strings.Map(func(r rune) rune {
		if replacement, ok := vietnameseChars[r]; ok {
			return replacement
		}
		return r
	}, file.Filename)

	// Then replace remaining special characters with underscores
	filename = strings.Map(func(r rune) rune {
		// Keep alphanumeric characters and dots
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		// Replace everything else with underscore
		return '_'
	}, filename)

	// Ensure filename is not empty after sanitization
	if filename == "" || filename == "." {
		// Generate random string of 6 characters
		randomStr, _ := generates.NewUUID()
		if randomStr == "" {
			randomStr = "default"
		}
		filename = fmt.Sprintf("file_%s", randomStr)
	}

	return fmt.Sprintf("%s_%s", id, filename)
}

// NewPermissionUseCase creates a new instance of PermissionUseCase
func NewUploadFileUsecase() UploadFileUsecase {
	return &uploadFileUsecase{}
}

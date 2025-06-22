package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type FileType string

const (
	FileTypeResume    FileType = "resume"
	FileTypePortfolio FileType = "portfolio"
)

type StorageService interface {
	UploadFile(userID uint, fileType FileType, file *multipart.FileHeader) (string, error)
	DeleteFile(url string) error
	GetFileURL(userID uint, fileType FileType, filename string) string
}

type FileService interface {
	UploadResume(userID uint, file *multipart.FileHeader) (string, error)
	UploadPortfolio(userID uint, file *multipart.FileHeader) (string, error)
	ValidateFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) error
}

type fileService struct {
	storage StorageService
}

func NewFileService(storage StorageService) FileService {
	return &fileService{
		storage: storage,
	}
}

func (s *fileService) UploadResume(userID uint, file *multipart.FileHeader) (string, error) {
	if err := s.ValidateFile(file, []string{".pdf"}, 10*1024*1024); err != nil {
		return "", err
	}
	
	return s.storage.UploadFile(userID, FileTypeResume, file)
}

func (s *fileService) UploadPortfolio(userID uint, file *multipart.FileHeader) (string, error) {
	if err := s.ValidateFile(file, []string{".pdf"}, 25*1024*1024); err != nil {
		return "", err
	}
	
	return s.storage.UploadFile(userID, FileTypePortfolio, file)
}

func (s *fileService) ValidateFile(file *multipart.FileHeader, allowedTypes []string, maxSize int64) error {
	if file.Size > maxSize {
		return fmt.Errorf("file size %d bytes exceeds maximum %d bytes", file.Size, maxSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			return nil
		}
	}

	return fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, allowedTypes)
}
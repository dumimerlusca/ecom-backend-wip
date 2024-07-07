package service

import (
	"context"
	"database/sql"
	"ecom-backend/internal/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

type UploadService struct {
	db        *sql.DB
	fileModel *model.FileModel
}

func NewUploadService(db *sql.DB, fileModel *model.FileModel) *UploadService {
	return &UploadService{db: db, fileModel: fileModel}
}

func (svc *UploadService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (*model.FileRecord, error) {
	file, err := fileHeader.Open()

	if err != nil {
		return nil, err

	}

	defer file.Close()

	fileRecord := &model.FileRecord{
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Extension:    path.Ext(fileHeader.Filename),
	}

	tx, err := svc.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	// Insert file record
	fileRecord, err = svc.fileModel.Insert(ctx, tx, fileRecord)

	if err != nil {
		return nil, err
	}

	// Check if uploads directory exists
	exists, isDir := directoryExists("uploads")

	if !exists || !isDir {
		// If directory does not exist, create it
		err := os.Mkdir("uploads", os.ModePerm)

		if err != nil {
			return nil, err
		}
	}

	// Save file to disk
	dst, err := os.Create("uploads/" + fileRecord.Id + fileRecord.Extension)

	if err != nil {
		return nil, err

	}

	defer dst.Close()

	_, err = io.Copy(dst, file)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return fileRecord, nil
}

func (s *UploadService) GetFileRecord(fileId string) (*model.FileRecord, error) {
	fileRecord, err := s.fileModel.FindById(context.Background(), s.db, fileId)

	if err != nil {
		return nil, err
	}

	return fileRecord, nil
}

func directoryExists(path string) (bool, bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false
	}
	if err != nil {
		fmt.Printf("Error checking directory: %v\n", err)
		return false, false
	}
	return true, info.IsDir()
}

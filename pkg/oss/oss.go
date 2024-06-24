package oss

import (
	"context"
	"mime/multipart"
)

// OSS 对象存储接口
type OSS interface {
	UploadFile(ctx context.Context, key string, file multipart.File, fsize int64) (string, string, error)
	DeleteFile(key string) error
}

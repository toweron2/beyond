package oss

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"mime/multipart"
	"os"
)

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

type Qiniu struct {
	Cfg       storage.Config
	AccessKey string
	SecretKey string
	Bucket    string
	Zone      *storage.Zone
}

// NewQiniu 返回一个配置好的 Qiniu 实例
func NewQiniu(accessKey, secretKey, bucket string, useHTTPS, useCdnDomains bool) *Qiniu {
	return &Qiniu{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Cfg: storage.Config{
			UseHTTPS:      useHTTPS,
			UseCdnDomains: useCdnDomains,
		},
	}
}

func (q *Qiniu) UploadFile(ctx context.Context, key string, file multipart.File, fsize int64) (string, string, error) {
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	putPolicy := storage.PutPolicy{
		Scope:      q.Bucket,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
	}
	upToken := putPolicy.UploadToken(mac)
	resumeUploader := storage.NewResumeUploaderV2(&q.Cfg)

	ret := MyPutRet{}
	putExtra := storage.RputV2Extra{
		/*Params: map[string]string{
			"x:name": "github logo",
		},*/
	}

	err := resumeUploader.Put(ctx, &ret, upToken, key, file, fsize, &putExtra)
	if err != nil {
		logx.Error("上传失败:", err)

	}

	logx.Debug(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)
	// })
	// wg.Wait()

	return ret.Key, ret.Name, nil
}

func (q *Qiniu) DeleteFile(key string) error {
	// TODO implement me
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	bm := storage.NewBucketManager(mac, &q.Cfg)

	// 删除文件
	err := bm.Delete(q.Bucket, key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	return nil
}

func (q *Qiniu) ReplaceFileUpload(key, filePath, oldKey string) (string, string, error) {
	mac := qbox.NewMac(q.AccessKey, q.SecretKey)
	wg := threading.NewRoutineGroup()

	// 删除旧数据
	wg.RunSafe(func() {
		// 删除开头  / --斜杆
		err := storage.NewBucketManager(mac, &q.Cfg).Delete(q.Bucket, oldKey[1:])
		if err != nil {
			logx.Error("oss图片删除失败:", err)
		} else {
			logx.Debug("oss删除图片成功")
		}
	})
	ret := MyPutRet{}
	// 上传oss
	wg.RunSafe(func() {
		// 使用 returnBody 自定义回复格式
		putPolicy := storage.PutPolicy{
			Scope:      q.Bucket,
			ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)"}`,
		}
		upToken := putPolicy.UploadToken(mac)
		formUploader := storage.NewFormUploader(&q.Cfg)
		ret := MyPutRet{}
		putExtra := storage.PutExtra{
			Params: map[string]string{
				"x:name": "github logo",
			},
		}
		err := formUploader.PutFile(context.Background(), &ret, upToken, key, filePath, &putExtra)
		// err = formUploader.Put(context.Background(), &ret, upToken, "test", bytes.NewReader(data), int64(len(data)), &putExtra)
		// err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
		if err != nil {
			logx.Error("上传失败:", err)

		}
		err = os.Remove(filePath)
		if err != nil {
			logx.Error("删除文件时出错:", err)

		}
		logx.Debug("文件删除成功:", filePath)
		logx.Debug(ret.Bucket, ret.Key, ret.Fsize, ret.Hash, ret.Name)
	})
	wg.Wait()

	return ret.Key, ret.Name, nil
}

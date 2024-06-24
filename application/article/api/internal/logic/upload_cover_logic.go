package logic

import (
	"beyond/application/article/api/code"
	"context"
	"fmt"
	"net/http"
	"time"

	"beyond/application/article/api/internal/svc"
	"beyond/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *http.Request) (resp *types.UploadCoverResp, err error) {
	_ = req.ParseMultipartForm(maxFileSize)
	file, header, err := req.FormFile("cover")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	objectKey := genFilename(header.Filename)

	uploadFile, s, err := l.svcCtx.Oss.UploadFile(l.ctx, objectKey, file, header.Size)
	if err != nil {
		// logx.Errorf("get bucket failed, err: %v", err)
		l.Logger.Errorf("put object failed, err: %v", err)
		return nil, code.PutBucketErr
	}
	fmt.Println(uploadFile, s)

	return &types.UploadCoverResp{CoverUrl: genFileURL(objectKey)}, nil
}

func genFilename(filename string) string {
	return fmt.Sprintf("%d_%s", time.Now().UnixMilli(), filename)
}
func genFileURL(objectKey string) string {
	// return fmt.Sprintf("https://beyond-article.oss-cn-shanghai.aliyuncs.com/%s", objectKey)
	return fmt.Sprintf("http://cdn.toweron.top/%s", objectKey)
}

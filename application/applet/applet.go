package main

import (
	"beyond/pkg/xcode"
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"beyond/application/applet/internal/config"
	"beyond/application/applet/internal/handler"
	"beyond/application/applet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/applet-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 使用 gorilla/handlers 包的 CompressionHandler 中间件开启数据压缩
	// compressedHandler := handlers.CompressHandlerLevel(http.HandlerFunc(handler), 5)
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return handlers.CompressHandlerLevel(next, 9).(http.HandlerFunc)
	})
	/*server.Use(func(next http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {
			// next(w, r)

			w.Header().Set("Accept-Encoding", "gzip")
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")
			w.Header().Set("Content-Type", "text/html; charset=utf-8")

			// 创建一个新的响应缓冲区
			buffer := new(bytes.Buffer)

			// 创建一个新的 ResponseWriter 代理，将响应写入到缓冲区
			proxyWriter := &responseProxyWriter{writer: w, buffer: buffer}

			// 创建 gzip writer
			gz := gzip.NewWriter(proxyWriter)

			// gz := gzip.NewWriter(w)
			defer gz.Close()
			gz.Write(buffer.Bytes())

		}
	})*/
	// 自定义错误处理方法
	httpx.SetErrorHandler(xcode.ErrHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// ResponseWriter 代理，将响应同时写入到原始 ResponseWriter 和缓冲区
type responseProxyWriter struct {
	writer http.ResponseWriter
	buffer *bytes.Buffer
}

func (rw *responseProxyWriter) Header() http.Header {
	return rw.writer.Header()
}

func (rw *responseProxyWriter) Write(data []byte) (int, error) {
	// 将数据同时写入到原始 ResponseWriter 和缓冲区
	rw.buffer.Write(data)
	return rw.writer.Write(data)
}

func (rw *responseProxyWriter) WriteHeader(statusCode int) {
	rw.writer.WriteHeader(statusCode)
}

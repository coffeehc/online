package api

import (
	"crypto/subtle"
	"fmt"
	"github.com/rs/cors"
	"net/http"
	"online/common/log"
	"online/server/dbm"
	"strings"

	"time"
)

type APIManager struct {
	staticDir string
	dbManager *dbm.Manager
}

func (a *APIManager) HookHTTPMiddleware(handler http.Handler) http.HandlerFunc {
	// 服务前端文件
	var feHandler http.Handler
	if a.staticDir != "" {
		feHandler = http.FileServer(http.Dir(a.staticDir))
	} else {
		feHandler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(404)
		})
	}
	ins := cors.Default()
	_ = ins

	//basicAuth := func(h http.HandlerFunc) http.HandlerFunc {
	//	if a.basicAuthPassword != "" {
	//		return BasicAuth(
	//			h, "falcon", a.basicAuthPassword, "Falcon",
	//		).ServeHTTP
	//	}
	//	return h
	//}
	//feHandler = basicAuth(feHandler.ServeHTTP)

	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		defer func() {
			end := time.Now()
			duration := end.Sub(start)

			if request.Response != nil {
				log.Infof("[%v] [%v] %v len:%v cost:%v response:[%v] body_len:[%v]",
					request.RemoteAddr, request.Method, request.RequestURI, request.ContentLength, duration.String(),
					request.Response.StatusCode, request.Response.ContentLength,
				)
			} else {
				log.Infof("[%v] [%v] %v len:%v cost:%v",
					request.RemoteAddr, request.Method, request.RequestURI, request.ContentLength, duration.String(),
				)
			}
		}()

		if strings.HasPrefix(request.RequestURI, "/api") || strings.HasPrefix(request.RequestURI, "/api") {
			//ins.Handler(handler).ServeHTTP(writer, request)

			// 如果是 Electron 模式一定要解决跨源问题，虽然这样很不安全
			//if a.electronMode {
			//	writer.Header().Set("Access-Control-Allow-Origin", "*")
			//	writer.Header().Set("Access-Control-Allow-Headers", "*")
			//	writer.Header().Set("Access-Control-Allow-Methods", "*")
			//	if request.Method == "OPTIONS" {
			//		writer.WriteHeader(200)
			//		return
			//	}
			//}

			//if strings.HasPrefix(request.RequestURI, "/api/license") {
			//	// 如果请求本来就是 license 相关的，就放行，不进行 license 检测
			//	handler.ServeHTTP(writer, request)
			//	return
			//}
			//rsp, err := xlic.LoadAndVerifyLicense(a.dbManager.DB)
			//if err != nil {
			//	writer.WriteHeader(402)
			//	_, _ = writer.Write([]byte(`<h2>当前系统使用许可 license 已失效，请联系销售人员购买新的许可证或申请试用</h2>`))
			//	return
			//}
			//_ = rsp
			handler.ServeHTTP(writer, request)
		} else {
			shouldServeStatic := strings.HasPrefix(
				request.RequestURI, "/static") ||
				request.RequestURI == "/" ||
				strings.HasPrefix(request.RequestURI, "/index.html")
			if shouldServeStatic {
				feHandler.ServeHTTP(writer, request)
				//feHandler.ServeHTTP(writer, request)
			} else {
				request.URL.Path = "/"
				request.RequestURI = "/"
				feHandler.ServeHTTP(writer, request)
				//feHandler.ServeHTTP(writer, request)
			}
		}
	}
}

// BasicAuth wraps a handler requiring HTTP basic auth for it using the given
// username and password and the specified realm, which shouldn't contain quotes.
//
// Most web browser display a dialog with something like:
//
//    The website says: "<realm>"
//
// Which is really stupid so you may want to set the realm to a message rather than
// an actual realm.
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func NewAPIManager(postgresParams string, frontendDir string) (*APIManager, error) {
	m, err := dbm.NewDBManager(postgresParams)
	if err != nil {
		return nil, fmt.Errorf("build database conn failed; %v", err)
	}
	return &APIManager{
		dbManager: m,
		staticDir: frontendDir,
	}, nil
}

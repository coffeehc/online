package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type WebHookServer struct {
	r    *mux.Router
	s    *http.Server
	addr string
}

func NewWebHookServerEx(port int, cb func(data interface{})) *WebHookServer {
	addr := fmt.Sprintf("127.0.0.1:%v", port)

	r := mux.NewRouter()
	r.HandleFunc("/webhook", func(writer http.ResponseWriter, request *http.Request) {
		defer writer.WriteHeader(200)
		target := fmt.Sprintf("http://%v/webhook", addr)
		request.URL, _ = url.Parse(target)
		if request.Body == nil {
			request.Body = http.NoBody
			request.GetBody = func() (io.ReadCloser, error) {
				return http.NoBody, nil
			}
		}
		cb(request)
	})
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		defer writer.WriteHeader(200)
		target := fmt.Sprintf("http://%v/", addr)
		request.URL, _ = url.Parse(target)
		if request.Body == nil {
			request.Body = http.NoBody
			request.GetBody = func() (io.ReadCloser, error) {
				return http.NoBody, nil
			}
		}
		cb(request)
	})
	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	return &WebHookServer{
		r: r, s: server,
		addr: addr,
	}
}

func NewWebHookServer(port int, cb func(data interface{})) *WebHookServer {
	r := mux.NewRouter()
	r.HandleFunc("/webhook", func(writer http.ResponseWriter, request *http.Request) {
		defer writer.WriteHeader(200)

		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return
		}
		cb(body)
	})
	addr := fmt.Sprintf("127.0.0.1:%v", port)
	server := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	return &WebHookServer{
		r: r, s: server,
		addr: addr,
	}

}

func (w *WebHookServer) Start() {
	go func() {
		err := w.s.ListenAndServe()
		if err != nil {
			//log.Errorf("serve failed: %s", err)
			//panic(err)
		}
	}()
}

func (w *WebHookServer) Shutdown() {
	_ = w.s.Shutdown(TimeoutContext(1 * time.Second))
}

func (w *WebHookServer) Addr() string {
	return fmt.Sprintf("http://%v/webhook", w.addr)
}

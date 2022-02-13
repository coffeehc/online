package api

import (
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
	"net/http"
	"online/common/log"
	"online/server/web/gen/restapi"
	"online/server/web/gen/restapi/operations"
	"time"
)

func GeneratePostgresParams(host string, port int, user string, password string) string {
	if host == "" {
		host = "127.0.0.1"
	}

	if port <= 0 {
		port = 5432
	}

	if user == "" {
		user = "root"
	}

	if password == "" {
		password = "password"
	}

	return fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=disable",
		host, port, user,
		"postgres", password,
	)
}

func StartServer(
	postgresParams string,
	port int,
	frontendDir string,
) error {
	log.Info("start to load web api...")
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return err
	}

	apiBase := operations.NewOnlineAPI(swaggerSpec)

	log.Info("start to create api manager...")
	manager, err := NewAPIManager(postgresParams, frontendDir)
	if err != nil {
		return err
	}

	apiBase.Logger = log.Infof
	log.Info("init api binding")
	manager.initAPIBase(apiBase)

	webServer := restapi.NewServer(apiBase)
	webServer.Port = port
	webServer.ConfigureAPI()
	webServer.SetHandler(
		manager.HookHTTPMiddleware(webServer.GetHandler()),
	)

	go func() {
		log.Infof("start to serve on: 0.0.0.0:%v", port)
		if err := webServer.Serve(); err != nil {
			log.Errorf("serve failed: %v", err)
		}
	}()

	ticker := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-ticker:
			url := fmt.Sprintf("http://127.0.0.1:%v", port)
			_, err := http.Get(url)
			if err != nil {
				continue
			}
			return nil
		case <-time.After(30 * time.Second):
			return errors.New("web service start timeout 30s!")
		}
	}
}

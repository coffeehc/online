package api

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"io/ioutil"
	"net/http"
	"net/url"
	"online/common/log"
	"online/server/api/webutil"
	"online/server/web/gen/restapi/operations"
	"strings"
	"time"
)

var secret string
var githubDailer = http.DefaultClient

func init() {
	raw, err := ioutil.ReadFile(".github-oauth-token")
	if err != nil {
		log.Errorf("load github-oauth-token failed: %s", err)
	}
	secret = strings.TrimSpace(string(raw))
}

func (a *APIManager) GetAuthFromGithubHandler() operations.GetAuthFromGithubHandlerFunc {
	apiName := "GetAuthFromGithubHandler"
	_ = apiName
	return func(params operations.GetAuthFromGithubParams) middleware.Responder {
		return nil
	}
}
func (a *APIManager) GetAuthFromGithubCallbackHandler() operations.GetAuthFromGithubCallbackHandlerFunc {
	apiName := "GetAuthFromGithubCallbackHandler"
	_ = apiName
	return func(params operations.GetAuthFromGithubCallbackParams) middleware.Responder {
		log.Infof("auth from code: %v", webutil.Str(params.Code))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		u, err := url.Parse(`https://github.com/login/oauth/access_token`)
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, err.Error())
		}

		values, err := url.ParseQuery("client_id=&client_secret=&code=")
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, err.Error())
		}
		values.Set("client_id", "ab2647aa8b16197389eb")
		values.Set("client_secret", secret)
		values.Set("code", webutil.Str(params.Code))
		u.RawQuery = values.Encode()
		req, err := http.NewRequestWithContext(ctx, "POST", u.String(), http.NoBody)
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, "auth to github failed: %s", err.Error())
		}
		rsp, err := githubDailer.Do(req)
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, err.Error())
		}
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, err.Error())
		}

		response, err := url.ParseQuery(string(body))
		if err != nil {
			return webutil.NewActionErrorResponder(apiName, err.Error())
		}

		if response.Get("error") != "" {
			return webutil.NewActionErrorResponder(apiName, "Github Oauth ERROR: %v Reason: %v", response.Get("error"), response.Get("error_description"))
		}

		// 登陆成功
		return webutil.NewActionSucceedResponder(apiName)
	}
}

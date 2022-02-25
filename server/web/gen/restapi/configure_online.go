// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"online/server/web/gen/models"
	"online/server/web/gen/restapi/operations"
)

//go:generate swagger generate server --target ../../gen --name Online --spec ../../../../swagger.yml --principal models.Principle

func configureFlags(api *operations.OnlineAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.OnlineAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	if api.TrustedAuth == nil {
		api.TrustedAuth = func(token string) (*models.Principle, error) {
			return nil, errors.NotImplemented("api key auth (trusted) Authorization from header param [Authorization] has not yet been implemented")
		}
	}
	// Applies when the "Authorization" header is set
	if api.UserAuth == nil {
		api.UserAuth = func(token string) (*models.Principle, error) {
			return nil, errors.NotImplemented("api key auth (user) Authorization from header param [Authorization] has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()

	if api.DeleteOperationHandler == nil {
		api.DeleteOperationHandler = operations.DeleteOperationHandlerFunc(func(params operations.DeleteOperationParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteOperation has not yet been implemented")
		})
	}
	if api.DeleteUserHandler == nil {
		api.DeleteUserHandler = operations.DeleteUserHandlerFunc(func(params operations.DeleteUserParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteUser has not yet been implemented")
		})
	}
	if api.DeleteYakitPluginHandler == nil {
		api.DeleteYakitPluginHandler = operations.DeleteYakitPluginHandlerFunc(func(params operations.DeleteYakitPluginParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteYakitPlugin has not yet been implemented")
		})
	}
	if api.GetAuthFromGithubHandler == nil {
		api.GetAuthFromGithubHandler = operations.GetAuthFromGithubHandlerFunc(func(params operations.GetAuthFromGithubParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetAuthFromGithub has not yet been implemented")
		})
	}
	if api.GetAuthFromGithubCallbackHandler == nil {
		api.GetAuthFromGithubCallbackHandler = operations.GetAuthFromGithubCallbackHandlerFunc(func(params operations.GetAuthFromGithubCallbackParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetAuthFromGithubCallback has not yet been implemented")
		})
	}
	if api.GetOperationHandler == nil {
		api.GetOperationHandler = operations.GetOperationHandlerFunc(func(params operations.GetOperationParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetOperation has not yet been implemented")
		})
	}
	if api.GetUserHandler == nil {
		api.GetUserHandler = operations.GetUserHandlerFunc(func(params operations.GetUserParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUser has not yet been implemented")
		})
	}
	if api.GetUserFetchHandler == nil {
		api.GetUserFetchHandler = operations.GetUserFetchHandlerFunc(func(params operations.GetUserFetchParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUserFetch has not yet been implemented")
		})
	}
	if api.GetUserTagsHandler == nil {
		api.GetUserTagsHandler = operations.GetUserTagsHandlerFunc(func(params operations.GetUserTagsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetUserTags has not yet been implemented")
		})
	}
	if api.GetYakitPluginHandler == nil {
		api.GetYakitPluginHandler = operations.GetYakitPluginHandlerFunc(func(params operations.GetYakitPluginParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetYakitPlugin has not yet been implemented")
		})
	}
	if api.GetYakitPluginFetchHandler == nil {
		api.GetYakitPluginFetchHandler = operations.GetYakitPluginFetchHandlerFunc(func(params operations.GetYakitPluginFetchParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetYakitPluginFetch has not yet been implemented")
		})
	}
	if api.GetYakitPluginTagsHandler == nil {
		api.GetYakitPluginTagsHandler = operations.GetYakitPluginTagsHandlerFunc(func(params operations.GetYakitPluginTagsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetYakitPluginTags has not yet been implemented")
		})
	}
	if api.PostOperationHandler == nil {
		api.PostOperationHandler = operations.PostOperationHandlerFunc(func(params operations.PostOperationParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostOperation has not yet been implemented")
		})
	}
	if api.PostUserTagsHandler == nil {
		api.PostUserTagsHandler = operations.PostUserTagsHandlerFunc(func(params operations.PostUserTagsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostUserTags has not yet been implemented")
		})
	}
	if api.PostYakitPluginHandler == nil {
		api.PostYakitPluginHandler = operations.PostYakitPluginHandlerFunc(func(params operations.PostYakitPluginParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostYakitPlugin has not yet been implemented")
		})
	}
	if api.PostYakitPluginTagsHandler == nil {
		api.PostYakitPluginTagsHandler = operations.PostYakitPluginTagsHandlerFunc(func(params operations.PostYakitPluginTagsParams, principal *models.Principle) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostYakitPluginTags has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

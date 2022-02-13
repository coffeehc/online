package api

import (
	"github.com/go-openapi/runtime/middleware"
	"online/server/web/gen/models"
	"online/server/web/gen/restapi/operations"
)

func (a *APIManager) GetYakitPluginHandler() operations.GetYakitPluginHandlerFunc {
	apiName := "GetYakitPluginHandler"
	_ = apiName
	return func(params operations.GetYakitPluginParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

func (a *APIManager) GetYakitPluginFetchHandler() operations.GetYakitPluginFetchHandlerFunc {
	apiName := "GetYakitPluginFetchHandler"
	_ = apiName
	return func(params operations.GetYakitPluginFetchParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

func (a *APIManager) GetYakitPluginTagsHandler() operations.GetYakitPluginTagsHandlerFunc {
	apiName := "GetYakitPluginTagsHandler"
	_ = apiName
	return func(params operations.GetYakitPluginTagsParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

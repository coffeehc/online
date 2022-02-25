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

func (a *APIManager) PostYakitPluginTagsHandler() operations.PostYakitPluginTagsHandlerFunc {
	apiName := "PostYakitPluginTagsHandler"
	_ = apiName
	return func(params operations.PostYakitPluginTagsParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

func (a *APIManager) PostYakitPluginHandler() operations.PostYakitPluginHandlerFunc {
	apiName := "PostYakitPluginHandler"
	_ = apiName
	return func(params operations.PostYakitPluginParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) DeleteYakitPluginHandler() operations.DeleteYakitPluginHandlerFunc {
	apiName := "DeleteYakitPluginHandler"
	_ = apiName
	return func(params operations.DeleteYakitPluginParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

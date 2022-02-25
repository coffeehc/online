package api

import (
	"github.com/go-openapi/runtime/middleware"
	"online/server/web/gen/models"
	"online/server/web/gen/restapi/operations"
)

func (a *APIManager) GetOperationHandler() operations.GetOperationHandlerFunc {
	apiName := "GetOperationHandler"
	_ = apiName
	return func(params operations.GetOperationParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) PostOperationHandler() operations.PostOperationHandlerFunc {
	apiName := "PostOperationHandler"
	_ = apiName
	return func(params operations.PostOperationParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) DeleteOperationHandler() operations.DeleteOperationHandlerFunc {
	apiName := "DeleteOperationHandler"
	_ = apiName
	return func(params operations.DeleteOperationParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

package api

import (
	"github.com/go-openapi/runtime/middleware"
	"online/server/web/gen/models"
	"online/server/web/gen/restapi/operations"
)

func (a *APIManager) GetUserFetchHandler() operations.GetUserFetchHandlerFunc {
	apiName := "GetUserFetchHandler"
	_ = apiName
	return func(params operations.GetUserFetchParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) GetUserHandler() operations.GetUserHandlerFunc {
	apiName := "GetUserHandler"
	_ = apiName
	return func(params operations.GetUserParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) GetUserTagsHandler() operations.GetUserTagsHandlerFunc {
	apiName := "GetUserTagsHandler"
	_ = apiName
	return func(params operations.GetUserTagsParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) PostUserTagsHandler() operations.PostUserTagsHandlerFunc {
	apiName := "PostUserTagsHandler"
	_ = apiName
	return func(params operations.PostUserTagsParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}
func (a *APIManager) DeleteUserHandler() operations.DeleteUserHandlerFunc {
	apiName := "DeleteUserHandler"
	_ = apiName
	return func(params operations.DeleteUserParams, principle *models.Principle) middleware.Responder {
		return nil
	}
}

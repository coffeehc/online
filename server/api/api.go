package api

import "online/server/web/gen/restapi/operations"

func (s *APIManager) initAPIBase(base *operations.OnlineAPI) {
	base.GetYakitPluginHandler = s.GetYakitPluginHandler()
	base.GetYakitPluginFetchHandler = s.GetYakitPluginFetchHandler()
	base.GetYakitPluginTagsHandler = s.GetYakitPluginTagsHandler()
	base.PostYakitPluginTagsHandler = s.PostYakitPluginTagsHandler()
	base.PostYakitPluginHandler = s.PostYakitPluginHandler()
	base.DeleteYakitPluginHandler = s.DeleteYakitPluginHandler()

	// 记录操作
	base.GetOperationHandler = s.GetOperationHandler()
	base.PostOperationHandler = s.PostOperationHandler()
	base.DeleteOperationHandler = s.DeleteOperationHandler()

	//
	base.GetUserFetchHandler = s.GetUserFetchHandler()
	base.GetUserHandler = s.GetUserHandler()
	base.GetUserTagsHandler = s.GetUserTagsHandler()
	base.PostUserTagsHandler = s.PostUserTagsHandler()
	base.DeleteUserHandler = s.DeleteUserHandler()

	// 认证
	base.GetAuthFromGithubHandler = s.GetAuthFromGithubHandler()
	base.GetAuthFromGithubCallbackHandler = s.GetAuthFromGithubCallbackHandler()
}

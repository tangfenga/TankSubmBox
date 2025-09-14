package rest

import (
	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/tool/result"
	"net/http"
)

type SubmissionController struct {
	BaseController
	submissionDao *SubmissionDao
	userProfileDao *UserProfileDao
}

func (this *SubmissionController) Init() {
	this.BaseController.Init()

	b := core.CONTEXT.GetBean(this.submissionDao)
	if b, ok := b.(*SubmissionDao); ok {
		this.submissionDao = b
	}

	b = core.CONTEXT.GetBean(this.userProfileDao)
	if b, ok := b.(*UserProfileDao); ok {
		this.userProfileDao = b
	}
}

func (this *SubmissionController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {
	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))
	routeMap["/api/submission/my"] = this.Wrap(this.MySubmission, USER_ROLE_USER)
	return routeMap
}

// 获取当前用户的提交
func (this *SubmissionController) MySubmission(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	user := this.checkUser(request)
	
	// 获取用户档案信息
	userProfile := this.userProfileDao.FindByUserUuid(user.Uuid)
	if userProfile == nil {
		return this.Success(nil)
	}
	
	// 根据学号查找提交
	submissions := this.submissionDao.FindByAuthorId(userProfile.StudentId)
	if len(submissions) == 0 {
		return this.Success(nil)
	}
	
	// 返回最新的提交
	return this.Success(submissions[0])
}
package rest

import (
	"net/http"
	"time"

	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/tool/result"
)

type SubmissionController struct {
	BaseController
	submissionDao *SubmissionDao
	userProfileDao *UserProfileDao
	collegeDao    *CollegeDao
	matterDao     *MatterDao
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

	b = core.CONTEXT.GetBean(this.collegeDao)
	if b, ok := b.(*CollegeDao); ok {
		this.collegeDao = b
	}

	b = core.CONTEXT.GetBean(this.matterDao)
	if b, ok := b.(*MatterDao); ok {
		this.matterDao = b
	}
}

func (this *SubmissionController) RegisterRoutes() map[string]func(writer http.ResponseWriter, request *http.Request) {
	routeMap := make(map[string]func(writer http.ResponseWriter, request *http.Request))
	routeMap["/api/submission/my"] = this.Wrap(this.MySubmission, USER_ROLE_USER)
	routeMap["/api/submission/recommend"] = this.Wrap(this.RecommendSubmission, USER_ROLE_COLLEGE_ADMIN)
	routeMap["/api/submission/by-matter"] = this.Wrap(this.GetSubmissionByMatter, USER_ROLE_USER)
	routeMap["/api/submission/recommended-by-college"] = this.Wrap(this.GetRecommendedByCollege, USER_ROLE_COLLEGE_ADMIN)
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

// 推荐作品
func (this *SubmissionController) RecommendSubmission(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	user := this.checkUser(request)
	
	// 从POST参数获取matterUuid
	matterUuid := request.FormValue("matterUuid")
	if matterUuid == "" {
		return result.BadRequest("matterUuid参数不能为空")
	}
	
	// 查找对应的提交
	submission := this.submissionDao.FindByMatterUuid(matterUuid)
	if submission == nil {
		return result.BadRequest("未找到对应的作品提交")
	}
	
	// 更新推荐状态
	submission.IsRecommended = true
	submission.RecommendedBy = user.Uuid
	submission.RecommendedAt = time.Now()
	
	this.submissionDao.Save(submission)
	
	return this.Success("推荐成功")
}

// 根据作品UUID获取提交信息
func (this *SubmissionController) GetSubmissionByMatter(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	matterUuid := request.FormValue("matterUuid")
	if matterUuid == "" {
		return result.BadRequest("matterUuid参数不能为空")
	}

	// 查找对应的提交
	submission := this.submissionDao.FindByMatterUuid(matterUuid)
	if submission == nil {
		return this.Success(nil)
	}

	return this.Success(submission)
}

// 获取学院管理员所在学院的推荐作品
func (this *SubmissionController) GetRecommendedByCollege(writer http.ResponseWriter, request *http.Request) *result.WebResult {
	user := this.checkUser(request)
	
	// 获取用户档案信息
	userProfile := this.userProfileDao.FindByUserUuid(user.Uuid)
	if userProfile == nil {
		return result.BadRequest("未找到用户档案信息")
	}
	
	// 获取学院管理员的学院名称
	if userProfile.College == "" {
		return result.BadRequest("用户未设置学院信息")
	}
	
	// 获取所有推荐的作品
	allRecommendedSubmissions := this.submissionDao.FindAllRecommended()
	
	// 提取学院管理员所在学院的推荐作品matterUuid
	matterUuids := make([]string, 0)
	for _, submission := range allRecommendedSubmissions {
		// 根据matterUuid获取matter
		matter := this.matterDao.FindByUuid(submission.MatterUuid)
		if matter == nil {
			continue
		}
		
		// 根据matter的用户获取用户档案
		matterUserProfile := this.userProfileDao.FindByUserUuid(matter.UserUuid)
		if matterUserProfile == nil {
			continue
		}
		
		// 检查是否属于学院管理员所在的学院
		if matterUserProfile.College == userProfile.College {
			matterUuids = append(matterUuids, submission.MatterUuid)
		}
	}
	
	return this.Success(matterUuids)
}
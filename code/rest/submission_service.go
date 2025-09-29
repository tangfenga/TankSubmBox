package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eyebluecn/tank/code/core"
	"github.com/eyebluecn/tank/code/tool/result"
)

type SubmissionService struct {
	BaseBean
	submissionDao     *SubmissionDao
	userDao           *UserDao
	spaceDao          *SpaceDao
	matterDao         *MatterDao
	matterService     *MatterService
	userProfileDao    *UserProfileDao
	trackDao          *TrackDao
	collegeDao        *CollegeDao
	userService       *UserService
	preferenceService *PreferenceService
}

func (this *SubmissionService) Init() {
	this.BaseBean.Init()

	b := core.CONTEXT.GetBean(this.submissionDao)
	if c, ok := b.(*SubmissionDao); ok {
		this.submissionDao = c
	}
	b = core.CONTEXT.GetBean(this.userDao)
	if c, ok := b.(*UserDao); ok {
		this.userDao = c
	}
	b = core.CONTEXT.GetBean(this.spaceDao)
	if c, ok := b.(*SpaceDao); ok {
		this.spaceDao = c
	}
	b = core.CONTEXT.GetBean(this.matterDao)
	if c, ok := b.(*MatterDao); ok {
		this.matterDao = c
	}
	b = core.CONTEXT.GetBean(this.matterService)
	if c, ok := b.(*MatterService); ok {
		this.matterService = c
	}
	b = core.CONTEXT.GetBean(this.userProfileDao)
	if c, ok := b.(*UserProfileDao); ok {
		this.userProfileDao = c
	}
	b = core.CONTEXT.GetBean(this.trackDao)
	if c, ok := b.(*TrackDao); ok {
		this.trackDao = c
	}
	b = core.CONTEXT.GetBean(this.collegeDao)
	if c, ok := b.(*CollegeDao); ok {
		this.collegeDao = c
	}
	b = core.CONTEXT.GetBean(this.userService)
	if c, ok := b.(*UserService); ok {
		this.userService = c
	}
	b = core.CONTEXT.GetBean(this.preferenceService)
	if c, ok := b.(*PreferenceService); ok {
		this.preferenceService = c
	}
}

func (this *SubmissionService) AdminCreateSubmissionFolder(
	request *http.Request,
	username string,
	password string,
	realName string,
	studentId string,
	collegeName string,
	userType string,
	phoneNumber string,
	trackId int64,
	workName string,
) *Matter {

	// 1. Find or create user
	user := this.userDao.FindByUsername(username)
	if user == nil {
		preference := this.preferenceService.Fetch()
		user = this.userService.CreateUser(request, username, -1, preference.DefaultTotalSizeLimit, password, USER_ROLE_USER, collegeName, realName, phoneNumber, userType, studentId)
	}

	// 2. Get user's space and profile
	space := this.spaceDao.CheckByUuid(user.SpaceUuid)
	userProfile := this.userProfileDao.FindByUserUuid(user.Uuid)
	if userProfile == nil {
		// If user existed but profile didn't, create it.
		userProfile = &UserProfile{
			UserUuid:    user.Uuid,
			RealName:    realName,
			College:     collegeName,
			PhoneNumber: phoneNumber,
			UserType:    userType,
			StudentId:   studentId,
		}
		this.userProfileDao.Create(userProfile)
	}

	// 3. Get track
	track := this.trackDao.Find(trackId)
	if track == nil {
		panic(result.BadRequest("Track not found"))
	}

	// 4. Create submission folder
	submissionFolderName := fmt.Sprintf("%s-%s-%s", track.Name, workName, userProfile.StudentId)
	rootMatter := this.matterDao.CheckWithRootByUuid(MATTER_ROOT, space)

	submissionFolder := this.matterService.AtomicCreateDirectory(request, rootMatter, submissionFolderName, user, space)

	// 5. Create/Update submission record
	submission := this.submissionDao.FindByMatterUuid(submissionFolder.Uuid)
	var collegeId int64 = 0
	if userProfile.College != "" {
		college := this.collegeDao.FindByName(userProfile.College)
		if college != nil {
			collegeId = college.Id
		}
	}

	if submission == nil {
		submission = &Submission{
			MatterUuid:    submissionFolder.Uuid,
			TrackId:       trackId,
			CollegeId:     collegeId,
			Title:         workName,
			AuthorName:    userProfile.RealName,
			AuthorId:      userProfile.StudentId,
			IsRecommended: false,
			CreateTime:    submissionFolder.CreateTime,
			UpdateTime:    submissionFolder.UpdateTime,
		}
	} else {
		submission.TrackId = trackId
		submission.Title = workName
		submission.UpdateTime = time.Now()
	}
	this.submissionDao.Save(submission)

	return submissionFolder
}

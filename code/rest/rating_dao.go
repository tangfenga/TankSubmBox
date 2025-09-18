package rest

import (
	"github.com/eyebluecn/tank/code/core"
)

// @Service
type RatingDao struct {
	BaseDao
}

func (this *RatingDao) Init() {
	this.BaseDao.Init()
}

func (this *RatingDao) Create(rating *Rating) *Rating {
	var db = core.CONTEXT.GetDB().Create(rating)
	this.PanicError(db.Error)
	return rating
}

func (this *RatingDao) Save(rating *Rating) *Rating {
	var db = core.CONTEXT.GetDB().Save(rating)
	this.PanicError(db.Error)
	return rating
}

func (this *RatingDao) FindBySubmissionAndJudge(submissionId int64, judgeUuid string) *Rating {
	var rating Rating
	db := core.CONTEXT.GetDB().Where("submission_id = ? AND judge_uuid = ?", submissionId, judgeUuid).First(&rating)
	if db.Error != nil {
		return nil
	}
	return &rating
}

func (this *RatingDao) FindScoredSubmissionIdsByJudge(judgeUuid string) []int64 {
	var submissionIds []int64
	db := core.CONTEXT.GetDB().Model(&Rating{}).
		Where("judge_uuid = ?", judgeUuid).
		Pluck("submission_id", &submissionIds)
	if db.Error != nil {
		return []int64{}
	}
	return submissionIds
}

// 获取所有提交的平均评分统计
type RatingStats struct {
	SubmissionId  int64             `json:"submissionId"`
	MatterName    string            `json:"matterName"`
	MatterUuid    string            `json:"matterUuid"`
	AuthorName    string            `json:"authorName"`    // 作者姓名
	CollegeName   string            `json:"collegeName"`   // 学院名称
	TrackName     string            `json:"trackName"`     // 赛道名称
	PhoneNumber   string            `json:"phoneNumber"`   // 手机号
	StudentId     string            `json:"studentId"`     // 学号
	CategoryScores map[string]float64 `json:"categoryScores"` // 各评分类别的平均分
	TotalScore    float64           `json:"totalScore"`    // 总分
	JudgeCount    int64             `json:"judgeCount"`    // 评委数量
}

func (this *RatingDao) FindRatingStats() []RatingStats {
	var stats []RatingStats
	
	// 获取所有有评分的提交
	var submissionIds []int64
	db := core.CONTEXT.GetDB().Model(&Rating{}).
		Distinct("submission_id").
		Pluck("submission_id", &submissionIds)
	if db.Error != nil {
		return []RatingStats{}
	}
	
	for _, submissionId := range submissionIds {
		// 获取提交信息
		submissionDao := &SubmissionDao{}
		submission := submissionDao.FindById(submissionId)
		if submission == nil {
			continue
		}
		
		// 获取该提交的所有评分
		var ratings []Rating
		db := core.CONTEXT.GetDB().Where("submission_id = ?", submissionId).Find(&ratings)
		if db.Error != nil {
			continue
		}
		
		// 按评分类别分组计算平均分
		categoryScores := make(map[string]float64)
		categoryCounts := make(map[string]int)
		
		for _, rating := range ratings {
			if rating.Comment != "" {
				categoryScores[rating.Comment] += float64(rating.Score)
				categoryCounts[rating.Comment]++
			}
		}
		
		// 计算每个类别的平均分
		categoryAverages := make(map[string]float64)
		totalScore := 0.0
		for category, sum := range categoryScores {
			if count := categoryCounts[category]; count > 0 {
				average := sum / float64(count)
				categoryAverages[category] = round(average, 2)
				totalScore += average
			}
		}
		
		// 获取matter信息
		matterDao := &MatterDao{}
		matter := matterDao.FindByUuid(submission.MatterUuid)
		
		// 获取赛道名称（从submission）
		trackDao := &TrackDao{}
		track := trackDao.Find(submission.TrackId)
		trackName := ""
		if track != nil {
			trackName = track.Name
		}
		
		// 获取用户信息（手机号、学号、学院）- 从matter的userUuid获取
		phoneNumber := ""
		studentId := ""
		collegeName := ""
		
		if matter != nil && matter.UserUuid != "" {
			userProfileDao := &UserProfileDao{}
			userProfile := userProfileDao.FindByUserUuid(matter.UserUuid)
			if userProfile != nil {
				phoneNumber = userProfile.PhoneNumber
				studentId = userProfile.StudentId
				collegeName = userProfile.College
			}
		}

		stats = append(stats, RatingStats{
			SubmissionId:   submissionId,
			MatterName:     submission.Title,
			MatterUuid:     submission.MatterUuid,
			AuthorName:     submission.AuthorName,
			CollegeName:    collegeName,
			TrackName:      trackName,
			PhoneNumber:    phoneNumber,
			StudentId:      studentId,
			CategoryScores: categoryAverages,
			TotalScore:     round(totalScore, 2),
			JudgeCount:     int64(len(ratings)),
		})
	}
	
	return stats
}

// 辅助函数：四舍五入到指定小数位
func round(value float64, precision int) float64 {
	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}
	return float64(int(value*multiplier+0.5)) / multiplier
}


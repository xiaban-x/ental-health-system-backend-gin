package models

import (
	"encoding/json"
	"errors"
	"time"
)

// PsychTest 心理测评模型
type PsychTest struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at" gorm:"index"`
	Title       string         `gorm:"size:100" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Category    string         `gorm:"size:50" json:"category"`
	Duration    int            `json:"duration"` // 预计完成时间（分钟）
	Status      string         `gorm:"size:20;default:'active'" json:"status"`
	CreatorID   uint           `json:"creator_id"`
	StartDate   *time.Time     `json:"start_date"`                   // 测评开始日期
	EndDate     *time.Time     `json:"end_date"`                     // 测评结束日期
	ScoreRules  string         `gorm:"type:text" json:"score_rules"` // JSON格式的评分规则
	Questions   []TestQuestion `gorm:"foreignKey:TestID" json:"questions,omitempty"`
	Results     []TestResult   `gorm:"foreignKey:TestID" json:"results,omitempty"`
}

// TestQuestion 测试问题模型
type TestQuestion struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
	TestID       uint       `json:"test_id"`
	QuestionNo   int        `json:"question_no"`               // 问题序号
	GroupID      uint       `json:"group_id"`                  // 问题分组ID
	GroupName    string     `gorm:"size:50" json:"group_name"` // 分组名称
	Content      string     `gorm:"type:text" json:"content"`
	Type         string     `gorm:"size:20" json:"type"` // 单选、多选、量表、开放式等
	Options      string     `gorm:"type:text" json:"options"`
	Score        int        `json:"score"` // 分值权重
	Required     bool       `gorm:"default:true" json:"required"`
	Dependencies string     `gorm:"type:text" json:"dependencies"` // JSON格式的问题依赖关系
}

// TestResult 测试结果模型
type TestResult struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at" gorm:"index"`
	TestID         uint       `json:"test_id"`
	StudentID      uint       `json:"student_id"`
	Student        Student    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Answers        string     `gorm:"type:text" json:"answers"`
	Score          int        `json:"score"`
	GroupScores    string     `gorm:"type:text" json:"group_scores"` // JSON格式的分组得分
	Result         string     `gorm:"type:text" json:"result"`
	Analysis       string     `gorm:"type:text" json:"analysis"`
	Recommendation string     `gorm:"type:text" json:"recommendation"`
	CompletedAt    time.Time  `json:"completed_at"`
}

// IsValid 检查测评是否在有效期内
func (pt *PsychTest) IsValid() bool {
	now := time.Now()
	if pt.StartDate != nil && now.Before(*pt.StartDate) {
		return false
	}
	if pt.EndDate != nil && now.After(*pt.EndDate) {
		return false
	}
	return pt.Status == "active"
}

// CalculateScore 计算测评得分
func (tr *TestResult) CalculateScore(test *PsychTest, answers map[string]interface{}) error {
	// 解析评分规则
	var scoreRules map[string]interface{}
	if err := json.Unmarshal([]byte(test.ScoreRules), &scoreRules); err != nil {
		return err
	}

	// 初始化分组得分
	groupScores := make(map[string]int)
	totalScore := 0

	// 计算每个问题的得分
	for _, question := range test.Questions {
		answerKey := string(question.ID)
		if answer, ok := answers[answerKey]; ok {
			// 根据问题类型和评分规则计算得分
			score := calculateQuestionScore(question, answer, scoreRules)
			totalScore += score

			// 更新分组得分
			if question.GroupID > 0 {
				groupScores[question.GroupName] += score
			}
		} else if question.Required {
			return errors.New("missing required answer")
		}
	}

	// 保存得分结果
	groupScoresJSON, err := json.Marshal(groupScores)
	if err != nil {
		return err
	}

	tr.Score = totalScore
	tr.GroupScores = string(groupScoresJSON)
	return nil
}

// calculateQuestionScore 计算单个问题的得分
func calculateQuestionScore(question TestQuestion, answer interface{}, scoreRules map[string]interface{}) int {
	switch question.Type {
	case "single":
		// 单选题得分计算
		return question.Score
	case "multiple":
		// 多选题得分计算
		if answers, ok := answer.([]interface{}); ok {
			return len(answers) * question.Score
		}
	case "scale":
		// 量表题得分计算
		if value, ok := answer.(float64); ok {
			return int(value) * question.Score
		}
	}
	return 0
}

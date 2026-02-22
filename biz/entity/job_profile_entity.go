package entity

import "time"

type CompetenciesJson struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Weight int    `json:"weight"`
}

type JobProfile struct {
	UserID                 string             //所属用户ID
	JobProfileID           string             //岗位画像ID
	JobTitle               string             //岗位名称
	Competencies           []CompetenciesJson //岗位能力要求
	RedLineCondition       []string           //岗位红线要求(自动淘汰规则)
	AiAdjustmentSuggestion string             //岗位调优建议
	CreatedAt              time.Time
}

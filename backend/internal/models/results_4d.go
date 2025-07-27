package models

import (
	"time"
)

type Result4D struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	DrawResultID  uint      `json:"draw_result_id" gorm:"not null;index"`
	PrizeCategory string    `json:"prize_category" gorm:"type:varchar(20);not null"` // '1st', '2nd', '3rd', 'Starter', 'Consolation'
	WinningNumber string    `json:"winning_number" gorm:"type:varchar(4);not null"`
	CreatedAt     time.Time `json:"created_at"`
	// PrizeAmount   float64   `json:"prize_amount" gorm:"type:decimal(10,2)"`

	DrawResult DrawResult `json:"-" gorm:"foreignKey:DrawResultID"`
}

func (Result4D) TableName() string {
	return "results_4d"
}

const (
	Prize1st         = "1st"
	Prize2nd         = "2nd"
	Prize3rd         = "3rd"
	PrizeStarter     = "Starter"
	PrizeConsolation = "Consolation"
)

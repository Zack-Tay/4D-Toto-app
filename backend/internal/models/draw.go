package models

import (
	"time"
)

type DrawResult struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	DrawType   string    `json:"draw_type" gorm:"type:varchar(10);not null;index"` // '4D' or 'TOTO'
	DrawDate   time.Time `json:"draw_date" gorm:"type:date;not null;index"`
	DrawNumber int       `json:"draw_number" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Results4D   []Result4D  `json:"results_4d,omitempty" gorm:"foreignKey:DrawResultID"`
	ResultsToto *ResultToto `json:"results_toto,omitempty" gorm:"foreignKey:DrawResultID"`
}

func (DrawResult) TableName() string {
	return "draw_results"
}

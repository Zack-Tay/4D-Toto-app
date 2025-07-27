package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type ResultToto struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	DrawResultID     uint      `json:"draw_result_id" gorm:"not null;uniqueIndex"`
	WinningNumbers   IntArray  `json:"winning_numbers" gorm:"type:jsonb;not null"` // 6 numbers
	AdditionalNumber uint      `json:"additional_number" gorm:"not null"`
	TotalPrizePool   float64   `json:"total_prize_pool" gorm:"type:decimal(12, 2)"`
	CreatedAt        time.Time `json:"created_at"`

	DrawResult  DrawResult       `json:"-" gorm:"foreignKey:DrawResultID"`
	PrizeGroups []TotoPrizeGroup `json:"prize_groups" gorm:"foreignKey:ResultTotoID"`
}

func (ResultToto) TableName() string {
	return "results_toto"
}

type TotoPrizeGroup struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ResultTotoID    uint      `json:"result_toto_id" gorm:"not null;index"`
	GroupNumber     int       `json:"group_number" gorm:"not null"`
	ShareAmount     float64   `json:"share_amount" gorm:"type:decimal(12, 2);not null"`
	WinningShares   int       `json:"winning_shares" gorm:"not null"`
	TotalAmount     float64   `json:"total_amount" gorm:"type:decimal(12, 2);not null"`
	WinningCriteria string    `json:"winning_criteria" gorm:"type:varchar(100)"`
	CreatedAt       time.Time `json:"created_at"`

	ResultToto ResultToto `json:"-" gorm:"foreignKey:ResultTotoID"`
}

func (TotoPrizeGroup) TableName() string {
	return "toto_prize_groups"
}

type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into IntArray", value)
	}

	return json.Unmarshal(bytes, a)
}

var TotoPrizeCriteria = map[int]string{
	1: "6 winning numbers",
	2: "5 winning numbers + additional number",
	3: "5 winning numbers",
	4: "4 winning numbers + additional number",
	5: "4 winning numbers",
	6: "3 winning numbers + additional number",
	7: "3 winning numbers",
}

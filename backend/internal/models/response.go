package models

import "time"

type Result4DReponse struct {
	DrawDate   time.Time      `json:"draw_date"`
	DrawNumber int            `json:"draw_number"`
	Prizes     Result4DPrizes `json:"prizes"`
}

type Result4DPrizes struct {
	First        string   `json:"first"`
	Second       string   `json:"second"`
	Third        string   `json:"third"`
	Starters     []string `json:"starters"`
	Consolations []string `json:"consolations"`
}

type ResultTotoResponse struct {
	DrawDate         time.Time            `json:"draw_date"`
	DrawNumber       int                  `json:"draw_number"`
	WinningNumbers   []int                `json:"winning_numbers"`
	AdditionalNumber int                  `json:"additional_number"`
	TotalPrizePool   float64              `json:"total_prize_pool"`
	PrizeBreakdown   []TotoPrizeBreakdown `json:"prize_breakdown"`
}

type TotoPrizeBreakdown struct {
	Group           int     `json:"group"`
	ShareAmount     float64 `json:"share_amount"`
	WinningShares   int     `json:"winning_shares"`
	TotalAmount     float64 `json:"total_amount"`
	WinningCriteria string  `json:"winning_criteria"`
}

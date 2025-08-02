package models

type Result4DReponse struct {
	DrawDate     string   `json:"draw_date"`
	DrawNumber   string   `json:"draw_number"`
	First        string   `json:"first"`
	Second       string   `json:"second"`
	Third        string   `json:"third"`
	Starters     []string `json:"starters"`
	Consolations []string `json:"consolations"`
}

type Result4DPrizes struct {
	First        string   `json:"first"`
	Second       string   `json:"second"`
	Third        string   `json:"third"`
	Starters     []string `json:"starters"`
	Consolations []string `json:"consolations"`
}

type ResultTotoResponse struct {
	DrawDate         string               `json:"draw_date"`
	DrawNumber       string               `json:"draw_number"`
	WinningNumbers   []string             `json:"winning_numbers"`
	AdditionalNumber string               `json:"additional_number"`
	Group1Prize      string               `json:"group1_prize"`
	PrizeBreakdown   []TotoPrizeBreakdown `json:"prize_breakdown"`
}

type TotoPrizeBreakdown struct {
	Group         string `json:"group"`
	ShareAmount   string `json:"share_amount"`
	WinningShares string `json:"winning_shares"`
	// WinningCriteria string `json:"winning_criteria"`
}

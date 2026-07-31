package domain

type Stat struct {
	Code         string `json:"code"`
	NameRU       string `json:"name_ru"`
	NameEN       string `json:"name_en"`
	Category     string `json:"category"`
	IsPercentage bool   `json:"is_percentage"`
	DisplayOrder int    `json:"display_order"`
}

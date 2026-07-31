package domain

type Item struct {
	ID               int      `json:"id"`
	GameID           int      `json:"game_id"`
	NameRU           string   `json:"name_ru"`
	NameEN           string   `json:"name_en"`
	Type             string   `json:"type"`
	Subtype          string   `json:"subtype"`
	LevelRequirement int      `json:"level_requirement"`
	ClassRestriction []string `json:"class_restriction"`
	BaseStats        Stats    `json:"base_stats"`
	GemSlots         int      `json:"gem_slots"`
	SetID            *int     `json:"set_id"`
	IconURL          string   `json:"icon_url"`
	Rarity           string   `json:"rarity"`
}

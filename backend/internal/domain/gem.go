package domain

type Gem struct {
	ID      int    `json:"id"`
	GameID  int    `json:"game_id"`
	NameRU  string `json:"name_ru"`
	NameEN  string `json:"name_en"`
	Level   int    `json:"level"`
	Bonuses Stats  `json:"bonuses"`
	IconURL string `json:"icon_url"`
}

package domain

type Class struct {
	ID                   int      `json:"id"`
	Code                 string   `json:"code"`
	NameRU               string   `json:"name_ru"`
	NameEN               string   `json:"name_en"`
	BaseStats            Stats    `json:"base_stats"`
	StatGrowthPerLevel   Stats    `json:"stat_growth_per_level"`
	AllowedEquipmentTypes []string `json:"allowed_equipment_types"`
	IconURL              string   `json:"icon_url"`
}

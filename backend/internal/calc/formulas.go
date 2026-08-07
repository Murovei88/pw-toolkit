package calc

// Константы для расчётов
const (
	PanguSoulMultiplier = 10.0  // Каждая душа Пань Гу даёт +10 к стату
	StarDiskMultiplier  = 5.0   // Каждый звёздный диск даёт +5 к стату
)

// GetRefineMultiplier возвращает множитель для уровня заточки
// Формулы взяты из community research (pwdev.ru, pwtools)
func GetRefineMultiplier(level int) float64 {
	switch {
	case level <= 0:
		return 1.0
	case level == 1:
		return 1.10
	case level == 2:
		return 1.20
	case level == 3:
		return 1.30
	case level == 4:
		return 1.40
	case level == 5:
		return 1.50
	case level == 6:
		return 1.60
	case level == 7:
		return 1.70
	case level == 8:
		return 1.80
	case level == 9:
		return 1.90
	case level == 10:
		return 2.00
	case level == 11:
		return 2.20
	case level == 12:
		return 2.50
	default:
		return 2.50 // max
	}
}

// StatCode — константы для названий характеристик
const (
	StatHP          = "hp"
	StatMP          = "mp"
	StatAttack      = "attack"
	StatDefense     = "defense"
	StatCritRate    = "crit_rate"
	StatCritDamage  = "crit_damage"
	StatDodge       = "dodge"
	StatAccuracy    = "accuracy"
	StatMagicAttack = "magic_attack"
	StatMagicDef    = "magic_defense"
)

package calc

// PropTypeToStat маппит ExtraPropType ID из XML на внутренние ключи статов
var PropTypeToStat = map[int]string{
	// Базовые характеристики (STR, DEX, INT, VIT)
	2066: "str", 2040: "dex", 2056: "int", 2079: "vit", 387: "vit",
	3317: "str", 3321: "hp", 3331: "crit_rate", 3335: "phys_pen", 3338: "mag_pen",

	// HP / MP
	2129: "hp", 636: "hp_pct",
	3520: "hp", 3425: "hp", 3427: "mp",

	// Атака / Защита
	1403: "phys_atk", 2275: "phys_atk", 1425: "mag_atk",
	1349: "phys_def", 1942: "phys_def", 3424: "phys_def", 3426: "mag_def",

	// Крит
	2110: "crit_rate", 217: "crit_dmg", 3234: "crit_dmg", 3473: "crit_dmg",
	3522: "crit_rate", 3527: "crit_rate", 3528: "crit_rate", 3454: "crit_rate",

	// Скорость атаки
	421: "atk_speed", 331: "atk_speed", 3525: "atk_speed", 3469: "atk_speed",

	// Пробивание / Снижение урона
	3694: "phys_pen", 3696: "mag_pen", 3697: "dmg_reduce",

	// Урон по мобам / Дух
	3688: "mob_dmg", 3691: "mob_dmg_pct", 3461: "mob_dmg", 3462: "mob_def",
	3230: "spirit", 638: "spirit", 3411: "spirit", 3529: "spirit",

	// Заточка (Refine) и прочее
	3519: "phys_atk", 3517: "phys_atk", 3524: "phys_atk",
	358: "skill_level", 391: "hp_pct", 2076: "int", 1510: "skill_level",
	3558: "skill_level", 3606: "skill_id",
}

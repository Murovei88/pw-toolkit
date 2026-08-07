package calc

import (
	"encoding/xml"
	"regexp"
	"fmt"
	"math"
)

// AllSimulatorInfos корневой элемент XML
type AllSimulatorInfos struct {
	XMLName    xml.Name `xml:"AllSimulatorInfos"`
	Version    string   `xml:"VERSION,attr"`
	Character  CharacterSimulator `xml:"CharacterSimulator"`
	Equipment  []EquipmentDetail `xml:"EquipmentSimulator>EquipDetail"`
}

type CharacterSimulator struct {
	Type       int `xml:"Type,attr"`
	RealmLevel int `xml:"RealmLevel,attr"`
}

type EquipmentDetail struct {
	EquipID    int    `xml:"EquipID,attr"`
	EquipLevel int    `xml:"EquipLevel,attr"`
	EquipPos   int    `xml:"EquipPos,attr"`
	ExtraInfo  []ExtraDetail `xml:"ExtraInfo>ExtraDetail"`
}

type ExtraDetail struct {
	ExtraPropType int   `xml:"ExtraPropType,attr"`
	IsLocal       int   `xml:"IsLocal,attr"`
	FirstParam    int64 `xml:"FirstParam,attr"`
}

// DecodeFloatFromInt декодирует IEEE 754 float из int32
func DecodeFloatFromInt(val int64) float64 {
	return float64(math.Float32frombits(uint32(val)))
}

// CalculateStats рассчитывает характеристики из XML
func CalculateStats(xmlData []byte) (map[string]float64, error) {
	var info AllSimulatorInfos
	if // Fix invalid XML attributes (e.g., SoulInfo00,="8" -> SoulInfo00="8")
 re := regexp.MustCompile(`,\s*="`)
 xmlData = re.ReplaceAll(xmlData, []byte(`="`))
  err := xml.Unmarshal(xmlData, &info); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	stats := make(map[string]float64)
	
	// Базовые характеристики для Воина (Type=0) на 105 уровне
	level := 105
	baseHP := 500.0 + float64(level-1) * 45.0  // ~45 HP за уровень
	baseMP := 100.0 + float64(level-1) * 5.0   // ~5 MP за уровень

	// Агрегируем статы от экипировки
	for _, equip := range info.Equipment {
		for _, extra := range equip.ExtraInfo {
			statKey := getStatKey(extra.ExtraPropType)
			if statKey == "" {
				continue
			}

			var val float64
			// Проверяем IsLocal для IEEE 754 float
			if extra.IsLocal == 1 {
				val = DecodeFloatFromInt(extra.FirstParam)
			} else {
				val = float64(extra.FirstParam)
			}

			// Обрабатываем проценты правильно
			switch statKey {
			case "hp_pct", "mp_pct":
				// Для IsLocal=1 значения уже в виде 0.1 = 10%, конвертируем в проценты
				if extra.IsLocal == 1 {
					val = val * 100.0
				}
			case "atk_speed":
				// Скорость атаки - это интервал между атаками
				// Конвертируем в атаки в секунду: 1 / interval
				if extra.IsLocal == 1 {
					val = 1.0 / val
				}
			}

			stats[statKey] += val
		}
	}

	// Добавляем базовые значения
	stats["hp"] = baseHP + stats["hp"]
	stats["mp"] = baseMP + stats["mp"]

	// Применяем процентные бонусы к HP/MP
	if hpPct, ok := stats["hp_pct"]; ok {
		stats["hp"] = stats["hp"] * (1.0 + hpPct/100.0)
	}
	if mpPct, ok := stats["mp_pct"]; ok {
		stats["mp"] = stats["mp"] * (1.0 + mpPct/100.0)
	}

	return stats, nil
}

// getStatKey возвращает ключ стата по ExtraPropType
func getStatKey(propType int) string {
	mapping := map[int]string{
		// Базовые статы
		2066: "str",
		2040: "dex",
		2056: "int",
		2079: "vit",
		387:  "vit",
		
		// HP/MP
		2129: "hp",
		636:  "hp_pct",
		637:  "mp_pct",
		
		// Атака/Защита
		1403: "phys_atk",
		2275: "phys_atk",
		1425: "mag_atk",
		1349: "phys_def",
		1942: "phys_def",
		3424: "phys_def",
		3426: "mag_def",
		
		// Крит
		2110: "crit_rate",
		217:  "crit_dmg",
		3234: "crit_dmg",
		3473: "crit_dmg",
		
		// Скорость атаки
		421:  "atk_speed",
		331:  "atk_speed",
		3525: "atk_speed",
		3469: "atk_speed",
		
		// Пробивание/Снижение урона
		3694: "phys_pen",
		3696: "mag_pen",
		3697: "dmg_reduce",
		
		// Урон по мобам/Дух
		3688: "mob_dmg",
		3691: "mob_dmg_pct",
		3461: "mob_dmg",
		3462: "mob_def",
		3230: "spirit",
		638:  "spirit",
		3411: "spirit",
		3529: "spirit",
		
		// Навыки
		358:  "skill_level",
		391:  "skill_level",
		1510: "skill_level",
		3558: "skill_level",
		3606: "skill_id",
	}
	
	return mapping[propType]
}

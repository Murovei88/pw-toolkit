package parser

import (
	"encoding/xml"
	"fmt"
	"math"
	"regexp"
	"strings"
)

// AllSimulatorInfos - корневая структура
type AllSimulatorInfos struct {
	XMLName            xml.Name           `xml:"AllSimulatorInfos"`
	Version            string             `xml:"VERSION,attr"`
	CharacterSimulator CharacterSimulator `xml:"CharacterSimulator"`
	EquipmentSimulator EquipmentSimulator `xml:"EquipmentSimulator"`
	// Остальные симуляторы добавим позже
}

type CharacterSimulator struct {
	Type        string `xml:"Type,attr"`
	RealmLevel  int    `xml:"RealmLevel,attr"`
	LevelInfo   string `xml:"LevelInfo"` // Временно как строка
	BaseProperty string `xml:"BaseProperty"`
}

type EquipmentSimulator struct {
	Type  string        `xml:"Type,attr"`
	Items []EquipDetail `xml:"EquipDetail"`
}

type EquipDetail struct {
	EquipID     int    `xml:"EquipID,attr"`
	EquipLevel  int    `xml:"EquipLevel,attr"`
	EquipHoles  int    `xml:"EquipHoles,attr"`
	EquipPos    int    `xml:"EquipPos,attr"`
	ExtraInfo   ExtraInfo `xml:"ExtraInfo"`
}

type ExtraInfo struct {
	Details []ExtraDetail `xml:"ExtraDetail"`
}

type ExtraDetail struct {
	ExtraPropType int   `xml:"ExtraPropType,attr"`
	IsLocal       int   `xml:"IsLocal,attr"`
	FirstParam    int64 `xml:"FirstParam,attr"`
	SecondParam   int64 `xml:"SecondParam,attr"`
}

// DecodeFloatFromInt декодирует int32, в котором запакован float32 (IsLocal="1")
func DecodeFloatFromInt(val int64) float64 {
	return float64(math.Float32frombits(uint32(val)))
}

// ParseSimulatorXML парсит XML-файл калькулятора
func ParseSimulatorXML(xmlData []byte) (*AllSimulatorInfos, error) {
	// Хак: нормализуем динамические теги игры в статические
	// <EquipDetail00> -> <EquipDetail>
	// <ExtraDetail01> -> <ExtraDetail>
	xmlStr := strings.ReplaceAll(string(xmlData), `,="`, `="`)
	
	reDynamicTags := regexp.MustCompile(`<(EquipDetail|ExtraDetail|StonePos|StoneID|BasePros|Level0|SoulInfo|CardInfo|ActivityProperty|PointValue|RuneID|SkillLevel|Reward|TitleID)\d+`)
	xmlStr = reDynamicTags.ReplaceAllStringFunc(xmlStr, func(match string) string {
		// Оставляем только базовое имя тега
		re := regexp.MustCompile(`\d+$`)
		return re.ReplaceAllString(match, "")
	})
	
	// То же самое для закрывающих тегов
	reCloseTags := regexp.MustCompile(`</(EquipDetail|ExtraDetail|StonePos|StoneID|BasePros|Level0|SoulInfo|CardInfo|ActivityProperty|PointValue|RuneID|SkillLevel|Reward|TitleID)\d+>`)
	xmlStr = reCloseTags.ReplaceAllStringFunc(xmlStr, func(match string) string {
		re := regexp.MustCompile(`\d+>`)
		return re.ReplaceAllString(match, ">")
	})

	var infos AllSimulatorInfos
	if err := xml.Unmarshal([]byte(xmlStr), &infos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	return &infos, nil
}

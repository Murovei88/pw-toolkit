package calc

import (
	"github.com/murovei88/pw-toolkit/internal/domain"
)

// StatCalculator вычисляет характеристики билда
type StatCalculator struct {
	// Зависимости для получения данных о предметах, классах, etc.
	// В MVP используем данные из самого билда (уже загружены из БД)
}

// NewStatCalculator создаёт новый калькулятор
func NewStatCalculator() *StatCalculator {
	return &StatCalculator{}
}

// Calculate возвращает итоговые характеристики билда
func (c *StatCalculator) Calculate(build *domain.Build, classes []*domain.Class, items []*domain.Item, gems []*domain.Gem) domain.Stats {
	stats := make(domain.Stats)

	// 1. Base stats from class
	classStats := c.calculateClassStats(build.ClassID, build.Level, classes)
	c.addStats(stats, classStats)

	// 2. Equipment stats (items + refine + gems)
	equipmentStats := c.calculateEquipmentStats(build.Equipment, items, gems)
	c.addStats(stats, equipmentStats)

	// 3. Cards, Books, Titles bonuses (упрощённо для MVP)
	// В будущем: загрузить из БД и суммировать бонусы
	// cardsStats := c.calculateCardsStats(build.Cards, cards)
	// booksStats := c.calculateBooksStats(build.Books, books)
	// titlesStats := c.calculateTitlesStats(build.Titles, titles)

	// 4. Genie bonuses (упрощённо)
	// genieStats := c.calculateGenieStats(build.GenieID, genies)

	// 5. Pangu Souls & Star Disks
	if build.PanguSouls != nil {
		for stat, value := range build.PanguSouls {
			stats[stat] += float64(value) * PanguSoulMultiplier
		}
	}

	if build.StarDisks != nil {
		for stat, value := range build.StarDisks {
			stats[stat] += float64(value) * StarDiskMultiplier
		}
	}

	// 6. Apply percentage modifiers (в будущем)
	// stats = c.applyPercentageModifiers(stats, percentageBonuses)

	return stats
}

// calculateClassStats вычисляет базовые статы от класса и уровня
func (c *StatCalculator) calculateClassStats(classID, level int, classes []*domain.Class) domain.Stats {
	stats := make(domain.Stats)

	// Находим класс
	var class *domain.Class
	for _, cl := range classes {
		if cl.ID == classID {
			class = cl
			break
		}
	}

	if class == nil {
		return stats
	}

	// Base stats + growth * (level - 1)
	for stat, baseValue := range class.BaseStats {
		growth := class.StatGrowthPerLevel[stat]
		stats[stat] = baseValue + growth*float64(level-1)
	}

	return stats
}

// calculateEquipmentStats вычисляет статы от экипировки
func (c *StatCalculator) calculateEquipmentStats(equipment domain.Equipment, items []*domain.Item, gems []*domain.Gem) domain.Stats {
	stats := make(domain.Stats)

	if equipment == nil {
		return stats
	}

	// Для каждого слота экипировки
	for _, equipped := range equipment {
		// Находим предмет
		var item *domain.Item
		for _, it := range items {
			if it.ID == equipped.ItemID {
				item = it
				break
			}
		}

		if item == nil {
			continue
		}

		// Base stats предмета * refine multiplier
		multiplier := GetRefineMultiplier(equipped.RefineLevel)
		for stat, value := range item.BaseStats {
			stats[stat] += value * multiplier
		}

		// Бонусы от камней
		for _, gemID := range equipped.Gems {
			if gemID == 0 {
				continue // пустой слот
			}

			var gem *domain.Gem
			for _, g := range gems {
				if g.ID == gemID {
					gem = g
					break
				}
			}

			if gem != nil {
				for stat, value := range gem.Bonuses {
					stats[stat] += value
				}
			}
		}
	}

	return stats
}

// addStats добавляет статы из source в target
func (c *StatCalculator) addStats(target, source domain.Stats) {
	for stat, value := range source {
		target[stat] += value
	}
}

// CalculatePreview вычисляет статы для preview (без сохранения в БД)
func (c *StatCalculator) CalculatePreview(req *PreviewRequest, classes []*domain.Class, items []*domain.Item, gems []*domain.Gem) domain.Stats {
	build := &domain.Build{
		ClassID:    req.ClassID,
		Level:      req.Level,
		Equipment:  req.Equipment,
		Cards:      req.Cards,
		Books:      req.Books,
		GenieID:    req.GenieID,
		PanguSouls: req.PanguSouls,
		StarDisks:  req.StarDisks,
		Titles:     req.Titles,
	}

	return c.Calculate(build, classes, items, gems)
}

// PreviewRequest — DTO для live-preview расчёта
type PreviewRequest struct {
	ClassID    int               `json:"class_id"`
	Level      int               `json:"level"`
	Equipment  domain.Equipment  `json:"equipment"`
	Cards      []int             `json:"cards"`
	Books      []int             `json:"books"`
	GenieID    *int              `json:"genie_id"`
	PanguSouls domain.PanguSouls `json:"pangu_souls"`
	StarDisks  domain.StarDisks  `json:"star_disks"`
	Titles     []int             `json:"titles"`
}

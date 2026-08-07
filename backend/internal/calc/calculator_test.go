package calc

import (
	"testing"

	"github.com/murovei88/pw-toolkit/internal/domain"
)

func TestStatCalculator_CalculateClassStats(t *testing.T) {
	calc := NewStatCalculator()

	classes := []*domain.Class{
		{
			ID:     1,
			Code:   "warrior",
			NameRU: "Воин",
			BaseStats: domain.Stats{
				"hp":      500,
				"mp":      100,
				"attack":  20,
				"defense": 15,
			},
			StatGrowthPerLevel: domain.Stats{
				"hp":      50,
				"mp":      5,
				"attack":  3,
				"defense": 2,
			},
		},
	}

	tests := []struct {
		name     string
		classID  int
		level    int
		expected domain.Stats
	}{
		{
			name:    "Level 1 warrior",
			classID: 1,
			level:   1,
			expected: domain.Stats{
				"hp":      500,
				"mp":      100,
				"attack":  20,
				"defense": 15,
			},
		},
		{
			name:    "Level 85 warrior",
			classID: 1,
			level:   85,
			expected: domain.Stats{
				"hp":      500 + 50*84,   // 4700
				"mp":      100 + 5*84,    // 520
				"attack":  20 + 3*84,     // 272
				"defense": 15 + 2*84,     // 183
			},
		},
		{
			name:     "Unknown class",
			classID:  999,
			level:    85,
			expected: domain.Stats{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := calc.calculateClassStats(tt.classID, tt.level, classes)
			
			for stat, expectedValue := range tt.expected {
				actualValue := stats[stat]
				if actualValue != expectedValue {
					t.Errorf("stat %s: expected %v, got %v", stat, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestGetRefineMultiplier(t *testing.T) {
	tests := []struct {
		level    int
		expected float64
	}{
		{0, 1.0},
		{1, 1.10},
		{5, 1.50},
		{8, 1.80},
		{10, 2.00},
		{12, 2.50},
		{15, 2.50}, // max
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := GetRefineMultiplier(tt.level)
			if actual != tt.expected {
				t.Errorf("level %d: expected %v, got %v", tt.level, tt.expected, actual)
			}
		})
	}
}

func TestStatCalculator_CalculateEquipmentStats(t *testing.T) {
	calc := NewStatCalculator()

	items := []*domain.Item{
		{
			ID: 100,
			BaseStats: domain.Stats{
				"attack": 100,
			},
		},
	}

	gems := []*domain.Gem{
		{
			ID: 201,
			Bonuses: domain.Stats{
				"attack": 10,
			},
		},
		{
			ID: 202,
			Bonuses: domain.Stats{
				"attack": 15,
			},
		},
	}

	equipment := domain.Equipment{
		"WEAPON": domain.EquippedItem{
			ItemID:      100,
			RefineLevel: 8,
			Gems:        []int{201, 202, 0, 0}, // 2 камня + 2 пустых слота
		},
	}

	stats := calc.calculateEquipmentStats(equipment, items, gems)

	// Expected: 100 * 1.80 (refine +8) + 10 (gem1) + 15 (gem2) = 180 + 25 = 205
	expectedAttack := 100*1.80 + 10 + 15
	if stats["attack"] != expectedAttack {
		t.Errorf("expected attack %v, got %v", expectedAttack, stats["attack"])
	}
}

func TestStatCalculator_Calculate(t *testing.T) {
	calc := NewStatCalculator()

	classes := []*domain.Class{
		{
			ID: 1,
			BaseStats: domain.Stats{
				"hp":     500,
				"attack": 20,
			},
			StatGrowthPerLevel: domain.Stats{
				"hp":     50,
				"attack": 3,
			},
		},
	}

	items := []*domain.Item{
		{
			ID: 100,
			BaseStats: domain.Stats{
				"attack": 100,
			},
		},
	}

	gems := []*domain.Gem{}

	build := &domain.Build{
		ClassID: 1,
		Level:   85,
		Equipment: domain.Equipment{
			"WEAPON": domain.EquippedItem{
				ItemID:      100,
				RefineLevel: 8,
				Gems:        []int{},
			},
		},
		PanguSouls: domain.PanguSouls{
			"attack": 5, // +50 attack (5 * 10)
		},
	}

	stats := calc.Calculate(build, classes, items, gems)

	// Class: hp = 500 + 50*84 = 4700
	// Class: attack = 20 + 3*84 = 272
	// Weapon: attack = 100 * 1.80 = 180
	// Pangu: attack = 5 * 10 = 50
	// Total attack = 272 + 180 + 50 = 502

	expectedHP := 4700.0
	expectedAttack := 272.0 + 180.0 + 50.0

	if stats["hp"] != expectedHP {
		t.Errorf("hp: expected %v, got %v", expectedHP, stats["hp"])
	}

	if stats["attack"] != expectedAttack {
		t.Errorf("attack: expected %v, got %v", expectedAttack, stats["attack"])
	}
}

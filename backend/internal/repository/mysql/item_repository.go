package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/murovei88/pw-toolkit/internal/domain"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) FindByID(ctx context.Context, id int) (*domain.Item, error) {
	query := `
		SELECT id, game_id, name_ru, name_en, type, subtype, level_requirement,
		       class_restriction, base_stats, gem_slots, set_id, icon_url, rarity
		FROM items
		WHERE id = ?
	`

	var item domain.Item
	var classRestriction, baseStats []byte
	var iconURL, rarity sql.NullString
	var setID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.GameID, &item.NameRU, &item.NameEN,
		&item.Type, &item.Subtype, &item.LevelRequirement,
		&classRestriction, &baseStats, &item.GemSlots,
		&setID, &iconURL, &rarity,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan item: %w", err)
	}

	json.Unmarshal(classRestriction, &item.ClassRestriction)
	json.Unmarshal(baseStats, &item.BaseStats)
	if setID.Valid {
		v := int(setID.Int64)
		item.SetID = &v
	}
	if iconURL.Valid {
		item.IconURL = iconURL.String
	}
	if rarity.Valid {
		item.Rarity = rarity.String
	}

	return &item, nil
}

func (r *ItemRepository) FindByIDs(ctx context.Context, ids []int) ([]*domain.Item, error) {
	if len(ids) == 0 {
		return []*domain.Item{}, nil
	}

	// Строим placeholder для IN clause
	placeholders := ""
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, game_id, name_ru, name_en, type, subtype, level_requirement,
		       class_restriction, base_stats, gem_slots, set_id, icon_url, rarity
		FROM items
		WHERE id IN (%s)
	`, placeholders)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		var classRestriction, baseStats []byte
		var iconURL, rarity sql.NullString
		var setID sql.NullInt64

		err := rows.Scan(
			&item.ID, &item.GameID, &item.NameRU, &item.NameEN,
			&item.Type, &item.Subtype, &item.LevelRequirement,
			&classRestriction, &baseStats, &item.GemSlots,
			&setID, &iconURL, &rarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item row: %w", err)
		}

		json.Unmarshal(classRestriction, &item.ClassRestriction)
		json.Unmarshal(baseStats, &item.BaseStats)
		if setID.Valid {
			v := int(setID.Int64)
			item.SetID = &v
		}
		if iconURL.Valid {
			item.IconURL = iconURL.String
		}
		if rarity.Valid {
			item.Rarity = rarity.String
		}

		items = append(items, &item)
	}

	return items, nil
}

func (r *ItemRepository) List(ctx context.Context, filter domain.ItemFilter) ([]*domain.Item, error) {
	// Упрощённая реализация для MVP
	query := `
		SELECT id, game_id, name_ru, name_en, type, subtype, level_requirement,
		       class_restriction, base_stats, gem_slots, set_id, icon_url, rarity
		FROM items
		ORDER BY id
		LIMIT 1000
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		var item domain.Item
		var classRestriction, baseStats []byte
		var iconURL, rarity sql.NullString
		var setID sql.NullInt64

		err := rows.Scan(
			&item.ID, &item.GameID, &item.NameRU, &item.NameEN,
			&item.Type, &item.Subtype, &item.LevelRequirement,
			&classRestriction, &baseStats, &item.GemSlots,
			&setID, &iconURL, &rarity,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item row: %w", err)
		}

		json.Unmarshal(classRestriction, &item.ClassRestriction)
		json.Unmarshal(baseStats, &item.BaseStats)
		if setID.Valid {
			v := int(setID.Int64)
			item.SetID = &v
		}
		if iconURL.Valid {
			item.IconURL = iconURL.String
		}
		if rarity.Valid {
			item.Rarity = rarity.String
		}

		items = append(items, &item)
	}

	return items, nil
}

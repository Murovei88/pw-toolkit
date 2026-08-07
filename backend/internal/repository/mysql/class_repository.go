package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/murovei88/pw-toolkit/internal/domain"
)

type ClassRepository struct {
	db *sql.DB
}

func NewClassRepository(db *sql.DB) *ClassRepository {
	return &ClassRepository{db: db}
}

func (r *ClassRepository) FindByID(ctx context.Context, id int) (*domain.Class, error) {
	query := `
		SELECT id, code, name_ru, name_en, base_stats, stat_growth_per_level, 
		       allowed_equipment_types, icon_url
		FROM classes
		WHERE id = ?
	`

	var c domain.Class
	var baseStats, statGrowth, allowedTypes []byte
	var iconURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Code, &c.NameRU, &c.NameEN,
		&baseStats, &statGrowth, &allowedTypes, &iconURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan class: %w", err)
	}

	json.Unmarshal(baseStats, &c.BaseStats)
	json.Unmarshal(statGrowth, &c.StatGrowthPerLevel)
	json.Unmarshal(allowedTypes, &c.AllowedEquipmentTypes)
	if iconURL.Valid {
		c.IconURL = iconURL.String
	}

	return &c, nil
}

func (r *ClassRepository) FindByCode(ctx context.Context, code string) (*domain.Class, error) {
	query := `
		SELECT id, code, name_ru, name_en, base_stats, stat_growth_per_level, 
		       allowed_equipment_types, icon_url
		FROM classes
		WHERE code = ?
	`

	var c domain.Class
	var baseStats, statGrowth, allowedTypes []byte
	var iconURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&c.ID, &c.Code, &c.NameRU, &c.NameEN,
		&baseStats, &statGrowth, &allowedTypes, &iconURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan class: %w", err)
	}

	json.Unmarshal(baseStats, &c.BaseStats)
	json.Unmarshal(statGrowth, &c.StatGrowthPerLevel)
	json.Unmarshal(allowedTypes, &c.AllowedEquipmentTypes)
	if iconURL.Valid {
		c.IconURL = iconURL.String
	}

	return &c, nil
}

func (r *ClassRepository) List(ctx context.Context) ([]*domain.Class, error) {
	query := `
		SELECT id, code, name_ru, name_en, base_stats, stat_growth_per_level, 
		       allowed_equipment_types, icon_url
		FROM classes
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query classes: %w", err)
	}
	defer rows.Close()

	var classes []*domain.Class
	for rows.Next() {
		var c domain.Class
		var baseStats, statGrowth, allowedTypes []byte
		var iconURL sql.NullString

		err := rows.Scan(
			&c.ID, &c.Code, &c.NameRU, &c.NameEN,
			&baseStats, &statGrowth, &allowedTypes, &iconURL,
		)
		if err != nil {
			return nil, fmt.Errorf("scan class row: %w", err)
		}

		json.Unmarshal(baseStats, &c.BaseStats)
		json.Unmarshal(statGrowth, &c.StatGrowthPerLevel)
		json.Unmarshal(allowedTypes, &c.AllowedEquipmentTypes)
		if iconURL.Valid {
			c.IconURL = iconURL.String
		}

		classes = append(classes, &c)
	}

	return classes, nil
}

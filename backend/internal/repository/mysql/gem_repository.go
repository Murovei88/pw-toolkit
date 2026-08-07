package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/murovei88/pw-toolkit/internal/domain"
)

type GemRepository struct {
	db *sql.DB
}

func NewGemRepository(db *sql.DB) *GemRepository {
	return &GemRepository{db: db}
}

func (r *GemRepository) FindByID(ctx context.Context, id int) (*domain.Gem, error) {
	query := `
		SELECT id, game_id, name_ru, name_en, level, bonuses, icon_url
		FROM gems
		WHERE id = ?
	`

	var gem domain.Gem
	var bonuses []byte
	var iconURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&gem.ID, &gem.GameID, &gem.NameRU, &gem.NameEN,
		&gem.Level, &bonuses, &iconURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan gem: %w", err)
	}

	json.Unmarshal(bonuses, &gem.Bonuses)
	if iconURL.Valid {
		gem.IconURL = iconURL.String
	}

	return &gem, nil
}

func (r *GemRepository) FindByIDs(ctx context.Context, ids []int) ([]*domain.Gem, error) {
	if len(ids) == 0 {
		return []*domain.Gem{}, nil
	}

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
		SELECT id, game_id, name_ru, name_en, level, bonuses, icon_url
		FROM gems
		WHERE id IN (%s)
	`, placeholders)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query gems: %w", err)
	}
	defer rows.Close()

	var gems []*domain.Gem
	for rows.Next() {
		var gem domain.Gem
		var bonuses []byte
		var iconURL sql.NullString

		err := rows.Scan(
			&gem.ID, &gem.GameID, &gem.NameRU, &gem.NameEN,
			&gem.Level, &bonuses, &iconURL,
		)
		if err != nil {
			return nil, fmt.Errorf("scan gem row: %w", err)
		}

		json.Unmarshal(bonuses, &gem.Bonuses)
		if iconURL.Valid {
			gem.IconURL = iconURL.String
		}

		gems = append(gems, &gem)
	}

	return gems, nil
}

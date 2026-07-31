package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/murovei88/pw-toolkit/internal/domain"
)

type BuildRepository struct {
	db *sql.DB
}

func NewBuildRepository(db *sql.DB) *BuildRepository {
	return &BuildRepository{db: db}
}

// buildRow — вспомогательная структура для чтения из БД
type buildRow struct {
	ID              string
	InternalID      int64
	Name            sql.NullString
	ClassID         int
	Level           int
	Equipment       []byte
	Cards           []byte
	Books           []byte
	GenieID         sql.NullInt64
	PanguSouls      []byte
	StarDisks       []byte
	Titles          []byte
	CalculatedStats []byte
	ViewCount       int
	LastViewedAt    time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (r *BuildRepository) Create(ctx context.Context, build *domain.Build) error {
	equipmentJSON, err := json.Marshal(build.Equipment)
	if err != nil {
		return fmt.Errorf("marshal equipment: %w", err)
	}

	cardsJSON, _ := json.Marshal(build.Cards)
	booksJSON, _ := json.Marshal(build.Books)
	panguSoulsJSON, _ := json.Marshal(build.PanguSouls)
	starDisksJSON, _ := json.Marshal(build.StarDisks)
	titlesJSON, _ := json.Marshal(build.Titles)
	calculatedStatsJSON, _ := json.Marshal(build.CalculatedStats)

	query := `
		INSERT INTO builds (
			id, name, class_id, level, equipment, cards, books, 
			genie_id, pangu_souls, star_disks, titles, calculated_stats
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		build.ID,
		build.Name,
		build.ClassID,
		build.Level,
		equipmentJSON,
		cardsJSON,
		booksJSON,
		build.GenieID,
		panguSoulsJSON,
		starDisksJSON,
		titlesJSON,
		calculatedStatsJSON,
	)
	if err != nil {
		return fmt.Errorf("insert build: %w", err)
	}

	return nil
}

func (r *BuildRepository) FindByID(ctx context.Context, id string) (*domain.Build, error) {
	query := `
		SELECT 
			id, internal_id, name, class_id, level, equipment, cards, books,
			genie_id, pangu_souls, star_disks, titles, calculated_stats,
			view_count, last_viewed_at, created_at, updated_at
		FROM builds
		WHERE id = ?
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var b buildRow
	err := row.Scan(
		&b.ID, &b.InternalID, &b.Name, &b.ClassID, &b.Level,
		&b.Equipment, &b.Cards, &b.Books, &b.GenieID,
		&b.PanguSouls, &b.StarDisks, &b.Titles, &b.CalculatedStats,
		&b.ViewCount, &b.LastViewedAt, &b.CreatedAt, &b.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan build: %w", err)
	}

	// Десериализуем JSON
	var equipment domain.Equipment
	if err := json.Unmarshal(b.Equipment, &equipment); err != nil {
		return nil, fmt.Errorf("unmarshal equipment: %w", err)
	}

	var cards []int
	if len(b.Cards) > 0 {
		json.Unmarshal(b.Cards, &cards)
	}

	var books []int
	if len(b.Books) > 0 {
		json.Unmarshal(b.Books, &books)
	}

	var panguSouls domain.PanguSouls
	if len(b.PanguSouls) > 0 {
		json.Unmarshal(b.PanguSouls, &panguSouls)
	}

	var starDisks domain.StarDisks
	if len(b.StarDisks) > 0 {
		json.Unmarshal(b.StarDisks, &starDisks)
	}

	var titles []int
	if len(b.Titles) > 0 {
		json.Unmarshal(b.Titles, &titles)
	}

	var calculatedStats domain.Stats
	if len(b.CalculatedStats) > 0 {
		json.Unmarshal(b.CalculatedStats, &calculatedStats)
	}

	var genieID *int
	if b.GenieID.Valid {
		v := int(b.GenieID.Int64)
		genieID = &v
	}

	var name string
	if b.Name.Valid {
		name = b.Name.String
	}

	return &domain.Build{
		ID:              b.ID,
		InternalID:      b.InternalID,
		Name:            name,
		ClassID:         b.ClassID,
		Level:           b.Level,
		Equipment:       equipment,
		Cards:           cards,
		Books:           books,
		GenieID:         genieID,
		PanguSouls:      panguSouls,
		StarDisks:       starDisks,
		Titles:          titles,
		CalculatedStats: calculatedStats,
		ViewCount:       b.ViewCount,
		LastViewedAt:    b.LastViewedAt,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
	}, nil
}

func (r *BuildRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `
		UPDATE builds 
		SET view_count = view_count + 1, last_viewed_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *BuildRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM builds WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("build not found")
	}
	return nil
}

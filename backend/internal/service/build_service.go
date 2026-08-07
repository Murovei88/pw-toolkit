package service

import (
	"context"
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/murovei88/pw-toolkit/internal/calc"
	"github.com/murovei88/pw-toolkit/internal/domain"
)

var (
	ErrBuildNotFound = errors.New("build not found")
	ErrInvalidBuild  = errors.New("invalid build data")
)

type BuildService struct {
	repo       domain.BuildRepository
	calc       *calc.StatCalculator
	classRepo  domain.ClassRepository
	itemRepo   domain.ItemRepository
	gemRepo    domain.GemRepository
}

func NewBuildService(
	repo domain.BuildRepository,
	classRepo domain.ClassRepository,
	itemRepo domain.ItemRepository,
	gemRepo domain.GemRepository,
) *BuildService {
	return &BuildService{
		repo:      repo,
		calc:      calc.NewStatCalculator(),
		classRepo: classRepo,
		itemRepo:  itemRepo,
		gemRepo:   gemRepo,
	}
}

func (s *BuildService) CreateBuild(ctx context.Context, build *domain.Build) error {
	if err := s.validateBuild(build); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBuild, err)
	}

	id, err := gonanoid.New(10)
	if err != nil {
		return fmt.Errorf("generate nanoid: %w", err)
	}
	build.ID = id

	if build.Equipment == nil {
		build.Equipment = make(domain.Equipment)
	}

	// Calculate stats before saving
	stats, err := s.calculateStats(ctx, build)
	if err != nil {
		return fmt.Errorf("calculate stats: %w", err)
	}
	build.CalculatedStats = stats

	return s.repo.Create(ctx, build)
}

func (s *BuildService) GetBuild(ctx context.Context, id string) (*domain.Build, error) {
	if len(id) != 10 {
		return nil, ErrBuildNotFound
	}

	build, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find build: %w", err)
	}
	if build == nil {
		return nil, ErrBuildNotFound
	}

	go func() {
		ctx := context.Background()
		_ = s.repo.IncrementViewCount(ctx, id)
	}()

	return build, nil
}

// CalculatePreview вычисляет статы без сохранения (для live-preview)
func (s *BuildService) CalculatePreview(ctx context.Context, req *calc.PreviewRequest) (domain.Stats, error) {
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

	return s.calculateStats(ctx, build)
}

func (s *BuildService) calculateStats(ctx context.Context, build *domain.Build) (domain.Stats, error) {
	// Загружаем классы
	classes, err := s.classRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("load classes: %w", err)
	}

	// Загружаем предметы из экипировки
	itemIDs := s.extractItemIDs(build.Equipment)
	items, err := s.itemRepo.FindByIDs(ctx, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}

	// Загружаем камни
	gemIDs := s.extractGemIDs(build.Equipment)
	gems, err := s.gemRepo.FindByIDs(ctx, gemIDs)
	if err != nil {
		return nil, fmt.Errorf("load gems: %w", err)
	}

	// Вычисляем статы
	stats := s.calc.Calculate(build, classes, items, gems)

	return stats, nil
}

func (s *BuildService) extractItemIDs(equipment domain.Equipment) []int {
	ids := make([]int, 0)
	for _, equipped := range equipment {
		ids = append(ids, equipped.ItemID)
	}
	return ids
}

func (s *BuildService) extractGemIDs(equipment domain.Equipment) []int {
	ids := make([]int, 0)
	for _, equipped := range equipment {
		for _, gemID := range equipped.Gems {
			if gemID > 0 {
				ids = append(ids, gemID)
			}
		}
	}
	return ids
}

func (s *BuildService) validateBuild(b *domain.Build) error {
	if b.ClassID <= 0 {
		return fmt.Errorf("class_id is required")
	}
	if b.Level < 1 || b.Level > 150 {
		return fmt.Errorf("level must be between 1 and 150")
	}
	return nil
}

package service

import (
	"context"
	"errors"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/murovei88/pw-toolkit/internal/domain"
)

var (
	ErrBuildNotFound = errors.New("build not found")
	ErrInvalidBuild  = errors.New("invalid build data")
)

type BuildService struct {
	repo domain.BuildRepository
}

func NewBuildService(repo domain.BuildRepository) *BuildService {
	return &BuildService{repo: repo}
}

// CreateBuild создаёт новый билд с уникальным публичным ID
func (s *BuildService) CreateBuild(ctx context.Context, build *domain.Build) error {
	// Валидация базовых полей
	if err := s.validateBuild(build); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidBuild, err)
	}

	// Генерация публичного ID (10 символов, ~59 bits entropy)
	// gonanoid.New(size) — использует стандартный URL-safe алфавит
	id, err := gonanoid.New(10)
	if err != nil {
		return fmt.Errorf("generate nanoid: %w", err)
	}
	build.ID = id

	// Инициализация пустых значений
	if build.Equipment == nil {
		build.Equipment = make(domain.Equipment)
	}
	if build.CalculatedStats == nil {
		build.CalculatedStats = make(domain.Stats)
	}

	// TODO: в Фазе 2.4 подключим CalculationService
	// build.CalculatedStats = s.calcService.Calculate(build)

	return s.repo.Create(ctx, build)
}

// GetBuild возвращает билд по ID и увеличивает счётчик просмотров
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

	// Увеличиваем счётчик просмотров (асинхронно, чтобы не замедлять ответ)
	go func() {
		ctx := context.Background()
		_ = s.repo.IncrementViewCount(ctx, id)
	}()

	return build, nil
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

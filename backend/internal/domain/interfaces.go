package domain

import "context"

type BuildRepository interface {
	Create(ctx context.Context, build *Build) error
	FindByID(ctx context.Context, id string) (*Build, error)
	Update(ctx context.Context, build *Build) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*Build, error)
	IncrementViewCount(ctx context.Context, id string) error
}

type ItemRepository interface {
	FindByID(ctx context.Context, id int) (*Item, error)
	FindByIDs(ctx context.Context, ids []int) ([]*Item, error)
	List(ctx context.Context, filter ItemFilter) ([]*Item, error)
}

type ClassRepository interface {
	FindByID(ctx context.Context, id int) (*Class, error)
	FindByCode(ctx context.Context, code string) (*Class, error)
	List(ctx context.Context) ([]*Class, error)
}

type GemRepository interface {
	FindByID(ctx context.Context, id int) (*Gem, error)
	FindByIDs(ctx context.Context, ids []int) ([]*Gem, error)
}

type ItemFilter struct {
	Type     string
	Subtype  string
	MinLevel int
	MaxLevel int
}

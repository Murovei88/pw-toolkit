package domain

import "time"

type Build struct {
	ID             string    `json:"id"`
	InternalID     int64     `json:"internal_id"`
	Name           string    `json:"name"`
	ClassID        int       `json:"class_id"`
	Level          int       `json:"level"`
	Equipment      Equipment `json:"equipment"`
	Cards          []int     `json:"cards"`
	Books          []int     `json:"books"`
	GenieID        *int      `json:"genie_id"`
	PanguSouls     PanguSouls `json:"pangu_souls"`
	StarDisks      StarDisks `json:"star_disks"`
	Titles         []int     `json:"titles"`
	CalculatedStats Stats    `json:"calculated_stats"`
	ViewCount      int       `json:"view_count"`
	LastViewedAt   time.Time `json:"last_viewed_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Equipment map[string]EquippedItem

type EquippedItem struct {
	ItemID      int   `json:"item_id"`
	RefineLevel int   `json:"refine_level"`
	Gems        []int `json:"gems"`
}

type PanguSouls map[string]int
type StarDisks map[string]int
type Stats map[string]float64

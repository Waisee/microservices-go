package record

import (
	"time"
)

type PartRecord struct {
	UUID          string    `db:"uuid"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	Price         int64     `db:"price"`
	PartType      PartType  `db:"part_type"`
	StockQuantity int64     `db:"stock_quantity"`
	CreatedAt     time.Time `db:"created_at"`
}

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

package model

import "time"

type Variant struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Value     int64     `db:"value"`
	Unit      string    `db:"unit"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type VariantRepositoryPriceDTO struct {
	ID        int64     `db:"variant_id"`
	Name      string    `db:"name"`
	Unit      string    `db:"unit"`
	Price     int64     `db:"price"`
	Value     int64     `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type VariantPriceDTO struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Unit  string `db:"unit"`
	Price int64  `db:"price"`
	Value int64  `db:"value"`
}

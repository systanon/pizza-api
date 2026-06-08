package model

import "time"

type AddonRow struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Price     int64     `db:"price"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CategoryAddonRow struct {
	CategoryID int64     `db:"category_id"`
	AddonID    int64     `db:"addon_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Addon struct {
	ID    int64
	Name  string
	Price int64
}

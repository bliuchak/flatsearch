package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	ext sqlx.ExtContext
}

func NewPostgres(ext sqlx.ExtContext) *Postgres {
	return &Postgres{
		ext: ext,
	}
}

func (p *Postgres) GetEstateByExtID(ctx context.Context, extID string) (Entity, error) {
	var estate Entity
	if err := sqlx.GetContext(ctx, p.ext, &estate, getEstateByExtID, extID); err != nil {
		return Entity{}, fmt.Errorf("select estate: %w", err)
	}

	return estate, nil
}

func (p *Postgres) InsertEstate(ctx context.Context, extID, title, url, price string) (string, error) {
	var estateID string
	if err := sqlx.GetContext(ctx, p.ext, &estateID, insertEstate, extID, title, url, price); err != nil {
		return "", fmt.Errorf("insert new estate: %w", err)
	}

	return estateID, nil
}

func (p *Postgres) UpdateEstate(ctx context.Context, extID, title, url, price string) (string, error) {
	var estateID string
	if _, err := p.ext.ExecContext(ctx, updateEstateByExtID, extID, title, url, price); err != nil {
		return "", fmt.Errorf("update new estate: %w", err)
	}

	return estateID, nil
}

var getEstateByExtID = `
SELECT id, ext_id, title, url, price
FROM estates
WHERE ext_id = $1
LIMIT 1
`

var insertEstate = `
INSERT INTO estates(ext_id, title, url, price) VALUES ($1, $2, $3, $4)
RETURNING id
`

var updateEstateByExtID = `
UPDATE estates SET title=$2, url=$3, price=$4, updated_at=now()
WHERE ext_id = $1
`

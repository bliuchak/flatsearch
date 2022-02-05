package storage

type Entity struct {
	ID    string  `db:"id"`
	ExtID string  `db:"ext_id"`
	Title string  `db:"title"`
	Url   string  `db:"url"`
	Price float64 `db:"price"`
}

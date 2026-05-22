package api

type Variant struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Price int64  `json:"price"`
	Value int64  `json:"value"`
}

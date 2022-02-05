package sreality

type Estate struct {
	Locality string   `json:"locality"`
	HashID   int      `json:"hash_id"`
	Price    float64  `json:"price"`
	Name     string   `json:"name"`
	Seo      Seo      `json:"seo"`
	Gps      Location `json:"gps"`
}

type Seo struct {
	CategoryMainCb int64  `json:"category_main_cb"`
	CategorySubCb  int64  `json:"category_sub_cb"`
	CategoryTypeCb int64  `json:"category_type_cb"`
	Locality       string `json:"locality"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Embedded struct {
	Estates []Estate `json:"estates"`
}

type ClientResponse struct {
	Title           string   `json:"title"`
	MetaDescription string   `json:"meta_description"`
	Locality        string   `json:"locality"`
	ResultSize      uint     `json:"result_size"`
	Embedded        Embedded `json:"_embedded"`
	PerPage         uint     `json:"per_page"`
	Page            uint     `json:"page"`
}

package sreality

type Response struct {
	Estates []Estate
}

type Estate struct {
	Locality string
	HashID   int
	Price    float64
	Name     string
	Seo      Seo
	Gps      Location
}

type Seo struct {
	CategoryMainCb int64
	CategorySubCb  int64
	CategoryTypeCb int64
	Locality       string
}

type Location struct {
	Lat float64
	Lon float64
}

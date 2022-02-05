package sreality

import (
	"math"

	"github.com/bliuchak/flatsearch/internal/platform/sreality"
	"github.com/bliuchak/flatsearch/internal/sreality/decoder"
)

type Decoder interface {
	Decode(data []byte, resp *sreality.ClientResponse) error
}

type SrealityAPI struct {
	client  *sreality.Client
	decoder Decoder
	perPage uint
}

func NewSrealityAPI(client *sreality.Client, json *decoder.JSON) *SrealityAPI {
	return &SrealityAPI{
		client:  client,
		decoder: json,
		perPage: 20,
	}
}

type Option func(s *SrealityAPI)

func WithPerPage(perPage uint) Option {
	return func(s *SrealityAPI) {
		s.perPage = perPage
	}
}

func (s *SrealityAPI) Opts(ops ...Option) *SrealityAPI {
	for _, o := range ops {
		o(s)
	}

	return s
}

func (s *SrealityAPI) GetAll() (Response, error) {
	currentPage := uint(1)

	output, err := s.get(currentPage)
	if err != nil {
		return Response{}, err
	}

	var estates []Estate
	for _, platformEstate := range output.Embedded.Estates {
		estates = append(estates, toEstate(platformEstate))
	}

	totalPages := uint(math.Ceil(float64(output.ResultSize) / float64(s.perPage)))

	for {
		currentPage++

		output, err := s.get(currentPage)
		if err != nil {
			return Response{}, err
		}

		for _, platformEstate := range output.Embedded.Estates {
			estates = append(estates, toEstate(platformEstate))
		}

		if currentPage == totalPages {
			break
		}
	}

	return Response{
		Estates: estates,
	}, nil
}

func (s *SrealityAPI) get(currentPage uint) (sreality.ClientResponse, error) {
	s.client.Country(sreality.CountryCzechRepublic).
		Region(sreality.RegionPrague).
		OperationType(sreality.OperationTypeRent).
		Category(sreality.CategoryElse).
		SubType(sreality.SubTypeElseGarage).
		PerPage(sreality.PerPage(s.perPage)).
		Page(sreality.Page(currentPage))

	resp, err := s.client.Do()
	if err != nil {
		return sreality.ClientResponse{}, err
	}

	var output sreality.ClientResponse

	if err := s.decoder.Decode(resp, &output); err != nil {
		return sreality.ClientResponse{}, err
	}

	return output, nil
}

func toEstate(platformEstate sreality.Estate) Estate {
	return Estate{
		Locality: platformEstate.Locality,
		HashID:   platformEstate.HashID,
		Price:    platformEstate.Price,
		Name:     platformEstate.Name,
		Seo: Seo{
			CategoryMainCb: platformEstate.Seo.CategoryMainCb,
			CategorySubCb:  platformEstate.Seo.CategorySubCb,
			CategoryTypeCb: platformEstate.Seo.CategoryTypeCb,
			Locality:       platformEstate.Seo.Locality,
		},
		Gps: Location{
			Lat: platformEstate.Gps.Lat,
			Lon: platformEstate.Gps.Lon,
		},
	}
}

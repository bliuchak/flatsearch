package sreality

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Client struct {
	httpClient *http.Client
	baseUrl    string

	category      Category
	operationType OperationType
	country       Country
	region        Region
	subType       SubType

	page    Page
	perPage PerPage
}

type Category string
type OperationType string
type Country string
type Region string
type SubType string
type Page uint
type PerPage uint

// https://www.sreality.cz/api/cs/v2/estates?category_main_cb=5&category_sub_cb=52&category_type_cb=2&locality_region_id=10&per_page=20
// https://www.sreality.cz/api/cs/v2/estates?category_main_cb=5category_type_cb=2locality_country_id=112locality_region_id=10category_sub_cb=52
// https://www.sreality.cz/api/cs/v2/estates/count?category_main_cb=5&category_sub_cb=52%7C34&category_type_cb=2&locality_country_id=112&locality_region_id=10

const (
	KeyCategoryMain = "category_main_cb"
	KeyCategoryType = "category_type_cb"
	KeyCountry      = "locality_country_id"
	KeyRegion       = "locality_region_id"
	KeySubType      = "category_sub_cb"
	KeyPage         = "page"
	KeyPerPage      = "per_page"

	// category_main_cb
	CategoryFlat     = Category("1")
	CategoryHouse    = Category("2")
	CategoryLand     = Category("3")
	CategoryCommerce = Category("4")
	CategoryElse     = Category("5")

	// category_type_cb
	OperationTypeSell    = OperationType("1")
	OperationTypeRent    = OperationType("2")
	OperationTypeAuction = OperationType("3")

	// category_sub_cb
	SubTypeElseGarage     = SubType("34")
	SubTypeElseParkingLot = SubType("52")

	// locality_country_id
	CountryCzechRepublic = Country("112")

	// locality_region_id
	RegionPrague                = Region("10")
	RegionCentralBohemianRegion = Region("11")
)

var (
	ErrEmptyCategory      = errors.New("empty category")
	ErrEmptyOperationType = errors.New("empty operation type")
	ErrEmptySubType       = errors.New("empty sub type")
	ErrEmptyCountry       = errors.New("empty country")
	ErrEmptyRegion        = errors.New("empty region")
)

type Option func(d *Client)

func NewClient(ops ...Option) *Client {
	c := Client{
		httpClient: http.DefaultClient,
		baseUrl:    "https://www.sreality.cz/api/cs/v2/estates",
		perPage:    PerPage(20),
	}

	for _, o := range ops {
		o(&c)
	}

	return &c
}

func (c *Client) Category(category Category) *Client {
	c.category = category
	return c
}

func (c *Client) OperationType(operationType OperationType) *Client {
	c.operationType = operationType
	return c
}

func (c *Client) Country(country Country) *Client {
	c.country = country
	return c
}

func (c *Client) Region(region Region) *Client {
	c.region = region
	return c
}

func (c *Client) SubType(subType SubType) *Client {
	c.subType = subType
	return c
}

func (c *Client) Page(page Page) *Client {
	c.page = page
	return c
}

func (c *Client) PerPage(perPage PerPage) *Client {
	c.perPage = perPage
	return c
}

func (c *Client) compose() (string, error) {
	var strs []string

	base := []string{c.baseUrl, "?"}

	if c.category == "" {
		return "", ErrEmptyCategory
	} else {
		strs = ampAppend(strs)
		strs = append(strs, KeyCategoryMain, "=", string(c.category))
	}

	if c.operationType == "" {
		return "", ErrEmptyOperationType
	} else {
		strs = ampAppend(strs)
		strs = append(strs, KeyCategoryType, "=", string(c.operationType))
	}

	if c.country == "" {
		return "", ErrEmptyCountry
	} else {
		strs = ampAppend(strs)
		strs = append(strs, KeyCountry, "=", string(c.country))
	}

	if c.region == "" {
		return "", ErrEmptyRegion
	} else {
		strs = ampAppend(strs)
		strs = append(strs, KeyRegion, "=", string(c.region))
	}

	if c.subType == "" {
		return "", ErrEmptySubType
	} else {
		strs = ampAppend(strs)
		strs = append(strs, KeySubType, "=", string(c.subType))
	}

	if c.page > 0 {
		strs = ampAppend(strs)
		strs = append(strs, KeyPage, "=", strconv.Itoa(int(c.page)))
	}

	if c.perPage > 0 {
		strs = ampAppend(strs)
		strs = append(strs, KeyPerPage, "=", strconv.Itoa(int(c.perPage)))
	}

	return join(append(base, strs...)...), nil
}

func (c *Client) Do() ([]byte, error) {
	url, err := c.compose()
	if err != nil {
		return nil, err
	}

	spew.Dump(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func WithHttpClient(client *http.Client) Option {
	return func(d *Client) {
		d.httpClient = client
	}
}

func WithBaseUrl(url string) Option {
	return func(d *Client) {
		d.baseUrl = url
	}
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

func ampAppend(strs []string) []string {
	if len(strs) > 0 && strs[len(strs)-1] != "&" {
		strs = append(strs, "&")
	}
	return strs
}

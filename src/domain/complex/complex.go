package complex

const (
	STATUS_ACTIVE   = "active"
	STATUS_DEACTIVE = "deactive"
)

type Complex struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	PriceFrom        float64   `json:"price_from"`
	PriceTo          float64   `json:"price_to"`
	Photo            string    `json:"photo"`
	Address          string    `json:"address"`
	GPS              gps       `json:"gps"`
	City             string    `json:"city"`
	NumberOfBuilding int64     `json:"number_of_building"`
	NumberOfVilla    int64     `json:"number_of_villa"`
	Services         services  `json:"services"`
	SupportedSales   []string  `json:"supported_sales"`
	Amenities        amenities `json:"amenities"`

	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
}

type amenities struct {
	CarParking        bool `json:"car_parking"`
	Gardens           bool `json:"gardens"`
	PrimarySchool     bool `json:"primary_school"`
	Kindergarten      bool `json:"kindergarten"`
	HealthCentre      bool `json:"health_centre"`
	ShoppingMall      bool `json:"shopping_mall"`
	Markets           bool `json:"markets"`
	Restaurant        bool `json:"restaurant"`
	Cafe              bool `json:"cafe"`
	Pharmacy          bool `json:"pharmacy"`
	Laundery          bool `json:"laundery"`
	Gym               bool `json:"gym"`
	Barbershop        bool `json:"barbershop"`
	BeautyCenter      bool `json:"beauty_center"`
	KidsEntertainment bool `json:"kids_entertainment"`
	Nursery           bool `json:"nursery"`
	Clinic            bool `json:"clinic"`
	Bakery            bool `json:"bakery"`
	GiftStore         bool `json:"gift_store"`
	ElectronicShops   bool `json:"electronic_shops"`
	Stationary        bool `json:"stationary"`
	SweetShop         bool `json:"sweet_shop"`
	Mosque            bool `json:"mosque"`
	School            bool `json:"school"`
	Hospital          bool `json:"hospital"`
	Shopping          bool `json:"shopping"`
	Garden            bool `json:"garden"`
}

type services struct {
	Electric       float64 `json:"electric"`
	WaterResources float64 `json:"water_resources"`
	Rubbish        float64 `json:"rubbish"`
	SecurityHours  string  `json:"security_hours"`
}

type gps struct {
	Long string `json:"long"`
	Lat  string `json:"lat"`
}

type Complexes []Complex

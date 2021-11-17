package complex

const (
	STATUS_ACTIVE   = "active"
	STATUS_DEACTIVE = "deactive"
)

type Complex struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Photo            string    `json:"photo"`
	Address          string    `json:"address"`
	GPS              gps       `json:"gps"`
	City             string    `json:"city"`
	NumberOfBuilding int64     `json:"number_of_building"`
	NumberOfVilla    int64     `json:"number_of_villa"`
	Services         services  `json:"services"`
	Amenities        amenities `json:"amenities"`

	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
}

type amenities struct {
	Mosque   bool `json:"mosque"`
	School   bool `json:"school"`
	Hospital bool `json:"hospital"`
	Gym      bool `json:"gym"`
	Shopping bool `json:"shopping"`
	Garden   bool `json:"garden"`
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

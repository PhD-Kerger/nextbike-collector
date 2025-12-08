package types

type Root struct {
	Countries []Country `json:"countries"`
}

// Country struct based on the provided Rust struct
// Note: You will need to define the City struct as well
type Country struct {
	Lat                   float64 `json:"lat" parquet:"name=lat, type=DOUBLE"`
	Lng                   float64 `json:"lng" parquet:"name=lng, type=DOUBLE"`
	Zoom                  uint8   `json:"zoom"`
	Name                  string  `json:"name"`
	Hotline               *string `json:"hotline,omitempty"`
	Domain                *string `json:"domain,omitempty"`
	Language              *string `json:"language,omitempty"`
	Email                 *string `json:"email,omitempty"`
	Timezone              string  `json:"timezone"`
	Currency              *string `json:"currency,omitempty"`
	CountryCallingCode    *string `json:"country_calling_code,omitempty"`
	SystemOperatorAddress *string `json:"system_operator_address,omitempty"`
	Country               string  `json:"country"`
	CountryName           string  `json:"country_name"`
	Terms                 *string `json:"terms,omitempty"`
	Policy                *string `json:"policy,omitempty"`
	Website               *string `json:"website,omitempty"`
	ShowBikeTypes         bool    `json:"show_bike_types"`
	ShowBikeTypeGroups    bool    `json:"show_bike_type_groups"`
	ShowFreeRacks         bool    `json:"show_free_racks"`
	BookedBikes           int16   `json:"booked_bikes"`
	SetPointBikes         uint16  `json:"set_point_bikes"`
	AvailableBikes        uint16  `json:"available_bikes"`
	CappedAvailableBikes  bool    `json:"capped_available_bikes"`
	NoRegistration        bool    `json:"no_registration"`
	Pricing               *string `json:"pricing,omitempty"`
	Vat                   *string `json:"vat,omitempty"`
	FaqURL                *string `json:"faq_url,omitempty"`
	StoreURIAndroid       string  `json:"store_uri_android"`
	StoreURIiOS           string  `json:"store_uri_ios"`
	Cities                []City  `json:"cities"`
}

type FlatCountry struct {
	Timestamp             uint32  `parquet:"timestamp"`
	Lat                   float64 `parquet:"lat"`
	Lng                   float64 `parquet:"lng"`
	Zoom                  uint32  `parquet:"zoom"`
	Name                  string  `parquet:"name"`
	Hotline               *string `parquet:"hotline,optional"`
	Domain                *string `parquet:"domain,optional"`
	Language              *string `parquet:"language,optional"`
	Email                 *string `parquet:"email,optional"`
	Timezone              string  `parquet:"timezone"`
	Currency              *string `parquet:"currency,optional"`
	CountryCallingCode    *string `parquet:"country_calling_code,optional"`
	SystemOperatorAddress *string `parquet:"system_operator_address,optional"`
	Country               string  `parquet:"country"`
	CountryName           string  `parquet:"country_name"`
	Terms                 *string `parquet:"terms,optional"`
	Policy                *string `parquet:"policy,optional"`
	Website               *string `parquet:"website,optional"`
	ShowBikeTypes         bool    `parquet:"show_bike_types"`
	ShowBikeTypeGroups    bool    `parquet:"show_bike_type_groups"`
	ShowFreeRacks         bool    `parquet:"show_free_racks"`
	BookedBikes           int32   `parquet:"booked_bikes"`
	SetPointBikes         uint32  `parquet:"set_point_bikes"`
	AvailableBikes        uint32  `parquet:"available_bikes"`
	CappedAvailableBikes  bool    `parquet:"capped_available_bikes"`
	NoRegistration        bool    `parquet:"no_registration"`
	Pricing               *string `parquet:"pricing,optional"`
	Vat                   *string `parquet:"vat,optional"`
	FaqURL                *string `parquet:"faq_url,optional"`
	StoreURIAndroid       string  `parquet:"store_uri_android"`
	StoreURIiOS           string  `parquet:"store_uri_ios"`
}

func (country *Country) ToFlat(timestamp uint32) FlatCountry {
	return FlatCountry{
		Timestamp:             timestamp,
		Lat:                   country.Lat,
		Lng:                   country.Lng,
		Zoom:                  uint32(country.Zoom),
		Name:                  country.Name,
		Hotline:               country.Hotline,
		Domain:                country.Domain,
		Language:              country.Language,
		Email:                 country.Email,
		Timezone:              country.Timezone,
		Currency:              country.Currency,
		CountryCallingCode:    country.CountryCallingCode,
		SystemOperatorAddress: country.SystemOperatorAddress,
		Country:               country.Country,
		CountryName:           country.CountryName,
		Terms:                 country.Terms,
		Policy:                country.Policy,
		Website:               country.Website,
		ShowBikeTypes:         country.ShowBikeTypes,
		ShowBikeTypeGroups:    country.ShowBikeTypeGroups,
		ShowFreeRacks:         country.ShowFreeRacks,
		BookedBikes:           int32(country.BookedBikes),
		SetPointBikes:         uint32(country.SetPointBikes),
		AvailableBikes:        uint32(country.AvailableBikes),
		CappedAvailableBikes:  country.CappedAvailableBikes,
		NoRegistration:        country.NoRegistration,
		Pricing:               country.Pricing,
		Vat:                   country.Vat,
		FaqURL:                country.FaqURL,
		StoreURIAndroid:       country.StoreURIAndroid,
		StoreURIiOS:           country.StoreURIiOS,
	}
}

// City struct based on the provided Rust struct
// Note: You will need to define Bounds and Place structs as well

type City struct {
	UID                  uint32            `json:"uid"`
	Lat                  float64           `json:"lat"`
	Lng                  float64           `json:"lng"`
	Zoom                 uint8             `json:"zoom"`
	MapsIcon             string            `json:"maps_icon"`
	Alias                string            `json:"alias"`
	BreakInfo            bool              `json:"break"`
	Name                 string            `json:"name"`
	NumPlaces            uint16            `json:"num_places"`
	RefreshRate          string            `json:"refresh_rate"`
	Bounds               Bounds            `json:"bounds"`
	BookedBikes          int32             `json:"booked_bikes"`
	SetPointBikes        uint16            `json:"set_point_bikes"`
	AvailableBikes       uint16            `json:"available_bikes"`
	ReturnToOfficialOnly bool              `json:"return_to_official_only"`
	BikeTypes            map[string]uint32 `json:"bike_types"`
	Website              *string           `json:"website,omitempty"`
	Places               []Place           `json:"places"`
}

type FlatCity struct {
	Timestamp            uint32            `parquet:"timestamp"`
	Parent_Country_Name  string            `parquet:"parent_country_name"`
	UID                  uint32            `parquet:"uid"`
	Lat                  float64           `parquet:"lat"`
	Lng                  float64           `parquet:"lng"`
	Zoom                 uint32            `parquet:"zoom"`
	MapsIcon             string            `parquet:"maps_icon"`
	Alias                string            `parquet:"alias"`
	BreakInfo            bool              `parquet:"break"`
	Name                 string            `parquet:"name"`
	NumPlaces            uint32            `parquet:"num_places"`
	RefreshRate          string            `parquet:"refresh_rate"`
	Bounds_SouthWest_Lat float64           `parquet:"bounds_south_west_lat"`
	Bounds_SouthWest_Lng float64           `parquet:"bounds_south_west_lng"`
	Bounds_NorthEast_Lat float64           `parquet:"bounds_north_east_lat"`
	Bounds_NorthEast_Lng float64           `parquet:"bounds_north_east_lng"`
	BookedBikes          int32             `parquet:"booked_bikes"`
	SetPointBikes        uint32            `parquet:"set_point_bikes"`
	AvailableBikes       uint32            `parquet:"available_bikes"`
	ReturnToOfficialOnly bool              `parquet:"return_to_official_only"`
	Website              *string           `parquet:"website,optional"`
	BikeTypes            map[string]uint32 `parquet:"bike_types"`
}

func (city *City) ToFlat(timestamp uint32, parentCountryName string) FlatCity {
	return FlatCity{
		Timestamp:            timestamp,
		Parent_Country_Name:  parentCountryName,
		UID:                  city.UID,
		Lat:                  city.Lat,
		Lng:                  city.Lng,
		Zoom:                 uint32(city.Zoom),
		MapsIcon:             city.MapsIcon,
		Alias:                city.Alias,
		BreakInfo:            city.BreakInfo,
		Name:                 city.Name,
		NumPlaces:            uint32(city.NumPlaces),
		RefreshRate:          city.RefreshRate,
		Bounds_SouthWest_Lat: city.Bounds.SouthWest.Lat,
		Bounds_SouthWest_Lng: city.Bounds.SouthWest.Lng,
		Bounds_NorthEast_Lat: city.Bounds.NorthEast.Lat,
		Bounds_NorthEast_Lng: city.Bounds.NorthEast.Lng,
		BookedBikes:          city.BookedBikes,
		SetPointBikes:        uint32(city.SetPointBikes),
		AvailableBikes:       uint32(city.AvailableBikes),
		ReturnToOfficialOnly: city.ReturnToOfficialOnly,
		Website:              city.Website,
		BikeTypes:            city.BikeTypes,
	}
}

// Bounds struct based on the provided Rust struct
// Coordinates struct for use in Bounds
// Place struct based on the provided Rust struct
// Note: You will need to define the Bike struct as well

type Bounds struct {
	SouthWest Coordinates `json:"south_west"`
	NorthEast Coordinates `json:"north_east"`
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Place struct {
	UID                  uint32            `json:"uid"`
	Lat                  float64           `json:"lat"`
	Lng                  float64           `json:"lng"`
	Bike                 bool              `json:"bike"`
	Name                 string            `json:"name"`
	Address              *string           `json:"address,omitempty"`
	Spot                 bool              `json:"spot"`
	Number               uint32            `json:"number"`
	BookedBikes          int32             `json:"booked_bikes"`
	Bikes                uint32            `json:"bikes"`
	BikesAvailableToRent uint32            `json:"bikes_available_to_rent"`
	BikeRacks            uint32            `json:"bike_racks"`
	FreeRacks            uint32            `json:"free_racks"`
	SpecialRacks         int32             `json:"special_racks"`
	FreeSpecialRacks     uint32            `json:"free_special_racks"`
	Maintenance          bool              `json:"maintenance"`
	TerminalType         string            `json:"terminal_type"`
	BikeList             []Bike            `json:"bike_list"`
	BikeNumbers          []string          `json:"bike_numbers"`
	BikeTypes            map[string]uint32 `json:"bike_types"`
	PlaceType            string            `json:"place_type"`
	RackLocks            bool              `json:"rack_locks"`
}

type FlatPlace struct {
	Timestamp            uint32            `parquet:"timestamp"`
	City_UID             uint32            `parquet:"city_uid"`
	UID                  uint32            `parquet:"uid"`
	Lat                  float64           `parquet:"lat"`
	Lng                  float64           `parquet:"lng"`
	Bike                 bool              `parquet:"bike"`
	Name                 string            `parquet:"name"`
	Address              *string           `parquet:"address,optional"`
	Spot                 bool              `parquet:"spot"`
	Number               uint32            `parquet:"number"`
	BookedBikes          int32             `parquet:"booked_bikes"`
	Bikes                uint32            `parquet:"bikes"`
	BikesAvailableToRent uint32            `parquet:"bikes_available_to_rent"`
	BikeRacks            uint32            `parquet:"bike_racks"`
	FreeRacks            uint32            `parquet:"free_racks"`
	FreeSpecialRacks     uint32            `parquet:"free_special_racks"`
	SpecialRacks         int32             `parquet:"special_racks"`
	Maintenance          bool              `parquet:"maintenance"`
	TerminalType         string            `parquet:"terminal_type"`
	BikeNumbers          []string          `parquet:"bike_numbers"`
	BikeTypes            map[string]uint32 `parquet:"bike_types"`
	PlaceType            string            `parquet:"place_type"`
	RackLocks            bool              `parquet:"rack_locks"`
}

func (place *Place) ToFlat(timestamp uint32, cityUID uint32) FlatPlace {
	return FlatPlace{
		Timestamp:            timestamp,
		City_UID:             cityUID,
		UID:                  place.UID,
		Lat:                  place.Lat,
		Lng:                  place.Lng,
		Bike:                 place.Bike,
		Name:                 place.Name,
		Address:              place.Address,
		Spot:                 place.Spot,
		Number:               place.Number,
		BookedBikes:          place.BookedBikes,
		Bikes:                place.Bikes,
		BikesAvailableToRent: place.BikesAvailableToRent,
		BikeRacks:            place.BikeRacks,
		FreeRacks:            place.FreeRacks,
		FreeSpecialRacks:     place.FreeSpecialRacks,
		SpecialRacks:         place.SpecialRacks,
		Maintenance:          place.Maintenance,
		TerminalType:         place.TerminalType,
		BikeNumbers:          place.BikeNumbers,
		BikeTypes:            place.BikeTypes,
		PlaceType:            place.PlaceType,
		RackLocks:            place.RackLocks,
	}
}

// Bike struct based on the provided Rust struct
type Bike struct {
	Number         string             `json:"number"`
	BikeType       uint16             `json:"bike_type"`
	LockTypes      []string           `json:"lock_types"`
	Active         bool               `json:"active"`
	State          *string            `json:"state,omitempty"`
	ElectricLock   bool               `json:"electric_lock"`
	Boardcomputer  uint64             `json:"boardcomputer"`
	PedelecBattery *uint8             `json:"pedelec_battery,omitempty"`
	BatteryPack    *map[string]uint32 `json:"battery_pack,omitempty"`
}

type FlatBike struct {
	Timestamp      uint32             `parquet:"timestamp"`
	Place_UID      uint32             `parquet:"place_uid"`
	Number         string             `parquet:"number"`
	BikeType       uint32             `parquet:"bike_type"`
	LockTypes      []string           `parquet:"lock_types"`
	Active         bool               `parquet:"active"`
	State          *string            `parquet:"state,optional"`
	ElectricLock   bool               `parquet:"electric_lock"`
	Boardcomputer  uint64             `parquet:"boardcomputer"`
	PedelecBattery *uint32            `parquet:"pedelec_battery,optional"`
	BatteryPack    *map[string]uint32 `parquet:"battery_pack,optional"`
}

func (bike *Bike) ToFlat(timestamp uint32, placeUID uint32) FlatBike {
	var pedelecBattery *uint32
	if bike.PedelecBattery != nil {
		val := uint32(*bike.PedelecBattery)
		pedelecBattery = &val
	}
	var batteryPackUint32 *map[string]uint32
	if bike.BatteryPack != nil {
		bp := make(map[string]uint32, len(*bike.BatteryPack))
		for k, v := range *bike.BatteryPack {
			bp[k] = v
		}
		batteryPackUint32 = &bp
	}
	return FlatBike{
		Timestamp:      timestamp,
		Place_UID:      placeUID,
		Number:         bike.Number,
		BikeType:       uint32(bike.BikeType),
		LockTypes:      bike.LockTypes,
		Active:         bike.Active,
		State:          bike.State,
		ElectricLock:   bike.ElectricLock,
		Boardcomputer:  bike.Boardcomputer,
		PedelecBattery: pedelecBattery,
		BatteryPack:    batteryPackUint32,
	}
}

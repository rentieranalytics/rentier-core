package apicalculation

import "github.com/rentieranalytics/rentier-core/domain"

type AVMCalculationRequest struct {
	Area                float64         `json:"area"`
	GeoPoint            domain.GeoPoint `json:"geo_point"`
	MarketType          string          `json:"market_type"`
	Rooms               *int            `json:"rooms,omitempty"`
	BuildingType        *string         `json:"building_type,omitempty"`
	BuildYear           *int            `json:"build_year,omitempty"`
	Floor               *int            `json:"floor,omitempty"`
	TotalFloors         *int            `json:"total_floors,omitempty"`
	Standard            *string         `json:"standard,omitempty"`
	WithParkPlace       *bool           `json:"with_park_place,omitempty"`
	WithElevator        *bool           `json:"with_elevator,omitempty"`
	WithBalcony         *bool           `json:"with_balcony,omitempty"`
	WithAdditionalSpace *bool           `json:"with_additional_space,omitempty"`
}

type AVMCalculationResponse struct {
	Price      string   `json:"price"`
	PriceRaw   string   `json:"price_raw"`
	PriceStats AVMStats `json:"price_stats"`
}

type AVMStats struct {
	AVMPriceEstimation AVMEstimation `json:"price_estimation"`
}

type AVMEstimation struct {
	Distance                      int             `json:"distance"`
	Points                        int             `json:"points"`
	RealEstates                   int             `json:"real_estates"`
	Days                          int             `json:"days"`
	Accuracy                      int             `json:"accuracy"`
	RealEstatesUsedForCalculation []AVMRealEstate `json:"real_estates_used_for_calculation"`
	AvgPrice                      string          `json:"avg_price"`
	AvgPriceM2                    string          `json:"avg_price_m2"`
	StandardDeviationPriceM2      string          `json:"standard_deviation_price_m2"`
	StandardDeviationPrice        string          `json:"standard_deviation_price"`
	DeviationPriceMin             string          `json:"deviation_price_min"`
	DeviationPriceMax             string          `json:"deviation_price_max"`
	DeviationPriceM2Min           string          `json:"deviation_price_m2_min"`
	DeviationPriceM2Max           string          `json:"deviation_price_m2_max"`
}

type AVMRealEstate struct {
	Address        string          `json:"address"`
	Area           float64         `json:"area"`
	BarnPriceM2    string          `json:"barn_price_m2"`
	BarnRawPriceM2 string          `json:"barn_raw_price_m2"`
	BuildYear      *int            `json:"build_year"`
	DomainName     string          `json:"domain__name"`
	Floor          *int32          `json:"floor"`
	GeoPoint       domain.GeoPoint `json:"geo_point"`
	ID             string          `json:"id"`
	IsActive       bool            `json:"is_active"`
	OfferTitle     string          `json:"offer_title"`
	Points         int             `json:"points"`
	Price          string          `json:"price"`
	PriceM2        string          `json:"price_m2"`
	Rooms          *int32          `json:"rooms"`
	URL            string          `json:"url"`
}

package calculations

type AVMResponse struct {
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
	Address        string              `json:"address"`
	Area           float64             `json:"area"`
	BarnPriceM2    string              `json:"barn_price_m2"`
	BarnRawPriceM2 string              `json:"barn_raw_price_m2"`
	BuildYear      *int                `json:"build_year"`
	DomainName     string              `json:"domain__name"`
	Floor          *int32              `json:"floor"`
	GeoPoint       RealEstatesGeoPoint `json:"geo_point"`
	ID             string              `json:"id"`
	IsActive       bool                `json:"is_active"`
	OfferTitle     string              `json:"offer_title"`
	Points         int                 `json:"points"`
	Price          string              `json:"price"`
	PriceM2        string              `json:"price_m2"`
	Rooms          *int32              `json:"rooms"`
	URL            string              `json:"url"`
}

type RealEstatesGeoPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

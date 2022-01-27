package endpoints

type GeoFilter struct {
	Action       GeoFilterActions `json:"action"`
	CountryCodes []string         `json:"countryCodes"`
	RelativePath string           `json:"relativePath"`
}

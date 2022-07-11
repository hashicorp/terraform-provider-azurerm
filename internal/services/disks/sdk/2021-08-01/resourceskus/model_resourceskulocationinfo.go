package resourceskus

type ResourceSkuLocationInfo struct {
	Location    *string                   `json:"location,omitempty"`
	ZoneDetails *[]ResourceSkuZoneDetails `json:"zoneDetails,omitempty"`
	Zones       *[]string                 `json:"zones,omitempty"`
}

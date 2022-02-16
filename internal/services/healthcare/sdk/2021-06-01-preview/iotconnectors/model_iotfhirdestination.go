package iotconnectors

type IotFhirDestination struct {
	Etag       *string                      `json:"etag,omitempty"`
	Id         *string                      `json:"id,omitempty"`
	Location   *string                      `json:"location,omitempty"`
	Name       *string                      `json:"name,omitempty"`
	Properties IotFhirDestinationProperties `json:"properties"`
	SystemData *SystemData                  `json:"systemData,omitempty"`
	Type       *string                      `json:"type,omitempty"`
}

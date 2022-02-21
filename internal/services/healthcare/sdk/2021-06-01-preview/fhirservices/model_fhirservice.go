package fhirservices

type FhirService struct {
	Etag       *string                         `json:"etag,omitempty"`
	Id         *string                         `json:"id,omitempty"`
	Identity   *ServiceManagedIdentityIdentity `json:"identity,omitempty"`
	Kind       *FhirServiceKind                `json:"kind,omitempty"`
	Location   *string                         `json:"location,omitempty"`
	Name       *string                         `json:"name,omitempty"`
	Properties *FhirServiceProperties          `json:"properties,omitempty"`
	SystemData *SystemData                     `json:"systemData,omitempty"`
	Tags       *map[string]string              `json:"tags,omitempty"`
	Type       *string                         `json:"type,omitempty"`
}

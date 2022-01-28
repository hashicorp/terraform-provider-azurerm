package dicomservices

type DicomService struct {
	Etag       *string                 `json:"etag,omitempty"`
	Id         *string                 `json:"id,omitempty"`
	Location   *string                 `json:"location,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *DicomServiceProperties `json:"properties,omitempty"`
	SystemData *SystemData             `json:"systemData,omitempty"`
	Tags       *map[string]string      `json:"tags,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}

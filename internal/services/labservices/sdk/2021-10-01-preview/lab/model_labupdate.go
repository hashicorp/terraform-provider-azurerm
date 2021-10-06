package lab

type LabUpdate struct {
	Properties *LabUpdateProperties `json:"properties,omitempty"`
	Tags       *[]string            `json:"tags,omitempty"`
}

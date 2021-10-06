package labplan

type LabPlanUpdate struct {
	Properties *LabPlanUpdateProperties `json:"properties,omitempty"`
	Tags       *[]string                `json:"tags,omitempty"`
}

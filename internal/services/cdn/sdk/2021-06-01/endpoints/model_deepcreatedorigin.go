package endpoints

type DeepCreatedOrigin struct {
	Name       string                       `json:"name"`
	Properties *DeepCreatedOriginProperties `json:"properties,omitempty"`
}

package endpoints

type DeepCreatedOriginGroup struct {
	Name       string                            `json:"name"`
	Properties *DeepCreatedOriginGroupProperties `json:"properties,omitempty"`
}

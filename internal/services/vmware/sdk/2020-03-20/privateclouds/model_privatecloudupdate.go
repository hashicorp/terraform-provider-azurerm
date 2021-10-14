package privateclouds

type PrivateCloudUpdate struct {
	Properties *PrivateCloudUpdateProperties `json:"properties,omitempty"`
	Tags       *map[string]string            `json:"tags,omitempty"`
}

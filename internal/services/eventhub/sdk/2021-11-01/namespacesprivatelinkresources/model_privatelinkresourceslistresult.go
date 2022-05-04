package namespacesprivatelinkresources

type PrivateLinkResourcesListResult struct {
	NextLink *string                `json:"nextLink,omitempty"`
	Value    *[]PrivateLinkResource `json:"value,omitempty"`
}

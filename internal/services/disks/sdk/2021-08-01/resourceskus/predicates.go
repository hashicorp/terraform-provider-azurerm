package resourceskus

type ResourceSkuInfoPredicate struct {
	ApiVersion   *string
	Name         *string
	ResourceType *string
	Tier         *string
}

func (p ResourceSkuInfoPredicate) Matches(input ResourceSkuInfo) bool {

	if p.ApiVersion != nil && (input.ApiVersion == nil && *p.ApiVersion != *input.ApiVersion) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.ResourceType != nil && (input.ResourceType == nil && *p.ResourceType != *input.ResourceType) {
		return false
	}

	if p.Tier != nil && (input.Tier == nil && *p.Tier != *input.Tier) {
		return false
	}

	return true
}

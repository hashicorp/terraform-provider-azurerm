package cognitiveservicesaccounts

type ResourceSkuRestrictions struct {
	ReasonCode      *ResourceSkuRestrictionsReasonCode `json:"reasonCode,omitempty"`
	RestrictionInfo *ResourceSkuRestrictionInfo        `json:"restrictionInfo,omitempty"`
	Type            *ResourceSkuRestrictionsType       `json:"type,omitempty"`
	Values          *[]string                          `json:"values,omitempty"`
}

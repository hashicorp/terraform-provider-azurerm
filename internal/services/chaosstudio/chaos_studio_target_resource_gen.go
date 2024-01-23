package chaosstudio

// NOTE: this placeholder below is just to keep the compiler happy until the upstream
// PR https://github.com/hashicorp/pandora/pull/3671 is merged

var _ = ChaosStudioTargetResource{} // to keep the unused linter happy

type ChaosStudioTargetResource struct{}

type ChaosStudioTargetResourceSchema struct {
	Location         string `tfschema:"location"`
	TargetType       string `tfschema:"target_type"`
	TargetResourceId string `tfschema:"target_resource_id"`
}

func (r ChaosStudioTargetResource) ResourceType() string {
	return ""
}

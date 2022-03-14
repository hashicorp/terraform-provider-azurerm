package vmhhostlist

type VMResourcesOperationPredicate struct {
	VmResourceId *string
}

func (p VMResourcesOperationPredicate) Matches(input VMResources) bool {

	if p.VmResourceId != nil && (input.VmResourceId == nil && *p.VmResourceId != *input.VmResourceId) {
		return false
	}

	return true
}

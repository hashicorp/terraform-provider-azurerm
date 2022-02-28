package vmhhostlist

type VMResourcesPredicate struct {
	VmResourceId *string
}

func (p VMResourcesPredicate) Matches(input VMResources) bool {

	if p.VmResourceId != nil && (input.VmResourceId == nil && *p.VmResourceId != *input.VmResourceId) {
		return false
	}

	return true
}

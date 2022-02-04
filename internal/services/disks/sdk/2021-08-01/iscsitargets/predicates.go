package iscsitargets

type IscsiTargetPredicate struct {
	Id        *string
	ManagedBy *string
	Name      *string
	Type      *string
}

func (p IscsiTargetPredicate) Matches(input IscsiTarget) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.ManagedBy != nil && (input.ManagedBy == nil && *p.ManagedBy != *input.ManagedBy) {
		return false
	}

	if p.Name != nil && (input.Name == nil && *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil && *p.Type != *input.Type) {
		return false
	}

	return true
}

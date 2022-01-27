package afdprofiles

type UsagePredicate struct {
	CurrentValue *int64
	Id           *string
	Limit        *int64
}

func (p UsagePredicate) Matches(input Usage) bool {

	if p.CurrentValue != nil && *p.CurrentValue != input.CurrentValue {
		return false
	}

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
		return false
	}

	if p.Limit != nil && *p.Limit != input.Limit {
		return false
	}

	return true
}

package rulesets

type RuleSetPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RuleSetPredicate) Matches(input RuleSet) bool {

	if p.Id != nil && (input.Id == nil && *p.Id != *input.Id) {
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

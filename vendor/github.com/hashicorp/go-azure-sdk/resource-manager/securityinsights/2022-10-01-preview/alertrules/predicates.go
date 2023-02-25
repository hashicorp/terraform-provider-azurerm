package alertrules

type AlertRuleOperationPredicate struct {
}

func (p AlertRuleOperationPredicate) Matches(input AlertRule) bool {

	return true
}

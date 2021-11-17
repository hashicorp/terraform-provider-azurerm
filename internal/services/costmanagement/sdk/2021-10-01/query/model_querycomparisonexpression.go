package query

type QueryComparisonExpression struct {
	Name     string            `json:"name"`
	Operator QueryOperatorType `json:"operator"`
	Values   []string          `json:"values"`
}

package views

type ReportConfigComparisonExpression struct {
	Name     string       `json:"name"`
	Operator OperatorType `json:"operator"`
	Values   []string     `json:"values"`
}

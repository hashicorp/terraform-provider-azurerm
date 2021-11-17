package views

type ReportConfigFilter struct {
	And       *[]ReportConfigFilter             `json:"and,omitempty"`
	Dimension *ReportConfigComparisonExpression `json:"dimension,omitempty"`
	Not       *ReportConfigFilter               `json:"not,omitempty"`
	Or        *[]ReportConfigFilter             `json:"or,omitempty"`
	Tag       *ReportConfigComparisonExpression `json:"tag,omitempty"`
}

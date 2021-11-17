package query

type QueryFilter struct {
	And       *[]QueryFilter             `json:"and,omitempty"`
	Dimension *QueryComparisonExpression `json:"dimension,omitempty"`
	Not       *QueryFilter               `json:"not,omitempty"`
	Or        *[]QueryFilter             `json:"or,omitempty"`
	Tag       *QueryComparisonExpression `json:"tag,omitempty"`
}

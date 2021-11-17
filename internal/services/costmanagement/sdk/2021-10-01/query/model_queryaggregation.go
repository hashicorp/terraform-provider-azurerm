package query

type QueryAggregation struct {
	Function FunctionType `json:"function"`
	Name     string       `json:"name"`
}

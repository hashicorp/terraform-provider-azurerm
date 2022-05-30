package alertsmanagement

type Condition struct {
	Field    *Field    `json:"field,omitempty"`
	Operator *Operator `json:"operator,omitempty"`
	Values   *[]string `json:"values,omitempty"`
}

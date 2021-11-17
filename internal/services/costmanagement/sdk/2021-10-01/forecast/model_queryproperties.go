package forecast

type QueryProperties struct {
	Columns  *[]QueryColumn   `json:"columns,omitempty"`
	NextLink *string          `json:"nextLink,omitempty"`
	Rows     *[][]interface{} `json:"rows,omitempty"`
}

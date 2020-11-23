package sdk

// TODO: make these more granular for the tests

type ExampleObj struct {
	Name        string            `tfschema:"name"`
	Number      int               `tfschema:"number"`
	Output      string            `tfschema:"output" computed:"true"`
	Enabled     bool              `tfschema:"enabled"`
	Networks    []string          `tfschema:"networks"`
	NetworksSet []string          `tfschema:"networks_set"`
	IntList     []int             `tfschema:"int_list"`
	IntSet      []int             `tfschema:"int_set"`
	FloatList   []float64         `tfschema:"float_list"`
	FloatSet    []float64         `tfschema:"float_set"`
	BoolList    []bool            `tfschema:"bool_list"`
	BoolSet     []bool            `tfschema:"bool_set"`
	List        []NetworkList     `tfschema:"list"`
	Set         []NetworkSet      `tfschema:"set"`
	Float       float64           `tfschema:"float"`
	Map         map[string]string `tfschema:"map"`
}

type NetworkList struct {
	Name  string         `tfschema:"name"`
	Inner []NetworkInner `tfschema:"inner"`
}

type NetworkListSet struct {
	Name string `tfschema:"name"`
}

type NetworkSet struct {
	Name  string       `tfschema:"name"`
	Inner []InnerInner `tfschema:"inner"`
}

type NetworkInner struct {
	Name  string           `tfschema:"name"`
	Inner []InnerInner     `tfschema:"inner"`
	Set   []NetworkListSet `tfschema:"set"`
}

type InnerInner struct {
	Name         string `tfschema:"name"`
	ShouldBeFine bool   `tfschema:"should_be_fine"`
}

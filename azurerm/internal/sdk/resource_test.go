package sdk

// TODO: make these more granular for the tests

type ExampleObj struct {
	Name        string            `hcl:"name"`
	Number      int               `hcl:"number"`
	Output      string            `hcl:"output" computed:"true"`
	Enabled     bool              `hcl:"enabled"`
	Networks    []string          `hcl:"networks"`
	NetworksSet []string          `hcl:"networks_set"`
	IntList     []int             `hcl:"int_list"`
	IntSet      []int             `hcl:"int_set"`
	FloatList   []float64         `hcl:"float_list"`
	FloatSet    []float64         `hcl:"float_set"`
	BoolList    []bool            `hcl:"bool_list"`
	BoolSet     []bool            `hcl:"bool_set"`
	List        []NetworkList     `hcl:"list"`
	Set         []NetworkSet      `hcl:"set"`
	Float       float64           `hcl:"float"`
	Map         map[string]string `hcl:"map"`
}

type NetworkList struct {
	Name  string         `hcl:"name"`
	Inner []NetworkInner `hcl:"inner"`
}

type NetworkListSet struct {
	Name string `hcl:"name"`
}

type NetworkSet struct {
	Name  string       `hcl:"name"`
	Inner []InnerInner `hcl:"inner"`
}

type NetworkInner struct {
	Name  string           `hcl:"name"`
	Inner []InnerInner     `hcl:"inner"`
	Set   []NetworkListSet `hcl:"set"`
}

type InnerInner struct {
	Name         string `hcl:"name"`
	ShouldBeFine bool   `hcl:"should_be_fine"`
}

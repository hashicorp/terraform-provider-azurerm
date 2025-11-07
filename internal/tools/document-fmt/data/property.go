package data

type Properties struct {
	Names   []string // Only really relevant to the documentation, could be used to track ordering in docs to compare against ordering we want
	Objects map[string]*Property
}

type Property struct {
	// Basic attributes
	Name        string
	Type        string
	Description string
	Required    bool
	Optional    bool
	Computed    bool
	ForceNew    bool
	Deprecated  bool

	PossibleValues []string
	DefaultValue   interface{} // Default value can be many types, TODO: convert func to cast from interface{} to string and change this field type to string

	// Block related attributes
	Nested          *Properties
	Block           bool
	BlockHasSection bool // TODO?

	// List or map related attributes
	NestedType string

	// Documentation related attributes
	AdditionalLines []string // Tracks any lines from docs beyond initial property, e.g. notes
	Count           int      // Property count, for doc parsing to detect duplicate entries
}

func NewProperties() *Properties {
	return &Properties{
		Names:   make([]string, 0),
		Objects: make(map[string]*Property),
	}
}

func (p *Property) String() string {
	return "TODO"
}

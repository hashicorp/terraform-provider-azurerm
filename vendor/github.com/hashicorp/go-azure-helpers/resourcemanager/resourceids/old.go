package resourceids

// Id defines a type for a ResourceId of some kind
type Id interface {
	// ID returns the fully formatted ID for this Resource ID
	ID() string

	// String returns a friendly description of the components of this Resource ID
	// which is suitable for use in error messages (for example 'MyThing %q / Resource Group %q')
	String() string
}

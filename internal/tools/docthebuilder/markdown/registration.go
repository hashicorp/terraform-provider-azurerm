package markdown

func RegisteredSections() []Section {
	return []Section{
		&FrontMatterSection{},
		&ExampleSection{},
		&ArgumentsSection{},
		&AttributesSection{},
		&TimeoutsSection{},
		&APISection{},
	}
}

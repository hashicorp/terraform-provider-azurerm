package environments

type Api interface {
	DomainSuffix() (*string, bool)
	Endpoint() (*string, bool)
	MicrosoftGraphAppId() (*string, bool)
	Name() string
	ResourceIdentifier() (*string, bool)
}

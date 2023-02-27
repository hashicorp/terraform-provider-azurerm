package environments

type Api interface {
	AppId() (*string, bool)
	DomainSuffix() (*string, bool)
	Endpoint() (*string, bool)
	Name() string
	ResourceIdentifier() (*string, bool)
}

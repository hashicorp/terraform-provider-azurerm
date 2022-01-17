package servergroups

type ServerNameItem struct {
	FullyQualifiedDomainName *string `json:"fullyQualifiedDomainName,omitempty"`
	Name                     *string `json:"name,omitempty"`
}

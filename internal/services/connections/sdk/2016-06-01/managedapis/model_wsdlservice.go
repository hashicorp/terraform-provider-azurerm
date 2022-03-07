package managedapis

type WsdlService struct {
	EndpointQualifiedNames *[]string `json:"endpointQualifiedNames,omitempty"`
	QualifiedName          string    `json:"qualifiedName"`
}

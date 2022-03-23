package objectreplicationpolicies

import "github.com/Azure/go-autorest/autorest"

type ObjectReplicationPoliciesClient struct {
	Client  autorest.Client
	baseUri string
}

func NewObjectReplicationPoliciesClientWithBaseURI(endpoint string) ObjectReplicationPoliciesClient {
	return ObjectReplicationPoliciesClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

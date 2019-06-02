package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2017-10-12/cdn"
)

type Client struct {
	CustomDomainsClient cdn.CustomDomainsClient
	EndpointsClient     cdn.EndpointsClient
	ProfilesClient      cdn.ProfilesClient
}

package virtualmachine

import "github.com/Azure/go-autorest/autorest"

type VirtualMachineClient struct {
	Client  autorest.Client
	baseUri string
}

func NewVirtualMachineClientWithBaseURI(endpoint string) VirtualMachineClient {
	return VirtualMachineClient{
		Client:  autorest.NewClientWithUserAgent(userAgent()),
		baseUri: endpoint,
	}
}

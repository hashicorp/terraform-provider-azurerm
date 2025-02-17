package workspace

import (
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

const defaultApiVersion = "2024-05-01"

type WorkspaceClient struct {
	Client *resourcemanager.Client
}

func NewWorkspaceClientWithBaseURI(sdkApi sdkEnv.Api) (*WorkspaceClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "workspace", defaultApiVersion)
	if err != nil {
		return nil, err
	}

	return &WorkspaceClient{
		Client: client,
	}, nil
}

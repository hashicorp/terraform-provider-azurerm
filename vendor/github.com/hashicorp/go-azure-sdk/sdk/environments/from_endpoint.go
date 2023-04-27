// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package environments

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/internal/metadata"
)

// FromEndpoint attempts to load an environment from the given Endpoint.
func FromEndpoint(ctx context.Context, endpoint, name string) (*Environment, error) {
	env := baseEnvironmentWithName("FromEnvironment")

	client := metadata.NewClientWithEndpoint(endpoint)
	config, err := client.GetMetaData(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving metadata from endpoint %q: %+v", endpoint, err)
	}

	if config.Name == "" {
		return nil, fmt.Errorf("retrieving metadata from endpoint: `name` was nil")
	}
	env.Name = config.Name

	if config.ResourceManagerEndpoint == "" {
		return nil, fmt.Errorf("retrieving metadata from endpoint: no `resourceManagerEndpoint` was returned")
	}
	env.ResourceManager = ResourceManagerAPI(config.ResourceManagerEndpoint)

	if config.ResourceIdentifiers.MicrosoftGraph == "" {
		return nil, fmt.Errorf("retrieving metdata from endpoint: no `microsoftGraphResourceId` was returned")
	}
	env.MicrosoftGraph = MicrosoftGraphAPI(config.ResourceIdentifiers.MicrosoftGraph)

	if err := env.updateFromMetaData(config); err != nil {
		return nil, fmt.Errorf("updating Environment from MetaData: %+v", err)
	}

	return &env, nil
}

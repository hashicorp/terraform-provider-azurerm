// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package load_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTestService.V20221201.LoadTests

			var config LoadTestResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := loadtests.NewLoadTestID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload loadtests.LoadTestResource
			if err := r.mapLoadTestResourceSchemaToLoadTestResource(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) mapLoadTestResourceSchemaToLoadTestResource(input LoadTestResourceSchema, output *loadtests.LoadTestResource) error {
	if err := r.EnsureEncryptionIdentityIDExistsInIdentity(input); err != nil {
		return err
	}

	identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(input.Identity)
	if err != nil {
		return fmt.Errorf("expanding Legacy SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Location = location.Normalize(input.Location)
	output.Tags = tags.Expand(input.Tags)

	if output.Properties == nil {
		output.Properties = &loadtests.LoadTestProperties{}
	}
	if err := r.mapLoadTestResourceSchemaToLoadTestProperties(input, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "LoadTestProperties", "Properties", err)
	}

	return nil
}

func (r Resource) mapLoadTestResourceSchemaToLoadTestResourceUpdate(input LoadTestResourceSchema, output *loadtests.LoadTestResourceUpdate) error {
	identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(input.Identity)
	if err != nil {
		return fmt.Errorf("expanding Legacy SystemAndUserAssigned Identity: %+v", err)
	}
	output.Identity = identity

	output.Tags = tags.Expand(input.Tags)

	output.Properties = &loadtests.LoadTestResourceUpdateProperties{
		Description: pointer.To(input.Description),
	}
	return nil
}

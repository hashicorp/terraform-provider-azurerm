// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package disks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagepool/2021-08-01/diskpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	disksValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/disks/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = DiskPoolResource{}
var _ sdk.ResourceWithDeprecationAndNoReplacement = DiskPoolResource{}

type DiskPoolResource struct{}

func (DiskPoolResource) DeprecationMessage() string {
	return "The `azurerm_disk_pool` resource is deprecated and will be removed in v4.0 of the AzureRM Provider."
}

type DiskPoolResourceModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	Sku               string                 `tfschema:"sku_name"`
	SubnetId          string                 `tfschema:"subnet_id"`
	Tags              map[string]interface{} `tfschema:"tags"`
	Zones             zones.Schema           `tfschema:"zones"`
}

func (DiskPoolResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: disksValidate.DiskPoolName(),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: disksValidate.DiskPoolSku(),
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleRequiredForceNew(),
	}
}

func (DiskPoolResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (DiskPoolResource) ModelObject() interface{} {
	return &DiskPoolResourceModel{}
}

func (DiskPoolResource) ResourceType() string {
	return "azurerm_disk_pool"
}

func (r DiskPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Disks.DiskPoolsClient

			m := DiskPoolResourceModel{}
			err := metadata.Decode(&m)
			if err != nil {
				return err
			}

			id := diskpools.NewDiskPoolID(subscriptionId, m.ResourceGroupName, m.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			createParameter := diskpools.DiskPoolCreate{
				Location: location.Normalize(m.Location),
				Properties: diskpools.DiskPoolCreateProperties{
					AvailabilityZones: &m.Zones,
					SubnetId:          m.SubnetId,
				},
				Sku:  expandDisksPoolSku(m.Sku),
				Tags: tags.Expand(m.Tags),
			}
			future, err := client.CreateOrUpdate(ctx, id, createParameter)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}

			//lintignore:R006
			return pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
				if err := r.retryError("waiting for creation", id.ID(), future.Poller.PollUntilDone(ctx)); err != nil {
					return err
				}
				metadata.SetID(id)
				return nil
			})
		},
	}
}

func (DiskPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			m := DiskPoolResourceModel{
				Name:              id.DiskPoolName,
				ResourceGroupName: id.ResourceGroupName,
			}
			if model := resp.Model; model != nil {
				if model.Sku != nil {
					m.Sku = model.Sku.Name
				}
				m.Tags = tags.Flatten(model.Tags)

				m.Location = location.Normalize(model.Location)
				m.SubnetId = model.Properties.SubnetId
				m.Zones = model.Properties.AvailabilityZones
			}

			return metadata.Encode(&m)
		},
	}
}

func (r DiskPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			future, err := client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id)
			}

			//lintignore:R006
			return pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
				return r.retryError("waiting for deletion", id.ID(), future.Poller.PollUntilDone(ctx))
			})
		},
	}
}

func (DiskPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return diskpools.ValidateDiskPoolID
}

func (r DiskPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Disks.DiskPoolsClient
			id, err := diskpools.ParseDiskPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(metadata.ResourceData.Id())
			defer locks.UnlockByID(metadata.ResourceData.Id())

			patch := diskpools.DiskPoolUpdate{}
			var m DiskPoolResourceModel
			if err = metadata.Decode(&m); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			if metadata.ResourceData.HasChange("sku") {
				sku := expandDisksPoolSku(m.Sku)
				patch.Sku = &sku
			}
			if metadata.ResourceData.HasChange("tags") {
				patch.Tags = tags.Expand(m.Tags)
			}

			future, err := client.Update(ctx, *id, patch)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("could not retrieve context deadline for %s", id.ID())
			}

			//lintignore:R006
			return pluginsdk.Retry(time.Until(deadline), func() *pluginsdk.RetryError {
				return r.retryError("waiting for update", id.ID(), future.Poller.PollUntilDone(ctx))
			})
		},
	}
}

func (DiskPoolResource) retryError(action string, id string, err error) *pluginsdk.RetryError {
	if err == nil {
		return nil
	}
	// according to https://docs.microsoft.com/en-us/azure/virtual-machines/disks-pools-troubleshoot#common-failure-codes-when-deploying-a-disk-pool the errors below are retryable.
	retryableErrors := []string{
		"DeploymentTimeout",
		"GoalStateApplicationTimeoutError",
		"OngoingOperationInProgress",
	}
	for _, retryableError := range retryableErrors {
		if strings.Contains(err.Error(), retryableError) {
			return pluginsdk.RetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
		}
	}
	return pluginsdk.NonRetryableError(fmt.Errorf("%s %s: %+v", action, id, err))
}

func expandDisksPoolSku(sku string) diskpools.Sku {
	parts := strings.Split(sku, "_")
	return diskpools.Sku{
		Name: sku,
		Tier: &parts[0],
	}
}

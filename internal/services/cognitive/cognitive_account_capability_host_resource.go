// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountcapabilityhost"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = CognitiveAccountCapabilityHostResource{}

type CognitiveAccountCapabilityHostModel struct {
	AiServicesConnections    []string `tfschema:"ai_services_connections"`
	CognitiveAccountId       string   `tfschema:"cognitive_account_id"`
	Name                     string   `tfschema:"name"`
	StorageConnections       []string `tfschema:"storage_connections"`
	SubnetId                 string   `tfschema:"subnet_id"`
	ThreadStorageConnections []string `tfschema:"thread_storage_connections"`
	VectorStoreConnections   []string `tfschema:"vector_store_connections"`
}

type CognitiveAccountCapabilityHostResource struct{}

func (r CognitiveAccountCapabilityHostResource) ResourceType() string {
	return "azurerm_cognitive_account_capability_host"
}

func (r CognitiveAccountCapabilityHostResource) ModelObject() interface{} {
	return &CognitiveAccountCapabilityHostModel{}
}

func (r CognitiveAccountCapabilityHostResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountcapabilityhost.ValidateCapabilityHostID
}

func (r CognitiveAccountCapabilityHostResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cognitiveservicesaccounts.AccountId{}),

		"ai_services_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"storage_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"thread_storage_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"vector_store_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r CognitiveAccountCapabilityHostResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountCapabilityHostResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CognitiveAccountCapabilityHostModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Cognitive.AccountCapabilityHostClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := accountcapabilityhost.NewCapabilityHostID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			resource := accountcapabilityhost.CapabilityHostResource{
				Properties: accountcapabilityhost.CapabilityHost{},
			}

			if model.SubnetId != "" {
				resource.Properties.CustomerSubnet = pointer.To(model.SubnetId)
			}

			if len(model.AiServicesConnections) > 0 {
				resource.Properties.AiServicesConnections = pointer.To(model.AiServicesConnections)
			}

			if len(model.StorageConnections) > 0 {
				resource.Properties.StorageConnections = pointer.To(model.StorageConnections)
			}

			if len(model.ThreadStorageConnections) > 0 {
				resource.Properties.ThreadStorageConnections = pointer.To(model.ThreadStorageConnections)
			}

			if len(model.VectorStoreConnections) > 0 {
				resource.Properties.VectorStoreConnections = pointer.To(model.VectorStoreConnections)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, resource); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveAccountCapabilityHostResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountCapabilityHostClient

			id, err := accountcapabilityhost.ParseCapabilityHostID(metadata.ResourceData.Id())
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

			state := CognitiveAccountCapabilityHostModel{
				Name:               id.CapabilityHostName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				state.AiServicesConnections = pointer.From(props.AiServicesConnections)
				state.StorageConnections = pointer.From(props.StorageConnections)
				state.ThreadStorageConnections = pointer.From(props.ThreadStorageConnections)
				state.VectorStoreConnections = pointer.From(props.VectorStoreConnections)

				if props.CustomerSubnet != nil && *props.CustomerSubnet != "" {
					subnetId, err := commonids.ParseSubnetIDInsensitively(*props.CustomerSubnet)
					if err != nil {
						return err
					}
					state.SubnetId = subnetId.ID()
				}

			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountCapabilityHostResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountCapabilityHostClient

			id, err := accountcapabilityhost.ParseCapabilityHostID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountcapabilityhost"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_capability_host -properties "name" -compare-values "subscription_id:cognitive_account_id,resource_group_name:cognitive_account_id,account_name:cognitive_account_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.Resource             = CognitiveAccountCapabilityHostResource{}
	_ sdk.ResourceWithIdentity = CognitiveAccountCapabilityHostResource{}
)

type CognitiveAccountCapabilityHostModel struct {
	Name                     string   `tfschema:"name"`
	CognitiveAccountId       string   `tfschema:"cognitive_account_id"`
	AiServicesConnections    []string `tfschema:"ai_services_connections"`
	StorageConnections       []string `tfschema:"storage_connections"`
	ThreadStorageConnections []string `tfschema:"thread_storage_connections"`
	VectorStoreConnections   []string `tfschema:"vector_store_connections"`
}

type CognitiveAccountCapabilityHostResource struct{}

func (r CognitiveAccountCapabilityHostResource) Identity() resourceids.ResourceId {
	return new(accountcapabilityhost.CapabilityHostId)
}

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
			ValidateFunc: validate.CapabilityHostName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cognitiveservicesaccounts.AccountId{}),

		// The Foundry Agent Service backend currently only accepts zero or exactly one non-empty
		// entry per connection field (see CreateCapabilityHostRequestDtoValidator), but the REST API
		// and SDK type these as `*[]string`. They are kept as lists here, with `MaxItems: 1` enforcing
		// the current backend restriction, so no breaking schema change is required if Azure later
		// allows more than one connection.
		"ai_services_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.AccountConnectionName(),
			},
		},

		"storage_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.AccountConnectionName(),
			},
		},

		"thread_storage_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.AccountConnectionName(),
			},
		},

		"vector_store_connections": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.AccountConnectionName(),
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			resource := accountcapabilityhost.CapabilityHost{
				Properties: accountcapabilityhost.CapabilityHostProperties{
					CapabilityHostKind: pointer.To(accountcapabilityhost.CapabilityHostKindAgents),
				},
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, resource, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

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

			return r.flatten(metadata, id, resp.Model)
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

func (r CognitiveAccountCapabilityHostResource) flatten(metadata sdk.ResourceMetaData, id *accountcapabilityhost.CapabilityHostId, model *accountcapabilityhost.CapabilityHost) error {
	state := CognitiveAccountCapabilityHostModel{
		Name:               id.CapabilityHostName,
		CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if model != nil {
		props := model.Properties

		state.AiServicesConnections = pointer.From(props.AiServicesConnections)
		state.StorageConnections = pointer.From(props.StorageConnections)
		state.ThreadStorageConnections = pointer.From(props.ThreadStorageConnections)
		state.VectorStoreConnections = pointer.From(props.VectorStoreConnections)
	}

	return metadata.Encode(&state)
}

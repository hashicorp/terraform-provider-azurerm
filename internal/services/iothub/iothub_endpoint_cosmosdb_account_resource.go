// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

type IotHubEndpointCosmosDBAccountResource struct{}

var (
	_ sdk.ResourceWithUpdate = IotHubEndpointCosmosDBAccountResource{}
)

type IotHubEndpointCosmosDBAccountModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroupName    string `tfschema:"resource_group_name"`
	AuthenticationType   string `tfschema:"authentication_type"`
	ContainerName        string `tfschema:"container_name"`
	DatabaseName         string `tfschema:"database_name"`
	EndpointUri          string `tfschema:"endpoint_uri"`
	IdentityId           string `tfschema:"identity_id"`
	IothubId             string `tfschema:"iothub_id"`
	PartitionKeyName     string `tfschema:"partition_key_name"`
	PartitionKeyTemplate string `tfschema:"partition_key_template"`
	PrimaryKey           string `tfschema:"primary_key"`
	SecondaryKey         string `tfschema:"secondary_key"`
}

func (r IotHubEndpointCosmosDBAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"database_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint_uri": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.AuthenticationTypeKeyBased),
				string(devices.AuthenticationTypeIdentityBased),
			}, false),
		},

		"identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"primary_key", "secondary_key"},
		},

		"partition_key_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"partition_key_template"},
		},

		"partition_key_template": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"partition_key_name"},
		},

		"primary_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"identity_id"},
			RequiredWith:  []string{"secondary_key"},
		},

		"secondary_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"identity_id"},
			RequiredWith:  []string{"primary_key"},
		},
	}
}

func (r IotHubEndpointCosmosDBAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r IotHubEndpointCosmosDBAccountResource) ResourceType() string {
	return "azurerm_iothub_endpoint_cosmosdb_account"
}

func (r IotHubEndpointCosmosDBAccountResource) ModelObject() interface{} {
	return &IotHubEndpointCosmosDBAccountResource{}
}

func (r IotHubEndpointCosmosDBAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.EndpointCosmosDBAccountID
}

func (r IotHubEndpointCosmosDBAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state IotHubEndpointCosmosDBAccountModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.IoTHub.ResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			iotHubId, err := parse.IotHubID(state.IothubId)
			if err != nil {
				return err
			}

			id := parse.NewEndpointCosmosDBAccountID(subscriptionId, iotHubId.ResourceGroup, iotHubId.Name, state.Name)

			locks.ByName(iotHubId.Name, IothubResourceName)
			defer locks.UnlockByName(iotHubId.Name, IothubResourceName)

			iothub, err := client.Get(ctx, iotHubId.ResourceGroup, iotHubId.Name)
			if err != nil {
				if utils.ResponseWasNotFound(iothub.Response) {
					return fmt.Errorf("%q was not found", iotHubId)
				}

				return fmt.Errorf("retrieving %q: %+v", iotHubId, err)
			}

			authenticationType := devices.AuthenticationType(state.AuthenticationType)
			cosmosDBAccountEndpoint := devices.RoutingCosmosDBSQLAPIProperties{
				Name:               pointer.To(id.EndpointName),
				SubscriptionID:     pointer.To(subscriptionId),
				ResourceGroup:      pointer.To(state.ResourceGroupName),
				AuthenticationType: authenticationType,
				CollectionName:     pointer.To(state.ContainerName),
				DatabaseName:       pointer.To(state.DatabaseName),
				EndpointURI:        pointer.To(state.EndpointUri),
			}

			if state.PartitionKeyName != "" {
				cosmosDBAccountEndpoint.PartitionKeyName = pointer.To(state.PartitionKeyName)
			}

			if state.PartitionKeyTemplate != "" {
				cosmosDBAccountEndpoint.PartitionKeyTemplate = pointer.To(state.PartitionKeyTemplate)
			}

			if authenticationType == devices.AuthenticationTypeKeyBased {
				if state.PrimaryKey == "" || state.SecondaryKey == "" {
					return fmt.Errorf("`primary_key` and `secondary_key` must be specified when `authentication_type` is `keyBased`")
				}
				cosmosDBAccountEndpoint.PrimaryKey = pointer.To(state.PrimaryKey)
				cosmosDBAccountEndpoint.SecondaryKey = pointer.To(state.SecondaryKey)
			} else {
				if state.PrimaryKey != "" || state.SecondaryKey != "" {
					return fmt.Errorf("`primary_key` or `secondary_key` cannot be specified when `authentication_type` is `identityBased`")
				}

				if state.IdentityId != "" {
					cosmosDBAccountEndpoint.Identity = &devices.ManagedIdentity{
						UserAssignedIdentity: pointer.To(state.IdentityId),
					}
				}
			}

			routing := iothub.Properties.Routing
			if routing == nil {
				routing = &devices.RoutingProperties{}
			}

			if routing.Endpoints == nil {
				routing.Endpoints = &devices.RoutingEndpoints{}
			}

			if routing.Endpoints.CosmosDBSQLCollections == nil {
				cosmosDBAccounts := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)
				routing.Endpoints.CosmosDBSQLCollections = &cosmosDBAccounts
			}

			endpoints := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)

			for _, existingEndpoint := range pointer.From(routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(existingEndpoint.Name), id.EndpointName) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
				endpoints = append(endpoints, existingEndpoint)
			}

			endpoints = append(endpoints, cosmosDBAccountEndpoint)
			routing.Endpoints.CosmosDBSQLCollections = &endpoints

			future, err := client.CreateOrUpdate(ctx, iotHubId.ResourceGroup, iotHubId.Name, iothub, "")
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the completion of the creation of %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotHubEndpointCosmosDBAccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var oldState IotHubEndpointCosmosDBAccountModel
			if err = metadata.Decode(&oldState); err != nil {
				return err
			}

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if utils.ResponseWasNotFound(iothub.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return metadata.MarkAsGone(id)
			}

			for _, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					state := &IotHubEndpointCosmosDBAccountModel{
						Name:                 id.EndpointName,
						ResourceGroupName:    pointer.From(endpoint.ResourceGroup),
						IothubId:             parse.NewIotHubID(id.SubscriptionId, id.ResourceGroup, id.IotHubName).ID(),
						ContainerName:        pointer.From(endpoint.CollectionName),
						DatabaseName:         pointer.From(endpoint.DatabaseName),
						EndpointUri:          pointer.From(endpoint.EndpointURI),
						PartitionKeyName:     pointer.From(endpoint.PartitionKeyName),
						PartitionKeyTemplate: pointer.From(endpoint.PartitionKeyTemplate),
						PrimaryKey:           oldState.PrimaryKey,
						SecondaryKey:         oldState.SecondaryKey,
					}

					authenticationType := string(devices.AuthenticationTypeKeyBased)
					if string(endpoint.AuthenticationType) != "" {
						authenticationType = string(endpoint.AuthenticationType)
					}
					state.AuthenticationType = authenticationType

					identityId := ""
					if endpoint.Identity != nil && endpoint.Identity.UserAssignedIdentity != nil {
						identityId = pointer.From(endpoint.Identity.UserAssignedIdentity)
					}
					state.IdentityId = identityId

					return metadata.Encode(state)
				}
			}

			return metadata.MarkAsGone(id)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r IotHubEndpointCosmosDBAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient

			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.IotHubName, IothubResourceName)
			defer locks.UnlockByName(id.IotHubName, IothubResourceName)

			var state IotHubEndpointCosmosDBAccountModel
			if err = metadata.Decode(&state); err != nil {
				return err
			}

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if utils.ResponseWasNotFound(iothub.Response) {
					return fmt.Errorf("%q was not found", id)
				}

				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return fmt.Errorf("%q was not found", id)
			}

			for i, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					authenticationType := devices.AuthenticationType(state.AuthenticationType)

					if authenticationType == devices.AuthenticationTypeKeyBased {
						if state.PrimaryKey == "" || state.SecondaryKey == "" {
							return fmt.Errorf("`primary_key` and `secondary_key` must be specified when `authentication_type` is `keyBased`")
						}

						endpoint.PrimaryKey = pointer.To(state.PrimaryKey)
						endpoint.SecondaryKey = pointer.To(state.SecondaryKey)
						endpoint.Identity = nil
					} else {
						if state.PrimaryKey != "" || state.SecondaryKey != "" {
							return fmt.Errorf("`primary_key` or `secondary_key` cannot be specified when `authentication_type` is `identityBased`")
						}

						endpoint.PrimaryKey = nil
						endpoint.SecondaryKey = nil
						if state.IdentityId != "" {
							endpoint.Identity = &devices.ManagedIdentity{
								UserAssignedIdentity: pointer.To(state.IdentityId),
							}
						} else {
							endpoint.Identity = nil
						}
					}

					if metadata.ResourceData.HasChange("authentication_type") {
						endpoint.AuthenticationType = authenticationType
					}

					if metadata.ResourceData.HasChange("partition_key_name") {
						if state.PartitionKeyName == "" {
							endpoint.PartitionKeyName = nil
						} else {
							endpoint.PartitionKeyName = pointer.To(state.PartitionKeyName)
						}
					}

					if metadata.ResourceData.HasChange("partition_key_template") {
						if state.PartitionKeyTemplate == "" {
							endpoint.PartitionKeyTemplate = nil
						} else {
							endpoint.PartitionKeyTemplate = pointer.To(state.PartitionKeyTemplate)
						}
					}

					(*iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections)[i] = endpoint

					future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
					if err != nil {
						return fmt.Errorf("updating %s: %+v", id, err)
					}

					if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for the completion of the update of %s: %+v", id, err)
					}

					return nil
				}
			}

			return fmt.Errorf("%q was not found", id)
		},
		Timeout: 30 * time.Minute,
	}
}

func (r IotHubEndpointCosmosDBAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.IotHubName, IothubResourceName)
			defer locks.UnlockByName(id.IotHubName, IothubResourceName)

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if utils.ResponseWasNotFound(iothub.Response) {
					return fmt.Errorf("%q was not found", id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return nil
			}

			updatedEndpoints := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)
			for _, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if !strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					updatedEndpoints = append(updatedEndpoints, endpoint)
				}
			}
			iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections = &updatedEndpoints

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

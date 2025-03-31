// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package extendedlocation

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// The resource name should be `azurerm_extended_location_custom_location` as in the resource doc,
// but in source code it was named `azurerm_extended_custom_location`.
// This resource is superseded by a new resource with the correct name, and will be removed in 5.0
func (r CustomLocationResource) DeprecatedInFavourOfResource() string {
	return "azurerm_extended_location_custom_location"
}

var (
	_ sdk.ResourceWithUpdate                = CustomLocationResource{}
	_ sdk.ResourceWithDeprecationReplacedBy = CustomLocationResource{}
)

type CustomLocationResource struct{}

func (r CustomLocationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z\d.\-_]*[A-Za-z\d]$`),
				"supported alphanumeric characters and periods, underscores, hyphens. Name should end with an alphanumeric character.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"namespace": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"namespace must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"host_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: arckubernetes.ValidateConnectedClusterID,
		},

		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem:     commonschema.ResourceIDReferenceElem(&extensions.ScopedExtensionId{}),
		},

		"host_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(customlocations.HostTypeKubernetes),
			}, false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsBase64,
					},
				},
			},
		},
	}
}

func (r CustomLocationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CustomLocationResource) ModelObject() interface{} {
	return &ExtendedLocationCustomLocationResourceModel{}
}

func (r CustomLocationResource) ResourceType() string {
	return "azurerm_extended_custom_location"
}

func (r CustomLocationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ExtendedLocationCustomLocationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.ExtendedLocation.CustomLocationsClient

			id := customlocations.NewCustomLocationID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customLocationProps := customlocations.CustomLocationProperties{
				ClusterExtensionIds: pointer.To(model.ClusterExtensionIds),
				DisplayName:         pointer.To(model.DisplayName),
				HostResourceId:      pointer.To(model.HostResourceId),
				HostType:            pointer.To(customlocations.HostType(model.HostType)),
				Namespace:           pointer.To(model.Namespace),
			}

			if len(model.Authentication) > 0 {
				auth := model.Authentication[0]
				customLocationProps.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
					Type:  pointer.To(auth.Type),
					Value: pointer.To(auth.Value),
				}
			}

			props := customlocations.CustomLocation{
				Id:         pointer.To(id.ID()),
				Location:   location.Normalize(model.Location),
				Name:       pointer.To(model.Name),
				Properties: pointer.To(customLocationProps),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CustomLocationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
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

			if model := resp.Model; model != nil {
				props := model.Properties

				state := ExtendedLocationCustomLocationResourceModel{
					Name:                id.CustomLocationName,
					ResourceGroupName:   id.ResourceGroupName,
					Location:            location.Normalize(model.Location),
					ClusterExtensionIds: pointer.From(props.ClusterExtensionIds),
					DisplayName:         pointer.From(props.DisplayName),
					HostResourceId:      pointer.From(props.HostResourceId),
					HostType:            string(pointer.From(props.HostType)),
					Namespace:           pointer.From(props.Namespace),
				}

				// API always returns an empty `authentication` block even it's not specified. Tracing the bug: https://github.com/Azure/azure-rest-api-specs/issues/30101
				if props.Authentication != nil && props.Authentication.Type != nil && props.Authentication.Value != nil {
					state.Authentication = []AuthModel{
						{
							Type:  pointer.From(props.Authentication.Type),
							Value: pointer.From(props.Authentication.Value),
						},
					}
				}

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r CustomLocationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r CustomLocationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ExtendedLocation.CustomLocationsClient
			id, err := customlocations.ParseCustomLocationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ExtendedLocationCustomLocationResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return err
			}

			model := existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retreiving properties for %s for update: %+v", *id, err)
			}
			d := metadata.ResourceData

			if d.HasChanges("authentication") {
				if len(state.Authentication) > 0 {
					auth := state.Authentication[0]
					model.Properties.Authentication = &customlocations.CustomLocationPropertiesAuthentication{
						Type:  pointer.To(auth.Type),
						Value: pointer.To(auth.Value),
					}
				}
			} else {
				model.Properties.Authentication = nil
			}

			if d.HasChange("display_name") {
				model.Properties.DisplayName = pointer.To(state.DisplayName)
			}

			if d.HasChange("cluster_extension_ids") {
				model.Properties.ClusterExtensionIds = pointer.To(state.ClusterExtensionIds)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r CustomLocationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return customlocations.ValidateCustomLocationID
}

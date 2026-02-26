// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package confluent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confluent/2024-07-01/connectorresources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/confluent/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ConfluentConnectorResource struct{}

type ConfluentConnectorResourceModel struct {
	ConnectorName     string   `tfschema:"connector_name"`
	ClusterId         string   `tfschema:"cluster_id"`
	EnvironmentId     string   `tfschema:"environment_id"`
	OrganizationId    string   `tfschema:"organization_id"`
	ResourceGroupName string   `tfschema:"resource_group_name"`
	ConnectorType     string   `tfschema:"connector_type"`
	ConnectorClass    string   `tfschema:"connector_class"`

	// Computed
	Id            string `tfschema:"id"`
	ConnectorId   string `tfschema:"connector_id"`
	ConnectorState string `tfschema:"connector_state"`
}

func (r ConfluentConnectorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"connector_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"organization_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"connector_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(connectorresources.PossibleValuesForConnectorType(), false),
		},

		"connector_class": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(connectorresources.PossibleValuesForConnectorClass(), false),
		},
	}
}

func (r ConfluentConnectorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connector_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"connector_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ConfluentConnectorResource) ModelObject() interface{} {
	return &ConfluentConnectorResourceModel{}
}

func (r ConfluentConnectorResource) ResourceType() string {
	return "azurerm_confluent_connector"
}

func (r ConfluentConnectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ConnectorClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ConfluentConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := connectorresources.NewConnectorID(subscriptionId, model.ResourceGroupName, model.OrganizationId, model.EnvironmentId, model.ClusterId, model.ConnectorName)

			existing, err := client.ConnectorGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_confluent_connector", id.ID())
			}

			payload := connectorresources.ConnectorResource{
				Name:       pointer.To(model.ConnectorName),
				Properties: expandConfluentConnectorProperties(model),
			}

			if _, err := client.ConnectorCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ConfluentConnectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ConnectorClient

			id, err := connectorresources.ParseConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ConnectorGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var state ConfluentConnectorResourceModel
			state.ConnectorName = id.ConnectorName
			state.ClusterId = id.ClusterId
			state.EnvironmentId = id.EnvironmentId
			state.OrganizationId = id.OrganizationName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Id = pointer.From(model.Id)

				props := model.Properties
				if props.ConnectorBasicInfo != nil {
					state.ConnectorId = pointer.From(props.ConnectorBasicInfo.ConnectorId)
					state.ConnectorState = string(pointer.From(props.ConnectorBasicInfo.ConnectorState))
					state.ConnectorType = string(pointer.From(props.ConnectorBasicInfo.ConnectorType))
					state.ConnectorClass = string(pointer.From(props.ConnectorBasicInfo.ConnectorClass))
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ConfluentConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Confluent.ConnectorClient

			id, err := connectorresources.ParseConnectorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ConnectorDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ConfluentConnectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ConnectorID
}

func expandConfluentConnectorProperties(model ConfluentConnectorResourceModel) connectorresources.ConnectorResourceProperties {
	// Note: The full implementation would require handling the polymorphic types
	// (ConnectorServiceTypeInfoBase and PartnerInfoBase) which have multiple implementations
	// For now, we create a minimal structure that allows the connector to be created
	// Users will need to configure the connector details through other means (e.g., Confluent portal/CLI)

	props := connectorresources.ConnectorResourceProperties{}

	if model.ConnectorType != "" || model.ConnectorClass != "" {
		info := &connectorresources.ConnectorInfoBase{}

		if model.ConnectorType != "" {
			info.ConnectorType = pointer.To(connectorresources.ConnectorType(model.ConnectorType))
		}

		if model.ConnectorClass != "" {
			info.ConnectorClass = pointer.To(connectorresources.ConnectorClass(model.ConnectorClass))
		}

		props.ConnectorBasicInfo = info
	}

	return props
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppDataSource struct{}

type ContainerAppDataSourceModel struct {
	Name                       string                                     `tfschema:"name"`
	ResourceGroup              string                                     `tfschema:"resource_group_name"`
	ManagedEnvironmentId       string                                     `tfschema:"container_app_environment_id"`
	Location                   string                                     `tfschema:"location"`
	RevisionMode               string                                     `tfschema:"revision_mode"`
	MaxInactiveRevisions       int64                                      `tfschema:"max_inactive_revisions"`
	Ingress                    []helpers.Ingress                          `tfschema:"ingress"`
	Registries                 []helpers.Registry                         `tfschema:"registry"`
	Secrets                    []helpers.Secret                           `tfschema:"secret"`
	Dapr                       []helpers.Dapr                             `tfschema:"dapr"`
	Template                   []helpers.ContainerTemplate                `tfschema:"template"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags                       map[string]interface{}                     `tfschema:"tags"`
	OutboundIpAddresses        []string                                   `tfschema:"outbound_ip_addresses"`
	LatestRevisionName         string                                     `tfschema:"latest_revision_name"`
	LatestRevisionFqdn         string                                     `tfschema:"latest_revision_fqdn"`
	CustomDomainVerificationId string                                     `tfschema:"custom_domain_verification_id"`
	WorkloadProfileName        string                                     `tfschema:"workload_profile_name"`
}

var _ sdk.DataSource = ContainerAppDataSource{}

func (r ContainerAppDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ContainerAppDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"template": helpers.ContainerTemplateSchemaComputed(),

		"revision_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ingress": helpers.ContainerAppIngressSchemaComputed(),

		"registry": helpers.ContainerAppRegistrySchemaComputed(),

		"secret": helpers.SecretsDataSourceSchema(),

		"dapr": helpers.ContainerDaprSchemaComputed(),

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"tags": commonschema.TagsDataSource(),

		"location": commonschema.LocationComputed(),

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"latest_revision_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The name of the latest Container Revision.",
		},

		"latest_revision_fqdn": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The fully qualified domain name of the latest Container App.",
		},

		"custom_domain_verification_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "The ID of the Custom Domain Verification for this Container App.",
		},

		"workload_profile_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"max_inactive_revisions": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

func (r ContainerAppDataSource) ModelObject() interface{} {
	return &ContainerAppDataSourceModel{}
}

func (r ContainerAppDataSource) ResourceType() string {
	return "azurerm_container_app"
}

func (r ContainerAppDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var containerApp ContainerAppDataSourceModel
			if err := metadata.Decode(&containerApp); err != nil {
				return err
			}

			id := containerapps.NewContainerAppID(subscriptionId, containerApp.ResourceGroup, containerApp.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			containerApp.Name = id.ContainerAppName
			containerApp.ResourceGroup = id.ResourceGroupName

			if model := existing.Model; model != nil {
				containerApp.Location = location.Normalize(model.Location)
				containerApp.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					envId, err := managedenvironments.ParseManagedEnvironmentIDInsensitively(pointer.From(props.ManagedEnvironmentId))
					if err != nil {
						return err
					}
					containerApp.ManagedEnvironmentId = envId.ID()
					containerApp.Template = helpers.FlattenContainerAppTemplate(props.Template)
					if config := props.Configuration; config != nil {
						if config.ActiveRevisionsMode != nil {
							if config.ActiveRevisionsMode != nil {
								containerApp.RevisionMode = string(pointer.From(config.ActiveRevisionsMode))
							}
							containerApp.Ingress = helpers.FlattenContainerAppIngress(config.Ingress, id.ContainerAppName)
							containerApp.Registries = helpers.FlattenContainerAppRegistries(config.Registries)
							containerApp.Dapr = helpers.FlattenContainerAppDapr(config.Dapr)
							containerApp.MaxInactiveRevisions = pointer.ToInt64(config.MaxInactiveRevisions)
						}
					}
					containerApp.LatestRevisionName = pointer.From(props.LatestRevisionName)
					containerApp.LatestRevisionFqdn = pointer.From(props.LatestRevisionFqdn)
					containerApp.CustomDomainVerificationId = pointer.From(props.CustomDomainVerificationId)
					containerApp.OutboundIpAddresses = pointer.From(props.OutboundIPAddresses)
					containerApp.WorkloadProfileName = pointer.From(props.WorkloadProfileName)
				}
			}

			secretsResp, err := client.ListSecrets(ctx, id)
			if err != nil {
				return fmt.Errorf("listing secrets for %s: %+v", id, err)
			}

			containerApp.Secrets = helpers.FlattenContainerAppSecrets(secretsResp.Model)
			metadata.SetID(id)

			return metadata.Encode(&containerApp)
		},
	}
}

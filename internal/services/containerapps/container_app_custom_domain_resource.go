// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppCustomDomainResource struct{}

var _ sdk.Resource = ContainerAppCustomDomainResource{}

type ContainerAppCustomDomainResourceModel struct {
	Name                 string `tfschema:"name"`
	ContainerAppId       string `tfschema:"container_app_id"`
	CertificateId        string `tfschema:"container_app_environment_certificate_id"`
	BindingType          string `tfschema:"certificate_binding_type"`
	ManagedCertificateId string `tfschema:"container_app_environment_managed_certificate_id"`
}

func (a ContainerAppCustomDomainResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The hostname of the Certificate. Must be the CN or a named SAN in the certificate.",
		},

		"container_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerapps.ValidateContainerAppID,
		},

		"container_app_environment_certificate_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"certificate_binding_type"},
			ValidateFunc: managedenvironments.ValidateCertificateID,
		},

		"certificate_binding_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForBindingType(), false),
			Description:  "The Binding type. Possible values include `Disabled` and `SniEnabled`.",
		},
	}
}

func (a ContainerAppCustomDomainResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_app_environment_managed_certificate_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (a ContainerAppCustomDomainResource) ModelObject() interface{} {
	return &ContainerAppCustomDomainResourceModel{}
}

func (a ContainerAppCustomDomainResource) ResourceType() string {
	return "azurerm_container_app_custom_domain"
}

func (a ContainerAppCustomDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ContainerAppCustomDomainId
}

func (a ContainerAppCustomDomainResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			model := ContainerAppCustomDomainResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			containerAppId, err := containerapps.ParseContainerAppID(model.ContainerAppId)
			if err != nil {
				return err
			}

			locks.ByID(containerAppId.ID())
			defer locks.UnlockByID(containerAppId.ID())

			id := parse.NewContainerAppCustomDomainId(containerAppId.SubscriptionId, containerAppId.ResourceGroupName, containerAppId.ContainerAppName, model.Name)

			var certificateId *managedenvironments.CertificateId
			if model.CertificateId != "" {
				certificateId, err = managedenvironments.ParseCertificateID(model.CertificateId)
				if err != nil {
					return err
				}
			}

			containerApp, err := client.Get(ctx, *containerAppId)
			if err != nil || containerApp.Model == nil {
				return fmt.Errorf("retrieving %s to create %s", containerAppId, id)
			}

			props := containerApp.Model.Properties
			if props == nil || props.Configuration == nil {
				return fmt.Errorf("could not retrieve properties of %s", containerAppId)
			}

			config := *props.Configuration

			if config.Ingress == nil {
				return fmt.Errorf("specified Container App (%s) has no Ingress configuration for Custom Domains", containerAppId)
			}

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, *containerAppId)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", *containerAppId, err)
				}
			}
			props.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			ingress := *config.Ingress

			customDomains := make([]containerapps.CustomDomain, 0)
			if existingCustomDomains := ingress.CustomDomains; existingCustomDomains != nil {
				for _, v := range *existingCustomDomains {
					if strings.EqualFold(v.Name, model.Name) {
						return metadata.ResourceRequiresImport(ContainerAppCustomDomainResource{}.ResourceType(), id)
					}
				}

				customDomains = *existingCustomDomains
			}

			customDomain := containerapps.CustomDomain{
				Name:        model.Name,
				BindingType: pointer.To(containerapps.BindingTypeDisabled),
			}

			if certificateId != nil {
				customDomain.CertificateId = pointer.To(certificateId.ID())
				customDomain.BindingType = pointer.To(containerapps.BindingType(model.BindingType))
			}

			customDomains = append(customDomains, customDomain)

			containerApp.Model.Properties.Configuration.Ingress.CustomDomains = pointer.To(customDomains)

			if err := client.CreateOrUpdateThenPoll(ctx, *containerAppId, *containerApp.Model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (a ContainerAppCustomDomainResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient
			state := ContainerAppCustomDomainResourceModel{}
			id, err := parse.ContainerAppCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerAppId := containerapps.NewContainerAppID(id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)

			containerApp, err := client.Get(ctx, containerAppId)
			if err != nil || containerApp.Model == nil {
				return fmt.Errorf("retrieving %s to read %s", containerAppId, id)
			}

			model := containerApp.Model

			if model.Properties == nil || model.Properties.Configuration == nil || model.Properties.Configuration.Ingress == nil {
				return fmt.Errorf("could not read Ingress configuration for %s", containerAppId)
			}

			ingress := *model.Properties.Configuration.Ingress
			found := false
			if customDomains := ingress.CustomDomains; customDomains != nil {
				for _, v := range *customDomains {
					if strings.EqualFold(v.Name, id.CustomDomainName) {
						found = true
						state.Name = id.CustomDomainName
						state.ContainerAppId = containerAppId.ID()
						if pointer.From(v.CertificateId) != "" {
							// The `v.CertificateId` returned from API has two possible values. when using an Azure created Managed Certificate,
							// its format is "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/managedCertificates/%s",
							// another format is "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/certificates/%s",
							// both cases are handled here to avoid parsing error.
							certId, err1 := managedenvironments.ParseCertificateIDInsensitively(pointer.From(v.CertificateId))
							if err1 != nil {
								managedCertId, err2 := managedenvironments.ParseManagedCertificateID(pointer.From(v.CertificateId))
								if err2 != nil {
									return err1
								}
								state.ManagedCertificateId = managedCertId.ID()
							} else {
								state.CertificateId = certId.ID()
							}
						}

						state.BindingType = string(pointer.From(v.BindingType))
					}
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			return metadata.Encode(&state)
		},
	}
}

func (a ContainerAppCustomDomainResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.ContainerAppClient

			id, err := parse.ContainerAppCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			containerAppId := containerapps.NewContainerAppID(id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)

			containerApp, err := client.Get(ctx, containerAppId)
			if err != nil || containerApp.Model == nil {
				return fmt.Errorf("retrieving %s to read %s", containerAppId, id)
			}

			model := containerApp.Model

			if model.Properties == nil || model.Properties.Configuration == nil || model.Properties.Configuration.Ingress == nil {
				return fmt.Errorf("could not read Ingress configuration for %s", containerAppId)
			}

			ingress := *model.Properties.Configuration.Ingress
			updatedCustomDomains := make([]containerapps.CustomDomain, 0)
			if customDomains := ingress.CustomDomains; customDomains != nil {
				for _, v := range *customDomains {
					if !strings.EqualFold(v.Name, id.CustomDomainName) {
						updatedCustomDomains = append(updatedCustomDomains, v)
					} else {
						// attempt to lock the cert if we have the ID
						certificateId := pointer.From(v.CertificateId)
						if certificateId != "" {
							locks.ByID(certificateId)
							defer locks.UnlockByID(certificateId)
						}
					}
				}
			}

			model.Properties.Configuration.Ingress.CustomDomains = pointer.To(updatedCustomDomains)

			// Delta-updates need the secrets back from the list API, or we'll end up removing them or erroring out.
			secretsResp, err := client.ListSecrets(ctx, containerAppId)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", containerAppId, err)
				}
			}
			model.Properties.Configuration.Secrets = helpers.UnpackContainerSecretsCollection(secretsResp.Model)

			if err := client.CreateOrUpdateThenPoll(ctx, containerAppId, *model); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

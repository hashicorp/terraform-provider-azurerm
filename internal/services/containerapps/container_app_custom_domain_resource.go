package containerapps

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppCustomDomainResource struct{}

var _ sdk.Resource = ContainerAppCustomDomainResource{}

type ContainerAppCustomDomainResourceModel struct {
	Name           string `tfschema:"name"`
	ContainerAppId string `tfschema:"container_app_id"`
	CertificateId  string `tfschema:"container_app_environment_certificate_id"`
	BindingType    string `tfschema:"certificate_binding_type"`
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
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateCertificateID,
		},

		"certificate_binding_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(containerapps.PossibleValuesForBindingType(), false),
			Description:  "The Binding type. Possible values include `Disabled` and `SniEnabled`.",
		},
	}
}

func (a ContainerAppCustomDomainResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
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

			certificateId, err := managedenvironments.ParseCertificateID(model.CertificateId)
			if err != nil {
				return err
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

			customDomains = append(customDomains, containerapps.CustomDomain{
				BindingType:   pointer.To(containerapps.BindingType(model.BindingType)),
				CertificateId: pointer.To(certificateId.ID()),
				Name:          model.Name,
			})

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
						certId, err := managedenvironments.ParseCertificateIDInsensitively(pointer.From(v.CertificateId))
						if err != nil {
							return err
						}
						state.CertificateId = certId.ID()
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

			// attempt to lock the cert if we have the ID
			if certIdRaw := metadata.ResourceData.Get("container_app_environment_certificate_id").(string); certIdRaw != "" {
				if certId, err := managedenvironments.ParseCertificateID(certIdRaw); err == nil {
					locks.ByID(certId.ID())
					defer locks.UnlockByID(certId.ID())
				}
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
					}
				}
			}

			model.Properties.Configuration.Ingress.CustomDomains = pointer.To(updatedCustomDomains)

			if err := client.CreateOrUpdateThenPoll(ctx, containerAppId, *model); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

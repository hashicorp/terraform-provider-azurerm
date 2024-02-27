package appservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticSiteCustomDomainResource struct{}

var _ sdk.Resource = StaticSiteCustomDomainResource{}

const (
	txtValidationType   = "dns-txt-token"
	cnameValidationType = "cname-delegation"
)

type StaticSiteCustomDomainResourceModel struct {
	DomainName      string `tfschema:"domain_name"`
	StaticSiteId    string `tfschema:"static_site_id"`
	ValidationType  string `tfschema:"validation_type"`
	ValidationToken string `tfschema:"validation_token"`
}

func (r StaticSiteCustomDomainResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"static_site_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: staticsites.ValidateStaticSiteID,
		},

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"validation_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  txtValidationType,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				txtValidationType,
				cnameValidationType,
			}, false),
		},
	}
}

func (r StaticSiteCustomDomainResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"validation_token": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r StaticSiteCustomDomainResource) ModelObject() interface{} {
	return &StaticSiteCustomDomainResource{}
}

func (r StaticSiteCustomDomainResource) ResourceType() string {
	return "azurerm_static_site_custom_domain"
}

func (r StaticSiteCustomDomainResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			model := StaticSiteCustomDomainResourceModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			siteId, err := staticsites.ParseStaticSiteID(model.StaticSiteId)
			if err != nil {
				return err
			}

			id := staticsites.NewCustomDomainID(siteId.SubscriptionId, siteId.ResourceGroupName, siteId.StaticSiteName, model.DomainName)

			existing, err := client.GetStaticSiteCustomDomain(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("failed checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customDomain := staticsites.StaticSiteCustomDomainRequestPropertiesARMResource{
				Properties: &staticsites.StaticSiteCustomDomainRequestPropertiesARMResourceProperties{
					ValidationMethod: pointer.To(model.ValidationType),
				},
			}

			if err = client.CreateOrUpdateStaticSiteCustomDomainThenPoll(ctx, id, customDomain); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StaticSiteCustomDomainResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Some values are not retrievable from the API so we try and load the config.
			state := StaticSiteCustomDomainResourceModel{}
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.GetStaticSiteCustomDomain(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state.DomainName = id.CustomDomainName
			state.StaticSiteId = staticsites.NewStaticSiteID(id.SubscriptionId, id.ResourceGroupName, id.StaticSiteName).ID()

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ValidationToken = pointer.From(props.ValidationToken)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r StaticSiteCustomDomainResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteStaticSiteCustomDomainThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StaticSiteCustomDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateCustomDomainID
}

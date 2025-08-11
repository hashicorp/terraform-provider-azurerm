// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/staticsites"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type StaticWebAppCustomDomainResource struct{}

var _ sdk.Resource = StaticWebAppCustomDomainResource{}

type StaticWebAppCustomDomainResourceModel struct {
	DomainName      string `tfschema:"domain_name"`
	StaticSiteId    string `tfschema:"static_web_app_id"`
	ValidationType  string `tfschema:"validation_type"`
	ValidationToken string `tfschema:"validation_token"`
}

func (r StaticWebAppCustomDomainResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"static_web_app_id": {
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
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				helpers.ValidationTypeTXT,
				helpers.ValidationTypeCName,
			}, false),
		},
	}
}

func (r StaticWebAppCustomDomainResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"validation_token": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},
	}
}

func (r StaticWebAppCustomDomainResource) ModelObject() interface{} {
	return &StaticWebAppCustomDomainResource{}
}

func (r StaticWebAppCustomDomainResource) ResourceType() string {
	return "azurerm_static_web_app_custom_domain"
}

func (r StaticWebAppCustomDomainResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			model := StaticWebAppCustomDomainResourceModel{}

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
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
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

			if strings.EqualFold(model.ValidationType, helpers.ValidationTypeCName) {
				if err = client.CreateOrUpdateStaticSiteCustomDomainThenPoll(ctx, id, customDomain); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
			} else {
				if _, err := client.CreateOrUpdateStaticSiteCustomDomain(ctx, id, customDomain); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
				deadline, ok := ctx.Deadline()
				if !ok {
					return fmt.Errorf("internal-error: context was missing a deadline")
				}
				stateConf := &pluginsdk.StateChangeConf{
					Pending: []string{
						string(staticsites.CustomDomainStatusRetrievingValidationToken),
					},
					Target: []string{
						string(staticsites.CustomDomainStatusValidating),
					},
					MinTimeout: 20 * time.Second,
					Timeout:    time.Until(deadline),
					Refresh: func() (interface{}, string, error) {
						domain, err := client.GetStaticSiteCustomDomain(ctx, id)
						if err != nil {
							return domain, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
						}

						if domain.Model == nil || domain.Model.Properties == nil {
							return nil, "Failed", fmt.Errorf("`properties` was missing from the response")
						}
						return domain, string(pointer.From(domain.Model.Properties.Status)), nil
					},
				}

				if _, err := stateConf.WaitForStateContext(ctx); err != nil {
					return fmt.Errorf("waiting for DNS Validation after Creation of %s %+v", id, err)
				}
			}

			// Once validated the token value is zeroed,
			domain, err := client.GetStaticSiteCustomDomain(ctx, id)
			if err != nil {
				return fmt.Errorf("reading validation token for %s: %+v", id, err)
			}
			if m := domain.Model; m != nil {
				if m.Properties != nil {
					if err = metadata.ResourceData.Set("validation_token", m.Properties.ValidationToken); err != nil {
						return fmt.Errorf("setting validation_toekn value for %s: %+v", id, err)
					}
				}
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StaticWebAppCustomDomainResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppService.StaticSitesClient

			id, err := staticsites.ParseCustomDomainID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Some values are not retrievable from the API so we try and load the config.
			state := StaticWebAppCustomDomainResourceModel{}
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

func (r StaticWebAppCustomDomainResource) Delete() sdk.ResourceFunc {
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

func (r StaticWebAppCustomDomainResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return staticsites.ValidateCustomDomainID
}

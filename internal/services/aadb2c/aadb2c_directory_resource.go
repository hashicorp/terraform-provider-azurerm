package aadb2c

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AadB2cDirectoryModel struct {
	BillingType           string            `tfschema:"billing_type"`
	CountryCode           string            `tfschema:"country_code"`
	DataResidencyLocation string            `tfschema:"data_residency_location"`
	DisplayName           string            `tfschema:"display_name"`
	DomainName            string            `tfschema:"domain_name"`
	EffectiveStartDate    string            `tfschema:"effective_start_date"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	Sku                   string            `tfschema:"sku_name"`
	Tags                  map[string]string `tfschema:"tags"`
	TenantId              string            `tfschema:"tenant_id"`
}

type AadB2cDirectoryResource struct{}

var (
	_ sdk.Resource           = AadB2cDirectoryResource{}
	_ sdk.ResourceWithUpdate = AadB2cDirectoryResource{}
)

func (r AadB2cDirectoryResource) ResourceType() string {
	return "azurerm_aadb2c_directory"
}

func (r AadB2cDirectoryResource) ModelObject() interface{} {
	return &AadB2cDirectoryModel{}
}

func (r AadB2cDirectoryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tenants.ValidateB2CDirectoryID
}

func (r AadB2cDirectoryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_name": {
			Description:  "Domain name of the B2C tenant, including onmicrosoft.com suffix.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"country_code": {
			Description:  "Country code of the B2C tenant. See https://aka.ms/B2CDataResidency for valid country codes.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"data_residency_location": {
			Description: "Location in which the B2C tenant is hosted and data resides. See https://aka.ms/B2CDataResidency for more information.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(tenants.LocationAsiaPacific),
				string(tenants.LocationAustralia),
				string(tenants.LocationEurope),
				string(tenants.LocationGlobal),
				string(tenants.LocationUnitedStates),
			}, false),
		},

		"display_name": {
			Description:  "The initial display name of the B2C tenant.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"sku_name": {
			Description: "Billing SKU for the B2C tenant. See https://aka.ms/b2cBilling for more information.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(tenants.SkuNamePremiumP1),
				string(tenants.SkuNamePremiumP2),
				// string(tenants.SkuNameStandard), // API doesn't seem to support "Standard", it's ignored and "PremiumP1" is used instead, even when patching
			}, false),
		},

		"tags": tags.Schema(),
	}
}

func (r AadB2cDirectoryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_type": {
			Description: "The type of billing for the B2C tenant. Possible values include: `MAU` or `Auths`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"effective_start_date": {
			Description: "The date from which the billing type took effect. May not be populated until after the first billing cycle.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tenant_id": {
			Description: "The Tenant ID for the B2C tenant.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r AadB2cDirectoryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AadB2cDirectoryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if model.CountryCode == "" {
				return fmt.Errorf("`country_code` is required when creating a new AADB2C directory")
			}
			if model.DisplayName == "" {
				return fmt.Errorf("`display_name` is required when creating a new AADB2C directory")
			}

			id := tenants.NewB2CDirectoryID(subscriptionId, model.ResourceGroup, model.DomainName)

			metadata.Logger.Infof("Import check for %s", id)
			existing, err := client.Get(ctx, id)
			if err != nil && existing.HttpResponse.StatusCode != http.StatusNotFound {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Domain name availability check for %s", id)
			availabilityResult, err := client.CheckNameAvailability(ctx, commonids.NewSubscriptionID(subscriptionId), tenants.CheckNameAvailabilityRequest{
				Name:        &model.DomainName,
				CountryCode: &model.CountryCode,
			})
			if err != nil {
				return fmt.Errorf("checking availability of `domain_name`: %v", err)
			}

			if availabilityResult.Model.NameAvailable == nil || !*availabilityResult.Model.NameAvailable {
				reason := "unknown reason"
				if availabilityResult.Model.Reason != nil {
					reason = *availabilityResult.Model.Reason
				}
				if availabilityResult.Model.Message != nil {
					reason = fmt.Sprintf("%s (%s)", reason, *availabilityResult.Model.Message)
				}
				return fmt.Errorf("checking availability of `domain_name`: the specified domain %q is unavailable: %s", model.DomainName, reason)
			}

			metadata.Logger.Infof("Creating %s", id)

			properties := tenants.CreateTenant{
				Location: tenants.Location(model.DataResidencyLocation),
				Properties: tenants.TenantPropertiesForCreate{
					CreateTenantProperties: tenants.CreateTenantProperties{
						CountryCode: model.CountryCode,
						DisplayName: model.DisplayName,
					},
				},
				Sku: tenants.Sku{
					Name: tenants.SkuName(model.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, properties); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AadB2cDirectoryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", id)
			var state AadB2cDirectoryModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			metadata.Logger.Infof("Updating %s", id)

			properties := tenants.UpdateTenant{
				Sku: tenants.Sku{
					Name: tenants.SkuName(state.Sku),
					Tier: tenants.SkuTierA0,
				},
				Tags: &state.Tags,
			}

			if _, err := client.Update(ctx, *id, properties); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r AadB2cDirectoryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Reading %s", id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if resp.HttpResponse.StatusCode == http.StatusNotFound {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := AadB2cDirectoryModel{
				DomainName:    id.DirectoryName,
				ResourceGroup: id.ResourceGroup,
				CountryCode:   metadata.ResourceData.Get("country_code").(string),
				DisplayName:   metadata.ResourceData.Get("display_name").(string),
			}

			if model.Location != nil {
				state.DataResidencyLocation = string(*model.Location)
			}

			if model.Sku != nil {
				state.Sku = string(model.Sku.Name)
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			if properties := model.Properties; properties != nil {
				if billingConfig := properties.BillingConfig; billingConfig != nil {
					if billingConfig.BillingType != nil {
						state.BillingType = string(*billingConfig.BillingType)
					}
					if billingConfig.EffectiveStartDateUtc != nil {
						state.EffectiveStartDate = *billingConfig.EffectiveStartDateUtc
					}
				}

				if properties.TenantId != nil {
					state.TenantId = *properties.TenantId
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r AadB2cDirectoryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient

			id, err := tenants.ParseB2CDirectoryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Deleting %s", id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

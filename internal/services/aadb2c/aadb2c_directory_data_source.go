package aadb2c

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview/tenants"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AadB2cDirectoryDataSourceModel struct {
	BillingType           string            `tfschema:"billing_type"`
	DataResidencyLocation string            `tfschema:"data_residency_location"`
	DomainName            string            `tfschema:"domain_name"`
	EffectiveStartDate    string            `tfschema:"effective_start_date"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	Sku                   string            `tfschema:"sku_name"`
	Tags                  map[string]string `tfschema:"tags"`
	TenantId              string            `tfschema:"tenant_id"`
}

type AadB2cDirectoryDataSource struct{}

var _ sdk.DataSource = AadB2cDirectoryDataSource{}

func (r AadB2cDirectoryDataSource) ResourceType() string {
	return "azurerm_aadb2c_directory"
}

func (r AadB2cDirectoryDataSource) ModelObject() interface{} {
	return &AadB2cDirectoryModel{}
}

func (r AadB2cDirectoryDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tenants.ValidateB2CDirectoryID
}

func (r AadB2cDirectoryDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"domain_name": {
			Description:  "Domain name of the B2C tenant, including onmicrosoft.com suffix.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r AadB2cDirectoryDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_type": {
			Description: "The type of billing for the B2C tenant. Possible values include: `MAU` or `Auths`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"data_residency_location": {
			Description: "Location in which the B2C tenant is hosted and data resides.",
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

		"sku_name": {
			Description: "Billing SKU for the B2C tenant.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r AadB2cDirectoryDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AadB2c.TenantsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state AadB2cDirectoryDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := tenants.NewB2CDirectoryID(subscriptionId, state.ResourceGroup, state.DomainName)

			metadata.Logger.Infof("Reading %s", id)
			resp, err := client.Get(ctx, id)
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

			state.DomainName = id.DirectoryName
			state.ResourceGroup = id.ResourceGroup

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

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = BillingAccountDataSource{}

type BillingAccountDataSource struct {
	Name          string                   `tfschema:"name"`
	AccountStatus string                   `tfschema:"account_status"`
	AccountType   string                   `tfschema:"account_type"`
	AgreementType string                   `tfschema:"agreement_type"`
	DisplayName   string                   `tfschema:"display_name"`
	HasReadAccess bool                     `tfschema:"has_read_access"`
	SoldTo        []helpers.AddressDetails `tfschema:"sold_to"`
}

func (BillingAccountDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (BillingAccountDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"account_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"agreement_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"has_read_access": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
		"sold_to": helpers.AddressDetailsSchema(),
	}
}

func (BillingAccountDataSource) ModelObject() interface{} {
	return nil
}

func (BillingAccountDataSource) ResourceType() string {
	return "azurerm_billing_account"
}

func (BillingAccountDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingAccountsClient

			var billingAccount BillingAccountDataSource
			if err := metadata.Decode(&billingAccount); err != nil {
				return err
			}

			name := metadata.ResourceData.Get("name").(string)
			id := billingaccounts.NewBillingAccountID(name)

			resp, err := client.Get(ctx, id, billingaccounts.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := BillingAccountDataSource{
				Name: id.BillingAccountName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AccountStatus = string(pointer.From(props.AccountStatus))
					state.AccountType = string(pointer.From(props.AccountType))
					state.AgreementType = string(pointer.From(props.AgreementType))
					state.DisplayName = string(pointer.From(props.DisplayName))
					state.HasReadAccess = bool(pointer.From(props.HasReadAccess))
					state.SoldTo = helpers.FlattenAddressDetails(props.SoldTo)
				}
			}

			metadata.SetID(id)

			if err := metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return nil
		},
	}
}

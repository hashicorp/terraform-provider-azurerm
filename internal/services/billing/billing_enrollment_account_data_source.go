// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2019-10-01-preview/enrollmentaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = BillingEnrollmentAccountDataSource{}

type BillingEnrollmentAccountDataSource struct {
	BillingAccountName    string `tfschema:"billing_account_name"`
	EnrollmentAccountName string `tfschema:"enrollment_account_name"`
	AccountName           string `tfschema:"account_name"`
	AccountOwner          string `tfschema:"account_owner"`
	CostCenter            string `tfschema:"cost_center"`
	EndDate               string `tfschema:"end_date"`
	StartDate             string `tfschema:"start_date"`
	Status                string `tfschema:"status"`
}

func (BillingEnrollmentAccountDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"billing_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"enrollment_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (BillingEnrollmentAccountDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"account_owner": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"cost_center": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"end_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"start_date": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (BillingEnrollmentAccountDataSource) ModelObject() interface{} {
	return nil
}

func (BillingEnrollmentAccountDataSource) ResourceType() string {
	return "azurerm_billing_enrollment_account"
}

func (BillingEnrollmentAccountDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.BillingEnrollmentAccountsClient

			var billingEnrollmentAccount BillingEnrollmentAccountDataSource
			if err := metadata.Decode(&billingEnrollmentAccount); err != nil {
				return err
			}

			billingAccountName := metadata.ResourceData.Get("billing_account_name").(string)
			enrollmentAccountName := metadata.ResourceData.Get("enrollment_account_name").(string)
			id := enrollmentaccounts.NewEnrollmentAccountID(billingAccountName, enrollmentAccountName)

			resp, err := client.GetByEnrollmentAccountId(ctx, id, enrollmentaccounts.DefaultGetByEnrollmentAccountIdOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := BillingEnrollmentAccountDataSource{
				EnrollmentAccountName: id.EnrollmentAccountName,
				BillingAccountName:    id.BillingAccountName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.AccountName = pointer.From(props.AccountName)
					state.AccountOwner = pointer.From(props.AccountOwner)
					state.CostCenter = pointer.From(props.CostCenter)
					state.EndDate = pointer.From(props.EndDate)
					state.StartDate = pointer.From(props.StartDate)
					state.Status = pointer.From(props.Status)
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

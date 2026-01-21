// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package codesigning

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/codesigning/2025-10-13/codesigningaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ArtifactsSigningAccountDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	AccountUri        string            `tfschema:"account_uri"`
	SkuName           string            `tfschema:"sku_name"`
	Tags              map[string]string `tfschema:"tags"`
}

type ArtifactsSigningAccountDataSource struct{}

var _ sdk.DataSource = ArtifactsSigningAccountDataSource{}

func (d ArtifactsSigningAccountDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 24),
				validation.StringMatch(
					regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*(?:-[A-Za-z0-9]+)*$"),
					"An account's name must be between 3-24 alphanumeric characters. The name must begin with a letter, end with a letter or digit, and not contain consecutive hyphens.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (d ArtifactsSigningAccountDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.LocationComputed(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (d ArtifactsSigningAccountDataSource) ModelObject() interface{} {
	return &ArtifactsSigningAccountDataSourceModel{}
}

func (d ArtifactsSigningAccountDataSource) ResourceType() string {
	return "azurerm_artifacts_signing_account"
}

func (d ArtifactsSigningAccountDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CodeSigning.Client.CodeSigningAccounts

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ArtifactsSigningAccountDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := codesigningaccounts.NewCodeSigningAccountID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state.Name = id.CodeSigningAccountName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					state.AccountUri = pointer.From(prop.AccountUri)
					if sku := prop.Sku; sku != nil {
						state.SkuName = string(sku.Name)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

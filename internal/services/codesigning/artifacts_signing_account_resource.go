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

type ArtifactSigningAccountModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	AccountUri        string            `tfschema:"account_uri"`
	SkuName           string            `tfschema:"sku_name"`
	Tags              map[string]string `tfschema:"tags"`
}

type ArtifactSigningAccountResource struct{}

var _ sdk.Resource = ArtifactSigningAccountResource{}

func (r ArtifactSigningAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 24),
				validation.StringMatch(
					regexp.MustCompile("^[A-Za-z][A-Za-z0-9]*(?:-[A-Za-z0-9]+)*$"),
					"An account's name must be between 3-24 alphanumeric characters. The name must begin with a letter, end with a letter or digit, and not contain consecutive hyphens.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice(
				codesigningaccounts.PossibleValuesForSkuName(),
				false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ArtifactSigningAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ArtifactSigningAccountResource) ModelObject() interface{} {
	return &ArtifactSigningAccountModel{}
}

func (r ArtifactSigningAccountResource) ResourceType() string {
	return "azurerm_artifact_signing_account"
}

func (r ArtifactSigningAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts

			var model ArtifactSigningAccountModel
			if err := meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionID := meta.Client.Account.SubscriptionId
			id := codesigningaccounts.NewCodeSigningAccountID(subscriptionID, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if !response.WasNotFound(existing.HttpResponse) {
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", id, err)
				}
				return meta.ResourceRequiresImport(r.ResourceType(), id)
			}

			req := codesigningaccounts.CodeSigningAccount{
				Name:     &model.Name,
				Location: location.Normalize(model.Location),
				Tags:     &model.Tags,
				Properties: &codesigningaccounts.CodeSigningAccountProperties{
					Sku: &codesigningaccounts.AccountSku{
						Name: codesigningaccounts.SkuName(model.SkuName),
					},
				},
			}

			err = client.CreateThenPoll(ctx, id, req)
			if err != nil {
				return fmt.Errorf("creating %s: %v", id, err)
			}

			meta.SetID(id)
			return nil
		},
	}
}

func (r ArtifactSigningAccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts

			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			result, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(result.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			output := ArtifactSigningAccountModel{
				Name:              id.CodeSigningAccountName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := result.Model; model != nil {
				output.Location = location.Normalize(model.Location)
				output.Tags = pointer.From(model.Tags)

				if prop := model.Properties; prop != nil {
					output.AccountUri = pointer.From(prop.AccountUri)
					if sku := prop.Sku; sku != nil {
						output.SkuName = string(sku.Name)
					}
				}
			}

			return meta.Encode(&output)
		},
	}
}

func (r ArtifactSigningAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts
			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}
			var model ArtifactSigningAccountModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var patch codesigningaccounts.CodeSigningAccountPatch
			if meta.ResourceData.HasChange("sku_name") {
				patch.Properties = pointer.To(codesigningaccounts.CodeSigningAccountPatchProperties{
					Sku: pointer.To(codesigningaccounts.AccountSkuPatch{
						Name: pointer.ToEnum[codesigningaccounts.SkuName](model.SkuName),
					}),
				})
			}
			if meta.ResourceData.HasChange("tags") {
				patch.Tags = pointer.To(model.Tags)
			}

			if err = client.UpdateThenPoll(ctx, *id, patch); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

func (r ArtifactSigningAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.CodeSigning.Client.CodeSigningAccounts

			id, err := codesigningaccounts.ParseCodeSigningAccountID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			meta.Logger.Infof("deleting %s", id)

			if _, err = client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %v", id, err)
			}
			return nil
		},
	}
}

func (r ArtifactSigningAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return codesigningaccounts.ValidateCodeSigningAccountID
}

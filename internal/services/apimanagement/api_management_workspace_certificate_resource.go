// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/certificate"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementWorkspaceCertificateModel struct {
	Name                         string `tfschema:"name"`
	ApiManagementWorkspaceId     string `tfschema:"api_management_workspace_id"`
	CertificateDataBase64        string `tfschema:"certificate_data_base64"`
	KeyVaultSecretId             string `tfschema:"key_vault_secret_id"`
	Password                     string `tfschema:"password"`
	UserAssignedIdentityClientId string `tfschema:"user_assigned_identity_client_id"`

	// Computed fields
	Expiration string `tfschema:"expiration"`
	Subject    string `tfschema:"subject"`
	Thumbprint string `tfschema:"thumbprint"`
}

type ApiManagementWorkspaceCertificateResource struct{}

var _ sdk.ResourceWithUpdate = ApiManagementWorkspaceCertificateResource{}

func (r ApiManagementWorkspaceCertificateResource) ResourceType() string {
	return "azurerm_api_management_workspace_certificate"
}

func (r ApiManagementWorkspaceCertificateResource) ModelObject() interface{} {
	return &ApiManagementWorkspaceCertificateModel{}
}

func (r ApiManagementWorkspaceCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificate.ValidateWorkspaceCertificateID
}

func (r ApiManagementWorkspaceCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,78}[a-zA-Z0-9])?$`),
				"The `name` must be 1–80 characters, using only letters, numbers, or hyphens, and not starting or ending with a hyphen."),
		},

		"api_management_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&workspace.WorkspaceId{}),

		"certificate_data_base64": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ExactlyOneOf:  []string{"certificate_data_base64", "key_vault_secret_id"},
			ConflictsWith: []string{"key_vault_secret_id", "user_assigned_identity_client_id"},
			ValidateFunc:  validation.StringIsBase64,
		},

		"key_vault_secret_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  validate.NestedItemIdWithOptionalVersion,
			ExactlyOneOf:  []string{"certificate_data_base64", "key_vault_secret_id"},
			ConflictsWith: []string{"certificate_data_base64", "password"},
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			RequiredWith: []string{"certificate_data_base64"},
		},

		"user_assigned_identity_client_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
			RequiredWith: []string{"key_vault_secret_id"},
		},
	}
}

func (r ApiManagementWorkspaceCertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"expiration": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subject": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"thumbprint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ApiManagementWorkspaceCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			var model ApiManagementWorkspaceCertificateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspace.ParseWorkspaceID(model.ApiManagementWorkspaceId)
			if err != nil {
				return err
			}

			id := certificate.NewWorkspaceCertificateID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.ServiceName, workspaceId.WorkspaceId, model.Name)
			existing, err := client.WorkspaceCertificateGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := certificate.CertificateCreateOrUpdateParameters{
				Properties: &certificate.CertificateCreateOrUpdateProperties{},
			}

			if model.KeyVaultSecretId != "" {
				parsedSecretId, err := parse.ParseOptionallyVersionedNestedItemID(model.KeyVaultSecretId)
				if err != nil {
					return err
				}

				parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{
					SecretIdentifier: pointer.To(parsedSecretId.ID()),
				}

				if model.UserAssignedIdentityClientId != "" {
					parameters.Properties.KeyVault.IdentityClientId = pointer.To(model.UserAssignedIdentityClientId)
				}
			}

			if model.CertificateDataBase64 != "" {
				parameters.Properties.Data = pointer.To(model.CertificateDataBase64)
				parameters.Properties.Password = pointer.To(model.Password)
			}

			if _, err := client.WorkspaceCertificateCreateOrUpdate(ctx, id, parameters, certificate.DefaultWorkspaceCertificateCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ApiManagementWorkspaceCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			id, err := certificate.ParseWorkspaceCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkspaceCertificateGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := ApiManagementWorkspaceCertificateModel{
				Name:                     id.CertificateId,
				ApiManagementWorkspaceId: workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.WorkspaceId).ID(),
			}

			if respModel := resp.Model; respModel != nil {
				if props := respModel.Properties; props != nil {
					model.Expiration = props.ExpirationDate
					model.Subject = props.Subject
					model.Thumbprint = props.Thumbprint
					// The API omitted `certificate_data_base64` from the response as it is considered sensitive.
					model.CertificateDataBase64 = metadata.ResourceData.Get("certificate_data_base64").(string)
					// The API omitted `password` from the response as it is considered sensitive.
					model.Password = metadata.ResourceData.Get("password").(string)

					if kv := props.KeyVault; kv != nil {
						model.KeyVaultSecretId = pointer.From(kv.SecretIdentifier)
						model.UserAssignedIdentityClientId = pointer.From(kv.IdentityClientId)
					}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ApiManagementWorkspaceCertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			var model ApiManagementWorkspaceCertificateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := certificate.ParseWorkspaceCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.WorkspaceCertificateGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := resp.Model
			parameters := certificate.CertificateCreateOrUpdateParameters{
				Properties: &certificate.CertificateCreateOrUpdateProperties{},
			}

			if kv := payload.Properties.KeyVault; kv != nil {
				parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{}
				if payload.Properties.KeyVault.SecretIdentifier != nil {
					parameters.Properties.KeyVault.SecretIdentifier = kv.SecretIdentifier
				}
				if payload.Properties.KeyVault.IdentityClientId != nil {
					parameters.Properties.KeyVault.IdentityClientId = kv.IdentityClientId
				}
			}

			if metadata.ResourceData.HasChange("user_assigned_identity_client_id") {
				if model.UserAssignedIdentityClientId != "" {
					if parameters.Properties.KeyVault == nil {
						parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{}
					}
					parameters.Properties.KeyVault.IdentityClientId = pointer.To(model.UserAssignedIdentityClientId)
				} else if parameters.Properties.KeyVault != nil {
					parameters.Properties.KeyVault.IdentityClientId = nil
				}
			}

			if metadata.ResourceData.HasChange("key_vault_secret_id") {
				if model.KeyVaultSecretId != "" {
					parsedSecretId, err := parse.ParseOptionallyVersionedNestedItemID(model.KeyVaultSecretId)
					if err != nil {
						return err
					}
					if parameters.Properties.KeyVault == nil {
						parameters.Properties.KeyVault = &certificate.KeyVaultContractCreateProperties{}
					}
					parameters.Properties.KeyVault.SecretIdentifier = pointer.To(parsedSecretId.ID())
				} else {
					parameters.Properties.KeyVault = nil
				}
			}

			if metadata.ResourceData.HasChange("certificate_data_base64") {
				if model.CertificateDataBase64 != "" {
					parameters.Properties.Data = pointer.To(model.CertificateDataBase64)
				} else {
					parameters.Properties.Data = nil
				}
			}

			if metadata.ResourceData.HasChange("password") {
				if model.Password != "" {
					parameters.Properties.Password = pointer.To(model.Password)
				} else {
					parameters.Properties.Password = nil
				}
			}

			if _, err := client.WorkspaceCertificateCreateOrUpdate(ctx, *id, parameters, certificate.DefaultWorkspaceCertificateCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementWorkspaceCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.CertificateClient_v2024_05_01

			id, err := certificate.ParseWorkspaceCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkspaceCertificateDelete(ctx, *id, certificate.DefaultWorkspaceCertificateDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

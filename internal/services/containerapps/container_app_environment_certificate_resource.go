// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentCertificateResource struct{}

type CertificateKeyVaultModel struct {
	Identity         string `tfschema:"identity"`
	KeyVaultSecretId string `tfschema:"key_vault_secret_id"`
}

type ContainerAppCertificateModel struct {
	Name                 string                 `tfschema:"name"`
	ManagedEnvironmentId string                 `tfschema:"container_app_environment_id"`
	Tags                 map[string]interface{} `tfschema:"tags"`

	// Write only?
	CertificatePassword string `tfschema:"certificate_password"`
	CertificateBlob     string `tfschema:"certificate_blob_base64"`

	CertificateKeyVault []CertificateKeyVaultModel `tfschema:"certificate_key_vault"`

	// Read Only
	SubjectName    string `tfschema:"subject_name"`
	Issuer         string `tfschema:"issuer"`
	IssueDate      string `tfschema:"issue_date"`
	ExpirationDate string `tfschema:"expiration_date"`
	Thumbprint     string `tfschema:"thumbprint"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentCertificateResource{}

func (r ContainerAppEnvironmentCertificateResource) ModelObject() interface{} {
	return &ContainerAppCertificateModel{}
}

func (r ContainerAppEnvironmentCertificateResource) ResourceType() string {
	return "azurerm_container_app_environment_certificate"
}

func (r ContainerAppEnvironmentCertificateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return certificates.ValidateCertificateID
}

func (r ContainerAppEnvironmentCertificateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
			Description:  "The name of the Container Apps Environment Certificate.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Certificate on.",
		},

		"certificate_blob_base64": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsBase64,
			Description:   "The Certificate Private Key as a base64 encoded PFX or PEM.",
			ExactlyOneOf: []string{"certificate_key_vault", "certificate_blob_base64"},
			RequiredWith: []string{"certificate_password"},
		},

		"certificate_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			Description:  "The password for the Certificate.",
			RequiredWith: []string{"certificate_blob_base64"},
		},

		"certificate_key_vault": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			MaxItems:      1,
			Description:   "Import Certificate from Key Vault",
			ExactlyOneOf:  []string{"certificate_key_vault", "certificate_blob_base64"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"identity": {
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "The Resource ID of a managed identity to authenticate with Azure Key Vault, or 'System' to use a system-assigned identity.",
						ValidateFunc: validation.Any(
							validation.StringInSlice([]string{"System"}, false),
							commonids.ValidateUserAssignedIdentityID,
						),
					},
					"key_vault_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						Description:  "The Base ID of the Key Vault Secret containing the certificate.",
						ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r ContainerAppEnvironmentCertificateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subject_name": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Subject Name for the Certificate.",
		},

		"issuer": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Certificate Issuer.",
		},

		"issue_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The date of issue for the Certificate.",
		},

		"expiration_date": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The expiration date for the Certificate.",
		},

		"thumbprint": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Thumbprint of the Certificate.",
		},
	}
}

func (r ContainerAppEnvironmentCertificateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.CertificatesClient
			environmentsClient := metadata.Client.ContainerApps.ManagedEnvironmentClient

			var cert ContainerAppCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			envId, err := managedenvironments.ParseManagedEnvironmentID(cert.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := certificates.NewCertificateID(metadata.Client.Account.SubscriptionId, envId.ResourceGroupName, envId.ManagedEnvironmentName, cert.Name)
			env, err := environmentsClient.Get(ctx, *envId)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", *envId, id, err)
			}

			model := certificates.Certificate{
				Location:   env.Model.Location,
				Name:       pointer.To(id.CertificateName),
				Properties: &certificates.CertificateProperties{},
				Tags:       tags.Expand(cert.Tags),
			}

			hasBlobAuth := cert.CertificateBlob != "" && cert.CertificatePassword != ""
			if hasBlobAuth {
				model.Properties.Password = pointer.To(cert.CertificatePassword)
				model.Properties.Value = pointer.To(cert.CertificateBlob)
			} else if len(cert.CertificateKeyVault) > 0 {
				kvConfig := cert.CertificateKeyVault[0]
				model.Properties.CertificateKeyVaultProperties = &certificates.CertificateKeyVaultProperties{
					Identity:    pointer.To(kvConfig.Identity),
					KeyVaultURL: pointer.To(kvConfig.KeyVaultSecretId),
				}
			}

			if _, err := client.CreateOrUpdate(ctx, id, model); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppEnvironmentCertificateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.CertificatesClient

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppCertificateModel

			state.Name = id.CertificateName
			state.ManagedEnvironmentId = certificates.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName).ID()

			if model := existing.Model; model != nil {
				state.Tags = tags.Flatten(model.Tags)

				// The Certificate Blob and Password are not retrievable in any way, so grab them back from config if we can. Imports will need `ignore_changes`.
				if certBlob, ok := metadata.ResourceData.GetOk("certificate_blob_base64"); ok {
					state.CertificateBlob = certBlob.(string)
				}

				if certPassword, ok := metadata.ResourceData.GetOk("certificate_password"); ok {
					state.CertificatePassword = certPassword.(string)
				}

				if props := model.Properties; props != nil {
					state.SubjectName = pointer.From(props.SubjectName)
					state.Issuer = pointer.From(props.Issuer)
					state.IssueDate = pointer.From(props.IssueDate)
					state.ExpirationDate = pointer.From(props.ExpirationDate)
					state.Thumbprint = pointer.From(props.Thumbprint)

					if kvProps := props.CertificateKeyVaultProperties; kvProps != nil {
						keyVaultConfig := CertificateKeyVaultModel{
							Identity:         pointer.From(kvProps.Identity),
							KeyVaultSecretId: pointer.From(kvProps.KeyVaultURL),
						}
						state.CertificateKeyVault = []CertificateKeyVaultModel{keyVaultConfig}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentCertificateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.CertificatesClient

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppEnvironmentCertificateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.CertificatesClient

			var cert ContainerAppCertificateModel

			if err := metadata.Decode(&cert); err != nil {
				return err
			}

			id, err := certificates.ParseCertificateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("tags") {
				patch := certificates.CertificatePatch{
					Tags: tags.Expand(cert.Tags),
				}

				if _, err = client.Update(ctx, *id, patch); err != nil {
					return fmt.Errorf("updating tags for %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}

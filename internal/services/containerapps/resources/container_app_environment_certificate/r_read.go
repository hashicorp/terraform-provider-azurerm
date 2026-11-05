// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
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
						state.CertificateKeyVault = []CertificateKeyVaultModel{
							{
								Identity:         pointer.From(kvProps.Identity),
								KeyVaultSecretId: pointer.From(kvProps.KeyVaultURL),
							},
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

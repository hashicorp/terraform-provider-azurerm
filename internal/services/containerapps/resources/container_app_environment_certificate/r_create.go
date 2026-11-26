// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if cert.CertificateBlob != "" {
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

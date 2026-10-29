// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
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

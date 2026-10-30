// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_certificate

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxcertificate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.Nginx.NginxCertificate
			id, err := nginxcertificate.ParseCertificateID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CertificateModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			// retrieve from GET
			existing, err := client.CertificatesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists when updating: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil when updating for %v", *id)
			}

			// have to pass all existing properties to update
			upd := existing.Model
			if meta.ResourceData.HasChange("key_virtual_path") {
				upd.Properties.KeyVirtualPath = pointer.To(model.KeyVirtualPath)
			}

			if meta.ResourceData.HasChange("certificate_virtual_path") {
				upd.Properties.CertificateVirtualPath = pointer.To(model.CertificateVirtualPath)
			}

			if meta.ResourceData.HasChange("key_vault_secret_id") {
				upd.Properties.KeyVaultSecretId = pointer.To(model.KeyVaultSecretId)
			}

			if err = client.CertificatesCreateOrUpdateThenPoll(ctx, *id, *upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}
			return nil
		},
	}
}

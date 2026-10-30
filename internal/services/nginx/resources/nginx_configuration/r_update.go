// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_configuration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-11-01-preview/nginxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Nginx.NginxConfiguration
			id, err := nginxconfiguration.ParseConfigurationID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ConfigurationModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			// retrieve from GET
			existing, err := client.ConfigurationsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving exists: +%v", *id)
			}
			if existing.Model == nil && existing.Model.Properties == nil {
				return fmt.Errorf("retrieving as nil for %v", *id)
			}

			// full update - fill in the existing fields from the API and then patch it
			upd := nginxconfiguration.NginxConfigurationRequest{
				Name: pointer.To(defaultConfigurationName),
				Properties: &nginxconfiguration.NginxConfigurationRequestProperties{
					RootFile: existing.Model.Properties.RootFile,
					Files:    existing.Model.Properties.Files,
					Package:  existing.Model.Properties.Package,
				},
			}

			if existing.Model.Properties.ProtectedFiles != nil {
				var pfs []nginxconfiguration.NginxConfigurationProtectedFileRequest
				for _, f := range *existing.Model.Properties.ProtectedFiles {
					pfs = append(pfs, nginxconfiguration.NginxConfigurationProtectedFileRequest{
						VirtualPath: f.VirtualPath,
					})
				}
				upd.Properties.ProtectedFiles = pointer.To(pfs)
			}

			if meta.ResourceData.HasChange("root_file") {
				upd.Properties.RootFile = pointer.To(model.RootFile)
			}

			if meta.ResourceData.HasChange("config_file") {
				upd.Properties.Files = model.toSDKFiles()
			}

			if meta.ResourceData.HasChange("protected_file") {
				upd.Properties.ProtectedFiles = model.toSDKProtectedFiles()
			}

			if meta.ResourceData.HasChange("package_data") {
				upd.Properties.Package = &nginxconfiguration.NginxConfigurationPackage{
					Data: pointer.To(model.PackageData),
				}
			}

			if err := client.ConfigurationsCreateOrUpdateThenPoll(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

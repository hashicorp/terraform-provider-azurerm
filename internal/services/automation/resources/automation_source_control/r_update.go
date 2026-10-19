// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_source_control

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/sourcecontrol"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.SourceControl

			id, err := sourcecontrol.ParseSourceControlID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SourceControlModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding err: %+v", err)
			}

			var upd sourcecontrol.SourceControlUpdateParameters
			prop := &sourcecontrol.SourceControlUpdateProperties{}
			if meta.ResourceData.HasChange("branch") {
				prop.Branch = pointer.To(model.Branch)
			}
			if meta.ResourceData.HasChange("folder_path") {
				prop.FolderPath = pointer.To(model.FolderPath)
			}
			if meta.ResourceData.HasChange("automatic_sync") {
				prop.AutoSync = pointer.To(model.AutoSync)
			}
			if meta.ResourceData.HasChange("folder_path") {
				prop.FolderPath = pointer.To(model.FolderPath)
			}
			if meta.ResourceData.HasChange("publish_runbook_enabled") {
				prop.PublishRunbook = pointer.To(model.PublishRunbook)
			}
			if meta.ResourceData.HasChange("description") {
				prop.Description = pointer.To(model.Description)
			}

			if meta.ResourceData.HasChange("security") {
				prop.SecurityToken = &sourcecontrol.SourceControlSecurityTokenProperties{
					AccessToken:  pointer.To(model.SecurityToken[0].TokenType),
					RefreshToken: pointer.To(model.SecurityToken[0].RefreshToken),
					TokenType:    pointer.ToEnum[sourcecontrol.TokenType](model.SecurityToken[0].TokenType),
				}
			}
			upd.Properties = prop
			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", *id, err)
			}

			return nil
		},
	}
}

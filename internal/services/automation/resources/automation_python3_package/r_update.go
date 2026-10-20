// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_python3_package

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/python3package"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (m Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) (err error) {
			client := meta.Client.Automation.Python3Package
			id, err := python3package.ParsePython3PackageID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			var model Python3PackageModel
			if err = meta.Decode(&model); err != nil {
				return fmt.Errorf("decoding %s: %+v", id, err)
			}

			var upd python3package.PythonPackageUpdateParameters
			if meta.ResourceData.HasChange("tags") {
				upd.Tags = &model.Tags
			}

			if _, err = client.Update(ctx, *id, upd); err != nil {
				return fmt.Errorf("updating %s: %v", id, err)
			}

			return nil
		},
	}
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func runTFVersionChecks(ctx context.Context, t testing.T, terraformVersion *version.Version, terraformVersionChecks []tfversion.TerraformVersionCheck) {
	t.Helper()

	for _, tfVersionCheck := range terraformVersionChecks {
		resp := tfversion.CheckTerraformVersionResponse{}
		tfVersionCheck.CheckTerraformVersion(ctx, tfversion.CheckTerraformVersionRequest{TerraformVersion: terraformVersion}, &resp)

		if resp.Error != nil {
			t.Fatalf(resp.Error.Error())
		}

		if resp.Skip != "" {
			t.Skip(resp.Skip)
		}
	}

}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto5server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// UpgradeResourceIdentity satisfies the tfprotov5.ProviderServer interface.
func (s *Server) UpgradeResourceIdentity(ctx context.Context, proto5Req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	panic("unimplemented") // TODO:ResourceIdentity: implement
}

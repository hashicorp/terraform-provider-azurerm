// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package proto6server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// UpgradeResourceIdentity satisfies the tfprotov6.ProviderServer interface.
func (s *Server) UpgradeResourceIdentity(ctx context.Context, proto6Req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	panic("unimplemented") // TODO:ResourceIdentity: implement
}

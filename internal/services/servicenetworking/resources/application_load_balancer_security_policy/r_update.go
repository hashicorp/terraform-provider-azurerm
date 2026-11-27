// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_security_policy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	servicenetworkingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
)

func (t Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := servicenetworkingClient.From(metadata.Client).SecurityPoliciesInterface

			id, err := securitypoliciesinterface.ParseSecurityPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config SecurityPoliciesModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			payload := securitypoliciesinterface.SecurityPolicyUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if _, err := client.Update(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

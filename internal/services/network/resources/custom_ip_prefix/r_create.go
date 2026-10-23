// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package custom_ip_prefix

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 9 * time.Hour,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			r.client = metadata.Client.Network.CustomIPPrefixes
			subscriptionId := metadata.Client.Account.SubscriptionId

			deadline, ok := ctx.Deadline()
			if !ok {
				return fmt.Errorf("internal-error: context has no deadline")
			}

			var model CustomIpPrefixModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := customipprefixes.NewCustomIPPrefixID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := r.client.Get(ctx, id, customipprefixes.DefaultGetOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			_, cidr, err := net.ParseCIDR(model.CIDR)
			if err != nil {
				return fmt.Errorf("parsing `cidr`: %+v", err)
			}

			results, err := r.client.ListAll(ctx, commonids.NewSubscriptionID(subscriptionId))
			if err != nil {
				return fmt.Errorf("listing existing %s: %+v", id, err)
			}

			if prefixes := results.Model; prefixes != nil {
				for _, prefix := range *prefixes {
					if prefix.Properties != nil && prefix.Properties.Cidr != nil {
						_, netw, err := net.ParseCIDR(*prefix.Properties.Cidr)
						if err != nil {
							// couldn't parse the existing custom prefix, so skip it
							continue
						}
						if cidr == netw {
							return metadata.ResourceRequiresImport(r.ResourceType(), id)
						}
					}
				}
			}

			payload := customipprefixes.CustomIPPrefix{
				Name:     &model.Name,
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     tags.Expand(model.Tags),
				Properties: &customipprefixes.CustomIPPrefixPropertiesFormat{
					Cidr:              &model.CIDR,
					CommissionedState: pointer.To(customipprefixes.CommissionedStateProvisioning),
				},
			}

			if model.ParentCustomIPPrefixID != "" {
				payload.Properties.CustomIPPrefixParent = &customipprefixes.SubResource{
					Id: &model.ParentCustomIPPrefixID,
				}
			}

			if model.WANValidationSignedMessage != "" {
				payload.Properties.SignedMessage = &model.WANValidationSignedMessage
			}

			if model.ROAValidityEndDate != "" {
				roaValidityEndDate, err := time.Parse("2006-01-02", model.ROAValidityEndDate)
				if err != nil {
					return err
				}
				authorizationMessage := fmt.Sprintf("%s|%s|%s", subscriptionId, model.CIDR, roaValidityEndDate.Format("20060102"))
				payload.Properties.AuthorizationMessage = &authorizationMessage
			}

			if len(model.Zones) > 0 {
				payload.Zones = &model.Zones
			}

			if err := r.client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{string(customipprefixes.ProvisioningStateUpdating)},
				Target:     []string{string(customipprefixes.ProvisioningStateSucceeded)},
				Refresh:    r.provisioningStateRefreshFunc(ctx, id),
				MinTimeout: 2 * time.Minute,
				Timeout:    time.Until(deadline),
			}
			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for ProvisioningState of %s: %+v", id, err)
			}

			desiredState := customipprefixes.CommissionedStateProvisioned
			if model.CommissioningEnabled {
				if model.InternetAdvertisingDisabled {
					desiredState = customipprefixes.CommissionedStateCommissionedNoInternetAdvertise
				} else {
					desiredState = customipprefixes.CommissionedStateCommissioned
				}
			}

			commissionedState, err := r.updateCommissionedState(ctx, id, desiredState)
			if err != nil {
				return err
			}
			if commissionedState == nil {
				return fmt.Errorf("waiting for CommissionedState: final commissionedState was nil")
			}

			log.Printf("[DEBUG] Final CommissionedState is %q for %s..", *commissionedState, id)
			return nil
		},
	}
}

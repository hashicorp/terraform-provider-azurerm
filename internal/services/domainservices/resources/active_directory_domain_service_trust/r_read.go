// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service_trust

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient
			id, err := parse.DomainServiceTrustID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", idsdk, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}

			props := model.Properties
			if props == nil {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing properties", idsdk)
			}

			existingTrusts := []domainservices.ForestTrust{}
			if fsettings := props.ResourceForestSettings; fsettings != nil {
				if settings := fsettings.Settings; settings != nil {
					existingTrusts = *settings
				}
			}
			var trust *domainservices.ForestTrust
			for _, setting := range existingTrusts {
				if setting.FriendlyName != nil && *setting.FriendlyName == id.TrustName {
					trust = pointer.To(setting)
				}
			}
			if trust == nil {
				return metadata.MarkAsGone(id)
			}

			// Retrieve the initial replica set id to construct the domain service id.
			replicaSets := helpers.FlattenDomainServiceReplicaSets(props.ReplicaSets)
			if len(replicaSets) == 0 {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing replica set details", idsdk)
			}
			initialReplicaSetId := replicaSets[0].(map[string]any)["id"].(string)
			dsid := parse.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, initialReplicaSetId)

			var state DomainServiceTrustModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			data := DomainServiceTrustModel{
				DomainServiceId: dsid.ID(),
				Name:            id.TrustName,
				// Setting the password from state as it is not returned by API.
				Password: state.Password,
			}

			if trust.TrustedDomainFqdn != nil {
				data.TrustedDomainFqdn = *trust.TrustedDomainFqdn
			}

			if trust.RemoteDnsIPs != nil {
				data.TrustedDomainDnsIPs = strings.Split(*trust.RemoteDnsIPs, ",")
			}

			return metadata.Encode(&data)
		},
	}
}

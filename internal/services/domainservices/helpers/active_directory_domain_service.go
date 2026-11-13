// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func DomainServiceControllerRefreshFunc(ctx context.Context, client *domainservices.DomainServicesClient, id parse.DomainServiceId, deleting bool) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Waiting for domain controllers to deploy...")
		idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.Name)
		resp, err := client.Get(ctx, idsdk)
		if err != nil {
			return nil, "error", err
		}
		if model := resp.Model; model == nil || model.Properties == nil || model.Properties.ReplicaSets == nil || len(*model.Properties.ReplicaSets) == 0 {
			return nil, "error", fmt.Errorf("API error: `replicaSets` was not returned")
		}
		// Loop through all replica sets and ensure they are running and each have two available domain controllers
		for _, repl := range *resp.Model.Properties.ReplicaSets {
			if repl.ServiceStatus == nil {
				return resp, "pending", nil
			}
			switch {
			case !deleting && strings.EqualFold(*repl.ServiceStatus, "TearingDown"):
				// Sometimes a service error will cause the replica set, or resource, to self destruct
				return resp, "error", fmt.Errorf("service error: a replica set is unexpectedly tearing down")
			case strings.EqualFold(*repl.ServiceStatus, "Failed"):
				// If a replica set enters a failed state, it needs manual intervention
				return resp, "error", fmt.Errorf("service error: a replica set has entered a Failed state and must be recovered or deleted manually")
			case !strings.EqualFold(*repl.ServiceStatus, "Running"):
				// If it's not yet running, it isn't ready
				return resp, "pending", nil
			case repl.DomainControllerIPAddress == nil || len(*repl.DomainControllerIPAddress) < 2:
				// When a domain controller is online, its IP address will be returned. We're looking for 2 active domain controllers.
				return resp, "pending", nil
			}
		}
		return resp, "available", nil
	}
}

func FlattenDomainServiceReplicaSets(input *[]domainservices.ReplicaSet) (ret []any) {
	if input == nil {
		return []any{}
	}

	for _, in := range *input {
		repl := map[string]any{
			"domain_controller_ip_addresses": make([]string, 0),
			"external_access_ip_address":     "",
			"location":                       location.NormalizeNilable(in.Location),
			"id":                             "",
			"service_status":                 "",
			"subnet_id":                      "",
		}
		if in.DomainControllerIPAddress != nil {
			repl["domain_controller_ip_addresses"] = *in.DomainControllerIPAddress
		}
		if in.ExternalAccessIPAddress != nil {
			repl["external_access_ip_address"] = *in.ExternalAccessIPAddress
		}
		if in.ReplicaSetId != nil {
			repl["id"] = *in.ReplicaSetId
		}
		if in.ServiceStatus != nil {
			repl["service_status"] = *in.ServiceStatus
		}
		if in.SubnetId != nil {
			repl["subnet_id"] = *in.SubnetId
		}
		ret = append(ret, repl)
	}

	return
}

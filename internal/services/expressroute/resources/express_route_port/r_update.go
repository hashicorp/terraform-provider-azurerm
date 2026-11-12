// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmExpressRoutePortUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRoutePorts
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressrouteports.ParseExpressRoutePortID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := existing.Model

	// Normalize the identity type - Azure API returns lowercase (e.g. "userAssigned")
	// but the SDK MarshalJSON expects exact case match (e.g. "UserAssigned")
	// Without normalization, MarshalJSON defaults to "None" and removes the identity
	if payload.Identity != nil {
		// Use string comparison to normalize the type
		switch strings.ToLower(string(payload.Identity.Type)) {
		case "systemassigned":
			payload.Identity.Type = identity.TypeSystemAssigned
		case "userassigned":
			payload.Identity.Type = identity.TypeUserAssigned
		case "systemassigned, userassigned", "systemassigned,userassigned":
			payload.Identity.Type = identity.TypeSystemAssignedUserAssigned
		case "none":
			payload.Identity.Type = identity.TypeNone
		}
		log.Printf("[DEBUG] Normalized identity type from %q to %q", existing.Model.Identity.Type, payload.Identity.Type)
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expandedIdentity
	}

	if d.HasChange("billing_type") {
		if v, ok := d.GetOk("billing_type"); ok {
			payload.Properties.BillingType = pointer.ToEnum[expressrouteports.ExpressRoutePortsBillingType](v.(string))
		}
	}

	if d.HasChanges("link1", "link2") {
		payload.Properties.Links = expandExpressRoutePortLinks(d.Get("link1").([]any), d.Get("link2").([]any))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// a lock is needed here for subresource express_route_port_authorization needs a lock.
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmExpressRoutePortRead(d, meta)
}

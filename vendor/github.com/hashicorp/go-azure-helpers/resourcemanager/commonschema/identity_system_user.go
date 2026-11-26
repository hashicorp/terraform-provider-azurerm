// Copyright IBM Corp. 2018, 2025
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NOTE: we intentionally don't have an Optional & Computed here for behavioural consistency.

// NOTE: there's two different types of SystemAssignedUserAssignedIdentity supported by Azure:
// The first (List) represents the IdentityIDs as a List of Strings
// The other (Map) represents the IdentityIDs as a Map of String : Object (containing Client/PrincipalID)
// from a users perspective however, these should both be represented using the same schema
// so we have a single schema and separate Expand/Flatten functions

// SystemAssignedUserAssignedIdentityRequired returns the System Assigned User Assigned Identity schema where this is Required
func SystemAssignedUserAssignedIdentityRequired() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "The type of Managed Service Identity to configure on this resource.",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
						string(identity.TypeSystemAssigned),
						string(identity.TypeSystemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Description: "A set of User Assigned Managed Identity IDs to assign to this resource.",
					Type:        schema.TypeSet,
					Optional:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
				"principal_id": {
					Description: "The Principal ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"tenant_id": {
					Description: "The Tenant ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

// SystemAssignedUserAssignedIdentityRequiredForceNew returns the System Assigned User Assigned Identity schema where this is Required and ForceNew
func SystemAssignedUserAssignedIdentityRequiredForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "The type of Managed Service Identity to configure on this resource.",
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
						string(identity.TypeSystemAssigned),
						string(identity.TypeSystemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Description: "A set of User Assigned Managed Identity IDs to assign to this resource.",
					Type:        schema.TypeSet,
					Optional:    true,
					ForceNew:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
				"principal_id": {
					Description: "The Principal ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"tenant_id": {
					Description: "The Tenant ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

// SystemAssignedUserAssignedIdentityOptional returns the System Assigned User Assigned Identity schema where this is Optional
func SystemAssignedUserAssignedIdentityOptional() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "The type of Managed Service Identity to configure on this resource.",
					Type:        schema.TypeString,
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
						string(identity.TypeSystemAssigned),
						string(identity.TypeSystemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Description: "A set of User Assigned Managed Identity IDs to assign to this resource.",
					Type:        schema.TypeSet,
					Optional:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
				"principal_id": {
					Description: "The Principal ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"tenant_id": {
					Description: "The Tenant ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

// SystemAssignedUserAssignedIdentityOptionalForceNew returns the System Assigned User Assigned Identity schema where this is Optional and ForceNew
func SystemAssignedUserAssignedIdentityOptionalForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "The type of Managed Service Identity to configure on this resource.",
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
						string(identity.TypeSystemAssigned),
						string(identity.TypeSystemAssignedUserAssigned),
					}, false),
				},
				"identity_ids": {
					Description: "A set of User Assigned Managed Identity IDs to assign to this resource.",
					Type:        schema.TypeSet,
					Optional:    true,
					ForceNew:    true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
				"principal_id": {
					Description: "The Principal ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"tenant_id": {
					Description: "The Tenant ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

// SystemAssignedUserAssignedIdentityComputed returns the System Assigned User Assigned Identity schema where this is Computed
func SystemAssignedUserAssignedIdentityComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Description: "The type of Managed Service Identity configured on this resource.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"identity_ids": {
					Description: "A list of User Assigned Managed Identity IDs assigned to this resource.",
					Type:        schema.TypeList,
					Computed:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"principal_id": {
					Description: "The Principal ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
				"tenant_id": {
					Description: "The Tenant ID associated with this Managed Service Identity.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

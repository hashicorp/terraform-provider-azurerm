// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// NOTE: we intentionally don't have an Optional & Computed here for behavioural consistency.

// UserAssignedIdentityRequired returns the User Assigned Identity schema where this is Required
func UserAssignedIdentityRequired() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
	}
}

// UserAssignedIdentityRequiredForceNew returns the User Assigned Identity schema where this is Required and ForceNew
func UserAssignedIdentityRequiredForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     schema.TypeSet,
					Required: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
	}
}

// UserAssignedIdentityOptional returns the User Assigned Identity schema where this is Optional
func UserAssignedIdentityOptional() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
	}
}

// UserAssignedIdentityOptionalForceNew returns the User Assigned Identity schema where this is Optional and ForceNew
func UserAssignedIdentityOptionalForceNew() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(identity.TypeUserAssigned),
					}, false),
				},
				"identity_ids": {
					Type:     schema.TypeSet,
					Required: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},
	}
}

// UserAssignedIdentityComputed returns the User Assigned Identity schema where this is Computed
func UserAssignedIdentityComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"identity_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	}
}

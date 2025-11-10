// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/accesscontrol"
)

func ExpandDataLakeGen2AceList(input []interface{}) (*accesscontrol.ACL, error) {
	if len(input) == 0 {
		return nil, nil
	}
	aceList := make([]accesscontrol.ACE, len(input))

	for i := 0; i < len(input); i++ {
		v := input[i].(map[string]interface{})

		isDefault := false
		if scopeRaw, ok := v["scope"]; ok {
			if scopeRaw.(string) == "default" {
				isDefault = true
			}
		}

		tagType := accesscontrol.TagType(v["type"].(string))

		var id *uuid.UUID
		if raw, ok := v["id"]; ok && raw != "" {
			idTemp, err := uuid.Parse(raw.(string))
			if err != nil {
				return nil, err
			}
			id = &idTemp
		}

		permissions := v["permissions"].(string)

		ace := accesscontrol.ACE{
			IsDefault:    isDefault,
			TagType:      tagType,
			TagQualifier: id,
			Permissions:  permissions,
		}
		aceList[i] = ace
	}

	return &accesscontrol.ACL{Entries: aceList}, nil
}

func FlattenDataLakeGen2AceList(d *pluginsdk.ResourceData, acl accesscontrol.ACL) []interface{} {
	existingACLs, _ := ExpandDataLakeGen2AceList(d.Get("ace").(*pluginsdk.Set).List())
	output := make([]interface{}, 0)

	for _, v := range acl.Entries {
		// Filter ACL defalt entries (ones without ID value, for scopes 'user', 'group', 'other', 'mask').
		//    Include default entries, only if use in a configuration, to match the state file.
		if v.TagQualifier == nil && existingACLs != nil && !isACLContainingEntry(existingACLs, v.TagType, v.TagQualifier, v.IsDefault) {
			continue
		}

		ace := make(map[string]interface{})

		scope := "access"
		if v.IsDefault {
			scope = "default"
		}
		ace["scope"] = scope
		ace["type"] = string(v.TagType)
		id := ""
		if v.TagQualifier != nil {
			id = v.TagQualifier.String()
		}
		ace["id"] = id
		ace["permissions"] = v.Permissions

		output = append(output, ace)
	}

	return output
}

func isACLContainingEntry(acl *accesscontrol.ACL, tagType accesscontrol.TagType, tagQualifier *uuid.UUID, isDefault bool) bool {
	if acl == nil || acl.Entries == nil || len(acl.Entries) == 0 {
		return false
	}

	for _, v := range acl.Entries {
		if v.TagType == tagType && v.TagQualifier == tagQualifier && v.IsDefault == isDefault {
			return true
		}
	}

	return false
}

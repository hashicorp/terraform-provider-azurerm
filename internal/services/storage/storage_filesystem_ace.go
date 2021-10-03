package storage

import (
	"github.com/google/uuid"
	"github.com/tombuildsstuff/giovanni/storage/accesscontrol"
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

func FlattenDataLakeGen2AceList(acl accesscontrol.ACL) []interface{} {
	output := make([]interface{}, len(acl.Entries))

	for i, v := range acl.Entries {
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

		output[i] = ace
	}
	return output
}

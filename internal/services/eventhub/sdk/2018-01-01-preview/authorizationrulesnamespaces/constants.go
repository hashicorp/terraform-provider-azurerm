package authorizationrulesnamespaces

import "strings"

type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

func PossibleValuesForAccessRights() []string {
	return []string{
		"Listen",
		"Manage",
		"Send",
	}
}

func parseAccessRights(input string) (*AccessRights, error) {
	vals := map[string]AccessRights{
		"listen": "Listen",
		"manage": "Manage",
		"send":   "Send",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := AccessRights(v)
	return &out, nil
}

type KeyType string

const (
	KeyTypePrimaryKey   KeyType = "PrimaryKey"
	KeyTypeSecondaryKey KeyType = "SecondaryKey"
)

func PossibleValuesForKeyType() []string {
	return []string{
		"PrimaryKey",
		"SecondaryKey",
	}
}

func parseKeyType(input string) (*KeyType, error) {
	vals := map[string]KeyType{
		"primarykey":   "PrimaryKey",
		"secondarykey": "SecondaryKey",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := KeyType(v)
	return &out, nil
}

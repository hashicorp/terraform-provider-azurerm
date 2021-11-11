package eventhubs

import "strings"

type AccessRights string

const (
	AccessRightsListen AccessRights = "Listen"
	AccessRightsManage AccessRights = "Manage"
	AccessRightsSend   AccessRights = "Send"
)

func PossibleValuesForAccessRights() []string {
	return []string{
		string(AccessRightsListen),
		string(AccessRightsManage),
		string(AccessRightsSend),
	}
}

func parseAccessRights(input string) (*AccessRights, error) {
	vals := map[string]AccessRights{
		"listen": AccessRightsListen,
		"manage": AccessRightsManage,
		"send":   AccessRightsSend,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessRights(input)
	return &out, nil
}

type CreatedByType string

const (
	CreatedByTypeApplication     CreatedByType = "Application"
	CreatedByTypeKey             CreatedByType = "Key"
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	CreatedByTypeUser            CreatedByType = "User"
)

func PossibleValuesForCreatedByType() []string {
	return []string{
		string(CreatedByTypeApplication),
		string(CreatedByTypeKey),
		string(CreatedByTypeManagedIdentity),
		string(CreatedByTypeUser),
	}
}

func parseCreatedByType(input string) (*CreatedByType, error) {
	vals := map[string]CreatedByType{
		"application":     CreatedByTypeApplication,
		"key":             CreatedByTypeKey,
		"managedidentity": CreatedByTypeManagedIdentity,
		"user":            CreatedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreatedByType(input)
	return &out, nil
}

type EncodingCaptureDescription string

const (
	EncodingCaptureDescriptionAvro        EncodingCaptureDescription = "Avro"
	EncodingCaptureDescriptionAvroDeflate EncodingCaptureDescription = "AvroDeflate"
)

func PossibleValuesForEncodingCaptureDescription() []string {
	return []string{
		string(EncodingCaptureDescriptionAvro),
		string(EncodingCaptureDescriptionAvroDeflate),
	}
}

func parseEncodingCaptureDescription(input string) (*EncodingCaptureDescription, error) {
	vals := map[string]EncodingCaptureDescription{
		"avro":        EncodingCaptureDescriptionAvro,
		"avrodeflate": EncodingCaptureDescriptionAvroDeflate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncodingCaptureDescription(input)
	return &out, nil
}

type EntityStatus string

const (
	EntityStatusActive          EntityStatus = "Active"
	EntityStatusCreating        EntityStatus = "Creating"
	EntityStatusDeleting        EntityStatus = "Deleting"
	EntityStatusDisabled        EntityStatus = "Disabled"
	EntityStatusReceiveDisabled EntityStatus = "ReceiveDisabled"
	EntityStatusRenaming        EntityStatus = "Renaming"
	EntityStatusRestoring       EntityStatus = "Restoring"
	EntityStatusSendDisabled    EntityStatus = "SendDisabled"
	EntityStatusUnknown         EntityStatus = "Unknown"
)

func PossibleValuesForEntityStatus() []string {
	return []string{
		string(EntityStatusActive),
		string(EntityStatusCreating),
		string(EntityStatusDeleting),
		string(EntityStatusDisabled),
		string(EntityStatusReceiveDisabled),
		string(EntityStatusRenaming),
		string(EntityStatusRestoring),
		string(EntityStatusSendDisabled),
		string(EntityStatusUnknown),
	}
}

func parseEntityStatus(input string) (*EntityStatus, error) {
	vals := map[string]EntityStatus{
		"active":          EntityStatusActive,
		"creating":        EntityStatusCreating,
		"deleting":        EntityStatusDeleting,
		"disabled":        EntityStatusDisabled,
		"receivedisabled": EntityStatusReceiveDisabled,
		"renaming":        EntityStatusRenaming,
		"restoring":       EntityStatusRestoring,
		"senddisabled":    EntityStatusSendDisabled,
		"unknown":         EntityStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EntityStatus(input)
	return &out, nil
}

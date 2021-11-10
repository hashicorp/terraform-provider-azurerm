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

type EncodingCaptureDescription string

const (
	EncodingCaptureDescriptionAvro        EncodingCaptureDescription = "Avro"
	EncodingCaptureDescriptionAvroDeflate EncodingCaptureDescription = "AvroDeflate"
)

func PossibleValuesForEncodingCaptureDescription() []string {
	return []string{
		"Avro",
		"AvroDeflate",
	}
}

func parseEncodingCaptureDescription(input string) (*EncodingCaptureDescription, error) {
	vals := map[string]EncodingCaptureDescription{
		"avro":        "Avro",
		"avrodeflate": "AvroDeflate",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := EncodingCaptureDescription(v)
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
		"Active",
		"Creating",
		"Deleting",
		"Disabled",
		"ReceiveDisabled",
		"Renaming",
		"Restoring",
		"SendDisabled",
		"Unknown",
	}
}

func parseEntityStatus(input string) (*EntityStatus, error) {
	vals := map[string]EntityStatus{
		"active":          "Active",
		"creating":        "Creating",
		"deleting":        "Deleting",
		"disabled":        "Disabled",
		"receivedisabled": "ReceiveDisabled",
		"renaming":        "Renaming",
		"restoring":       "Restoring",
		"senddisabled":    "SendDisabled",
		"unknown":         "Unknown",
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// it could be a new value - best effort convert this
	v := input

	out := EntityStatus(v)
	return &out, nil
}

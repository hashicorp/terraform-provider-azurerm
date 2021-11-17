package domains

import "strings"

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

type DomainProvisioningState string

const (
	DomainProvisioningStateCanceled  DomainProvisioningState = "Canceled"
	DomainProvisioningStateCreating  DomainProvisioningState = "Creating"
	DomainProvisioningStateDeleting  DomainProvisioningState = "Deleting"
	DomainProvisioningStateFailed    DomainProvisioningState = "Failed"
	DomainProvisioningStateSucceeded DomainProvisioningState = "Succeeded"
	DomainProvisioningStateUpdating  DomainProvisioningState = "Updating"
)

func PossibleValuesForDomainProvisioningState() []string {
	return []string{
		string(DomainProvisioningStateCanceled),
		string(DomainProvisioningStateCreating),
		string(DomainProvisioningStateDeleting),
		string(DomainProvisioningStateFailed),
		string(DomainProvisioningStateSucceeded),
		string(DomainProvisioningStateUpdating),
	}
}

func parseDomainProvisioningState(input string) (*DomainProvisioningState, error) {
	vals := map[string]DomainProvisioningState{
		"canceled":  DomainProvisioningStateCanceled,
		"creating":  DomainProvisioningStateCreating,
		"deleting":  DomainProvisioningStateDeleting,
		"failed":    DomainProvisioningStateFailed,
		"succeeded": DomainProvisioningStateSucceeded,
		"updating":  DomainProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainProvisioningState(input)
	return &out, nil
}

type InputSchema string

const (
	InputSchemaCloudEventSchemaVOneZero InputSchema = "CloudEventSchemaV1_0"
	InputSchemaCustomEventSchema        InputSchema = "CustomEventSchema"
	InputSchemaEventGridSchema          InputSchema = "EventGridSchema"
)

func PossibleValuesForInputSchema() []string {
	return []string{
		string(InputSchemaCloudEventSchemaVOneZero),
		string(InputSchemaCustomEventSchema),
		string(InputSchemaEventGridSchema),
	}
}

func parseInputSchema(input string) (*InputSchema, error) {
	vals := map[string]InputSchema{
		"cloudeventschemav1_0": InputSchemaCloudEventSchemaVOneZero,
		"customeventschema":    InputSchemaCustomEventSchema,
		"eventgridschema":      InputSchemaEventGridSchema,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InputSchema(input)
	return &out, nil
}

type InputSchemaMappingType string

const (
	InputSchemaMappingTypeJson InputSchemaMappingType = "Json"
)

func PossibleValuesForInputSchemaMappingType() []string {
	return []string{
		string(InputSchemaMappingTypeJson),
	}
}

func parseInputSchemaMappingType(input string) (*InputSchemaMappingType, error) {
	vals := map[string]InputSchemaMappingType{
		"json": InputSchemaMappingTypeJson,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InputSchemaMappingType(input)
	return &out, nil
}

type IpActionType string

const (
	IpActionTypeAllow IpActionType = "Allow"
)

func PossibleValuesForIpActionType() []string {
	return []string{
		string(IpActionTypeAllow),
	}
}

func parseIpActionType(input string) (*IpActionType, error) {
	vals := map[string]IpActionType{
		"allow": IpActionTypeAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IpActionType(input)
	return &out, nil
}

type PersistedConnectionStatus string

const (
	PersistedConnectionStatusApproved     PersistedConnectionStatus = "Approved"
	PersistedConnectionStatusDisconnected PersistedConnectionStatus = "Disconnected"
	PersistedConnectionStatusPending      PersistedConnectionStatus = "Pending"
	PersistedConnectionStatusRejected     PersistedConnectionStatus = "Rejected"
)

func PossibleValuesForPersistedConnectionStatus() []string {
	return []string{
		string(PersistedConnectionStatusApproved),
		string(PersistedConnectionStatusDisconnected),
		string(PersistedConnectionStatusPending),
		string(PersistedConnectionStatusRejected),
	}
}

func parsePersistedConnectionStatus(input string) (*PersistedConnectionStatus, error) {
	vals := map[string]PersistedConnectionStatus{
		"approved":     PersistedConnectionStatusApproved,
		"disconnected": PersistedConnectionStatusDisconnected,
		"pending":      PersistedConnectionStatusPending,
		"rejected":     PersistedConnectionStatusRejected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PersistedConnectionStatus(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type ResourceProvisioningState string

const (
	ResourceProvisioningStateCanceled  ResourceProvisioningState = "Canceled"
	ResourceProvisioningStateCreating  ResourceProvisioningState = "Creating"
	ResourceProvisioningStateDeleting  ResourceProvisioningState = "Deleting"
	ResourceProvisioningStateFailed    ResourceProvisioningState = "Failed"
	ResourceProvisioningStateSucceeded ResourceProvisioningState = "Succeeded"
	ResourceProvisioningStateUpdating  ResourceProvisioningState = "Updating"
)

func PossibleValuesForResourceProvisioningState() []string {
	return []string{
		string(ResourceProvisioningStateCanceled),
		string(ResourceProvisioningStateCreating),
		string(ResourceProvisioningStateDeleting),
		string(ResourceProvisioningStateFailed),
		string(ResourceProvisioningStateSucceeded),
		string(ResourceProvisioningStateUpdating),
	}
}

func parseResourceProvisioningState(input string) (*ResourceProvisioningState, error) {
	vals := map[string]ResourceProvisioningState{
		"canceled":  ResourceProvisioningStateCanceled,
		"creating":  ResourceProvisioningStateCreating,
		"deleting":  ResourceProvisioningStateDeleting,
		"failed":    ResourceProvisioningStateFailed,
		"succeeded": ResourceProvisioningStateSucceeded,
		"updating":  ResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceProvisioningState(input)
	return &out, nil
}

type Sku string

const (
	SkuBasic   Sku = "Basic"
	SkuPremium Sku = "Premium"
)

func PossibleValuesForSku() []string {
	return []string{
		string(SkuBasic),
		string(SkuPremium),
	}
}

func parseSku(input string) (*Sku, error) {
	vals := map[string]Sku{
		"basic":   SkuBasic,
		"premium": SkuPremium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Sku(input)
	return &out, nil
}

package account

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

type LastModifiedByType string

const (
	LastModifiedByTypeApplication     LastModifiedByType = "Application"
	LastModifiedByTypeKey             LastModifiedByType = "Key"
	LastModifiedByTypeManagedIdentity LastModifiedByType = "ManagedIdentity"
	LastModifiedByTypeUser            LastModifiedByType = "User"
)

func PossibleValuesForLastModifiedByType() []string {
	return []string{
		string(LastModifiedByTypeApplication),
		string(LastModifiedByTypeKey),
		string(LastModifiedByTypeManagedIdentity),
		string(LastModifiedByTypeUser),
	}
}

func parseLastModifiedByType(input string) (*LastModifiedByType, error) {
	vals := map[string]LastModifiedByType{
		"application":     LastModifiedByTypeApplication,
		"key":             LastModifiedByTypeKey,
		"managedidentity": LastModifiedByTypeManagedIdentity,
		"user":            LastModifiedByTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LastModifiedByType(input)
	return &out, nil
}

type Name string

const (
	NameStandard Name = "Standard"
)

func PossibleValuesForName() []string {
	return []string{
		string(NameStandard),
	}
}

func parseName(input string) (*Name, error) {
	vals := map[string]Name{
		"standard": NameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Name(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateMoving       ProvisioningState = "Moving"
	ProvisioningStateSoftDeleted  ProvisioningState = "SoftDeleted"
	ProvisioningStateSoftDeleting ProvisioningState = "SoftDeleting"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUnknown      ProvisioningState = "Unknown"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateSoftDeleted),
		string(ProvisioningStateSoftDeleting),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"moving":       ProvisioningStateMoving,
		"softdeleted":  ProvisioningStateSoftDeleted,
		"softdeleting": ProvisioningStateSoftDeleting,
		"succeeded":    ProvisioningStateSucceeded,
		"unknown":      ProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled     PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled      PublicNetworkAccess = "Enabled"
	PublicNetworkAccessNotSpecified PublicNetworkAccess = "NotSpecified"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
		string(PublicNetworkAccessNotSpecified),
	}
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled":     PublicNetworkAccessDisabled,
		"enabled":      PublicNetworkAccessEnabled,
		"notspecified": PublicNetworkAccessNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type Status string

const (
	StatusApproved     Status = "Approved"
	StatusDisconnected Status = "Disconnected"
	StatusPending      Status = "Pending"
	StatusRejected     Status = "Rejected"
	StatusUnknown      Status = "Unknown"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusApproved),
		string(StatusDisconnected),
		string(StatusPending),
		string(StatusRejected),
		string(StatusUnknown),
	}
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"approved":     StatusApproved,
		"disconnected": StatusDisconnected,
		"pending":      StatusPending,
		"rejected":     StatusRejected,
		"unknown":      StatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

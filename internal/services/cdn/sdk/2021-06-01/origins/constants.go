package origins

import "strings"

type IdentityType string

const (
	IdentityTypeApplication     IdentityType = "application"
	IdentityTypeKey             IdentityType = "key"
	IdentityTypeManagedIdentity IdentityType = "managedIdentity"
	IdentityTypeUser            IdentityType = "user"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeApplication),
		string(IdentityTypeKey),
		string(IdentityTypeManagedIdentity),
		string(IdentityTypeUser),
	}
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"application":     IdentityTypeApplication,
		"key":             IdentityTypeKey,
		"managedidentity": IdentityTypeManagedIdentity,
		"user":            IdentityTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type OriginResourceState string

const (
	OriginResourceStateActive   OriginResourceState = "Active"
	OriginResourceStateCreating OriginResourceState = "Creating"
	OriginResourceStateDeleting OriginResourceState = "Deleting"
)

func PossibleValuesForOriginResourceState() []string {
	return []string{
		string(OriginResourceStateActive),
		string(OriginResourceStateCreating),
		string(OriginResourceStateDeleting),
	}
}

func parseOriginResourceState(input string) (*OriginResourceState, error) {
	vals := map[string]OriginResourceState{
		"active":   OriginResourceStateActive,
		"creating": OriginResourceStateCreating,
		"deleting": OriginResourceStateDeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OriginResourceState(input)
	return &out, nil
}

type PrivateEndpointStatus string

const (
	PrivateEndpointStatusApproved     PrivateEndpointStatus = "Approved"
	PrivateEndpointStatusDisconnected PrivateEndpointStatus = "Disconnected"
	PrivateEndpointStatusPending      PrivateEndpointStatus = "Pending"
	PrivateEndpointStatusRejected     PrivateEndpointStatus = "Rejected"
	PrivateEndpointStatusTimeout      PrivateEndpointStatus = "Timeout"
)

func PossibleValuesForPrivateEndpointStatus() []string {
	return []string{
		string(PrivateEndpointStatusApproved),
		string(PrivateEndpointStatusDisconnected),
		string(PrivateEndpointStatusPending),
		string(PrivateEndpointStatusRejected),
		string(PrivateEndpointStatusTimeout),
	}
}

func parsePrivateEndpointStatus(input string) (*PrivateEndpointStatus, error) {
	vals := map[string]PrivateEndpointStatus{
		"approved":     PrivateEndpointStatusApproved,
		"disconnected": PrivateEndpointStatusDisconnected,
		"pending":      PrivateEndpointStatusPending,
		"rejected":     PrivateEndpointStatusRejected,
		"timeout":      PrivateEndpointStatusTimeout,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateEndpointStatus(input)
	return &out, nil
}

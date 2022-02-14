package kubernetes

import "strings"

type AuthenticationMethod string

const (
	AuthenticationMethodAAD   AuthenticationMethod = "AAD"
	AuthenticationMethodToken AuthenticationMethod = "Token"
)

func PossibleValuesForAuthenticationMethod() []string {
	return []string{
		string(AuthenticationMethodAAD),
		string(AuthenticationMethodToken),
	}
}

func parseAuthenticationMethod(input string) (*AuthenticationMethod, error) {
	vals := map[string]AuthenticationMethod{
		"aad":   AuthenticationMethodAAD,
		"token": AuthenticationMethodToken,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMethod(input)
	return &out, nil
}

type ConnectivityStatus string

const (
	ConnectivityStatusConnected  ConnectivityStatus = "Connected"
	ConnectivityStatusConnecting ConnectivityStatus = "Connecting"
	ConnectivityStatusExpired    ConnectivityStatus = "Expired"
	ConnectivityStatusOffline    ConnectivityStatus = "Offline"
)

func PossibleValuesForConnectivityStatus() []string {
	return []string{
		string(ConnectivityStatusConnected),
		string(ConnectivityStatusConnecting),
		string(ConnectivityStatusExpired),
		string(ConnectivityStatusOffline),
	}
}

func parseConnectivityStatus(input string) (*ConnectivityStatus, error) {
	vals := map[string]ConnectivityStatus{
		"connected":  ConnectivityStatusConnected,
		"connecting": ConnectivityStatusConnecting,
		"expired":    ConnectivityStatusExpired,
		"offline":    ConnectivityStatusOffline,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityStatus(input)
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

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

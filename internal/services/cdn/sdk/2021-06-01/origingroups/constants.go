package origingroups

import "strings"

type HealthProbeRequestType string

const (
	HealthProbeRequestTypeGET    HealthProbeRequestType = "GET"
	HealthProbeRequestTypeHEAD   HealthProbeRequestType = "HEAD"
	HealthProbeRequestTypeNotSet HealthProbeRequestType = "NotSet"
)

func PossibleValuesForHealthProbeRequestType() []string {
	return []string{
		string(HealthProbeRequestTypeGET),
		string(HealthProbeRequestTypeHEAD),
		string(HealthProbeRequestTypeNotSet),
	}
}

func parseHealthProbeRequestType(input string) (*HealthProbeRequestType, error) {
	vals := map[string]HealthProbeRequestType{
		"get":    HealthProbeRequestTypeGET,
		"head":   HealthProbeRequestTypeHEAD,
		"notset": HealthProbeRequestTypeNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HealthProbeRequestType(input)
	return &out, nil
}

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

type OriginGroupResourceState string

const (
	OriginGroupResourceStateActive   OriginGroupResourceState = "Active"
	OriginGroupResourceStateCreating OriginGroupResourceState = "Creating"
	OriginGroupResourceStateDeleting OriginGroupResourceState = "Deleting"
)

func PossibleValuesForOriginGroupResourceState() []string {
	return []string{
		string(OriginGroupResourceStateActive),
		string(OriginGroupResourceStateCreating),
		string(OriginGroupResourceStateDeleting),
	}
}

func parseOriginGroupResourceState(input string) (*OriginGroupResourceState, error) {
	vals := map[string]OriginGroupResourceState{
		"active":   OriginGroupResourceStateActive,
		"creating": OriginGroupResourceStateCreating,
		"deleting": OriginGroupResourceStateDeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OriginGroupResourceState(input)
	return &out, nil
}

type ProbeProtocol string

const (
	ProbeProtocolHttp   ProbeProtocol = "Http"
	ProbeProtocolHttps  ProbeProtocol = "Https"
	ProbeProtocolNotSet ProbeProtocol = "NotSet"
)

func PossibleValuesForProbeProtocol() []string {
	return []string{
		string(ProbeProtocolHttp),
		string(ProbeProtocolHttps),
		string(ProbeProtocolNotSet),
	}
}

func parseProbeProtocol(input string) (*ProbeProtocol, error) {
	vals := map[string]ProbeProtocol{
		"http":   ProbeProtocolHttp,
		"https":  ProbeProtocolHttps,
		"notset": ProbeProtocolNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProbeProtocol(input)
	return &out, nil
}

type ResponseBasedDetectedErrorTypes string

const (
	ResponseBasedDetectedErrorTypesNone             ResponseBasedDetectedErrorTypes = "None"
	ResponseBasedDetectedErrorTypesTcpAndHttpErrors ResponseBasedDetectedErrorTypes = "TcpAndHttpErrors"
	ResponseBasedDetectedErrorTypesTcpErrorsOnly    ResponseBasedDetectedErrorTypes = "TcpErrorsOnly"
)

func PossibleValuesForResponseBasedDetectedErrorTypes() []string {
	return []string{
		string(ResponseBasedDetectedErrorTypesNone),
		string(ResponseBasedDetectedErrorTypesTcpAndHttpErrors),
		string(ResponseBasedDetectedErrorTypesTcpErrorsOnly),
	}
}

func parseResponseBasedDetectedErrorTypes(input string) (*ResponseBasedDetectedErrorTypes, error) {
	vals := map[string]ResponseBasedDetectedErrorTypes{
		"none":             ResponseBasedDetectedErrorTypesNone,
		"tcpandhttperrors": ResponseBasedDetectedErrorTypesTcpAndHttpErrors,
		"tcperrorsonly":    ResponseBasedDetectedErrorTypesTcpErrorsOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResponseBasedDetectedErrorTypes(input)
	return &out, nil
}

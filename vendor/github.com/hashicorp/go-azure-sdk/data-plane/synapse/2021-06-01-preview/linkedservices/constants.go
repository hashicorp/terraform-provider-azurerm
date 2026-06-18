package linkedservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Db2AuthenticationType string

const (
	Db2AuthenticationTypeBasic Db2AuthenticationType = "Basic"
)

func PossibleValuesForDb2AuthenticationType() []string {
	return []string{
		string(Db2AuthenticationTypeBasic),
	}
}

func (s *Db2AuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDb2AuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDb2AuthenticationType(input string) (*Db2AuthenticationType, error) {
	vals := map[string]Db2AuthenticationType{
		"basic": Db2AuthenticationTypeBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Db2AuthenticationType(input)
	return &out, nil
}

type FtpAuthenticationType string

const (
	FtpAuthenticationTypeAnonymous FtpAuthenticationType = "Anonymous"
	FtpAuthenticationTypeBasic     FtpAuthenticationType = "Basic"
)

func PossibleValuesForFtpAuthenticationType() []string {
	return []string{
		string(FtpAuthenticationTypeAnonymous),
		string(FtpAuthenticationTypeBasic),
	}
}

func (s *FtpAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFtpAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFtpAuthenticationType(input string) (*FtpAuthenticationType, error) {
	vals := map[string]FtpAuthenticationType{
		"anonymous": FtpAuthenticationTypeAnonymous,
		"basic":     FtpAuthenticationTypeBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FtpAuthenticationType(input)
	return &out, nil
}

type GoogleAdWordsAuthenticationType string

const (
	GoogleAdWordsAuthenticationTypeServiceAuthentication GoogleAdWordsAuthenticationType = "ServiceAuthentication"
	GoogleAdWordsAuthenticationTypeUserAuthentication    GoogleAdWordsAuthenticationType = "UserAuthentication"
)

func PossibleValuesForGoogleAdWordsAuthenticationType() []string {
	return []string{
		string(GoogleAdWordsAuthenticationTypeServiceAuthentication),
		string(GoogleAdWordsAuthenticationTypeUserAuthentication),
	}
}

func (s *GoogleAdWordsAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGoogleAdWordsAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGoogleAdWordsAuthenticationType(input string) (*GoogleAdWordsAuthenticationType, error) {
	vals := map[string]GoogleAdWordsAuthenticationType{
		"serviceauthentication": GoogleAdWordsAuthenticationTypeServiceAuthentication,
		"userauthentication":    GoogleAdWordsAuthenticationTypeUserAuthentication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GoogleAdWordsAuthenticationType(input)
	return &out, nil
}

type GoogleBigQueryAuthenticationType string

const (
	GoogleBigQueryAuthenticationTypeServiceAuthentication GoogleBigQueryAuthenticationType = "ServiceAuthentication"
	GoogleBigQueryAuthenticationTypeUserAuthentication    GoogleBigQueryAuthenticationType = "UserAuthentication"
)

func PossibleValuesForGoogleBigQueryAuthenticationType() []string {
	return []string{
		string(GoogleBigQueryAuthenticationTypeServiceAuthentication),
		string(GoogleBigQueryAuthenticationTypeUserAuthentication),
	}
}

func (s *GoogleBigQueryAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGoogleBigQueryAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGoogleBigQueryAuthenticationType(input string) (*GoogleBigQueryAuthenticationType, error) {
	vals := map[string]GoogleBigQueryAuthenticationType{
		"serviceauthentication": GoogleBigQueryAuthenticationTypeServiceAuthentication,
		"userauthentication":    GoogleBigQueryAuthenticationTypeUserAuthentication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GoogleBigQueryAuthenticationType(input)
	return &out, nil
}

type HBaseAuthenticationType string

const (
	HBaseAuthenticationTypeAnonymous HBaseAuthenticationType = "Anonymous"
	HBaseAuthenticationTypeBasic     HBaseAuthenticationType = "Basic"
)

func PossibleValuesForHBaseAuthenticationType() []string {
	return []string{
		string(HBaseAuthenticationTypeAnonymous),
		string(HBaseAuthenticationTypeBasic),
	}
}

func (s *HBaseAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHBaseAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHBaseAuthenticationType(input string) (*HBaseAuthenticationType, error) {
	vals := map[string]HBaseAuthenticationType{
		"anonymous": HBaseAuthenticationTypeAnonymous,
		"basic":     HBaseAuthenticationTypeBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HBaseAuthenticationType(input)
	return &out, nil
}

type HTTPAuthenticationType string

const (
	HTTPAuthenticationTypeAnonymous         HTTPAuthenticationType = "Anonymous"
	HTTPAuthenticationTypeBasic             HTTPAuthenticationType = "Basic"
	HTTPAuthenticationTypeClientCertificate HTTPAuthenticationType = "ClientCertificate"
	HTTPAuthenticationTypeDigest            HTTPAuthenticationType = "Digest"
	HTTPAuthenticationTypeWindows           HTTPAuthenticationType = "Windows"
)

func PossibleValuesForHTTPAuthenticationType() []string {
	return []string{
		string(HTTPAuthenticationTypeAnonymous),
		string(HTTPAuthenticationTypeBasic),
		string(HTTPAuthenticationTypeClientCertificate),
		string(HTTPAuthenticationTypeDigest),
		string(HTTPAuthenticationTypeWindows),
	}
}

func (s *HTTPAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPAuthenticationType(input string) (*HTTPAuthenticationType, error) {
	vals := map[string]HTTPAuthenticationType{
		"anonymous":         HTTPAuthenticationTypeAnonymous,
		"basic":             HTTPAuthenticationTypeBasic,
		"clientcertificate": HTTPAuthenticationTypeClientCertificate,
		"digest":            HTTPAuthenticationTypeDigest,
		"windows":           HTTPAuthenticationTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPAuthenticationType(input)
	return &out, nil
}

type HiveAuthenticationType string

const (
	HiveAuthenticationTypeAnonymous                    HiveAuthenticationType = "Anonymous"
	HiveAuthenticationTypeUsername                     HiveAuthenticationType = "Username"
	HiveAuthenticationTypeUsernameAndPassword          HiveAuthenticationType = "UsernameAndPassword"
	HiveAuthenticationTypeWindowsAzureHDInsightService HiveAuthenticationType = "WindowsAzureHDInsightService"
)

func PossibleValuesForHiveAuthenticationType() []string {
	return []string{
		string(HiveAuthenticationTypeAnonymous),
		string(HiveAuthenticationTypeUsername),
		string(HiveAuthenticationTypeUsernameAndPassword),
		string(HiveAuthenticationTypeWindowsAzureHDInsightService),
	}
}

func (s *HiveAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHiveAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHiveAuthenticationType(input string) (*HiveAuthenticationType, error) {
	vals := map[string]HiveAuthenticationType{
		"anonymous":                    HiveAuthenticationTypeAnonymous,
		"username":                     HiveAuthenticationTypeUsername,
		"usernameandpassword":          HiveAuthenticationTypeUsernameAndPassword,
		"windowsazurehdinsightservice": HiveAuthenticationTypeWindowsAzureHDInsightService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HiveAuthenticationType(input)
	return &out, nil
}

type HiveServerType string

const (
	HiveServerTypeHiveServerOne    HiveServerType = "HiveServer1"
	HiveServerTypeHiveServerTwo    HiveServerType = "HiveServer2"
	HiveServerTypeHiveThriftServer HiveServerType = "HiveThriftServer"
)

func PossibleValuesForHiveServerType() []string {
	return []string{
		string(HiveServerTypeHiveServerOne),
		string(HiveServerTypeHiveServerTwo),
		string(HiveServerTypeHiveThriftServer),
	}
}

func (s *HiveServerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHiveServerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHiveServerType(input string) (*HiveServerType, error) {
	vals := map[string]HiveServerType{
		"hiveserver1":      HiveServerTypeHiveServerOne,
		"hiveserver2":      HiveServerTypeHiveServerTwo,
		"hivethriftserver": HiveServerTypeHiveThriftServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HiveServerType(input)
	return &out, nil
}

type HiveThriftTransportProtocol string

const (
	HiveThriftTransportProtocolBinary HiveThriftTransportProtocol = "Binary"
	HiveThriftTransportProtocolHTTP   HiveThriftTransportProtocol = "HTTP "
	HiveThriftTransportProtocolSASL   HiveThriftTransportProtocol = "SASL"
)

func PossibleValuesForHiveThriftTransportProtocol() []string {
	return []string{
		string(HiveThriftTransportProtocolBinary),
		string(HiveThriftTransportProtocolHTTP),
		string(HiveThriftTransportProtocolSASL),
	}
}

func (s *HiveThriftTransportProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHiveThriftTransportProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHiveThriftTransportProtocol(input string) (*HiveThriftTransportProtocol, error) {
	vals := map[string]HiveThriftTransportProtocol{
		"binary": HiveThriftTransportProtocolBinary,
		"http ":  HiveThriftTransportProtocolHTTP,
		"sasl":   HiveThriftTransportProtocolSASL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HiveThriftTransportProtocol(input)
	return &out, nil
}

type ImpalaAuthenticationType string

const (
	ImpalaAuthenticationTypeAnonymous           ImpalaAuthenticationType = "Anonymous"
	ImpalaAuthenticationTypeSASLUsername        ImpalaAuthenticationType = "SASLUsername"
	ImpalaAuthenticationTypeUsernameAndPassword ImpalaAuthenticationType = "UsernameAndPassword"
)

func PossibleValuesForImpalaAuthenticationType() []string {
	return []string{
		string(ImpalaAuthenticationTypeAnonymous),
		string(ImpalaAuthenticationTypeSASLUsername),
		string(ImpalaAuthenticationTypeUsernameAndPassword),
	}
}

func (s *ImpalaAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImpalaAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImpalaAuthenticationType(input string) (*ImpalaAuthenticationType, error) {
	vals := map[string]ImpalaAuthenticationType{
		"anonymous":           ImpalaAuthenticationTypeAnonymous,
		"saslusername":        ImpalaAuthenticationTypeSASLUsername,
		"usernameandpassword": ImpalaAuthenticationTypeUsernameAndPassword,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImpalaAuthenticationType(input)
	return &out, nil
}

type IntegrationRuntimeReferenceType string

const (
	IntegrationRuntimeReferenceTypeIntegrationRuntimeReference IntegrationRuntimeReferenceType = "IntegrationRuntimeReference"
)

func PossibleValuesForIntegrationRuntimeReferenceType() []string {
	return []string{
		string(IntegrationRuntimeReferenceTypeIntegrationRuntimeReference),
	}
}

func (s *IntegrationRuntimeReferenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIntegrationRuntimeReferenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIntegrationRuntimeReferenceType(input string) (*IntegrationRuntimeReferenceType, error) {
	vals := map[string]IntegrationRuntimeReferenceType{
		"integrationruntimereference": IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IntegrationRuntimeReferenceType(input)
	return &out, nil
}

type MongoDbAuthenticationType string

const (
	MongoDbAuthenticationTypeAnonymous MongoDbAuthenticationType = "Anonymous"
	MongoDbAuthenticationTypeBasic     MongoDbAuthenticationType = "Basic"
)

func PossibleValuesForMongoDbAuthenticationType() []string {
	return []string{
		string(MongoDbAuthenticationTypeAnonymous),
		string(MongoDbAuthenticationTypeBasic),
	}
}

func (s *MongoDbAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMongoDbAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMongoDbAuthenticationType(input string) (*MongoDbAuthenticationType, error) {
	vals := map[string]MongoDbAuthenticationType{
		"anonymous": MongoDbAuthenticationTypeAnonymous,
		"basic":     MongoDbAuthenticationTypeBasic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MongoDbAuthenticationType(input)
	return &out, nil
}

type ODataAadServicePrincipalCredentialType string

const (
	ODataAadServicePrincipalCredentialTypeServicePrincipalCert ODataAadServicePrincipalCredentialType = "ServicePrincipalCert"
	ODataAadServicePrincipalCredentialTypeServicePrincipalKey  ODataAadServicePrincipalCredentialType = "ServicePrincipalKey"
)

func PossibleValuesForODataAadServicePrincipalCredentialType() []string {
	return []string{
		string(ODataAadServicePrincipalCredentialTypeServicePrincipalCert),
		string(ODataAadServicePrincipalCredentialTypeServicePrincipalKey),
	}
}

func (s *ODataAadServicePrincipalCredentialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseODataAadServicePrincipalCredentialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseODataAadServicePrincipalCredentialType(input string) (*ODataAadServicePrincipalCredentialType, error) {
	vals := map[string]ODataAadServicePrincipalCredentialType{
		"serviceprincipalcert": ODataAadServicePrincipalCredentialTypeServicePrincipalCert,
		"serviceprincipalkey":  ODataAadServicePrincipalCredentialTypeServicePrincipalKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ODataAadServicePrincipalCredentialType(input)
	return &out, nil
}

type ODataAuthenticationType string

const (
	ODataAuthenticationTypeAadServicePrincipal    ODataAuthenticationType = "AadServicePrincipal"
	ODataAuthenticationTypeAnonymous              ODataAuthenticationType = "Anonymous"
	ODataAuthenticationTypeBasic                  ODataAuthenticationType = "Basic"
	ODataAuthenticationTypeManagedServiceIdentity ODataAuthenticationType = "ManagedServiceIdentity"
	ODataAuthenticationTypeWindows                ODataAuthenticationType = "Windows"
)

func PossibleValuesForODataAuthenticationType() []string {
	return []string{
		string(ODataAuthenticationTypeAadServicePrincipal),
		string(ODataAuthenticationTypeAnonymous),
		string(ODataAuthenticationTypeBasic),
		string(ODataAuthenticationTypeManagedServiceIdentity),
		string(ODataAuthenticationTypeWindows),
	}
}

func (s *ODataAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseODataAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseODataAuthenticationType(input string) (*ODataAuthenticationType, error) {
	vals := map[string]ODataAuthenticationType{
		"aadserviceprincipal":    ODataAuthenticationTypeAadServicePrincipal,
		"anonymous":              ODataAuthenticationTypeAnonymous,
		"basic":                  ODataAuthenticationTypeBasic,
		"managedserviceidentity": ODataAuthenticationTypeManagedServiceIdentity,
		"windows":                ODataAuthenticationTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ODataAuthenticationType(input)
	return &out, nil
}

type ParameterType string

const (
	ParameterTypeArray        ParameterType = "Array"
	ParameterTypeBool         ParameterType = "Bool"
	ParameterTypeFloat        ParameterType = "Float"
	ParameterTypeInt          ParameterType = "Int"
	ParameterTypeObject       ParameterType = "Object"
	ParameterTypeSecureString ParameterType = "SecureString"
	ParameterTypeString       ParameterType = "String"
)

func PossibleValuesForParameterType() []string {
	return []string{
		string(ParameterTypeArray),
		string(ParameterTypeBool),
		string(ParameterTypeFloat),
		string(ParameterTypeInt),
		string(ParameterTypeObject),
		string(ParameterTypeSecureString),
		string(ParameterTypeString),
	}
}

func (s *ParameterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseParameterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseParameterType(input string) (*ParameterType, error) {
	vals := map[string]ParameterType{
		"array":        ParameterTypeArray,
		"bool":         ParameterTypeBool,
		"float":        ParameterTypeFloat,
		"int":          ParameterTypeInt,
		"object":       ParameterTypeObject,
		"securestring": ParameterTypeSecureString,
		"string":       ParameterTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ParameterType(input)
	return &out, nil
}

type PhoenixAuthenticationType string

const (
	PhoenixAuthenticationTypeAnonymous                    PhoenixAuthenticationType = "Anonymous"
	PhoenixAuthenticationTypeUsernameAndPassword          PhoenixAuthenticationType = "UsernameAndPassword"
	PhoenixAuthenticationTypeWindowsAzureHDInsightService PhoenixAuthenticationType = "WindowsAzureHDInsightService"
)

func PossibleValuesForPhoenixAuthenticationType() []string {
	return []string{
		string(PhoenixAuthenticationTypeAnonymous),
		string(PhoenixAuthenticationTypeUsernameAndPassword),
		string(PhoenixAuthenticationTypeWindowsAzureHDInsightService),
	}
}

func (s *PhoenixAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePhoenixAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePhoenixAuthenticationType(input string) (*PhoenixAuthenticationType, error) {
	vals := map[string]PhoenixAuthenticationType{
		"anonymous":                    PhoenixAuthenticationTypeAnonymous,
		"usernameandpassword":          PhoenixAuthenticationTypeUsernameAndPassword,
		"windowsazurehdinsightservice": PhoenixAuthenticationTypeWindowsAzureHDInsightService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PhoenixAuthenticationType(input)
	return &out, nil
}

type PrestoAuthenticationType string

const (
	PrestoAuthenticationTypeAnonymous PrestoAuthenticationType = "Anonymous"
	PrestoAuthenticationTypeLDAP      PrestoAuthenticationType = "LDAP"
)

func PossibleValuesForPrestoAuthenticationType() []string {
	return []string{
		string(PrestoAuthenticationTypeAnonymous),
		string(PrestoAuthenticationTypeLDAP),
	}
}

func (s *PrestoAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrestoAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrestoAuthenticationType(input string) (*PrestoAuthenticationType, error) {
	vals := map[string]PrestoAuthenticationType{
		"anonymous": PrestoAuthenticationTypeAnonymous,
		"ldap":      PrestoAuthenticationTypeLDAP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrestoAuthenticationType(input)
	return &out, nil
}

type RestServiceAuthenticationType string

const (
	RestServiceAuthenticationTypeAadServicePrincipal      RestServiceAuthenticationType = "AadServicePrincipal"
	RestServiceAuthenticationTypeAnonymous                RestServiceAuthenticationType = "Anonymous"
	RestServiceAuthenticationTypeBasic                    RestServiceAuthenticationType = "Basic"
	RestServiceAuthenticationTypeManagedServiceIdentity   RestServiceAuthenticationType = "ManagedServiceIdentity"
	RestServiceAuthenticationTypeOAuthTwoClientCredential RestServiceAuthenticationType = "OAuth2ClientCredential"
)

func PossibleValuesForRestServiceAuthenticationType() []string {
	return []string{
		string(RestServiceAuthenticationTypeAadServicePrincipal),
		string(RestServiceAuthenticationTypeAnonymous),
		string(RestServiceAuthenticationTypeBasic),
		string(RestServiceAuthenticationTypeManagedServiceIdentity),
		string(RestServiceAuthenticationTypeOAuthTwoClientCredential),
	}
}

func (s *RestServiceAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRestServiceAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRestServiceAuthenticationType(input string) (*RestServiceAuthenticationType, error) {
	vals := map[string]RestServiceAuthenticationType{
		"aadserviceprincipal":    RestServiceAuthenticationTypeAadServicePrincipal,
		"anonymous":              RestServiceAuthenticationTypeAnonymous,
		"basic":                  RestServiceAuthenticationTypeBasic,
		"managedserviceidentity": RestServiceAuthenticationTypeManagedServiceIdentity,
		"oauth2clientcredential": RestServiceAuthenticationTypeOAuthTwoClientCredential,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RestServiceAuthenticationType(input)
	return &out, nil
}

type SapHanaAuthenticationType string

const (
	SapHanaAuthenticationTypeBasic   SapHanaAuthenticationType = "Basic"
	SapHanaAuthenticationTypeWindows SapHanaAuthenticationType = "Windows"
)

func PossibleValuesForSapHanaAuthenticationType() []string {
	return []string{
		string(SapHanaAuthenticationTypeBasic),
		string(SapHanaAuthenticationTypeWindows),
	}
}

func (s *SapHanaAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSapHanaAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSapHanaAuthenticationType(input string) (*SapHanaAuthenticationType, error) {
	vals := map[string]SapHanaAuthenticationType{
		"basic":   SapHanaAuthenticationTypeBasic,
		"windows": SapHanaAuthenticationTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapHanaAuthenticationType(input)
	return &out, nil
}

type ServiceNowAuthenticationType string

const (
	ServiceNowAuthenticationTypeBasic    ServiceNowAuthenticationType = "Basic"
	ServiceNowAuthenticationTypeOAuthTwo ServiceNowAuthenticationType = "OAuth2"
)

func PossibleValuesForServiceNowAuthenticationType() []string {
	return []string{
		string(ServiceNowAuthenticationTypeBasic),
		string(ServiceNowAuthenticationTypeOAuthTwo),
	}
}

func (s *ServiceNowAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceNowAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceNowAuthenticationType(input string) (*ServiceNowAuthenticationType, error) {
	vals := map[string]ServiceNowAuthenticationType{
		"basic":  ServiceNowAuthenticationTypeBasic,
		"oauth2": ServiceNowAuthenticationTypeOAuthTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceNowAuthenticationType(input)
	return &out, nil
}

type SftpAuthenticationType string

const (
	SftpAuthenticationTypeBasic        SftpAuthenticationType = "Basic"
	SftpAuthenticationTypeSshPublicKey SftpAuthenticationType = "SshPublicKey"
)

func PossibleValuesForSftpAuthenticationType() []string {
	return []string{
		string(SftpAuthenticationTypeBasic),
		string(SftpAuthenticationTypeSshPublicKey),
	}
}

func (s *SftpAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSftpAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSftpAuthenticationType(input string) (*SftpAuthenticationType, error) {
	vals := map[string]SftpAuthenticationType{
		"basic":        SftpAuthenticationTypeBasic,
		"sshpublickey": SftpAuthenticationTypeSshPublicKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SftpAuthenticationType(input)
	return &out, nil
}

type SparkAuthenticationType string

const (
	SparkAuthenticationTypeAnonymous                    SparkAuthenticationType = "Anonymous"
	SparkAuthenticationTypeUsername                     SparkAuthenticationType = "Username"
	SparkAuthenticationTypeUsernameAndPassword          SparkAuthenticationType = "UsernameAndPassword"
	SparkAuthenticationTypeWindowsAzureHDInsightService SparkAuthenticationType = "WindowsAzureHDInsightService"
)

func PossibleValuesForSparkAuthenticationType() []string {
	return []string{
		string(SparkAuthenticationTypeAnonymous),
		string(SparkAuthenticationTypeUsername),
		string(SparkAuthenticationTypeUsernameAndPassword),
		string(SparkAuthenticationTypeWindowsAzureHDInsightService),
	}
}

func (s *SparkAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSparkAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSparkAuthenticationType(input string) (*SparkAuthenticationType, error) {
	vals := map[string]SparkAuthenticationType{
		"anonymous":                    SparkAuthenticationTypeAnonymous,
		"username":                     SparkAuthenticationTypeUsername,
		"usernameandpassword":          SparkAuthenticationTypeUsernameAndPassword,
		"windowsazurehdinsightservice": SparkAuthenticationTypeWindowsAzureHDInsightService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SparkAuthenticationType(input)
	return &out, nil
}

type SparkServerType string

const (
	SparkServerTypeSharkServer       SparkServerType = "SharkServer"
	SparkServerTypeSharkServerTwo    SparkServerType = "SharkServer2"
	SparkServerTypeSparkThriftServer SparkServerType = "SparkThriftServer"
)

func PossibleValuesForSparkServerType() []string {
	return []string{
		string(SparkServerTypeSharkServer),
		string(SparkServerTypeSharkServerTwo),
		string(SparkServerTypeSparkThriftServer),
	}
}

func (s *SparkServerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSparkServerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSparkServerType(input string) (*SparkServerType, error) {
	vals := map[string]SparkServerType{
		"sharkserver":       SparkServerTypeSharkServer,
		"sharkserver2":      SparkServerTypeSharkServerTwo,
		"sparkthriftserver": SparkServerTypeSparkThriftServer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SparkServerType(input)
	return &out, nil
}

type SparkThriftTransportProtocol string

const (
	SparkThriftTransportProtocolBinary SparkThriftTransportProtocol = "Binary"
	SparkThriftTransportProtocolHTTP   SparkThriftTransportProtocol = "HTTP "
	SparkThriftTransportProtocolSASL   SparkThriftTransportProtocol = "SASL"
)

func PossibleValuesForSparkThriftTransportProtocol() []string {
	return []string{
		string(SparkThriftTransportProtocolBinary),
		string(SparkThriftTransportProtocolHTTP),
		string(SparkThriftTransportProtocolSASL),
	}
}

func (s *SparkThriftTransportProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSparkThriftTransportProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSparkThriftTransportProtocol(input string) (*SparkThriftTransportProtocol, error) {
	vals := map[string]SparkThriftTransportProtocol{
		"binary": SparkThriftTransportProtocolBinary,
		"http ":  SparkThriftTransportProtocolHTTP,
		"sasl":   SparkThriftTransportProtocolSASL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SparkThriftTransportProtocol(input)
	return &out, nil
}

type SybaseAuthenticationType string

const (
	SybaseAuthenticationTypeBasic   SybaseAuthenticationType = "Basic"
	SybaseAuthenticationTypeWindows SybaseAuthenticationType = "Windows"
)

func PossibleValuesForSybaseAuthenticationType() []string {
	return []string{
		string(SybaseAuthenticationTypeBasic),
		string(SybaseAuthenticationTypeWindows),
	}
}

func (s *SybaseAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSybaseAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSybaseAuthenticationType(input string) (*SybaseAuthenticationType, error) {
	vals := map[string]SybaseAuthenticationType{
		"basic":   SybaseAuthenticationTypeBasic,
		"windows": SybaseAuthenticationTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SybaseAuthenticationType(input)
	return &out, nil
}

type TeradataAuthenticationType string

const (
	TeradataAuthenticationTypeBasic   TeradataAuthenticationType = "Basic"
	TeradataAuthenticationTypeWindows TeradataAuthenticationType = "Windows"
)

func PossibleValuesForTeradataAuthenticationType() []string {
	return []string{
		string(TeradataAuthenticationTypeBasic),
		string(TeradataAuthenticationTypeWindows),
	}
}

func (s *TeradataAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTeradataAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTeradataAuthenticationType(input string) (*TeradataAuthenticationType, error) {
	vals := map[string]TeradataAuthenticationType{
		"basic":   TeradataAuthenticationTypeBasic,
		"windows": TeradataAuthenticationTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TeradataAuthenticationType(input)
	return &out, nil
}

type Type string

const (
	TypeLinkedServiceReference Type = "LinkedServiceReference"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeLinkedServiceReference),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"linkedservicereference": TypeLinkedServiceReference,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type WebAuthenticationType string

const (
	WebAuthenticationTypeAnonymous         WebAuthenticationType = "Anonymous"
	WebAuthenticationTypeBasic             WebAuthenticationType = "Basic"
	WebAuthenticationTypeClientCertificate WebAuthenticationType = "ClientCertificate"
)

func PossibleValuesForWebAuthenticationType() []string {
	return []string{
		string(WebAuthenticationTypeAnonymous),
		string(WebAuthenticationTypeBasic),
		string(WebAuthenticationTypeClientCertificate),
	}
}

func (s *WebAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWebAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWebAuthenticationType(input string) (*WebAuthenticationType, error) {
	vals := map[string]WebAuthenticationType{
		"anonymous":         WebAuthenticationTypeAnonymous,
		"basic":             WebAuthenticationTypeBasic,
		"clientcertificate": WebAuthenticationTypeClientCertificate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WebAuthenticationType(input)
	return &out, nil
}

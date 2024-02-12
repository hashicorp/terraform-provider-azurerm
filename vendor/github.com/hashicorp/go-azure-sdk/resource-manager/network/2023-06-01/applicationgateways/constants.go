package applicationgateways

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendHealthServerHealth string

const (
	ApplicationGatewayBackendHealthServerHealthDown     ApplicationGatewayBackendHealthServerHealth = "Down"
	ApplicationGatewayBackendHealthServerHealthDraining ApplicationGatewayBackendHealthServerHealth = "Draining"
	ApplicationGatewayBackendHealthServerHealthPartial  ApplicationGatewayBackendHealthServerHealth = "Partial"
	ApplicationGatewayBackendHealthServerHealthUnknown  ApplicationGatewayBackendHealthServerHealth = "Unknown"
	ApplicationGatewayBackendHealthServerHealthUp       ApplicationGatewayBackendHealthServerHealth = "Up"
)

func PossibleValuesForApplicationGatewayBackendHealthServerHealth() []string {
	return []string{
		string(ApplicationGatewayBackendHealthServerHealthDown),
		string(ApplicationGatewayBackendHealthServerHealthDraining),
		string(ApplicationGatewayBackendHealthServerHealthPartial),
		string(ApplicationGatewayBackendHealthServerHealthUnknown),
		string(ApplicationGatewayBackendHealthServerHealthUp),
	}
}

func (s *ApplicationGatewayBackendHealthServerHealth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayBackendHealthServerHealth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayBackendHealthServerHealth(input string) (*ApplicationGatewayBackendHealthServerHealth, error) {
	vals := map[string]ApplicationGatewayBackendHealthServerHealth{
		"down":     ApplicationGatewayBackendHealthServerHealthDown,
		"draining": ApplicationGatewayBackendHealthServerHealthDraining,
		"partial":  ApplicationGatewayBackendHealthServerHealthPartial,
		"unknown":  ApplicationGatewayBackendHealthServerHealthUnknown,
		"up":       ApplicationGatewayBackendHealthServerHealthUp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayBackendHealthServerHealth(input)
	return &out, nil
}

type ApplicationGatewayClientRevocationOptions string

const (
	ApplicationGatewayClientRevocationOptionsNone ApplicationGatewayClientRevocationOptions = "None"
	ApplicationGatewayClientRevocationOptionsOCSP ApplicationGatewayClientRevocationOptions = "OCSP"
)

func PossibleValuesForApplicationGatewayClientRevocationOptions() []string {
	return []string{
		string(ApplicationGatewayClientRevocationOptionsNone),
		string(ApplicationGatewayClientRevocationOptionsOCSP),
	}
}

func (s *ApplicationGatewayClientRevocationOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayClientRevocationOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayClientRevocationOptions(input string) (*ApplicationGatewayClientRevocationOptions, error) {
	vals := map[string]ApplicationGatewayClientRevocationOptions{
		"none": ApplicationGatewayClientRevocationOptionsNone,
		"ocsp": ApplicationGatewayClientRevocationOptionsOCSP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayClientRevocationOptions(input)
	return &out, nil
}

type ApplicationGatewayCookieBasedAffinity string

const (
	ApplicationGatewayCookieBasedAffinityDisabled ApplicationGatewayCookieBasedAffinity = "Disabled"
	ApplicationGatewayCookieBasedAffinityEnabled  ApplicationGatewayCookieBasedAffinity = "Enabled"
)

func PossibleValuesForApplicationGatewayCookieBasedAffinity() []string {
	return []string{
		string(ApplicationGatewayCookieBasedAffinityDisabled),
		string(ApplicationGatewayCookieBasedAffinityEnabled),
	}
}

func (s *ApplicationGatewayCookieBasedAffinity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayCookieBasedAffinity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayCookieBasedAffinity(input string) (*ApplicationGatewayCookieBasedAffinity, error) {
	vals := map[string]ApplicationGatewayCookieBasedAffinity{
		"disabled": ApplicationGatewayCookieBasedAffinityDisabled,
		"enabled":  ApplicationGatewayCookieBasedAffinityEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayCookieBasedAffinity(input)
	return &out, nil
}

type ApplicationGatewayCustomErrorStatusCode string

const (
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveHundred   ApplicationGatewayCustomErrorStatusCode = "HttpStatus500"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroFour  ApplicationGatewayCustomErrorStatusCode = "HttpStatus504"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroThree ApplicationGatewayCustomErrorStatusCode = "HttpStatus503"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroTwo   ApplicationGatewayCustomErrorStatusCode = "HttpStatus502"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourHundred   ApplicationGatewayCustomErrorStatusCode = "HttpStatus400"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroEight ApplicationGatewayCustomErrorStatusCode = "HttpStatus408"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFive  ApplicationGatewayCustomErrorStatusCode = "HttpStatus405"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFour  ApplicationGatewayCustomErrorStatusCode = "HttpStatus404"
	ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroThree ApplicationGatewayCustomErrorStatusCode = "HttpStatus403"
)

func PossibleValuesForApplicationGatewayCustomErrorStatusCode() []string {
	return []string{
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveHundred),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroFour),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroThree),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroTwo),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourHundred),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroEight),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFive),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFour),
		string(ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroThree),
	}
}

func (s *ApplicationGatewayCustomErrorStatusCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayCustomErrorStatusCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayCustomErrorStatusCode(input string) (*ApplicationGatewayCustomErrorStatusCode, error) {
	vals := map[string]ApplicationGatewayCustomErrorStatusCode{
		"httpstatus500": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveHundred,
		"httpstatus504": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroFour,
		"httpstatus503": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroThree,
		"httpstatus502": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFiveZeroTwo,
		"httpstatus400": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourHundred,
		"httpstatus408": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroEight,
		"httpstatus405": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFive,
		"httpstatus404": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroFour,
		"httpstatus403": ApplicationGatewayCustomErrorStatusCodeHTTPStatusFourZeroThree,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayCustomErrorStatusCode(input)
	return &out, nil
}

type ApplicationGatewayFirewallMode string

const (
	ApplicationGatewayFirewallModeDetection  ApplicationGatewayFirewallMode = "Detection"
	ApplicationGatewayFirewallModePrevention ApplicationGatewayFirewallMode = "Prevention"
)

func PossibleValuesForApplicationGatewayFirewallMode() []string {
	return []string{
		string(ApplicationGatewayFirewallModeDetection),
		string(ApplicationGatewayFirewallModePrevention),
	}
}

func (s *ApplicationGatewayFirewallMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayFirewallMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayFirewallMode(input string) (*ApplicationGatewayFirewallMode, error) {
	vals := map[string]ApplicationGatewayFirewallMode{
		"detection":  ApplicationGatewayFirewallModeDetection,
		"prevention": ApplicationGatewayFirewallModePrevention,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayFirewallMode(input)
	return &out, nil
}

type ApplicationGatewayLoadDistributionAlgorithm string

const (
	ApplicationGatewayLoadDistributionAlgorithmIPHash           ApplicationGatewayLoadDistributionAlgorithm = "IpHash"
	ApplicationGatewayLoadDistributionAlgorithmLeastConnections ApplicationGatewayLoadDistributionAlgorithm = "LeastConnections"
	ApplicationGatewayLoadDistributionAlgorithmRoundRobin       ApplicationGatewayLoadDistributionAlgorithm = "RoundRobin"
)

func PossibleValuesForApplicationGatewayLoadDistributionAlgorithm() []string {
	return []string{
		string(ApplicationGatewayLoadDistributionAlgorithmIPHash),
		string(ApplicationGatewayLoadDistributionAlgorithmLeastConnections),
		string(ApplicationGatewayLoadDistributionAlgorithmRoundRobin),
	}
}

func (s *ApplicationGatewayLoadDistributionAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayLoadDistributionAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayLoadDistributionAlgorithm(input string) (*ApplicationGatewayLoadDistributionAlgorithm, error) {
	vals := map[string]ApplicationGatewayLoadDistributionAlgorithm{
		"iphash":           ApplicationGatewayLoadDistributionAlgorithmIPHash,
		"leastconnections": ApplicationGatewayLoadDistributionAlgorithmLeastConnections,
		"roundrobin":       ApplicationGatewayLoadDistributionAlgorithmRoundRobin,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayLoadDistributionAlgorithm(input)
	return &out, nil
}

type ApplicationGatewayOperationalState string

const (
	ApplicationGatewayOperationalStateRunning  ApplicationGatewayOperationalState = "Running"
	ApplicationGatewayOperationalStateStarting ApplicationGatewayOperationalState = "Starting"
	ApplicationGatewayOperationalStateStopped  ApplicationGatewayOperationalState = "Stopped"
	ApplicationGatewayOperationalStateStopping ApplicationGatewayOperationalState = "Stopping"
)

func PossibleValuesForApplicationGatewayOperationalState() []string {
	return []string{
		string(ApplicationGatewayOperationalStateRunning),
		string(ApplicationGatewayOperationalStateStarting),
		string(ApplicationGatewayOperationalStateStopped),
		string(ApplicationGatewayOperationalStateStopping),
	}
}

func (s *ApplicationGatewayOperationalState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayOperationalState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayOperationalState(input string) (*ApplicationGatewayOperationalState, error) {
	vals := map[string]ApplicationGatewayOperationalState{
		"running":  ApplicationGatewayOperationalStateRunning,
		"starting": ApplicationGatewayOperationalStateStarting,
		"stopped":  ApplicationGatewayOperationalStateStopped,
		"stopping": ApplicationGatewayOperationalStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayOperationalState(input)
	return &out, nil
}

type ApplicationGatewayProtocol string

const (
	ApplicationGatewayProtocolHTTP  ApplicationGatewayProtocol = "Http"
	ApplicationGatewayProtocolHTTPS ApplicationGatewayProtocol = "Https"
	ApplicationGatewayProtocolTcp   ApplicationGatewayProtocol = "Tcp"
	ApplicationGatewayProtocolTls   ApplicationGatewayProtocol = "Tls"
)

func PossibleValuesForApplicationGatewayProtocol() []string {
	return []string{
		string(ApplicationGatewayProtocolHTTP),
		string(ApplicationGatewayProtocolHTTPS),
		string(ApplicationGatewayProtocolTcp),
		string(ApplicationGatewayProtocolTls),
	}
}

func (s *ApplicationGatewayProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayProtocol(input string) (*ApplicationGatewayProtocol, error) {
	vals := map[string]ApplicationGatewayProtocol{
		"http":  ApplicationGatewayProtocolHTTP,
		"https": ApplicationGatewayProtocolHTTPS,
		"tcp":   ApplicationGatewayProtocolTcp,
		"tls":   ApplicationGatewayProtocolTls,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayProtocol(input)
	return &out, nil
}

type ApplicationGatewayRedirectType string

const (
	ApplicationGatewayRedirectTypeFound     ApplicationGatewayRedirectType = "Found"
	ApplicationGatewayRedirectTypePermanent ApplicationGatewayRedirectType = "Permanent"
	ApplicationGatewayRedirectTypeSeeOther  ApplicationGatewayRedirectType = "SeeOther"
	ApplicationGatewayRedirectTypeTemporary ApplicationGatewayRedirectType = "Temporary"
)

func PossibleValuesForApplicationGatewayRedirectType() []string {
	return []string{
		string(ApplicationGatewayRedirectTypeFound),
		string(ApplicationGatewayRedirectTypePermanent),
		string(ApplicationGatewayRedirectTypeSeeOther),
		string(ApplicationGatewayRedirectTypeTemporary),
	}
}

func (s *ApplicationGatewayRedirectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayRedirectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayRedirectType(input string) (*ApplicationGatewayRedirectType, error) {
	vals := map[string]ApplicationGatewayRedirectType{
		"found":     ApplicationGatewayRedirectTypeFound,
		"permanent": ApplicationGatewayRedirectTypePermanent,
		"seeother":  ApplicationGatewayRedirectTypeSeeOther,
		"temporary": ApplicationGatewayRedirectTypeTemporary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayRedirectType(input)
	return &out, nil
}

type ApplicationGatewayRequestRoutingRuleType string

const (
	ApplicationGatewayRequestRoutingRuleTypeBasic            ApplicationGatewayRequestRoutingRuleType = "Basic"
	ApplicationGatewayRequestRoutingRuleTypePathBasedRouting ApplicationGatewayRequestRoutingRuleType = "PathBasedRouting"
)

func PossibleValuesForApplicationGatewayRequestRoutingRuleType() []string {
	return []string{
		string(ApplicationGatewayRequestRoutingRuleTypeBasic),
		string(ApplicationGatewayRequestRoutingRuleTypePathBasedRouting),
	}
}

func (s *ApplicationGatewayRequestRoutingRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayRequestRoutingRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayRequestRoutingRuleType(input string) (*ApplicationGatewayRequestRoutingRuleType, error) {
	vals := map[string]ApplicationGatewayRequestRoutingRuleType{
		"basic":            ApplicationGatewayRequestRoutingRuleTypeBasic,
		"pathbasedrouting": ApplicationGatewayRequestRoutingRuleTypePathBasedRouting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayRequestRoutingRuleType(input)
	return &out, nil
}

type ApplicationGatewaySkuName string

const (
	ApplicationGatewaySkuNameBasic          ApplicationGatewaySkuName = "Basic"
	ApplicationGatewaySkuNameStandardLarge  ApplicationGatewaySkuName = "Standard_Large"
	ApplicationGatewaySkuNameStandardMedium ApplicationGatewaySkuName = "Standard_Medium"
	ApplicationGatewaySkuNameStandardSmall  ApplicationGatewaySkuName = "Standard_Small"
	ApplicationGatewaySkuNameStandardVTwo   ApplicationGatewaySkuName = "Standard_v2"
	ApplicationGatewaySkuNameWAFLarge       ApplicationGatewaySkuName = "WAF_Large"
	ApplicationGatewaySkuNameWAFMedium      ApplicationGatewaySkuName = "WAF_Medium"
	ApplicationGatewaySkuNameWAFVTwo        ApplicationGatewaySkuName = "WAF_v2"
)

func PossibleValuesForApplicationGatewaySkuName() []string {
	return []string{
		string(ApplicationGatewaySkuNameBasic),
		string(ApplicationGatewaySkuNameStandardLarge),
		string(ApplicationGatewaySkuNameStandardMedium),
		string(ApplicationGatewaySkuNameStandardSmall),
		string(ApplicationGatewaySkuNameStandardVTwo),
		string(ApplicationGatewaySkuNameWAFLarge),
		string(ApplicationGatewaySkuNameWAFMedium),
		string(ApplicationGatewaySkuNameWAFVTwo),
	}
}

func (s *ApplicationGatewaySkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewaySkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewaySkuName(input string) (*ApplicationGatewaySkuName, error) {
	vals := map[string]ApplicationGatewaySkuName{
		"basic":           ApplicationGatewaySkuNameBasic,
		"standard_large":  ApplicationGatewaySkuNameStandardLarge,
		"standard_medium": ApplicationGatewaySkuNameStandardMedium,
		"standard_small":  ApplicationGatewaySkuNameStandardSmall,
		"standard_v2":     ApplicationGatewaySkuNameStandardVTwo,
		"waf_large":       ApplicationGatewaySkuNameWAFLarge,
		"waf_medium":      ApplicationGatewaySkuNameWAFMedium,
		"waf_v2":          ApplicationGatewaySkuNameWAFVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewaySkuName(input)
	return &out, nil
}

type ApplicationGatewaySslCipherSuite string

const (
	ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHA                  ApplicationGatewaySslCipherSuite = "TLS_DHE_DSS_WITH_AES_128_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHATwoFiveSix        ApplicationGatewaySslCipherSuite = "TLS_DHE_DSS_WITH_AES_128_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHA                   ApplicationGatewaySslCipherSuite = "TLS_DHE_DSS_WITH_AES_256_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHATwoFiveSix         ApplicationGatewaySslCipherSuite = "TLS_DHE_DSS_WITH_AES_256_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHThreeDESEDECBCSHA                     ApplicationGatewaySslCipherSuite = "TLS_DHE_DSS_WITH_3DES_EDE_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightCBCSHA                  ApplicationGatewaySslCipherSuite = "TLS_DHE_RSA_WITH_AES_128_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix        ApplicationGatewaySslCipherSuite = "TLS_DHE_RSA_WITH_AES_128_GCM_SHA256"
	ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixCBCSHA                   ApplicationGatewaySslCipherSuite = "TLS_DHE_RSA_WITH_AES_256_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour     ApplicationGatewaySslCipherSuite = "TLS_DHE_RSA_WITH_AES_256_GCM_SHA384"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHA              ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHATwoFiveSix    ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightGCMSHATwoFiveSix    ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHA               ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHAThreeEightFour ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384"
	ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixGCMSHAThreeEightFour ApplicationGatewaySslCipherSuite = "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHA                ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHATwoFiveSix      ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix      ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHA                 ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHAThreeEightFour   ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384"
	ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour   ApplicationGatewaySslCipherSuite = "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHA                     ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_128_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHATwoFiveSix           ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_128_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightGCMSHATwoFiveSix           ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_128_GCM_SHA256"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHA                      ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_256_CBC_SHA"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHATwoFiveSix            ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_256_CBC_SHA256"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixGCMSHAThreeEightFour        ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_AES_256_GCM_SHA384"
	ApplicationGatewaySslCipherSuiteTLSRSAWITHThreeDESEDECBCSHA                        ApplicationGatewaySslCipherSuite = "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
)

func PossibleValuesForApplicationGatewaySslCipherSuite() []string {
	return []string{
		string(ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHThreeDESEDECBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightGCMSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixGCMSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightGCMSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHA),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHATwoFiveSix),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixGCMSHAThreeEightFour),
		string(ApplicationGatewaySslCipherSuiteTLSRSAWITHThreeDESEDECBCSHA),
	}
}

func (s *ApplicationGatewaySslCipherSuite) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewaySslCipherSuite(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewaySslCipherSuite(input string) (*ApplicationGatewaySslCipherSuite, error) {
	vals := map[string]ApplicationGatewaySslCipherSuite{
		"tls_dhe_dss_with_aes_128_cbc_sha":        ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHA,
		"tls_dhe_dss_with_aes_128_cbc_sha256":     ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESOneTwoEightCBCSHATwoFiveSix,
		"tls_dhe_dss_with_aes_256_cbc_sha":        ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHA,
		"tls_dhe_dss_with_aes_256_cbc_sha256":     ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHAESTwoFiveSixCBCSHATwoFiveSix,
		"tls_dhe_dss_with_3des_ede_cbc_sha":       ApplicationGatewaySslCipherSuiteTLSDHEDSSWITHThreeDESEDECBCSHA,
		"tls_dhe_rsa_with_aes_128_cbc_sha":        ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightCBCSHA,
		"tls_dhe_rsa_with_aes_128_gcm_sha256":     ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix,
		"tls_dhe_rsa_with_aes_256_cbc_sha":        ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixCBCSHA,
		"tls_dhe_rsa_with_aes_256_gcm_sha384":     ApplicationGatewaySslCipherSuiteTLSDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour,
		"tls_ecdhe_ecdsa_with_aes_128_cbc_sha":    ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHA,
		"tls_ecdhe_ecdsa_with_aes_128_cbc_sha256": ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightCBCSHATwoFiveSix,
		"tls_ecdhe_ecdsa_with_aes_128_gcm_sha256": ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESOneTwoEightGCMSHATwoFiveSix,
		"tls_ecdhe_ecdsa_with_aes_256_cbc_sha":    ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHA,
		"tls_ecdhe_ecdsa_with_aes_256_cbc_sha384": ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixCBCSHAThreeEightFour,
		"tls_ecdhe_ecdsa_with_aes_256_gcm_sha384": ApplicationGatewaySslCipherSuiteTLSECDHEECDSAWITHAESTwoFiveSixGCMSHAThreeEightFour,
		"tls_ecdhe_rsa_with_aes_128_cbc_sha":      ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHA,
		"tls_ecdhe_rsa_with_aes_128_cbc_sha256":   ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightCBCSHATwoFiveSix,
		"tls_ecdhe_rsa_with_aes_128_gcm_sha256":   ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESOneTwoEightGCMSHATwoFiveSix,
		"tls_ecdhe_rsa_with_aes_256_cbc_sha":      ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHA,
		"tls_ecdhe_rsa_with_aes_256_cbc_sha384":   ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixCBCSHAThreeEightFour,
		"tls_ecdhe_rsa_with_aes_256_gcm_sha384":   ApplicationGatewaySslCipherSuiteTLSECDHERSAWITHAESTwoFiveSixGCMSHAThreeEightFour,
		"tls_rsa_with_aes_128_cbc_sha":            ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHA,
		"tls_rsa_with_aes_128_cbc_sha256":         ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightCBCSHATwoFiveSix,
		"tls_rsa_with_aes_128_gcm_sha256":         ApplicationGatewaySslCipherSuiteTLSRSAWITHAESOneTwoEightGCMSHATwoFiveSix,
		"tls_rsa_with_aes_256_cbc_sha":            ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHA,
		"tls_rsa_with_aes_256_cbc_sha256":         ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixCBCSHATwoFiveSix,
		"tls_rsa_with_aes_256_gcm_sha384":         ApplicationGatewaySslCipherSuiteTLSRSAWITHAESTwoFiveSixGCMSHAThreeEightFour,
		"tls_rsa_with_3des_ede_cbc_sha":           ApplicationGatewaySslCipherSuiteTLSRSAWITHThreeDESEDECBCSHA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewaySslCipherSuite(input)
	return &out, nil
}

type ApplicationGatewaySslPolicyName string

const (
	ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneFiveZeroFiveZeroOne   ApplicationGatewaySslPolicyName = "AppGwSslPolicy20150501"
	ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOne  ApplicationGatewaySslPolicyName = "AppGwSslPolicy20170401"
	ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOneS ApplicationGatewaySslPolicyName = "AppGwSslPolicy20170401S"
	ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOne     ApplicationGatewaySslPolicyName = "AppGwSslPolicy20220101"
	ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOneS    ApplicationGatewaySslPolicyName = "AppGwSslPolicy20220101S"
)

func PossibleValuesForApplicationGatewaySslPolicyName() []string {
	return []string{
		string(ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneFiveZeroFiveZeroOne),
		string(ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOne),
		string(ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOneS),
		string(ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOne),
		string(ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOneS),
	}
}

func (s *ApplicationGatewaySslPolicyName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewaySslPolicyName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewaySslPolicyName(input string) (*ApplicationGatewaySslPolicyName, error) {
	vals := map[string]ApplicationGatewaySslPolicyName{
		"appgwsslpolicy20150501":  ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneFiveZeroFiveZeroOne,
		"appgwsslpolicy20170401":  ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOne,
		"appgwsslpolicy20170401s": ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroOneSevenZeroFourZeroOneS,
		"appgwsslpolicy20220101":  ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOne,
		"appgwsslpolicy20220101s": ApplicationGatewaySslPolicyNameAppGwSslPolicyTwoZeroTwoTwoZeroOneZeroOneS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewaySslPolicyName(input)
	return &out, nil
}

type ApplicationGatewaySslPolicyType string

const (
	ApplicationGatewaySslPolicyTypeCustom     ApplicationGatewaySslPolicyType = "Custom"
	ApplicationGatewaySslPolicyTypeCustomVTwo ApplicationGatewaySslPolicyType = "CustomV2"
	ApplicationGatewaySslPolicyTypePredefined ApplicationGatewaySslPolicyType = "Predefined"
)

func PossibleValuesForApplicationGatewaySslPolicyType() []string {
	return []string{
		string(ApplicationGatewaySslPolicyTypeCustom),
		string(ApplicationGatewaySslPolicyTypeCustomVTwo),
		string(ApplicationGatewaySslPolicyTypePredefined),
	}
}

func (s *ApplicationGatewaySslPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewaySslPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewaySslPolicyType(input string) (*ApplicationGatewaySslPolicyType, error) {
	vals := map[string]ApplicationGatewaySslPolicyType{
		"custom":     ApplicationGatewaySslPolicyTypeCustom,
		"customv2":   ApplicationGatewaySslPolicyTypeCustomVTwo,
		"predefined": ApplicationGatewaySslPolicyTypePredefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewaySslPolicyType(input)
	return &out, nil
}

type ApplicationGatewaySslProtocol string

const (
	ApplicationGatewaySslProtocolTLSvOneOne   ApplicationGatewaySslProtocol = "TLSv1_1"
	ApplicationGatewaySslProtocolTLSvOneThree ApplicationGatewaySslProtocol = "TLSv1_3"
	ApplicationGatewaySslProtocolTLSvOneTwo   ApplicationGatewaySslProtocol = "TLSv1_2"
	ApplicationGatewaySslProtocolTLSvOneZero  ApplicationGatewaySslProtocol = "TLSv1_0"
)

func PossibleValuesForApplicationGatewaySslProtocol() []string {
	return []string{
		string(ApplicationGatewaySslProtocolTLSvOneOne),
		string(ApplicationGatewaySslProtocolTLSvOneThree),
		string(ApplicationGatewaySslProtocolTLSvOneTwo),
		string(ApplicationGatewaySslProtocolTLSvOneZero),
	}
}

func (s *ApplicationGatewaySslProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewaySslProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewaySslProtocol(input string) (*ApplicationGatewaySslProtocol, error) {
	vals := map[string]ApplicationGatewaySslProtocol{
		"tlsv1_1": ApplicationGatewaySslProtocolTLSvOneOne,
		"tlsv1_3": ApplicationGatewaySslProtocolTLSvOneThree,
		"tlsv1_2": ApplicationGatewaySslProtocolTLSvOneTwo,
		"tlsv1_0": ApplicationGatewaySslProtocolTLSvOneZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewaySslProtocol(input)
	return &out, nil
}

type ApplicationGatewayTier string

const (
	ApplicationGatewayTierBasic        ApplicationGatewayTier = "Basic"
	ApplicationGatewayTierStandard     ApplicationGatewayTier = "Standard"
	ApplicationGatewayTierStandardVTwo ApplicationGatewayTier = "Standard_v2"
	ApplicationGatewayTierWAF          ApplicationGatewayTier = "WAF"
	ApplicationGatewayTierWAFVTwo      ApplicationGatewayTier = "WAF_v2"
)

func PossibleValuesForApplicationGatewayTier() []string {
	return []string{
		string(ApplicationGatewayTierBasic),
		string(ApplicationGatewayTierStandard),
		string(ApplicationGatewayTierStandardVTwo),
		string(ApplicationGatewayTierWAF),
		string(ApplicationGatewayTierWAFVTwo),
	}
}

func (s *ApplicationGatewayTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayTier(input string) (*ApplicationGatewayTier, error) {
	vals := map[string]ApplicationGatewayTier{
		"basic":       ApplicationGatewayTierBasic,
		"standard":    ApplicationGatewayTierStandard,
		"standard_v2": ApplicationGatewayTierStandardVTwo,
		"waf":         ApplicationGatewayTierWAF,
		"waf_v2":      ApplicationGatewayTierWAFVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayTier(input)
	return &out, nil
}

type ApplicationGatewayTierTypes string

const (
	ApplicationGatewayTierTypesStandard     ApplicationGatewayTierTypes = "Standard"
	ApplicationGatewayTierTypesStandardVTwo ApplicationGatewayTierTypes = "Standard_v2"
	ApplicationGatewayTierTypesWAF          ApplicationGatewayTierTypes = "WAF"
	ApplicationGatewayTierTypesWAFVTwo      ApplicationGatewayTierTypes = "WAF_v2"
)

func PossibleValuesForApplicationGatewayTierTypes() []string {
	return []string{
		string(ApplicationGatewayTierTypesStandard),
		string(ApplicationGatewayTierTypesStandardVTwo),
		string(ApplicationGatewayTierTypesWAF),
		string(ApplicationGatewayTierTypesWAFVTwo),
	}
}

func (s *ApplicationGatewayTierTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayTierTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayTierTypes(input string) (*ApplicationGatewayTierTypes, error) {
	vals := map[string]ApplicationGatewayTierTypes{
		"standard":    ApplicationGatewayTierTypesStandard,
		"standard_v2": ApplicationGatewayTierTypesStandardVTwo,
		"waf":         ApplicationGatewayTierTypesWAF,
		"waf_v2":      ApplicationGatewayTierTypesWAFVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayTierTypes(input)
	return &out, nil
}

type ApplicationGatewayWafRuleActionTypes string

const (
	ApplicationGatewayWafRuleActionTypesAllow          ApplicationGatewayWafRuleActionTypes = "Allow"
	ApplicationGatewayWafRuleActionTypesAnomalyScoring ApplicationGatewayWafRuleActionTypes = "AnomalyScoring"
	ApplicationGatewayWafRuleActionTypesBlock          ApplicationGatewayWafRuleActionTypes = "Block"
	ApplicationGatewayWafRuleActionTypesLog            ApplicationGatewayWafRuleActionTypes = "Log"
	ApplicationGatewayWafRuleActionTypesNone           ApplicationGatewayWafRuleActionTypes = "None"
)

func PossibleValuesForApplicationGatewayWafRuleActionTypes() []string {
	return []string{
		string(ApplicationGatewayWafRuleActionTypesAllow),
		string(ApplicationGatewayWafRuleActionTypesAnomalyScoring),
		string(ApplicationGatewayWafRuleActionTypesBlock),
		string(ApplicationGatewayWafRuleActionTypesLog),
		string(ApplicationGatewayWafRuleActionTypesNone),
	}
}

func (s *ApplicationGatewayWafRuleActionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayWafRuleActionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayWafRuleActionTypes(input string) (*ApplicationGatewayWafRuleActionTypes, error) {
	vals := map[string]ApplicationGatewayWafRuleActionTypes{
		"allow":          ApplicationGatewayWafRuleActionTypesAllow,
		"anomalyscoring": ApplicationGatewayWafRuleActionTypesAnomalyScoring,
		"block":          ApplicationGatewayWafRuleActionTypesBlock,
		"log":            ApplicationGatewayWafRuleActionTypesLog,
		"none":           ApplicationGatewayWafRuleActionTypesNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayWafRuleActionTypes(input)
	return &out, nil
}

type ApplicationGatewayWafRuleStateTypes string

const (
	ApplicationGatewayWafRuleStateTypesDisabled ApplicationGatewayWafRuleStateTypes = "Disabled"
	ApplicationGatewayWafRuleStateTypesEnabled  ApplicationGatewayWafRuleStateTypes = "Enabled"
)

func PossibleValuesForApplicationGatewayWafRuleStateTypes() []string {
	return []string{
		string(ApplicationGatewayWafRuleStateTypesDisabled),
		string(ApplicationGatewayWafRuleStateTypesEnabled),
	}
}

func (s *ApplicationGatewayWafRuleStateTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationGatewayWafRuleStateTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationGatewayWafRuleStateTypes(input string) (*ApplicationGatewayWafRuleStateTypes, error) {
	vals := map[string]ApplicationGatewayWafRuleStateTypes{
		"disabled": ApplicationGatewayWafRuleStateTypesDisabled,
		"enabled":  ApplicationGatewayWafRuleStateTypesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationGatewayWafRuleStateTypes(input)
	return &out, nil
}

type DdosSettingsProtectionMode string

const (
	DdosSettingsProtectionModeDisabled                DdosSettingsProtectionMode = "Disabled"
	DdosSettingsProtectionModeEnabled                 DdosSettingsProtectionMode = "Enabled"
	DdosSettingsProtectionModeVirtualNetworkInherited DdosSettingsProtectionMode = "VirtualNetworkInherited"
)

func PossibleValuesForDdosSettingsProtectionMode() []string {
	return []string{
		string(DdosSettingsProtectionModeDisabled),
		string(DdosSettingsProtectionModeEnabled),
		string(DdosSettingsProtectionModeVirtualNetworkInherited),
	}
}

func (s *DdosSettingsProtectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosSettingsProtectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosSettingsProtectionMode(input string) (*DdosSettingsProtectionMode, error) {
	vals := map[string]DdosSettingsProtectionMode{
		"disabled":                DdosSettingsProtectionModeDisabled,
		"enabled":                 DdosSettingsProtectionModeEnabled,
		"virtualnetworkinherited": DdosSettingsProtectionModeVirtualNetworkInherited,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosSettingsProtectionMode(input)
	return &out, nil
}

type DeleteOptions string

const (
	DeleteOptionsDelete DeleteOptions = "Delete"
	DeleteOptionsDetach DeleteOptions = "Detach"
)

func PossibleValuesForDeleteOptions() []string {
	return []string{
		string(DeleteOptionsDelete),
		string(DeleteOptionsDetach),
	}
}

func (s *DeleteOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteOptions(input string) (*DeleteOptions, error) {
	vals := map[string]DeleteOptions{
		"delete": DeleteOptionsDelete,
		"detach": DeleteOptionsDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteOptions(input)
	return &out, nil
}

type FlowLogFormatType string

const (
	FlowLogFormatTypeJSON FlowLogFormatType = "JSON"
)

func PossibleValuesForFlowLogFormatType() []string {
	return []string{
		string(FlowLogFormatTypeJSON),
	}
}

func (s *FlowLogFormatType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFlowLogFormatType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFlowLogFormatType(input string) (*FlowLogFormatType, error) {
	vals := map[string]FlowLogFormatType{
		"json": FlowLogFormatTypeJSON,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FlowLogFormatType(input)
	return &out, nil
}

type GatewayLoadBalancerTunnelInterfaceType string

const (
	GatewayLoadBalancerTunnelInterfaceTypeExternal GatewayLoadBalancerTunnelInterfaceType = "External"
	GatewayLoadBalancerTunnelInterfaceTypeInternal GatewayLoadBalancerTunnelInterfaceType = "Internal"
	GatewayLoadBalancerTunnelInterfaceTypeNone     GatewayLoadBalancerTunnelInterfaceType = "None"
)

func PossibleValuesForGatewayLoadBalancerTunnelInterfaceType() []string {
	return []string{
		string(GatewayLoadBalancerTunnelInterfaceTypeExternal),
		string(GatewayLoadBalancerTunnelInterfaceTypeInternal),
		string(GatewayLoadBalancerTunnelInterfaceTypeNone),
	}
}

func (s *GatewayLoadBalancerTunnelInterfaceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayLoadBalancerTunnelInterfaceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayLoadBalancerTunnelInterfaceType(input string) (*GatewayLoadBalancerTunnelInterfaceType, error) {
	vals := map[string]GatewayLoadBalancerTunnelInterfaceType{
		"external": GatewayLoadBalancerTunnelInterfaceTypeExternal,
		"internal": GatewayLoadBalancerTunnelInterfaceTypeInternal,
		"none":     GatewayLoadBalancerTunnelInterfaceTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayLoadBalancerTunnelInterfaceType(input)
	return &out, nil
}

type GatewayLoadBalancerTunnelProtocol string

const (
	GatewayLoadBalancerTunnelProtocolNative GatewayLoadBalancerTunnelProtocol = "Native"
	GatewayLoadBalancerTunnelProtocolNone   GatewayLoadBalancerTunnelProtocol = "None"
	GatewayLoadBalancerTunnelProtocolVXLAN  GatewayLoadBalancerTunnelProtocol = "VXLAN"
)

func PossibleValuesForGatewayLoadBalancerTunnelProtocol() []string {
	return []string{
		string(GatewayLoadBalancerTunnelProtocolNative),
		string(GatewayLoadBalancerTunnelProtocolNone),
		string(GatewayLoadBalancerTunnelProtocolVXLAN),
	}
}

func (s *GatewayLoadBalancerTunnelProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayLoadBalancerTunnelProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayLoadBalancerTunnelProtocol(input string) (*GatewayLoadBalancerTunnelProtocol, error) {
	vals := map[string]GatewayLoadBalancerTunnelProtocol{
		"native": GatewayLoadBalancerTunnelProtocolNative,
		"none":   GatewayLoadBalancerTunnelProtocolNone,
		"vxlan":  GatewayLoadBalancerTunnelProtocolVXLAN,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayLoadBalancerTunnelProtocol(input)
	return &out, nil
}

type IPAllocationMethod string

const (
	IPAllocationMethodDynamic IPAllocationMethod = "Dynamic"
	IPAllocationMethodStatic  IPAllocationMethod = "Static"
)

func PossibleValuesForIPAllocationMethod() []string {
	return []string{
		string(IPAllocationMethodDynamic),
		string(IPAllocationMethodStatic),
	}
}

func (s *IPAllocationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPAllocationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPAllocationMethod(input string) (*IPAllocationMethod, error) {
	vals := map[string]IPAllocationMethod{
		"dynamic": IPAllocationMethodDynamic,
		"static":  IPAllocationMethodStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAllocationMethod(input)
	return &out, nil
}

type IPVersion string

const (
	IPVersionIPvFour IPVersion = "IPv4"
	IPVersionIPvSix  IPVersion = "IPv6"
)

func PossibleValuesForIPVersion() []string {
	return []string{
		string(IPVersionIPvFour),
		string(IPVersionIPvSix),
	}
}

func (s *IPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPVersion(input string) (*IPVersion, error) {
	vals := map[string]IPVersion{
		"ipv4": IPVersionIPvFour,
		"ipv6": IPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPVersion(input)
	return &out, nil
}

type LoadBalancerBackendAddressAdminState string

const (
	LoadBalancerBackendAddressAdminStateDown LoadBalancerBackendAddressAdminState = "Down"
	LoadBalancerBackendAddressAdminStateNone LoadBalancerBackendAddressAdminState = "None"
	LoadBalancerBackendAddressAdminStateUp   LoadBalancerBackendAddressAdminState = "Up"
)

func PossibleValuesForLoadBalancerBackendAddressAdminState() []string {
	return []string{
		string(LoadBalancerBackendAddressAdminStateDown),
		string(LoadBalancerBackendAddressAdminStateNone),
		string(LoadBalancerBackendAddressAdminStateUp),
	}
}

func (s *LoadBalancerBackendAddressAdminState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerBackendAddressAdminState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerBackendAddressAdminState(input string) (*LoadBalancerBackendAddressAdminState, error) {
	vals := map[string]LoadBalancerBackendAddressAdminState{
		"down": LoadBalancerBackendAddressAdminStateDown,
		"none": LoadBalancerBackendAddressAdminStateNone,
		"up":   LoadBalancerBackendAddressAdminStateUp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerBackendAddressAdminState(input)
	return &out, nil
}

type NatGatewaySkuName string

const (
	NatGatewaySkuNameStandard NatGatewaySkuName = "Standard"
)

func PossibleValuesForNatGatewaySkuName() []string {
	return []string{
		string(NatGatewaySkuNameStandard),
	}
}

func (s *NatGatewaySkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNatGatewaySkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNatGatewaySkuName(input string) (*NatGatewaySkuName, error) {
	vals := map[string]NatGatewaySkuName{
		"standard": NatGatewaySkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NatGatewaySkuName(input)
	return &out, nil
}

type NetworkInterfaceAuxiliaryMode string

const (
	NetworkInterfaceAuxiliaryModeAcceleratedConnections NetworkInterfaceAuxiliaryMode = "AcceleratedConnections"
	NetworkInterfaceAuxiliaryModeFloating               NetworkInterfaceAuxiliaryMode = "Floating"
	NetworkInterfaceAuxiliaryModeMaxConnections         NetworkInterfaceAuxiliaryMode = "MaxConnections"
	NetworkInterfaceAuxiliaryModeNone                   NetworkInterfaceAuxiliaryMode = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliaryMode() []string {
	return []string{
		string(NetworkInterfaceAuxiliaryModeAcceleratedConnections),
		string(NetworkInterfaceAuxiliaryModeFloating),
		string(NetworkInterfaceAuxiliaryModeMaxConnections),
		string(NetworkInterfaceAuxiliaryModeNone),
	}
}

func (s *NetworkInterfaceAuxiliaryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliaryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliaryMode(input string) (*NetworkInterfaceAuxiliaryMode, error) {
	vals := map[string]NetworkInterfaceAuxiliaryMode{
		"acceleratedconnections": NetworkInterfaceAuxiliaryModeAcceleratedConnections,
		"floating":               NetworkInterfaceAuxiliaryModeFloating,
		"maxconnections":         NetworkInterfaceAuxiliaryModeMaxConnections,
		"none":                   NetworkInterfaceAuxiliaryModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliaryMode(input)
	return &out, nil
}

type NetworkInterfaceAuxiliarySku string

const (
	NetworkInterfaceAuxiliarySkuAEight NetworkInterfaceAuxiliarySku = "A8"
	NetworkInterfaceAuxiliarySkuAFour  NetworkInterfaceAuxiliarySku = "A4"
	NetworkInterfaceAuxiliarySkuAOne   NetworkInterfaceAuxiliarySku = "A1"
	NetworkInterfaceAuxiliarySkuATwo   NetworkInterfaceAuxiliarySku = "A2"
	NetworkInterfaceAuxiliarySkuNone   NetworkInterfaceAuxiliarySku = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliarySku() []string {
	return []string{
		string(NetworkInterfaceAuxiliarySkuAEight),
		string(NetworkInterfaceAuxiliarySkuAFour),
		string(NetworkInterfaceAuxiliarySkuAOne),
		string(NetworkInterfaceAuxiliarySkuATwo),
		string(NetworkInterfaceAuxiliarySkuNone),
	}
}

func (s *NetworkInterfaceAuxiliarySku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliarySku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliarySku(input string) (*NetworkInterfaceAuxiliarySku, error) {
	vals := map[string]NetworkInterfaceAuxiliarySku{
		"a8":   NetworkInterfaceAuxiliarySkuAEight,
		"a4":   NetworkInterfaceAuxiliarySkuAFour,
		"a1":   NetworkInterfaceAuxiliarySkuAOne,
		"a2":   NetworkInterfaceAuxiliarySkuATwo,
		"none": NetworkInterfaceAuxiliarySkuNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliarySku(input)
	return &out, nil
}

type NetworkInterfaceMigrationPhase string

const (
	NetworkInterfaceMigrationPhaseAbort     NetworkInterfaceMigrationPhase = "Abort"
	NetworkInterfaceMigrationPhaseCommit    NetworkInterfaceMigrationPhase = "Commit"
	NetworkInterfaceMigrationPhaseCommitted NetworkInterfaceMigrationPhase = "Committed"
	NetworkInterfaceMigrationPhaseNone      NetworkInterfaceMigrationPhase = "None"
	NetworkInterfaceMigrationPhasePrepare   NetworkInterfaceMigrationPhase = "Prepare"
)

func PossibleValuesForNetworkInterfaceMigrationPhase() []string {
	return []string{
		string(NetworkInterfaceMigrationPhaseAbort),
		string(NetworkInterfaceMigrationPhaseCommit),
		string(NetworkInterfaceMigrationPhaseCommitted),
		string(NetworkInterfaceMigrationPhaseNone),
		string(NetworkInterfaceMigrationPhasePrepare),
	}
}

func (s *NetworkInterfaceMigrationPhase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceMigrationPhase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceMigrationPhase(input string) (*NetworkInterfaceMigrationPhase, error) {
	vals := map[string]NetworkInterfaceMigrationPhase{
		"abort":     NetworkInterfaceMigrationPhaseAbort,
		"commit":    NetworkInterfaceMigrationPhaseCommit,
		"committed": NetworkInterfaceMigrationPhaseCommitted,
		"none":      NetworkInterfaceMigrationPhaseNone,
		"prepare":   NetworkInterfaceMigrationPhasePrepare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceMigrationPhase(input)
	return &out, nil
}

type NetworkInterfaceNicType string

const (
	NetworkInterfaceNicTypeElastic  NetworkInterfaceNicType = "Elastic"
	NetworkInterfaceNicTypeStandard NetworkInterfaceNicType = "Standard"
)

func PossibleValuesForNetworkInterfaceNicType() []string {
	return []string{
		string(NetworkInterfaceNicTypeElastic),
		string(NetworkInterfaceNicTypeStandard),
	}
}

func (s *NetworkInterfaceNicType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceNicType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceNicType(input string) (*NetworkInterfaceNicType, error) {
	vals := map[string]NetworkInterfaceNicType{
		"elastic":  NetworkInterfaceNicTypeElastic,
		"standard": NetworkInterfaceNicTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceNicType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicIPAddressDnsSettingsDomainNameLabelScope string

const (
	PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse            PublicIPAddressDnsSettingsDomainNameLabelScope = "NoReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse PublicIPAddressDnsSettingsDomainNameLabelScope = "ResourceGroupReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse  PublicIPAddressDnsSettingsDomainNameLabelScope = "SubscriptionReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse        PublicIPAddressDnsSettingsDomainNameLabelScope = "TenantReuse"
)

func PossibleValuesForPublicIPAddressDnsSettingsDomainNameLabelScope() []string {
	return []string{
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse),
	}
}

func (s *PublicIPAddressDnsSettingsDomainNameLabelScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressDnsSettingsDomainNameLabelScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressDnsSettingsDomainNameLabelScope(input string) (*PublicIPAddressDnsSettingsDomainNameLabelScope, error) {
	vals := map[string]PublicIPAddressDnsSettingsDomainNameLabelScope{
		"noreuse":            PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse,
		"resourcegroupreuse": PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse,
		"subscriptionreuse":  PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse,
		"tenantreuse":        PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressDnsSettingsDomainNameLabelScope(input)
	return &out, nil
}

type PublicIPAddressMigrationPhase string

const (
	PublicIPAddressMigrationPhaseAbort     PublicIPAddressMigrationPhase = "Abort"
	PublicIPAddressMigrationPhaseCommit    PublicIPAddressMigrationPhase = "Commit"
	PublicIPAddressMigrationPhaseCommitted PublicIPAddressMigrationPhase = "Committed"
	PublicIPAddressMigrationPhaseNone      PublicIPAddressMigrationPhase = "None"
	PublicIPAddressMigrationPhasePrepare   PublicIPAddressMigrationPhase = "Prepare"
)

func PossibleValuesForPublicIPAddressMigrationPhase() []string {
	return []string{
		string(PublicIPAddressMigrationPhaseAbort),
		string(PublicIPAddressMigrationPhaseCommit),
		string(PublicIPAddressMigrationPhaseCommitted),
		string(PublicIPAddressMigrationPhaseNone),
		string(PublicIPAddressMigrationPhasePrepare),
	}
}

func (s *PublicIPAddressMigrationPhase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressMigrationPhase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressMigrationPhase(input string) (*PublicIPAddressMigrationPhase, error) {
	vals := map[string]PublicIPAddressMigrationPhase{
		"abort":     PublicIPAddressMigrationPhaseAbort,
		"commit":    PublicIPAddressMigrationPhaseCommit,
		"committed": PublicIPAddressMigrationPhaseCommitted,
		"none":      PublicIPAddressMigrationPhaseNone,
		"prepare":   PublicIPAddressMigrationPhasePrepare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressMigrationPhase(input)
	return &out, nil
}

type PublicIPAddressSkuName string

const (
	PublicIPAddressSkuNameBasic    PublicIPAddressSkuName = "Basic"
	PublicIPAddressSkuNameStandard PublicIPAddressSkuName = "Standard"
)

func PossibleValuesForPublicIPAddressSkuName() []string {
	return []string{
		string(PublicIPAddressSkuNameBasic),
		string(PublicIPAddressSkuNameStandard),
	}
}

func (s *PublicIPAddressSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuName(input string) (*PublicIPAddressSkuName, error) {
	vals := map[string]PublicIPAddressSkuName{
		"basic":    PublicIPAddressSkuNameBasic,
		"standard": PublicIPAddressSkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuName(input)
	return &out, nil
}

type PublicIPAddressSkuTier string

const (
	PublicIPAddressSkuTierGlobal   PublicIPAddressSkuTier = "Global"
	PublicIPAddressSkuTierRegional PublicIPAddressSkuTier = "Regional"
)

func PossibleValuesForPublicIPAddressSkuTier() []string {
	return []string{
		string(PublicIPAddressSkuTierGlobal),
		string(PublicIPAddressSkuTierRegional),
	}
}

func (s *PublicIPAddressSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuTier(input string) (*PublicIPAddressSkuTier, error) {
	vals := map[string]PublicIPAddressSkuTier{
		"global":   PublicIPAddressSkuTierGlobal,
		"regional": PublicIPAddressSkuTierRegional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuTier(input)
	return &out, nil
}

type RouteNextHopType string

const (
	RouteNextHopTypeInternet              RouteNextHopType = "Internet"
	RouteNextHopTypeNone                  RouteNextHopType = "None"
	RouteNextHopTypeVirtualAppliance      RouteNextHopType = "VirtualAppliance"
	RouteNextHopTypeVirtualNetworkGateway RouteNextHopType = "VirtualNetworkGateway"
	RouteNextHopTypeVnetLocal             RouteNextHopType = "VnetLocal"
)

func PossibleValuesForRouteNextHopType() []string {
	return []string{
		string(RouteNextHopTypeInternet),
		string(RouteNextHopTypeNone),
		string(RouteNextHopTypeVirtualAppliance),
		string(RouteNextHopTypeVirtualNetworkGateway),
		string(RouteNextHopTypeVnetLocal),
	}
}

func (s *RouteNextHopType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRouteNextHopType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRouteNextHopType(input string) (*RouteNextHopType, error) {
	vals := map[string]RouteNextHopType{
		"internet":              RouteNextHopTypeInternet,
		"none":                  RouteNextHopTypeNone,
		"virtualappliance":      RouteNextHopTypeVirtualAppliance,
		"virtualnetworkgateway": RouteNextHopTypeVirtualNetworkGateway,
		"vnetlocal":             RouteNextHopTypeVnetLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RouteNextHopType(input)
	return &out, nil
}

type SecurityRuleAccess string

const (
	SecurityRuleAccessAllow SecurityRuleAccess = "Allow"
	SecurityRuleAccessDeny  SecurityRuleAccess = "Deny"
)

func PossibleValuesForSecurityRuleAccess() []string {
	return []string{
		string(SecurityRuleAccessAllow),
		string(SecurityRuleAccessDeny),
	}
}

func (s *SecurityRuleAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleAccess(input string) (*SecurityRuleAccess, error) {
	vals := map[string]SecurityRuleAccess{
		"allow": SecurityRuleAccessAllow,
		"deny":  SecurityRuleAccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleAccess(input)
	return &out, nil
}

type SecurityRuleDirection string

const (
	SecurityRuleDirectionInbound  SecurityRuleDirection = "Inbound"
	SecurityRuleDirectionOutbound SecurityRuleDirection = "Outbound"
)

func PossibleValuesForSecurityRuleDirection() []string {
	return []string{
		string(SecurityRuleDirectionInbound),
		string(SecurityRuleDirectionOutbound),
	}
}

func (s *SecurityRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleDirection(input string) (*SecurityRuleDirection, error) {
	vals := map[string]SecurityRuleDirection{
		"inbound":  SecurityRuleDirectionInbound,
		"outbound": SecurityRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleDirection(input)
	return &out, nil
}

type SecurityRuleProtocol string

const (
	SecurityRuleProtocolAh   SecurityRuleProtocol = "Ah"
	SecurityRuleProtocolAny  SecurityRuleProtocol = "*"
	SecurityRuleProtocolEsp  SecurityRuleProtocol = "Esp"
	SecurityRuleProtocolIcmp SecurityRuleProtocol = "Icmp"
	SecurityRuleProtocolTcp  SecurityRuleProtocol = "Tcp"
	SecurityRuleProtocolUdp  SecurityRuleProtocol = "Udp"
)

func PossibleValuesForSecurityRuleProtocol() []string {
	return []string{
		string(SecurityRuleProtocolAh),
		string(SecurityRuleProtocolAny),
		string(SecurityRuleProtocolEsp),
		string(SecurityRuleProtocolIcmp),
		string(SecurityRuleProtocolTcp),
		string(SecurityRuleProtocolUdp),
	}
}

func (s *SecurityRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleProtocol(input string) (*SecurityRuleProtocol, error) {
	vals := map[string]SecurityRuleProtocol{
		"ah":   SecurityRuleProtocolAh,
		"*":    SecurityRuleProtocolAny,
		"esp":  SecurityRuleProtocolEsp,
		"icmp": SecurityRuleProtocolIcmp,
		"tcp":  SecurityRuleProtocolTcp,
		"udp":  SecurityRuleProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleProtocol(input)
	return &out, nil
}

type SyncMode string

const (
	SyncModeAutomatic SyncMode = "Automatic"
	SyncModeManual    SyncMode = "Manual"
)

func PossibleValuesForSyncMode() []string {
	return []string{
		string(SyncModeAutomatic),
		string(SyncModeManual),
	}
}

func (s *SyncMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyncMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyncMode(input string) (*SyncMode, error) {
	vals := map[string]SyncMode{
		"automatic": SyncModeAutomatic,
		"manual":    SyncModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncMode(input)
	return &out, nil
}

type TransportProtocol string

const (
	TransportProtocolAll TransportProtocol = "All"
	TransportProtocolTcp TransportProtocol = "Tcp"
	TransportProtocolUdp TransportProtocol = "Udp"
)

func PossibleValuesForTransportProtocol() []string {
	return []string{
		string(TransportProtocolAll),
		string(TransportProtocolTcp),
		string(TransportProtocolUdp),
	}
}

func (s *TransportProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransportProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransportProtocol(input string) (*TransportProtocol, error) {
	vals := map[string]TransportProtocol{
		"all": TransportProtocolAll,
		"tcp": TransportProtocolTcp,
		"udp": TransportProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransportProtocol(input)
	return &out, nil
}

type VirtualNetworkPrivateEndpointNetworkPolicies string

const (
	VirtualNetworkPrivateEndpointNetworkPoliciesDisabled VirtualNetworkPrivateEndpointNetworkPolicies = "Disabled"
	VirtualNetworkPrivateEndpointNetworkPoliciesEnabled  VirtualNetworkPrivateEndpointNetworkPolicies = "Enabled"
)

func PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies() []string {
	return []string{
		string(VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
		string(VirtualNetworkPrivateEndpointNetworkPoliciesEnabled),
	}
}

func (s *VirtualNetworkPrivateEndpointNetworkPolicies) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPrivateEndpointNetworkPolicies(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPrivateEndpointNetworkPolicies(input string) (*VirtualNetworkPrivateEndpointNetworkPolicies, error) {
	vals := map[string]VirtualNetworkPrivateEndpointNetworkPolicies{
		"disabled": VirtualNetworkPrivateEndpointNetworkPoliciesDisabled,
		"enabled":  VirtualNetworkPrivateEndpointNetworkPoliciesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPrivateEndpointNetworkPolicies(input)
	return &out, nil
}

type VirtualNetworkPrivateLinkServiceNetworkPolicies string

const (
	VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled VirtualNetworkPrivateLinkServiceNetworkPolicies = "Disabled"
	VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled  VirtualNetworkPrivateLinkServiceNetworkPolicies = "Enabled"
)

func PossibleValuesForVirtualNetworkPrivateLinkServiceNetworkPolicies() []string {
	return []string{
		string(VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled),
		string(VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled),
	}
}

func (s *VirtualNetworkPrivateLinkServiceNetworkPolicies) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPrivateLinkServiceNetworkPolicies(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPrivateLinkServiceNetworkPolicies(input string) (*VirtualNetworkPrivateLinkServiceNetworkPolicies, error) {
	vals := map[string]VirtualNetworkPrivateLinkServiceNetworkPolicies{
		"disabled": VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled,
		"enabled":  VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPrivateLinkServiceNetworkPolicies(input)
	return &out, nil
}

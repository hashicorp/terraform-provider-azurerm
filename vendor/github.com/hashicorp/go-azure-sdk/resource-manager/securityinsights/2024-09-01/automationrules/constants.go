package automationrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionType string

const (
	ActionTypeAddIncidentTask  ActionType = "AddIncidentTask"
	ActionTypeModifyProperties ActionType = "ModifyProperties"
	ActionTypeRunPlaybook      ActionType = "RunPlaybook"
)

func PossibleValuesForActionType() []string {
	return []string{
		string(ActionTypeAddIncidentTask),
		string(ActionTypeModifyProperties),
		string(ActionTypeRunPlaybook),
	}
}

func (s *ActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionType(input string) (*ActionType, error) {
	vals := map[string]ActionType{
		"addincidenttask":  ActionTypeAddIncidentTask,
		"modifyproperties": ActionTypeModifyProperties,
		"runplaybook":      ActionTypeRunPlaybook,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionType(input)
	return &out, nil
}

type AutomationRuleBooleanConditionSupportedOperator string

const (
	AutomationRuleBooleanConditionSupportedOperatorAnd AutomationRuleBooleanConditionSupportedOperator = "And"
	AutomationRuleBooleanConditionSupportedOperatorOr  AutomationRuleBooleanConditionSupportedOperator = "Or"
)

func PossibleValuesForAutomationRuleBooleanConditionSupportedOperator() []string {
	return []string{
		string(AutomationRuleBooleanConditionSupportedOperatorAnd),
		string(AutomationRuleBooleanConditionSupportedOperatorOr),
	}
}

func (s *AutomationRuleBooleanConditionSupportedOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRuleBooleanConditionSupportedOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRuleBooleanConditionSupportedOperator(input string) (*AutomationRuleBooleanConditionSupportedOperator, error) {
	vals := map[string]AutomationRuleBooleanConditionSupportedOperator{
		"and": AutomationRuleBooleanConditionSupportedOperatorAnd,
		"or":  AutomationRuleBooleanConditionSupportedOperatorOr,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRuleBooleanConditionSupportedOperator(input)
	return &out, nil
}

type AutomationRulePropertyArrayChangedConditionSupportedArrayType string

const (
	AutomationRulePropertyArrayChangedConditionSupportedArrayTypeAlerts   AutomationRulePropertyArrayChangedConditionSupportedArrayType = "Alerts"
	AutomationRulePropertyArrayChangedConditionSupportedArrayTypeComments AutomationRulePropertyArrayChangedConditionSupportedArrayType = "Comments"
	AutomationRulePropertyArrayChangedConditionSupportedArrayTypeLabels   AutomationRulePropertyArrayChangedConditionSupportedArrayType = "Labels"
	AutomationRulePropertyArrayChangedConditionSupportedArrayTypeTactics  AutomationRulePropertyArrayChangedConditionSupportedArrayType = "Tactics"
)

func PossibleValuesForAutomationRulePropertyArrayChangedConditionSupportedArrayType() []string {
	return []string{
		string(AutomationRulePropertyArrayChangedConditionSupportedArrayTypeAlerts),
		string(AutomationRulePropertyArrayChangedConditionSupportedArrayTypeComments),
		string(AutomationRulePropertyArrayChangedConditionSupportedArrayTypeLabels),
		string(AutomationRulePropertyArrayChangedConditionSupportedArrayTypeTactics),
	}
}

func (s *AutomationRulePropertyArrayChangedConditionSupportedArrayType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyArrayChangedConditionSupportedArrayType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyArrayChangedConditionSupportedArrayType(input string) (*AutomationRulePropertyArrayChangedConditionSupportedArrayType, error) {
	vals := map[string]AutomationRulePropertyArrayChangedConditionSupportedArrayType{
		"alerts":   AutomationRulePropertyArrayChangedConditionSupportedArrayTypeAlerts,
		"comments": AutomationRulePropertyArrayChangedConditionSupportedArrayTypeComments,
		"labels":   AutomationRulePropertyArrayChangedConditionSupportedArrayTypeLabels,
		"tactics":  AutomationRulePropertyArrayChangedConditionSupportedArrayTypeTactics,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyArrayChangedConditionSupportedArrayType(input)
	return &out, nil
}

type AutomationRulePropertyArrayChangedConditionSupportedChangeType string

const (
	AutomationRulePropertyArrayChangedConditionSupportedChangeTypeAdded AutomationRulePropertyArrayChangedConditionSupportedChangeType = "Added"
)

func PossibleValuesForAutomationRulePropertyArrayChangedConditionSupportedChangeType() []string {
	return []string{
		string(AutomationRulePropertyArrayChangedConditionSupportedChangeTypeAdded),
	}
}

func (s *AutomationRulePropertyArrayChangedConditionSupportedChangeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyArrayChangedConditionSupportedChangeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyArrayChangedConditionSupportedChangeType(input string) (*AutomationRulePropertyArrayChangedConditionSupportedChangeType, error) {
	vals := map[string]AutomationRulePropertyArrayChangedConditionSupportedChangeType{
		"added": AutomationRulePropertyArrayChangedConditionSupportedChangeTypeAdded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyArrayChangedConditionSupportedChangeType(input)
	return &out, nil
}

type AutomationRulePropertyArrayConditionSupportedArrayConditionType string

const (
	AutomationRulePropertyArrayConditionSupportedArrayConditionTypeAnyItem AutomationRulePropertyArrayConditionSupportedArrayConditionType = "AnyItem"
)

func PossibleValuesForAutomationRulePropertyArrayConditionSupportedArrayConditionType() []string {
	return []string{
		string(AutomationRulePropertyArrayConditionSupportedArrayConditionTypeAnyItem),
	}
}

func (s *AutomationRulePropertyArrayConditionSupportedArrayConditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyArrayConditionSupportedArrayConditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyArrayConditionSupportedArrayConditionType(input string) (*AutomationRulePropertyArrayConditionSupportedArrayConditionType, error) {
	vals := map[string]AutomationRulePropertyArrayConditionSupportedArrayConditionType{
		"anyitem": AutomationRulePropertyArrayConditionSupportedArrayConditionTypeAnyItem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyArrayConditionSupportedArrayConditionType(input)
	return &out, nil
}

type AutomationRulePropertyArrayConditionSupportedArrayType string

const (
	AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetailValues AutomationRulePropertyArrayConditionSupportedArrayType = "CustomDetailValues"
	AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetails      AutomationRulePropertyArrayConditionSupportedArrayType = "CustomDetails"
)

func PossibleValuesForAutomationRulePropertyArrayConditionSupportedArrayType() []string {
	return []string{
		string(AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetailValues),
		string(AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetails),
	}
}

func (s *AutomationRulePropertyArrayConditionSupportedArrayType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyArrayConditionSupportedArrayType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyArrayConditionSupportedArrayType(input string) (*AutomationRulePropertyArrayConditionSupportedArrayType, error) {
	vals := map[string]AutomationRulePropertyArrayConditionSupportedArrayType{
		"customdetailvalues": AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetailValues,
		"customdetails":      AutomationRulePropertyArrayConditionSupportedArrayTypeCustomDetails,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyArrayConditionSupportedArrayType(input)
	return &out, nil
}

type AutomationRulePropertyChangedConditionSupportedChangedType string

const (
	AutomationRulePropertyChangedConditionSupportedChangedTypeChangedFrom AutomationRulePropertyChangedConditionSupportedChangedType = "ChangedFrom"
	AutomationRulePropertyChangedConditionSupportedChangedTypeChangedTo   AutomationRulePropertyChangedConditionSupportedChangedType = "ChangedTo"
)

func PossibleValuesForAutomationRulePropertyChangedConditionSupportedChangedType() []string {
	return []string{
		string(AutomationRulePropertyChangedConditionSupportedChangedTypeChangedFrom),
		string(AutomationRulePropertyChangedConditionSupportedChangedTypeChangedTo),
	}
}

func (s *AutomationRulePropertyChangedConditionSupportedChangedType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyChangedConditionSupportedChangedType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyChangedConditionSupportedChangedType(input string) (*AutomationRulePropertyChangedConditionSupportedChangedType, error) {
	vals := map[string]AutomationRulePropertyChangedConditionSupportedChangedType{
		"changedfrom": AutomationRulePropertyChangedConditionSupportedChangedTypeChangedFrom,
		"changedto":   AutomationRulePropertyChangedConditionSupportedChangedTypeChangedTo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyChangedConditionSupportedChangedType(input)
	return &out, nil
}

type AutomationRulePropertyChangedConditionSupportedPropertyType string

const (
	AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentOwner    AutomationRulePropertyChangedConditionSupportedPropertyType = "IncidentOwner"
	AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentSeverity AutomationRulePropertyChangedConditionSupportedPropertyType = "IncidentSeverity"
	AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentStatus   AutomationRulePropertyChangedConditionSupportedPropertyType = "IncidentStatus"
)

func PossibleValuesForAutomationRulePropertyChangedConditionSupportedPropertyType() []string {
	return []string{
		string(AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentOwner),
		string(AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentSeverity),
		string(AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentStatus),
	}
}

func (s *AutomationRulePropertyChangedConditionSupportedPropertyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyChangedConditionSupportedPropertyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyChangedConditionSupportedPropertyType(input string) (*AutomationRulePropertyChangedConditionSupportedPropertyType, error) {
	vals := map[string]AutomationRulePropertyChangedConditionSupportedPropertyType{
		"incidentowner":    AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentOwner,
		"incidentseverity": AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentSeverity,
		"incidentstatus":   AutomationRulePropertyChangedConditionSupportedPropertyTypeIncidentStatus,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyChangedConditionSupportedPropertyType(input)
	return &out, nil
}

type AutomationRulePropertyConditionSupportedOperator string

const (
	AutomationRulePropertyConditionSupportedOperatorContains      AutomationRulePropertyConditionSupportedOperator = "Contains"
	AutomationRulePropertyConditionSupportedOperatorEndsWith      AutomationRulePropertyConditionSupportedOperator = "EndsWith"
	AutomationRulePropertyConditionSupportedOperatorEquals        AutomationRulePropertyConditionSupportedOperator = "Equals"
	AutomationRulePropertyConditionSupportedOperatorNotContains   AutomationRulePropertyConditionSupportedOperator = "NotContains"
	AutomationRulePropertyConditionSupportedOperatorNotEndsWith   AutomationRulePropertyConditionSupportedOperator = "NotEndsWith"
	AutomationRulePropertyConditionSupportedOperatorNotEquals     AutomationRulePropertyConditionSupportedOperator = "NotEquals"
	AutomationRulePropertyConditionSupportedOperatorNotStartsWith AutomationRulePropertyConditionSupportedOperator = "NotStartsWith"
	AutomationRulePropertyConditionSupportedOperatorStartsWith    AutomationRulePropertyConditionSupportedOperator = "StartsWith"
)

func PossibleValuesForAutomationRulePropertyConditionSupportedOperator() []string {
	return []string{
		string(AutomationRulePropertyConditionSupportedOperatorContains),
		string(AutomationRulePropertyConditionSupportedOperatorEndsWith),
		string(AutomationRulePropertyConditionSupportedOperatorEquals),
		string(AutomationRulePropertyConditionSupportedOperatorNotContains),
		string(AutomationRulePropertyConditionSupportedOperatorNotEndsWith),
		string(AutomationRulePropertyConditionSupportedOperatorNotEquals),
		string(AutomationRulePropertyConditionSupportedOperatorNotStartsWith),
		string(AutomationRulePropertyConditionSupportedOperatorStartsWith),
	}
}

func (s *AutomationRulePropertyConditionSupportedOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyConditionSupportedOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyConditionSupportedOperator(input string) (*AutomationRulePropertyConditionSupportedOperator, error) {
	vals := map[string]AutomationRulePropertyConditionSupportedOperator{
		"contains":      AutomationRulePropertyConditionSupportedOperatorContains,
		"endswith":      AutomationRulePropertyConditionSupportedOperatorEndsWith,
		"equals":        AutomationRulePropertyConditionSupportedOperatorEquals,
		"notcontains":   AutomationRulePropertyConditionSupportedOperatorNotContains,
		"notendswith":   AutomationRulePropertyConditionSupportedOperatorNotEndsWith,
		"notequals":     AutomationRulePropertyConditionSupportedOperatorNotEquals,
		"notstartswith": AutomationRulePropertyConditionSupportedOperatorNotStartsWith,
		"startswith":    AutomationRulePropertyConditionSupportedOperatorStartsWith,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyConditionSupportedOperator(input)
	return &out, nil
}

type AutomationRulePropertyConditionSupportedProperty string

const (
	AutomationRulePropertyConditionSupportedPropertyAccountAadTenantId             AutomationRulePropertyConditionSupportedProperty = "AccountAadTenantId"
	AutomationRulePropertyConditionSupportedPropertyAccountAadUserId               AutomationRulePropertyConditionSupportedProperty = "AccountAadUserId"
	AutomationRulePropertyConditionSupportedPropertyAccountNTDomain                AutomationRulePropertyConditionSupportedProperty = "AccountNTDomain"
	AutomationRulePropertyConditionSupportedPropertyAccountName                    AutomationRulePropertyConditionSupportedProperty = "AccountName"
	AutomationRulePropertyConditionSupportedPropertyAccountObjectGuid              AutomationRulePropertyConditionSupportedProperty = "AccountObjectGuid"
	AutomationRulePropertyConditionSupportedPropertyAccountPUID                    AutomationRulePropertyConditionSupportedProperty = "AccountPUID"
	AutomationRulePropertyConditionSupportedPropertyAccountSid                     AutomationRulePropertyConditionSupportedProperty = "AccountSid"
	AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix               AutomationRulePropertyConditionSupportedProperty = "AccountUPNSuffix"
	AutomationRulePropertyConditionSupportedPropertyAlertAnalyticRuleIds           AutomationRulePropertyConditionSupportedProperty = "AlertAnalyticRuleIds"
	AutomationRulePropertyConditionSupportedPropertyAlertProductNames              AutomationRulePropertyConditionSupportedProperty = "AlertProductNames"
	AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceId        AutomationRulePropertyConditionSupportedProperty = "AzureResourceResourceId"
	AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionId    AutomationRulePropertyConditionSupportedProperty = "AzureResourceSubscriptionId"
	AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppId          AutomationRulePropertyConditionSupportedProperty = "CloudApplicationAppId"
	AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName        AutomationRulePropertyConditionSupportedProperty = "CloudApplicationAppName"
	AutomationRulePropertyConditionSupportedPropertyDNSDomainName                  AutomationRulePropertyConditionSupportedProperty = "DNSDomainName"
	AutomationRulePropertyConditionSupportedPropertyFileDirectory                  AutomationRulePropertyConditionSupportedProperty = "FileDirectory"
	AutomationRulePropertyConditionSupportedPropertyFileHashValue                  AutomationRulePropertyConditionSupportedProperty = "FileHashValue"
	AutomationRulePropertyConditionSupportedPropertyFileName                       AutomationRulePropertyConditionSupportedProperty = "FileName"
	AutomationRulePropertyConditionSupportedPropertyHostAzureID                    AutomationRulePropertyConditionSupportedProperty = "HostAzureID"
	AutomationRulePropertyConditionSupportedPropertyHostNTDomain                   AutomationRulePropertyConditionSupportedProperty = "HostNTDomain"
	AutomationRulePropertyConditionSupportedPropertyHostName                       AutomationRulePropertyConditionSupportedProperty = "HostName"
	AutomationRulePropertyConditionSupportedPropertyHostNetBiosName                AutomationRulePropertyConditionSupportedProperty = "HostNetBiosName"
	AutomationRulePropertyConditionSupportedPropertyHostOSVersion                  AutomationRulePropertyConditionSupportedProperty = "HostOSVersion"
	AutomationRulePropertyConditionSupportedPropertyIPAddress                      AutomationRulePropertyConditionSupportedProperty = "IPAddress"
	AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsKey       AutomationRulePropertyConditionSupportedProperty = "IncidentCustomDetailsKey"
	AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsValue     AutomationRulePropertyConditionSupportedProperty = "IncidentCustomDetailsValue"
	AutomationRulePropertyConditionSupportedPropertyIncidentDescription            AutomationRulePropertyConditionSupportedProperty = "IncidentDescription"
	AutomationRulePropertyConditionSupportedPropertyIncidentLabel                  AutomationRulePropertyConditionSupportedProperty = "IncidentLabel"
	AutomationRulePropertyConditionSupportedPropertyIncidentProviderName           AutomationRulePropertyConditionSupportedProperty = "IncidentProviderName"
	AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds AutomationRulePropertyConditionSupportedProperty = "IncidentRelatedAnalyticRuleIds"
	AutomationRulePropertyConditionSupportedPropertyIncidentSeverity               AutomationRulePropertyConditionSupportedProperty = "IncidentSeverity"
	AutomationRulePropertyConditionSupportedPropertyIncidentStatus                 AutomationRulePropertyConditionSupportedProperty = "IncidentStatus"
	AutomationRulePropertyConditionSupportedPropertyIncidentTactics                AutomationRulePropertyConditionSupportedProperty = "IncidentTactics"
	AutomationRulePropertyConditionSupportedPropertyIncidentTitle                  AutomationRulePropertyConditionSupportedProperty = "IncidentTitle"
	AutomationRulePropertyConditionSupportedPropertyIncidentUpdatedBySource        AutomationRulePropertyConditionSupportedProperty = "IncidentUpdatedBySource"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceId                    AutomationRulePropertyConditionSupportedProperty = "IoTDeviceId"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel                 AutomationRulePropertyConditionSupportedProperty = "IoTDeviceModel"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceName                  AutomationRulePropertyConditionSupportedProperty = "IoTDeviceName"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem       AutomationRulePropertyConditionSupportedProperty = "IoTDeviceOperatingSystem"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceType                  AutomationRulePropertyConditionSupportedProperty = "IoTDeviceType"
	AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor                AutomationRulePropertyConditionSupportedProperty = "IoTDeviceVendor"
	AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction      AutomationRulePropertyConditionSupportedProperty = "MailMessageDeliveryAction"
	AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation    AutomationRulePropertyConditionSupportedProperty = "MailMessageDeliveryLocation"
	AutomationRulePropertyConditionSupportedPropertyMailMessagePOneSender          AutomationRulePropertyConditionSupportedProperty = "MailMessageP1Sender"
	AutomationRulePropertyConditionSupportedPropertyMailMessagePTwoSender          AutomationRulePropertyConditionSupportedProperty = "MailMessageP2Sender"
	AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient           AutomationRulePropertyConditionSupportedProperty = "MailMessageRecipient"
	AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP            AutomationRulePropertyConditionSupportedProperty = "MailMessageSenderIP"
	AutomationRulePropertyConditionSupportedPropertyMailMessageSubject             AutomationRulePropertyConditionSupportedProperty = "MailMessageSubject"
	AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName             AutomationRulePropertyConditionSupportedProperty = "MailboxDisplayName"
	AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress          AutomationRulePropertyConditionSupportedProperty = "MailboxPrimaryAddress"
	AutomationRulePropertyConditionSupportedPropertyMailboxUPN                     AutomationRulePropertyConditionSupportedProperty = "MailboxUPN"
	AutomationRulePropertyConditionSupportedPropertyMalwareCategory                AutomationRulePropertyConditionSupportedProperty = "MalwareCategory"
	AutomationRulePropertyConditionSupportedPropertyMalwareName                    AutomationRulePropertyConditionSupportedProperty = "MalwareName"
	AutomationRulePropertyConditionSupportedPropertyProcessCommandLine             AutomationRulePropertyConditionSupportedProperty = "ProcessCommandLine"
	AutomationRulePropertyConditionSupportedPropertyProcessId                      AutomationRulePropertyConditionSupportedProperty = "ProcessId"
	AutomationRulePropertyConditionSupportedPropertyRegistryKey                    AutomationRulePropertyConditionSupportedProperty = "RegistryKey"
	AutomationRulePropertyConditionSupportedPropertyRegistryValueData              AutomationRulePropertyConditionSupportedProperty = "RegistryValueData"
	AutomationRulePropertyConditionSupportedPropertyURL                            AutomationRulePropertyConditionSupportedProperty = "Url"
)

func PossibleValuesForAutomationRulePropertyConditionSupportedProperty() []string {
	return []string{
		string(AutomationRulePropertyConditionSupportedPropertyAccountAadTenantId),
		string(AutomationRulePropertyConditionSupportedPropertyAccountAadUserId),
		string(AutomationRulePropertyConditionSupportedPropertyAccountNTDomain),
		string(AutomationRulePropertyConditionSupportedPropertyAccountName),
		string(AutomationRulePropertyConditionSupportedPropertyAccountObjectGuid),
		string(AutomationRulePropertyConditionSupportedPropertyAccountPUID),
		string(AutomationRulePropertyConditionSupportedPropertyAccountSid),
		string(AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix),
		string(AutomationRulePropertyConditionSupportedPropertyAlertAnalyticRuleIds),
		string(AutomationRulePropertyConditionSupportedPropertyAlertProductNames),
		string(AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceId),
		string(AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionId),
		string(AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppId),
		string(AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName),
		string(AutomationRulePropertyConditionSupportedPropertyDNSDomainName),
		string(AutomationRulePropertyConditionSupportedPropertyFileDirectory),
		string(AutomationRulePropertyConditionSupportedPropertyFileHashValue),
		string(AutomationRulePropertyConditionSupportedPropertyFileName),
		string(AutomationRulePropertyConditionSupportedPropertyHostAzureID),
		string(AutomationRulePropertyConditionSupportedPropertyHostNTDomain),
		string(AutomationRulePropertyConditionSupportedPropertyHostName),
		string(AutomationRulePropertyConditionSupportedPropertyHostNetBiosName),
		string(AutomationRulePropertyConditionSupportedPropertyHostOSVersion),
		string(AutomationRulePropertyConditionSupportedPropertyIPAddress),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsKey),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsValue),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentDescription),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentLabel),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentProviderName),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentSeverity),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentStatus),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentTactics),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentTitle),
		string(AutomationRulePropertyConditionSupportedPropertyIncidentUpdatedBySource),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceId),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceName),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceType),
		string(AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessagePOneSender),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessagePTwoSender),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP),
		string(AutomationRulePropertyConditionSupportedPropertyMailMessageSubject),
		string(AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName),
		string(AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress),
		string(AutomationRulePropertyConditionSupportedPropertyMailboxUPN),
		string(AutomationRulePropertyConditionSupportedPropertyMalwareCategory),
		string(AutomationRulePropertyConditionSupportedPropertyMalwareName),
		string(AutomationRulePropertyConditionSupportedPropertyProcessCommandLine),
		string(AutomationRulePropertyConditionSupportedPropertyProcessId),
		string(AutomationRulePropertyConditionSupportedPropertyRegistryKey),
		string(AutomationRulePropertyConditionSupportedPropertyRegistryValueData),
		string(AutomationRulePropertyConditionSupportedPropertyURL),
	}
}

func (s *AutomationRulePropertyConditionSupportedProperty) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutomationRulePropertyConditionSupportedProperty(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutomationRulePropertyConditionSupportedProperty(input string) (*AutomationRulePropertyConditionSupportedProperty, error) {
	vals := map[string]AutomationRulePropertyConditionSupportedProperty{
		"accountaadtenantid":             AutomationRulePropertyConditionSupportedPropertyAccountAadTenantId,
		"accountaaduserid":               AutomationRulePropertyConditionSupportedPropertyAccountAadUserId,
		"accountntdomain":                AutomationRulePropertyConditionSupportedPropertyAccountNTDomain,
		"accountname":                    AutomationRulePropertyConditionSupportedPropertyAccountName,
		"accountobjectguid":              AutomationRulePropertyConditionSupportedPropertyAccountObjectGuid,
		"accountpuid":                    AutomationRulePropertyConditionSupportedPropertyAccountPUID,
		"accountsid":                     AutomationRulePropertyConditionSupportedPropertyAccountSid,
		"accountupnsuffix":               AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix,
		"alertanalyticruleids":           AutomationRulePropertyConditionSupportedPropertyAlertAnalyticRuleIds,
		"alertproductnames":              AutomationRulePropertyConditionSupportedPropertyAlertProductNames,
		"azureresourceresourceid":        AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceId,
		"azureresourcesubscriptionid":    AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionId,
		"cloudapplicationappid":          AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppId,
		"cloudapplicationappname":        AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName,
		"dnsdomainname":                  AutomationRulePropertyConditionSupportedPropertyDNSDomainName,
		"filedirectory":                  AutomationRulePropertyConditionSupportedPropertyFileDirectory,
		"filehashvalue":                  AutomationRulePropertyConditionSupportedPropertyFileHashValue,
		"filename":                       AutomationRulePropertyConditionSupportedPropertyFileName,
		"hostazureid":                    AutomationRulePropertyConditionSupportedPropertyHostAzureID,
		"hostntdomain":                   AutomationRulePropertyConditionSupportedPropertyHostNTDomain,
		"hostname":                       AutomationRulePropertyConditionSupportedPropertyHostName,
		"hostnetbiosname":                AutomationRulePropertyConditionSupportedPropertyHostNetBiosName,
		"hostosversion":                  AutomationRulePropertyConditionSupportedPropertyHostOSVersion,
		"ipaddress":                      AutomationRulePropertyConditionSupportedPropertyIPAddress,
		"incidentcustomdetailskey":       AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsKey,
		"incidentcustomdetailsvalue":     AutomationRulePropertyConditionSupportedPropertyIncidentCustomDetailsValue,
		"incidentdescription":            AutomationRulePropertyConditionSupportedPropertyIncidentDescription,
		"incidentlabel":                  AutomationRulePropertyConditionSupportedPropertyIncidentLabel,
		"incidentprovidername":           AutomationRulePropertyConditionSupportedPropertyIncidentProviderName,
		"incidentrelatedanalyticruleids": AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds,
		"incidentseverity":               AutomationRulePropertyConditionSupportedPropertyIncidentSeverity,
		"incidentstatus":                 AutomationRulePropertyConditionSupportedPropertyIncidentStatus,
		"incidenttactics":                AutomationRulePropertyConditionSupportedPropertyIncidentTactics,
		"incidenttitle":                  AutomationRulePropertyConditionSupportedPropertyIncidentTitle,
		"incidentupdatedbysource":        AutomationRulePropertyConditionSupportedPropertyIncidentUpdatedBySource,
		"iotdeviceid":                    AutomationRulePropertyConditionSupportedPropertyIoTDeviceId,
		"iotdevicemodel":                 AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel,
		"iotdevicename":                  AutomationRulePropertyConditionSupportedPropertyIoTDeviceName,
		"iotdeviceoperatingsystem":       AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem,
		"iotdevicetype":                  AutomationRulePropertyConditionSupportedPropertyIoTDeviceType,
		"iotdevicevendor":                AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor,
		"mailmessagedeliveryaction":      AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction,
		"mailmessagedeliverylocation":    AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation,
		"mailmessagep1sender":            AutomationRulePropertyConditionSupportedPropertyMailMessagePOneSender,
		"mailmessagep2sender":            AutomationRulePropertyConditionSupportedPropertyMailMessagePTwoSender,
		"mailmessagerecipient":           AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient,
		"mailmessagesenderip":            AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP,
		"mailmessagesubject":             AutomationRulePropertyConditionSupportedPropertyMailMessageSubject,
		"mailboxdisplayname":             AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName,
		"mailboxprimaryaddress":          AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress,
		"mailboxupn":                     AutomationRulePropertyConditionSupportedPropertyMailboxUPN,
		"malwarecategory":                AutomationRulePropertyConditionSupportedPropertyMalwareCategory,
		"malwarename":                    AutomationRulePropertyConditionSupportedPropertyMalwareName,
		"processcommandline":             AutomationRulePropertyConditionSupportedPropertyProcessCommandLine,
		"processid":                      AutomationRulePropertyConditionSupportedPropertyProcessId,
		"registrykey":                    AutomationRulePropertyConditionSupportedPropertyRegistryKey,
		"registryvaluedata":              AutomationRulePropertyConditionSupportedPropertyRegistryValueData,
		"url":                            AutomationRulePropertyConditionSupportedPropertyURL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutomationRulePropertyConditionSupportedProperty(input)
	return &out, nil
}

type ConditionType string

const (
	ConditionTypeBoolean              ConditionType = "Boolean"
	ConditionTypeProperty             ConditionType = "Property"
	ConditionTypePropertyArray        ConditionType = "PropertyArray"
	ConditionTypePropertyArrayChanged ConditionType = "PropertyArrayChanged"
	ConditionTypePropertyChanged      ConditionType = "PropertyChanged"
)

func PossibleValuesForConditionType() []string {
	return []string{
		string(ConditionTypeBoolean),
		string(ConditionTypeProperty),
		string(ConditionTypePropertyArray),
		string(ConditionTypePropertyArrayChanged),
		string(ConditionTypePropertyChanged),
	}
}

func (s *ConditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConditionType(input string) (*ConditionType, error) {
	vals := map[string]ConditionType{
		"boolean":              ConditionTypeBoolean,
		"property":             ConditionTypeProperty,
		"propertyarray":        ConditionTypePropertyArray,
		"propertyarraychanged": ConditionTypePropertyArrayChanged,
		"propertychanged":      ConditionTypePropertyChanged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionType(input)
	return &out, nil
}

type IncidentClassification string

const (
	IncidentClassificationBenignPositive IncidentClassification = "BenignPositive"
	IncidentClassificationFalsePositive  IncidentClassification = "FalsePositive"
	IncidentClassificationTruePositive   IncidentClassification = "TruePositive"
	IncidentClassificationUndetermined   IncidentClassification = "Undetermined"
)

func PossibleValuesForIncidentClassification() []string {
	return []string{
		string(IncidentClassificationBenignPositive),
		string(IncidentClassificationFalsePositive),
		string(IncidentClassificationTruePositive),
		string(IncidentClassificationUndetermined),
	}
}

func (s *IncidentClassification) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIncidentClassification(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIncidentClassification(input string) (*IncidentClassification, error) {
	vals := map[string]IncidentClassification{
		"benignpositive": IncidentClassificationBenignPositive,
		"falsepositive":  IncidentClassificationFalsePositive,
		"truepositive":   IncidentClassificationTruePositive,
		"undetermined":   IncidentClassificationUndetermined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IncidentClassification(input)
	return &out, nil
}

type IncidentClassificationReason string

const (
	IncidentClassificationReasonInaccurateData        IncidentClassificationReason = "InaccurateData"
	IncidentClassificationReasonIncorrectAlertLogic   IncidentClassificationReason = "IncorrectAlertLogic"
	IncidentClassificationReasonSuspiciousActivity    IncidentClassificationReason = "SuspiciousActivity"
	IncidentClassificationReasonSuspiciousButExpected IncidentClassificationReason = "SuspiciousButExpected"
)

func PossibleValuesForIncidentClassificationReason() []string {
	return []string{
		string(IncidentClassificationReasonInaccurateData),
		string(IncidentClassificationReasonIncorrectAlertLogic),
		string(IncidentClassificationReasonSuspiciousActivity),
		string(IncidentClassificationReasonSuspiciousButExpected),
	}
}

func (s *IncidentClassificationReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIncidentClassificationReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIncidentClassificationReason(input string) (*IncidentClassificationReason, error) {
	vals := map[string]IncidentClassificationReason{
		"inaccuratedata":        IncidentClassificationReasonInaccurateData,
		"incorrectalertlogic":   IncidentClassificationReasonIncorrectAlertLogic,
		"suspiciousactivity":    IncidentClassificationReasonSuspiciousActivity,
		"suspiciousbutexpected": IncidentClassificationReasonSuspiciousButExpected,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IncidentClassificationReason(input)
	return &out, nil
}

type IncidentLabelType string

const (
	IncidentLabelTypeAutoAssigned IncidentLabelType = "AutoAssigned"
	IncidentLabelTypeUser         IncidentLabelType = "User"
)

func PossibleValuesForIncidentLabelType() []string {
	return []string{
		string(IncidentLabelTypeAutoAssigned),
		string(IncidentLabelTypeUser),
	}
}

func (s *IncidentLabelType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIncidentLabelType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIncidentLabelType(input string) (*IncidentLabelType, error) {
	vals := map[string]IncidentLabelType{
		"autoassigned": IncidentLabelTypeAutoAssigned,
		"user":         IncidentLabelTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IncidentLabelType(input)
	return &out, nil
}

type IncidentSeverity string

const (
	IncidentSeverityHigh          IncidentSeverity = "High"
	IncidentSeverityInformational IncidentSeverity = "Informational"
	IncidentSeverityLow           IncidentSeverity = "Low"
	IncidentSeverityMedium        IncidentSeverity = "Medium"
)

func PossibleValuesForIncidentSeverity() []string {
	return []string{
		string(IncidentSeverityHigh),
		string(IncidentSeverityInformational),
		string(IncidentSeverityLow),
		string(IncidentSeverityMedium),
	}
}

func (s *IncidentSeverity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIncidentSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIncidentSeverity(input string) (*IncidentSeverity, error) {
	vals := map[string]IncidentSeverity{
		"high":          IncidentSeverityHigh,
		"informational": IncidentSeverityInformational,
		"low":           IncidentSeverityLow,
		"medium":        IncidentSeverityMedium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IncidentSeverity(input)
	return &out, nil
}

type IncidentStatus string

const (
	IncidentStatusActive IncidentStatus = "Active"
	IncidentStatusClosed IncidentStatus = "Closed"
	IncidentStatusNew    IncidentStatus = "New"
)

func PossibleValuesForIncidentStatus() []string {
	return []string{
		string(IncidentStatusActive),
		string(IncidentStatusClosed),
		string(IncidentStatusNew),
	}
}

func (s *IncidentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIncidentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIncidentStatus(input string) (*IncidentStatus, error) {
	vals := map[string]IncidentStatus{
		"active": IncidentStatusActive,
		"closed": IncidentStatusClosed,
		"new":    IncidentStatusNew,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IncidentStatus(input)
	return &out, nil
}

type OwnerType string

const (
	OwnerTypeGroup   OwnerType = "Group"
	OwnerTypeUnknown OwnerType = "Unknown"
	OwnerTypeUser    OwnerType = "User"
)

func PossibleValuesForOwnerType() []string {
	return []string{
		string(OwnerTypeGroup),
		string(OwnerTypeUnknown),
		string(OwnerTypeUser),
	}
}

func (s *OwnerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOwnerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOwnerType(input string) (*OwnerType, error) {
	vals := map[string]OwnerType{
		"group":   OwnerTypeGroup,
		"unknown": OwnerTypeUnknown,
		"user":    OwnerTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OwnerType(input)
	return &out, nil
}

type TriggersOn string

const (
	TriggersOnAlerts    TriggersOn = "Alerts"
	TriggersOnIncidents TriggersOn = "Incidents"
)

func PossibleValuesForTriggersOn() []string {
	return []string{
		string(TriggersOnAlerts),
		string(TriggersOnIncidents),
	}
}

func (s *TriggersOn) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggersOn(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggersOn(input string) (*TriggersOn, error) {
	vals := map[string]TriggersOn{
		"alerts":    TriggersOnAlerts,
		"incidents": TriggersOnIncidents,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggersOn(input)
	return &out, nil
}

type TriggersWhen string

const (
	TriggersWhenCreated TriggersWhen = "Created"
	TriggersWhenUpdated TriggersWhen = "Updated"
)

func PossibleValuesForTriggersWhen() []string {
	return []string{
		string(TriggersWhenCreated),
		string(TriggersWhenUpdated),
	}
}

func (s *TriggersWhen) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggersWhen(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggersWhen(input string) (*TriggersWhen, error) {
	vals := map[string]TriggersWhen{
		"created": TriggersWhenCreated,
		"updated": TriggersWhenUpdated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggersWhen(input)
	return &out, nil
}

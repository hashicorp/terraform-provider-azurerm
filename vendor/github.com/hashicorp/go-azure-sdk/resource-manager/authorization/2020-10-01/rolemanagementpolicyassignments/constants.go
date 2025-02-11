package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApprovalMode string

const (
	ApprovalModeNoApproval  ApprovalMode = "NoApproval"
	ApprovalModeParallel    ApprovalMode = "Parallel"
	ApprovalModeSerial      ApprovalMode = "Serial"
	ApprovalModeSingleStage ApprovalMode = "SingleStage"
)

func PossibleValuesForApprovalMode() []string {
	return []string{
		string(ApprovalModeNoApproval),
		string(ApprovalModeParallel),
		string(ApprovalModeSerial),
		string(ApprovalModeSingleStage),
	}
}

func (s *ApprovalMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApprovalMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApprovalMode(input string) (*ApprovalMode, error) {
	vals := map[string]ApprovalMode{
		"noapproval":  ApprovalModeNoApproval,
		"parallel":    ApprovalModeParallel,
		"serial":      ApprovalModeSerial,
		"singlestage": ApprovalModeSingleStage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApprovalMode(input)
	return &out, nil
}

type EnablementRules string

const (
	EnablementRulesJustification             EnablementRules = "Justification"
	EnablementRulesMultiFactorAuthentication EnablementRules = "MultiFactorAuthentication"
	EnablementRulesTicketing                 EnablementRules = "Ticketing"
)

func PossibleValuesForEnablementRules() []string {
	return []string{
		string(EnablementRulesJustification),
		string(EnablementRulesMultiFactorAuthentication),
		string(EnablementRulesTicketing),
	}
}

func (s *EnablementRules) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnablementRules(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnablementRules(input string) (*EnablementRules, error) {
	vals := map[string]EnablementRules{
		"justification":             EnablementRulesJustification,
		"multifactorauthentication": EnablementRulesMultiFactorAuthentication,
		"ticketing":                 EnablementRulesTicketing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnablementRules(input)
	return &out, nil
}

type NotificationDeliveryMechanism string

const (
	NotificationDeliveryMechanismEmail NotificationDeliveryMechanism = "Email"
)

func PossibleValuesForNotificationDeliveryMechanism() []string {
	return []string{
		string(NotificationDeliveryMechanismEmail),
	}
}

func (s *NotificationDeliveryMechanism) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotificationDeliveryMechanism(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotificationDeliveryMechanism(input string) (*NotificationDeliveryMechanism, error) {
	vals := map[string]NotificationDeliveryMechanism{
		"email": NotificationDeliveryMechanismEmail,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotificationDeliveryMechanism(input)
	return &out, nil
}

type NotificationLevel string

const (
	NotificationLevelAll      NotificationLevel = "All"
	NotificationLevelCritical NotificationLevel = "Critical"
	NotificationLevelNone     NotificationLevel = "None"
)

func PossibleValuesForNotificationLevel() []string {
	return []string{
		string(NotificationLevelAll),
		string(NotificationLevelCritical),
		string(NotificationLevelNone),
	}
}

func (s *NotificationLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotificationLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotificationLevel(input string) (*NotificationLevel, error) {
	vals := map[string]NotificationLevel{
		"all":      NotificationLevelAll,
		"critical": NotificationLevelCritical,
		"none":     NotificationLevelNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotificationLevel(input)
	return &out, nil
}

type RecipientType string

const (
	RecipientTypeAdmin     RecipientType = "Admin"
	RecipientTypeApprover  RecipientType = "Approver"
	RecipientTypeRequestor RecipientType = "Requestor"
)

func PossibleValuesForRecipientType() []string {
	return []string{
		string(RecipientTypeAdmin),
		string(RecipientTypeApprover),
		string(RecipientTypeRequestor),
	}
}

func (s *RecipientType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecipientType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecipientType(input string) (*RecipientType, error) {
	vals := map[string]RecipientType{
		"admin":     RecipientTypeAdmin,
		"approver":  RecipientTypeApprover,
		"requestor": RecipientTypeRequestor,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecipientType(input)
	return &out, nil
}

type RoleManagementPolicyRuleType string

const (
	RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule              RoleManagementPolicyRuleType = "RoleManagementPolicyApprovalRule"
	RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule RoleManagementPolicyRuleType = "RoleManagementPolicyAuthenticationContextRule"
	RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule            RoleManagementPolicyRuleType = "RoleManagementPolicyEnablementRule"
	RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule            RoleManagementPolicyRuleType = "RoleManagementPolicyExpirationRule"
	RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule          RoleManagementPolicyRuleType = "RoleManagementPolicyNotificationRule"
)

func PossibleValuesForRoleManagementPolicyRuleType() []string {
	return []string{
		string(RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule),
		string(RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule),
		string(RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule),
		string(RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule),
		string(RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule),
	}
}

func (s *RoleManagementPolicyRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoleManagementPolicyRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoleManagementPolicyRuleType(input string) (*RoleManagementPolicyRuleType, error) {
	vals := map[string]RoleManagementPolicyRuleType{
		"rolemanagementpolicyapprovalrule":              RoleManagementPolicyRuleTypeRoleManagementPolicyApprovalRule,
		"rolemanagementpolicyauthenticationcontextrule": RoleManagementPolicyRuleTypeRoleManagementPolicyAuthenticationContextRule,
		"rolemanagementpolicyenablementrule":            RoleManagementPolicyRuleTypeRoleManagementPolicyEnablementRule,
		"rolemanagementpolicyexpirationrule":            RoleManagementPolicyRuleTypeRoleManagementPolicyExpirationRule,
		"rolemanagementpolicynotificationrule":          RoleManagementPolicyRuleTypeRoleManagementPolicyNotificationRule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoleManagementPolicyRuleType(input)
	return &out, nil
}

type UserType string

const (
	UserTypeGroup UserType = "Group"
	UserTypeUser  UserType = "User"
)

func PossibleValuesForUserType() []string {
	return []string{
		string(UserTypeGroup),
		string(UserTypeUser),
	}
}

func (s *UserType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUserType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUserType(input string) (*UserType, error) {
	vals := map[string]UserType{
		"group": UserTypeGroup,
		"user":  UserTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UserType(input)
	return &out, nil
}

package rolemanagementpolicyassignments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

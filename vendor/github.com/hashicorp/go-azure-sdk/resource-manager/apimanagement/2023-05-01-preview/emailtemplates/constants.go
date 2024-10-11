package emailtemplates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TemplateName string

const (
	TemplateNameAccountClosedDeveloper                            TemplateName = "accountClosedDeveloper"
	TemplateNameApplicationApprovedNotificationMessage            TemplateName = "applicationApprovedNotificationMessage"
	TemplateNameConfirmSignUpIdentityDefault                      TemplateName = "confirmSignUpIdentityDefault"
	TemplateNameEmailChangeIdentityDefault                        TemplateName = "emailChangeIdentityDefault"
	TemplateNameInviteUserNotificationMessage                     TemplateName = "inviteUserNotificationMessage"
	TemplateNameNewCommentNotificationMessage                     TemplateName = "newCommentNotificationMessage"
	TemplateNameNewDeveloperNotificationMessage                   TemplateName = "newDeveloperNotificationMessage"
	TemplateNameNewIssueNotificationMessage                       TemplateName = "newIssueNotificationMessage"
	TemplateNamePasswordResetByAdminNotificationMessage           TemplateName = "passwordResetByAdminNotificationMessage"
	TemplateNamePasswordResetIdentityDefault                      TemplateName = "passwordResetIdentityDefault"
	TemplateNamePurchaseDeveloperNotificationMessage              TemplateName = "purchaseDeveloperNotificationMessage"
	TemplateNameQuotaLimitApproachingDeveloperNotificationMessage TemplateName = "quotaLimitApproachingDeveloperNotificationMessage"
	TemplateNameRejectDeveloperNotificationMessage                TemplateName = "rejectDeveloperNotificationMessage"
	TemplateNameRequestDeveloperNotificationMessage               TemplateName = "requestDeveloperNotificationMessage"
)

func PossibleValuesForTemplateName() []string {
	return []string{
		string(TemplateNameAccountClosedDeveloper),
		string(TemplateNameApplicationApprovedNotificationMessage),
		string(TemplateNameConfirmSignUpIdentityDefault),
		string(TemplateNameEmailChangeIdentityDefault),
		string(TemplateNameInviteUserNotificationMessage),
		string(TemplateNameNewCommentNotificationMessage),
		string(TemplateNameNewDeveloperNotificationMessage),
		string(TemplateNameNewIssueNotificationMessage),
		string(TemplateNamePasswordResetByAdminNotificationMessage),
		string(TemplateNamePasswordResetIdentityDefault),
		string(TemplateNamePurchaseDeveloperNotificationMessage),
		string(TemplateNameQuotaLimitApproachingDeveloperNotificationMessage),
		string(TemplateNameRejectDeveloperNotificationMessage),
		string(TemplateNameRequestDeveloperNotificationMessage),
	}
}

func (s *TemplateName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTemplateName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTemplateName(input string) (*TemplateName, error) {
	vals := map[string]TemplateName{
		"accountcloseddeveloper":                            TemplateNameAccountClosedDeveloper,
		"applicationapprovednotificationmessage":            TemplateNameApplicationApprovedNotificationMessage,
		"confirmsignupidentitydefault":                      TemplateNameConfirmSignUpIdentityDefault,
		"emailchangeidentitydefault":                        TemplateNameEmailChangeIdentityDefault,
		"inviteusernotificationmessage":                     TemplateNameInviteUserNotificationMessage,
		"newcommentnotificationmessage":                     TemplateNameNewCommentNotificationMessage,
		"newdevelopernotificationmessage":                   TemplateNameNewDeveloperNotificationMessage,
		"newissuenotificationmessage":                       TemplateNameNewIssueNotificationMessage,
		"passwordresetbyadminnotificationmessage":           TemplateNamePasswordResetByAdminNotificationMessage,
		"passwordresetidentitydefault":                      TemplateNamePasswordResetIdentityDefault,
		"purchasedevelopernotificationmessage":              TemplateNamePurchaseDeveloperNotificationMessage,
		"quotalimitapproachingdevelopernotificationmessage": TemplateNameQuotaLimitApproachingDeveloperNotificationMessage,
		"rejectdevelopernotificationmessage":                TemplateNameRejectDeveloperNotificationMessage,
		"requestdevelopernotificationmessage":               TemplateNameRequestDeveloperNotificationMessage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TemplateName(input)
	return &out, nil
}

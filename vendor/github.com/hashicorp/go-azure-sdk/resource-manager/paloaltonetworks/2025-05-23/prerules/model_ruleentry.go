package prerules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleEntry struct {
	ActionType                   *ActionEnum             `json:"actionType,omitempty"`
	Applications                 *[]string               `json:"applications,omitempty"`
	AuditComment                 *string                 `json:"auditComment,omitempty"`
	Category                     *Category               `json:"category,omitempty"`
	DecryptionRuleType           *DecryptionRuleTypeEnum `json:"decryptionRuleType,omitempty"`
	Description                  *string                 `json:"description,omitempty"`
	Destination                  *DestinationAddr        `json:"destination,omitempty"`
	EnableLogging                *StateEnum              `json:"enableLogging,omitempty"`
	Etag                         *string                 `json:"etag,omitempty"`
	InboundInspectionCertificate *string                 `json:"inboundInspectionCertificate,omitempty"`
	NegateDestination            *BooleanEnum            `json:"negateDestination,omitempty"`
	NegateSource                 *BooleanEnum            `json:"negateSource,omitempty"`
	Priority                     *int64                  `json:"priority,omitempty"`
	Protocol                     *string                 `json:"protocol,omitempty"`
	ProtocolPortList             *[]string               `json:"protocolPortList,omitempty"`
	ProvisioningState            *ProvisioningState      `json:"provisioningState,omitempty"`
	RuleName                     string                  `json:"ruleName"`
	RuleState                    *StateEnum              `json:"ruleState,omitempty"`
	Source                       *SourceAddr             `json:"source,omitempty"`
	Tags                         *[]TagInfo              `json:"tags,omitempty"`
}

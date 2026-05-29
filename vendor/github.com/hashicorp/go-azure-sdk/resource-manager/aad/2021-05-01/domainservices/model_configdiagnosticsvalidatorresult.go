package domainservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigDiagnosticsValidatorResult struct {
	Issues                      *[]ConfigDiagnosticsValidatorResultIssue `json:"issues,omitempty"`
	ReplicaSetSubnetDisplayName *string                                  `json:"replicaSetSubnetDisplayName,omitempty"`
	Status                      *Status                                  `json:"status,omitempty"`
	ValidatorId                 *string                                  `json:"validatorId,omitempty"`
}

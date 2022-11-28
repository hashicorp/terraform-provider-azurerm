package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPolicyAction struct {
	BaseBlob *ManagementPolicyBaseBlob `json:"baseBlob"`
	Snapshot *ManagementPolicySnapShot `json:"snapshot"`
	Version  *ManagementPolicyVersion  `json:"version"`
}

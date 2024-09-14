package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TroubleshootingAdditionalProperties struct {
	UnhealthyReplicaInfo *UnhealthyReplicaInfo `json:"unhealthyReplicaInfo,omitempty"`
}

package instancefailovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstanceFailoverGroupReadWriteEndpoint struct {
	FailoverPolicy                         ReadWriteEndpointFailoverPolicy `json:"failoverPolicy"`
	FailoverWithDataLossGracePeriodMinutes *int64                          `json:"failoverWithDataLossGracePeriodMinutes,omitempty"`
}

package failovergroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FailoverGroupReadWriteEndpoint struct {
	FailoverPolicy                         ReadWriteEndpointFailoverPolicy `json:"failoverPolicy"`
	FailoverWithDataLossGracePeriodMinutes *int64                          `json:"failoverWithDataLossGracePeriodMinutes,omitempty"`
}

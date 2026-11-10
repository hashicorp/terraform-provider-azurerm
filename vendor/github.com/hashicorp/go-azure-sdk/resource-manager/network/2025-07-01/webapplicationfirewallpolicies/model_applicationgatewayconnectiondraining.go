package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayConnectionDraining struct {
	DrainTimeoutInSec int64 `json:"drainTimeoutInSec"`
	Enabled           bool  `json:"enabled"`
}

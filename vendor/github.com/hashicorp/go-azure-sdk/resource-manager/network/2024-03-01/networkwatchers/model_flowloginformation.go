package networkwatchers

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FlowLogInformation struct {
	FlowAnalyticsConfiguration *TrafficAnalyticsProperties        `json:"flowAnalyticsConfiguration,omitempty"`
	Identity                   *identity.SystemAndUserAssignedMap `json:"identity,omitempty"`
	Properties                 FlowLogProperties                  `json:"properties"`
	TargetResourceId           string                             `json:"targetResourceId"`
}

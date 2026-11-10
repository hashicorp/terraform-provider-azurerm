package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorEndpointScope struct {
	Exclude *[]ConnectionMonitorEndpointScopeItem `json:"exclude,omitempty"`
	Include *[]ConnectionMonitorEndpointScopeItem `json:"include,omitempty"`
}

package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MultiRegionSettings struct {
	Regions       *[]RegionSetting `json:"regions,omitempty"`
	RoutingMethod *RoutingMethods  `json:"routingMethod,omitempty"`
}

package communicationsgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceRegionProperties struct {
	Name                    string                  `json:"name"`
	PrimaryRegionProperties PrimaryRegionProperties `json:"primaryRegionProperties"`
}

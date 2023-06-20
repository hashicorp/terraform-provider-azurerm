package contactprofile

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfilesProperties struct {
	AutoTrackingConfiguration    *AutoTrackingConfiguration                    `json:"autoTrackingConfiguration,omitempty"`
	EventHubUri                  *string                                       `json:"eventHubUri,omitempty"`
	Links                        []ContactProfileLink                          `json:"links"`
	MinimumElevationDegrees      *float64                                      `json:"minimumElevationDegrees,omitempty"`
	MinimumViableContactDuration *string                                       `json:"minimumViableContactDuration,omitempty"`
	NetworkConfiguration         ContactProfilesPropertiesNetworkConfiguration `json:"networkConfiguration"`
	ProvisioningState            *ProvisioningState                            `json:"provisioningState,omitempty"`
	ThirdPartyConfigurations     *[]ContactProfileThirdPartyConfiguration      `json:"thirdPartyConfigurations,omitempty"`
}

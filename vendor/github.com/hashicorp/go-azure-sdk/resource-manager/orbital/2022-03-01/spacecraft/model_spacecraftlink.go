package spacecraft

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SpacecraftLink struct {
	Authorizations     *[]AuthorizedGroundstation `json:"authorizations,omitempty"`
	BandwidthMHz       float64                    `json:"bandwidthMHz"`
	CenterFrequencyMHz float64                    `json:"centerFrequencyMHz"`
	Direction          Direction                  `json:"direction"`
	Name               string                     `json:"name"`
	Polarization       Polarization               `json:"polarization"`
}

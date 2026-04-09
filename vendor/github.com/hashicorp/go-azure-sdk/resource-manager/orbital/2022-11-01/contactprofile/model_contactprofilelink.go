package contactprofile

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfileLink struct {
	Channels            []ContactProfileLinkChannel `json:"channels"`
	Direction           Direction                   `json:"direction"`
	EirpdBW             *float64                    `json:"eirpdBW,omitempty"`
	GainOverTemperature *float64                    `json:"gainOverTemperature,omitempty"`
	Name                string                      `json:"name"`
	Polarization        Polarization                `json:"polarization"`
}

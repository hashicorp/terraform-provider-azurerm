package contactprofile

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactProfileLinkChannel struct {
	BandwidthMHz              float64  `json:"bandwidthMHz"`
	CenterFrequencyMHz        float64  `json:"centerFrequencyMHz"`
	DecodingConfiguration     *string  `json:"decodingConfiguration,omitempty"`
	DemodulationConfiguration *string  `json:"demodulationConfiguration,omitempty"`
	EncodingConfiguration     *string  `json:"encodingConfiguration,omitempty"`
	EndPoint                  EndPoint `json:"endPoint"`
	ModulationConfiguration   *string  `json:"modulationConfiguration,omitempty"`
	Name                      string   `json:"name"`
}

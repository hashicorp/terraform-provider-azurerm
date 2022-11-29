package contentkeypolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyPlayReadyPlayRight struct {
	AgcAndColorStripeRestriction                       *int64                                                        `json:"agcAndColorStripeRestriction,omitempty"`
	AllowPassingVideoContentToUnknownOutput            ContentKeyPolicyPlayReadyUnknownOutputPassingOption           `json:"allowPassingVideoContentToUnknownOutput"`
	AnalogVideoOpl                                     *int64                                                        `json:"analogVideoOpl,omitempty"`
	CompressedDigitalAudioOpl                          *int64                                                        `json:"compressedDigitalAudioOpl,omitempty"`
	CompressedDigitalVideoOpl                          *int64                                                        `json:"compressedDigitalVideoOpl,omitempty"`
	DigitalVideoOnlyContentRestriction                 bool                                                          `json:"digitalVideoOnlyContentRestriction"`
	ExplicitAnalogTelevisionOutputRestriction          *ContentKeyPolicyPlayReadyExplicitAnalogTelevisionRestriction `json:"explicitAnalogTelevisionOutputRestriction,omitempty"`
	FirstPlayExpiration                                *string                                                       `json:"firstPlayExpiration,omitempty"`
	ImageConstraintForAnalogComponentVideoRestriction  bool                                                          `json:"imageConstraintForAnalogComponentVideoRestriction"`
	ImageConstraintForAnalogComputerMonitorRestriction bool                                                          `json:"imageConstraintForAnalogComputerMonitorRestriction"`
	ScmsRestriction                                    *int64                                                        `json:"scmsRestriction,omitempty"`
	UncompressedDigitalAudioOpl                        *int64                                                        `json:"uncompressedDigitalAudioOpl,omitempty"`
	UncompressedDigitalVideoOpl                        *int64                                                        `json:"uncompressedDigitalVideoOpl,omitempty"`
}

package expressroutecrossconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvertisedPublicPrefixProperties struct {
	Prefix          *string                                          `json:"prefix,omitempty"`
	Signature       *string                                          `json:"signature,omitempty"`
	ValidationId    *string                                          `json:"validationId,omitempty"`
	ValidationState *AdvertisedPublicPrefixPropertiesValidationState `json:"validationState,omitempty"`
}

package hubs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationResult struct {
	ApplicationPlatform *string `json:"applicationPlatform,omitempty"`
	Outcome             *string `json:"outcome,omitempty"`
	PnsHandle           *string `json:"pnsHandle,omitempty"`
	RegistrationId      *string `json:"registrationId,omitempty"`
}

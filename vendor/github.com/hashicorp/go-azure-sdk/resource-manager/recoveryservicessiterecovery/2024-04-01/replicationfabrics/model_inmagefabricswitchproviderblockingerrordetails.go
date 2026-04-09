package replicationfabrics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageFabricSwitchProviderBlockingErrorDetails struct {
	ErrorCode              *string            `json:"errorCode,omitempty"`
	ErrorMessage           *string            `json:"errorMessage,omitempty"`
	ErrorMessageParameters *map[string]string `json:"errorMessageParameters,omitempty"`
	ErrorTags              *map[string]string `json:"errorTags,omitempty"`
	PossibleCauses         *string            `json:"possibleCauses,omitempty"`
	RecommendedAction      *string            `json:"recommendedAction,omitempty"`
}

package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyPoolProfileDefinition struct {
	FailContainerGroupCreateOnReuseFailure *bool   `json:"failContainerGroupCreateOnReuseFailure,omitempty"`
	Id                                     *string `json:"id,omitempty"`
}

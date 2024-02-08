package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigServerSettingsValidateResult struct {
	Details *[]ConfigServerSettingsErrorRecord `json:"details,omitempty"`
	IsValid *bool                              `json:"isValid,omitempty"`
}

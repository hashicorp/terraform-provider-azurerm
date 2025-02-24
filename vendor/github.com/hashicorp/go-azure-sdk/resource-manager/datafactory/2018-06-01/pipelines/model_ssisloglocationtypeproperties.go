package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSISLogLocationTypeProperties struct {
	AccessCredential   *SSISAccessCredential `json:"accessCredential,omitempty"`
	LogRefreshInterval *string               `json:"logRefreshInterval,omitempty"`
}

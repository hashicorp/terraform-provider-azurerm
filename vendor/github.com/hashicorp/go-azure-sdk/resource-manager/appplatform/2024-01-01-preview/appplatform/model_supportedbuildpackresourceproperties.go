package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupportedBuildpackResourceProperties struct {
	BuildpackId *string `json:"buildpackId,omitempty"`
	Version     *string `json:"version,omitempty"`
}

package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyComputeScaleProperties struct {
	DataIntegrationUnit *int64 `json:"dataIntegrationUnit,omitempty"`
	TimeToLive          *int64 `json:"timeToLive,omitempty"`
}

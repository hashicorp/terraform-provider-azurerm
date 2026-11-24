package agentpools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScaleProfile struct {
	Autoscale *[]AutoScaleProfile   `json:"autoscale,omitempty"`
	Manual    *[]ManualScaleProfile `json:"manual,omitempty"`
}

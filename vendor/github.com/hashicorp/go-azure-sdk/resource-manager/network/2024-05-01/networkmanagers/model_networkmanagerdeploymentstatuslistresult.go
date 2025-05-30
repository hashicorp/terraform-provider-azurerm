package networkmanagers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerDeploymentStatusListResult struct {
	SkipToken *string                           `json:"skipToken,omitempty"`
	Value     *[]NetworkManagerDeploymentStatus `json:"value,omitempty"`
}

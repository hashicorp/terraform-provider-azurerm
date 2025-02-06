package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleConfigurationUpdateParameter struct {
	Autoscale *Autoscale `json:"autoscale,omitempty"`
}

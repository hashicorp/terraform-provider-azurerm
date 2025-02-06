package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExcludedServicesConfig struct {
	ExcludedServicesConfigId *string `json:"excludedServicesConfigId,omitempty"`
	ExcludedServicesList     *string `json:"excludedServicesList,omitempty"`
}

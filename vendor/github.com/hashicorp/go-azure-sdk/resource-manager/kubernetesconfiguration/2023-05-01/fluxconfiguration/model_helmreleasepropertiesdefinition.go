package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HelmReleasePropertiesDefinition struct {
	FailureCount        *int64                     `json:"failureCount,omitempty"`
	HelmChartRef        *ObjectReferenceDefinition `json:"helmChartRef,omitempty"`
	InstallFailureCount *int64                     `json:"installFailureCount,omitempty"`
	LastRevisionApplied *int64                     `json:"lastRevisionApplied,omitempty"`
	UpgradeFailureCount *int64                     `json:"upgradeFailureCount,omitempty"`
}

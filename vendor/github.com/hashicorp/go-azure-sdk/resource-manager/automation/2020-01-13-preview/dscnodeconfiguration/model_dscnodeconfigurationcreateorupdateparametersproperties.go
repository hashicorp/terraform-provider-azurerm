package dscnodeconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscNodeConfigurationCreateOrUpdateParametersProperties struct {
	Configuration                   DscConfigurationAssociationProperty `json:"configuration"`
	IncrementNodeConfigurationBuild *bool                               `json:"incrementNodeConfigurationBuild,omitempty"`
	Source                          ContentSource                       `json:"source"`
}

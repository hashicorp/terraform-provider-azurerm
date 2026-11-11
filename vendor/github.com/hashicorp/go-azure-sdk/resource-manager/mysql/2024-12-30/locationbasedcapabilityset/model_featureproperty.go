package locationbasedcapabilityset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FeatureProperty struct {
	FeatureName  *string `json:"featureName,omitempty"`
	FeatureValue *string `json:"featureValue,omitempty"`
}

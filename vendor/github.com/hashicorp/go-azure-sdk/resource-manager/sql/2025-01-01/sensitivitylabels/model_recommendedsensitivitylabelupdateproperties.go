package sensitivitylabels

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecommendedSensitivityLabelUpdateProperties struct {
	Column string                                `json:"column"`
	Op     RecommendedSensitivityLabelUpdateKind `json:"op"`
	Schema string                                `json:"schema"`
	Table  string                                `json:"table"`
}

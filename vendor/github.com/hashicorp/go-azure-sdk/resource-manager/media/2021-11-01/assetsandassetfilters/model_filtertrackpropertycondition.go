package assetsandassetfilters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FilterTrackPropertyCondition struct {
	Operation FilterTrackPropertyCompareOperation `json:"operation"`
	Property  FilterTrackPropertyType             `json:"property"`
	Value     string                              `json:"value"`
}

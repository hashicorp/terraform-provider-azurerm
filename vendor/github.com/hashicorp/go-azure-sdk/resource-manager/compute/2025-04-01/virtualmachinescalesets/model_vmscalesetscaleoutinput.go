package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMScaleSetScaleOutInput struct {
	Capacity   int64                              `json:"capacity"`
	Properties *VMScaleSetScaleOutInputProperties `json:"properties,omitempty"`
}

package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModelDeprecationInfo struct {
	FineTune  *string `json:"fineTune,omitempty"`
	Inference *string `json:"inference,omitempty"`
}

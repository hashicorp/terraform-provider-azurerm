package diagnostic

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataMaskingEntity struct {
	Mode  *DataMaskingMode `json:"mode,omitempty"`
	Value *string          `json:"value,omitempty"`
}

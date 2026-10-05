package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DWCopyCommandSettings struct {
	AdditionalOptions *map[string]string           `json:"additionalOptions,omitempty"`
	DefaultValues     *[]DWCopyCommandDefaultValue `json:"defaultValues,omitempty"`
}

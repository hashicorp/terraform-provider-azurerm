package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SitePhpErrorLogFlagProperties struct {
	LocalLogErrors           *string `json:"localLogErrors,omitempty"`
	LocalLogErrorsMaxLength  *string `json:"localLogErrorsMaxLength,omitempty"`
	MasterLogErrors          *string `json:"masterLogErrors,omitempty"`
	MasterLogErrorsMaxLength *string `json:"masterLogErrorsMaxLength,omitempty"`
}

package connectionanalyzers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtocolSettings struct {
	HTTPConfiguration *HTTPConfiguration `json:"httpConfiguration,omitempty"`
	Protocol          *Protocol          `json:"protocol,omitempty"`
}

package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebActivityTypeProperties struct {
	Authentication        *WebActivityAuthentication   `json:"authentication,omitempty"`
	Body                  *interface{}                 `json:"body,omitempty"`
	ConnectVia            *IntegrationRuntimeReference `json:"connectVia,omitempty"`
	Datasets              *[]DatasetReference          `json:"datasets,omitempty"`
	DisableCertValidation *bool                        `json:"disableCertValidation,omitempty"`
	HTTPRequestTimeout    *interface{}                 `json:"httpRequestTimeout,omitempty"`
	Headers               *map[string]interface{}      `json:"headers,omitempty"`
	LinkedServices        *[]LinkedServiceReference    `json:"linkedServices,omitempty"`
	Method                WebActivityMethod            `json:"method"`
	TurnOffAsync          *bool                        `json:"turnOffAsync,omitempty"`
	Url                   interface{}                  `json:"url"`
}

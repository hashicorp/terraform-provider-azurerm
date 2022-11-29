package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiResourceProperties struct {
	ApiDefinitionUrl     *string                         `json:"apiDefinitionUrl,omitempty"`
	ApiDefinitions       *ApiResourceDefinitions         `json:"apiDefinitions"`
	BackendService       *ApiResourceBackendService      `json:"backendService"`
	Capabilities         *[]string                       `json:"capabilities,omitempty"`
	ConnectionParameters *map[string]ConnectionParameter `json:"connectionParameters,omitempty"`
	GeneralInformation   *ApiResourceGeneralInformation  `json:"generalInformation"`
	Metadata             *ApiResourceMetadata            `json:"metadata"`
	Name                 *string                         `json:"name,omitempty"`
	Policies             *ApiResourcePolicies            `json:"policies"`
	RuntimeUrls          *[]string                       `json:"runtimeUrls,omitempty"`
	Swagger              *interface{}                    `json:"swagger,omitempty"`
}

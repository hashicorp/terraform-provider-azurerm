package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdmCredentialProperties struct {
	AuthTokenURL string `json:"authTokenUrl"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

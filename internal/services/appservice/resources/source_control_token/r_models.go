// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package source_control_token

type AppServiceSourceControlTokenModel struct {
	Token       string `tfschema:"token"`
	TokenSecret string `tfschema:"token_secret"`
	Type        string `tfschema:"type"`
}

func (r Resource) ModelObject() any {
	return &AppServiceSourceControlTokenModel{}
}

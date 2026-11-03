// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package nginx_api_key

type APIKeyModel struct {
	Name              string `tfschema:"name"`
	NginxDeploymentId string `tfschema:"nginx_deployment_id"`
	SecretText        string `tfschema:"secret_text"`
	Hint              string `tfschema:"hint"`
	EndDateTime       string `tfschema:"end_date_time"`
}

func (m Resource) ModelObject() any {
	return &APIKeyModel{}
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_blocklist

type cognitiveRaiBlocklistModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Description        string            `tfschema:"description"`
	Tags               map[string]string `tfschema:"tags"`
}

func (c Resource) ModelObject() any {
	return &cognitiveRaiBlocklistModel{}
}

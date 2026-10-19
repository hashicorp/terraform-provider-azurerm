// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_policy

type AccountRaiPolicyResourceModel struct {
	Name           string                          `tfschema:"name"`
	AccountId      string                          `tfschema:"cognitive_account_id"`
	BasePolicyName string                          `tfschema:"base_policy_name"`
	ContentFilter  []AccountRaiPolicyContentFilter `tfschema:"content_filter"`
	Mode           string                          `tfschema:"mode"`
	Tags           map[string]string               `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AccountRaiPolicyResourceModel{}
}

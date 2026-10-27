// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule_association

type DataCollectionRuleAssociationModel struct {
	Name                     string `tfschema:"name"`
	TargetResourceId         string `tfschema:"target_resource_id"`
	DataCollectionEndpointId string `tfschema:"data_collection_endpoint_id"`
	DataCollectionRuleId     string `tfschema:"data_collection_rule_id"`
	Description              string `tfschema:"description"`
}

func (r Resource) ModelObject() any {
	return &DataCollectionRuleAssociationModel{}
}

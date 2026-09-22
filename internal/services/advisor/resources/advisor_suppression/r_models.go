// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package advisor_suppression

type AdvisorSuppressionResourceModel struct {
	Name             string `tfschema:"name"`
	SuppressionID    string `tfschema:"suppression_id"`
	RecommendationID string `tfschema:"recommendation_id"`
	ResourceID       string `tfschema:"resource_id"`
	TTL              string `tfschema:"ttl"`
}

func (Resource) ModelObject() interface{} {
	return &AdvisorSuppressionResourceModel{}
}

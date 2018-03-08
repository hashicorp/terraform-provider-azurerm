package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var streamAnalyticsTransformationType = "Microsoft.StreamAnalytics/streamingjobs/transformations"

func streamAnalyticsTransformationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(3, 64),
				},
				"streaming_units": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"query": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

func streamAnalyticsTransformationFromSchema(transSchema interface{}) *streamanalytics.Transformation {
	transMap := transSchema.(map[string]interface{})
	name := transMap["name"].(string)
	streamingUnits := transMap["streaming_units"].(int) // defaults to 1 so no validation needed
	query := transMap["query"].(string)

	streamingUnits32 := int32(streamingUnits)
	return &streamanalytics.Transformation{
		Name: &name,
		Type: &streamAnalyticsTransformationType,
		TransformationProperties: &streamanalytics.TransformationProperties{
			StreamingUnits: &streamingUnits32,
			Query:          &query,
		},
	}
}

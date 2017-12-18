package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	DelimSpace       = " "
	DelimComma       = ","
	DelimTab         = "\t"
	DelimSemiColon   = ";"
	DelimVerticalBar = "|"
)

// Allbut returns a list all possible datasource other than the specified arg
func Allbut(source string) (ret []string) {
	sources := []string{
		"blob",
		"table",
		"event_hub",
		"sql_database",
		"service_bus_queues",
		"service_bus_topics",
		"documentdb",
	}
	for _, val := range sources {
		if val != source {
			ret = append(ret, val)
		}
	}

	return

}

func streamAnalyticsOutputSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Elem: map[string]*schema.Schema{
			"serialization": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.TypeCsv),
								string(streamanalytics.TypeAvro),
								string(streamanalytics.TypeJSON),
							}, false),
						},
						"field_delimiter": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								DelimTab,
								DelimComma,
								DelimSemiColon,
								DelimVerticalBar,
								DelimSpace,
							}, false),
						},
						"encoding": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.Array),
								string(streamanalytics.LineSeparated),
							}, false),
						},
					},
				},
			},
			"datasource": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"blob": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("blob"),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_account_name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"storage_account_key": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"container": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"path_pattern": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"date_format": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_format": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_partition_count": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"table": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("table"),
							Elem: map[string]*schema.Schema{
								"account_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"account_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"table": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"partition_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"row_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"columns_to_remove": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
								"batch_size": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
									Default:  100,
								},
							},
						},
						"event_hub": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("event_hub"),
							Elem: map[string]*schema.Schema{

								"namespace": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"event_hub_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"partition_key": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
						"sql_database": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("sql_database"),
							Elem: map[string]*schema.Schema{
								"server": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"database": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"user": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"password": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"table": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
						"service_bus_queues": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("service_bus_queues"),
							Elem: map[string]*schema.Schema{
								"namespace": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"queue_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"property_columns": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
						"service_bus_topics": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("service_bus_topics"),
							Elem: map[string]*schema.Schema{
								"namespace": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"shared_access_policy_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"topic_name": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"property_columns": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
						"documentdb": &schema.Schema{
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: Allbut("documentdb"),
							Elem: map[string]*schema.Schema{
								"account_id": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"account_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"database": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"collection_name_pattern": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"partition_key": &schema.Schema{
									Type:     schema.TypeString,
									Required: true,
								},
								"document_id": &schema.Schema{
									Type:     schema.TypeString,
									Required: false,
								},
							},
						},
					},
				},
			},
		},
	}
}

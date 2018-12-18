package azure

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// SchemaBatchPoolImageReference returns the schema for a Batch pool image reference
func SchemaBatchPoolImageReference() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"publisher": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"offer": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"sku": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"version": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},
		},
		Set: resourceArmBatchPoolImageReferenceHash,
	}
}

// SchemaBatchPoolImageReferenceForDataSource returns the schema for a Batch pool image reference data source
func SchemaBatchPoolImageReferenceForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"publisher": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"offer": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"sku": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"version": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
		Set: resourceArmBatchPoolImageReferenceHash,
	}
}

// SchemaBatchPoolFixedScale returns the schema for the Batch pool fixed scale settings
func SchemaBatchPoolFixedScale() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"target_dedicated_nodes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"target_low_priority_nodes": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"resize_timeout": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "PT15M",
				},
			},
		},
	}
}

// SchemaBatchPoolFixedScaleForDataSource returns the schema for the Batch pool fixed scale settings data source
func SchemaBatchPoolFixedScaleForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"target_dedicated_nodes": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"target_low_priority_nodes": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"resize_timeout": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// SchemaBatchPoolAutoScale returns the schema for the Batch pool auto scale settings
func SchemaBatchPoolAutoScale() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"evaluation_interval": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "PT15M",
				},
				"formula": {
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}

// SchemaBatchPoolAutoScaleForDataSource returns the schema for the Batch pool auto scale settings data source
func SchemaBatchPoolAutoScaleForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"evaluation_interval": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"formula": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// SchemaBatchPoolStartTask returns the schema for a Batch pool start task
func SchemaBatchPoolStartTask() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"command_line": {
					Type:     schema.TypeString,
					Required: true,
				},

				"max_task_retry_count": {
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},

				"wait_for_success": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"environment": {
					Type:     schema.TypeMap,
					Optional: true,
				},

				"user_identity": schemaBatchPoolStartTaskUserIdentity(),
			},
		},
	}
}

// SchemaBatchPoolStartTaskForDataSource returns the schema for a Batch pool start task
func SchemaBatchPoolStartTaskForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"command_line": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"max_task_retry_count": {
					Type:     schema.TypeInt,
					Computed: true,
				},

				"wait_for_success": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"environment": {
					Type:     schema.TypeMap,
					Computed: true,
				},

				"user_identity": schemaBatchPoolStartTaskUserIdentityForDataSource(),
			},
		},
	}
}

// FlattenBatchPoolAutoScaleSettings flattens the auto scale settings for a Batch pool
func FlattenBatchPoolAutoScaleSettings(settings *batch.AutoScaleSettings) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result["evaluation_interval"] = settings.EvaluationInterval
	result["formula"] = settings.Formula

	return append(results, result)
}

// FlattenBatchPoolFixedScaleSettings flattens the fixed scale settings for a Batch pool
func FlattenBatchPoolFixedScaleSettings(settings *batch.FixedScaleSettings) []interface{} {
	results := make([]interface{}, 0)
	result := make(map[string]interface{})

	if settings == nil {
		log.Printf("[DEBUG] settings is nil")
		return results
	}

	result["target_dedicated_nodes"] = settings.TargetDedicatedNodes
	result["target_low_priority_nodes"] = settings.TargetLowPriorityNodes
	result["resize_timeout"] = settings.ResizeTimeout

	return append(results, result)
}

// ValidateAzureRMBatchPoolName validates the name of a Batch pool
func ValidateAzureRMBatchPoolName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"any combination of alphanumeric characters including hyphens and underscores are allowed in %q: %q", k, value))
	}

	if 1 > len(value) {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 64 characters: %q %d", k, value, len(value)))
	}

	return warnings, errors
}

// FlattenBatchPoolImageReference flattens the Batch pool image reference
func FlattenBatchPoolImageReference(image *batch.ImageReference) *schema.Set {
	result := make(map[string]interface{})
	if image.Publisher != nil {
		result["publisher"] = *image.Publisher
	}
	if image.Offer != nil {
		result["offer"] = *image.Offer
	}
	if image.Sku != nil {
		result["sku"] = *image.Sku
	}
	if image.Version != nil {
		result["version"] = *image.Version
	}
	if image.ID != nil {
		result["id"] = *image.ID
	}

	return schema.NewSet(resourceArmBatchPoolImageReferenceHash, []interface{}{result})
}

// ExpandBatchPoolImageReference expands a Set into an image reference
func ExpandBatchPoolImageReference(set *schema.Set) (*batch.ImageReference, error) {
	if set == nil || set.Len() == 0 {
		return nil, fmt.Errorf("Error: storage image reference should be defined")
	}

	storageImageRef := set.List()[0].(map[string]interface{})

	storageImageRefOffer, storageImageRefOfferOk := storageImageRef["offer"].(string)
	if !storageImageRefOfferOk {
		return nil, fmt.Errorf("Error: storage image reference offer should be defined")
	}

	storageImageRefPublisher, storageImageRefPublisherOK := storageImageRef["publisher"].(string)
	if !storageImageRefPublisherOK {
		return nil, fmt.Errorf("Error: storage image reference publisher should be defined")
	}

	storageImageRefSku, storageImageRefSkuOK := storageImageRef["sku"].(string)
	if !storageImageRefSkuOK {
		return nil, fmt.Errorf("Error: storage image reference sku should be defined")
	}

	storageImageRefVersion, storageImageRefVersionOK := storageImageRef["version"].(string)
	if !storageImageRefVersionOK {
		return nil, fmt.Errorf("Error: storage image reference version should be defined")
	}

	imageRef := &batch.ImageReference{
		Offer:     &storageImageRefOffer,
		Publisher: &storageImageRefPublisher,
		Sku:       &storageImageRefSku,
		Version:   &storageImageRefVersion,
	}

	return imageRef, nil
}

func schemaBatchPoolStartTaskUserIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user_name": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"auto_user": {
					Type:     schema.TypeSet,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"elevation_level": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(batch.NonAdmin),
								ValidateFunc: validation.StringInSlice([]string{
									string(batch.NonAdmin),
									string(batch.Admin),
								}, false),
							},
							"scope": {
								Type:     schema.TypeString,
								Optional: true,
								Default:  string(batch.AutoUserScopeTask),
								ValidateFunc: validation.StringInSlice([]string{
									string(batch.AutoUserScopeTask),
									string(batch.AutoUserScopePool),
								}, false),
							},
						},
					},
				},
			},
		},
	}
}

func schemaBatchPoolStartTaskUserIdentityForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"user_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"auto_user": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"elevation_level": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"scope": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmBatchPoolImageReferenceHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["publisher"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["offer"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["sku"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["id"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
		if v, ok := m["version"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return hashcode.String(buf.String())
}

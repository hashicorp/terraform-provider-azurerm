package iotoperations

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/!azure/azure-sdk-for-go/sdk/resourcemanager/iotoperations/armiotoperations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func resourceIotOperationsDataflowProfile() *schema.Resource {
    return &schema.Resource{
        Create: resourceIotOperationsDataflowProfileCreate,
        Read:   resourceIotOperationsDataflowProfileRead,
        Update: resourceIotOperationsDataflowProfileUpdate,
        Delete: resourceIotOperationsDataflowProfileDelete,
        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(30 * time.Minute),
            Update: schema.DefaultTimeout(30 * time.Minute),
            Delete: schema.DefaultTimeout(30 * time.Minute),
        },
        Schema: map[string]*schema.Schema{
            "dataflow_profile_name": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
                MinLen:       3,
                MaxLen:       63,
            },
            "instance_name": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
                MinLen:       3,
                MaxLen:       63,
            },
            "resource_group_name": {
                Type:         schema.TypeString,
                Required:     true,
                MinLen:       1,
                MaxLen:       90,
            },
            "subscription_id": {
                Type:         schema.TypeString,
                Required:     true,
                ValidateFunc: validation.IsUUID,
            },
            "api_version": {
                Type:         schema.TypeString,
                Required:     true,
                MinLen:       1,
            },
            "extended_location": {
                Type:     schema.TypeList,
                Optional: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "name": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "type": {
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                    },
                },
            },
            "properties": {
                Type:     schema.TypeList,
                Optional: true,
                MaxItems: 1,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "instance_count": {
                            Type:     schema.TypeInt,
                            Optional: true,
                        },
                        "diagnostics": {
                            Type:     schema.TypeList,
                            Optional: true,
                            MaxItems: 1,
                            Elem: &schema.Resource{
                                Schema: map[string]*schema.Schema{
                                    "logs": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "level": {
                                                    Type:     schema.TypeString,
                                                    Optional: true,
                                                },
                                            },
                                        },
                                    },
                                    "metrics": {
                                        Type:     schema.TypeList,
                                        Optional: true,
                                        MaxItems: 1,
                                        Elem: &schema.Resource{
                                            Schema: map[string]*schema.Schema{
                                                "prometheus_port": {
                                                    Type:     schema.TypeInt,
                                                    Optional: true,
                                                },
                                            },
                                        },
                                    },
                                },
                            },
                        },
                        "provisioning_state": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                    },
                },
            },
            "id": {
                Type:     schema.TypeString,
                Computed: true,
            },
            "system_data": {
                Type:     schema.TypeMap,
                Computed: true,
            },
        },
    }
}

func resourceIotOperationsDataflowProfileCreate(d *schema.ResourceData, meta interface{}) error {
    ctx := context.Background()

    name := d.Get("name").(string)
    rg := d.Get("resource_group_name").(string)
    instance := d.Get("instance_name").(string)

    top := meta.(*clients.Client)
    if top.IoTOperations == nil || top.IoTOperations.DataflowProfileClient == nil {
        return fmt.Errorf("iotoperations client is not configured")
    }
    svc := top.IoTOperations

    // build resource
    profile := armiotoperations.DataflowProfileResource{}
    if v, ok := d.GetOk("properties"); ok {
        if props, err := expandDataflowProfileProperties(v.([]interface{})); err == nil && props != nil {
            profile.Properties = props
        } else if err != nil {
            return fmt.Errorf("expanding properties: %+v", err)
        }
    }

    // extended location
    if v, ok := d.GetOk("extended_location"); ok && len(v.([]interface{})) > 0 {
        el := v.([]interface{})[0].(map[string]interface{})
        profile.ExtendedLocation = &armiotoperations.ExtendedLocation{}
        if v2, ok2 := el["name"].(string); ok2 && v2 != "" {
            profile.ExtendedLocation.Name = to.Ptr(v2)
        }
        if v2, ok2 := el["type"].(string); ok2 && v2 != "" {
            profile.ExtendedLocation.Type = to.Ptr(armiotoperations.ExtendedLocationType(v2))
        }
    }

    poller, err := svc.DataflowProfileClient.BeginCreateOrUpdate(ctx, rg, instance, name, profile, nil)
    if err != nil {
        return fmt.Errorf("starting DataflowProfile create/update: %+v", err)
    }
    res, err := poller.PollUntilDone(ctx, nil)
    if err != nil {
        return fmt.Errorf("waiting for DataflowProfile create/update: %+v", err)
    }

    // set ID
    if res.DataflowProfileResource.ID != nil {
        d.SetId(*res.DataflowProfileResource.ID)
    } else {
        d.SetId(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s/dataflowProfiles/%s", top.SubscriptionID(), rg, instance, name))
    }

    if err := d.Set("properties", []interface{}{flattenDataflowProfileProperties(res.DataflowProfileResource.Properties)}); err != nil {
        return fmt.Errorf("setting properties: %+v", err)
    }

    return nil
}

func resourceIotOperationsDataflowProfileRead(d *schema.ResourceData, meta interface{}) error {
    ctx := context.Background()

    name := d.Get("name").(string)
    rg := d.Get("resource_group_name").(string)
    instance := d.Get("instance_name").(string)

    top := meta.(*clients.Client)
    if top.IoTOperations == nil || top.IoTOperations.DataflowProfileClient == nil {
        return fmt.Errorf("iotoperations client is not configured")
    }
    svc := top.IoTOperations

    res, err := svc.DataflowProfileClient.Get(ctx, rg, instance, name, nil)
    if err != nil {
        return fmt.Errorf("failed to get DataflowProfile: %+v", err)
    }

    if err := d.Set("properties", []interface{}{flattenDataflowProfileProperties(res.DataflowProfileResource.Properties)}); err != nil {
        return fmt.Errorf("setting properties: %+v", err)
    }

    return nil
}

func resourceIotOperationsDataflowProfileUpdate(d *schema.ResourceData, meta interface{}) error {
    ctx := context.Background()

    name := d.Get("name").(string)
    rg := d.Get("resource_group_name").(string)
    instance := d.Get("instance_name").(string)

    top := meta.(*clients.Client)
    if top.IoTOperations == nil || top.IoTOperations.DataflowProfileClient == nil {
        return fmt.Errorf("iotoperations client is not configured")
    }
    svc := top.IoTOperations

    profile := armiotoperations.DataflowProfileResource{}
    if v, ok := d.GetOk("properties"); ok {
        if props, err := expandDataflowProfileProperties(v.([]interface{})); err == nil && props != nil {
            profile.Properties = props
        } else if err != nil {
            return fmt.Errorf("expanding properties: %+v", err)
        }
    }

    // extended location
    if v, ok := d.GetOk("extended_location"); ok && len(v.([]interface{})) > 0 {
        el := v.([]interface{})[0].(map[string]interface{})
        profile.ExtendedLocation = &armiotoperations.ExtendedLocation{}
        if v2, ok2 := el["name"].(string); ok2 && v2 != "" {
            profile.ExtendedLocation.Name = to.Ptr(v2)
        }
        if v2, ok2 := el["type"].(string); ok2 && v2 != "" {
            profile.ExtendedLocation.Type = to.Ptr(armiotoperations.ExtendedLocationType(v2))
        }
    }

    poller, err := svc.DataflowProfileClient.BeginCreateOrUpdate(ctx, rg, instance, name, profile, nil)
    if err != nil {
        return fmt.Errorf("starting DataflowProfile update: %+v", err)
    }
    res, err := poller.PollUntilDone(ctx, nil)
    if err != nil {
        return fmt.Errorf("waiting for DataflowProfile update: %+v", err)
    }

    if err := d.Set("properties", []interface{}{flattenDataflowProfileProperties(res.DataflowProfileResource.Properties)}); err != nil {
        return fmt.Errorf("setting properties: %+v", err)
    }

    return nil
}

func resourceIotOperationsDataflowProfileDelete(d *schema.ResourceData, meta interface{}) error {
    ctx := context.Background()

    name := d.Get("name").(string)
    rg := d.Get("resource_group_name").(string)
    instance := d.Get("instance_name").(string)

    top := meta.(*clients.Client)
    if top.IoTOperations == nil || top.IoTOperations.DataflowProfileClient == nil {
        return fmt.Errorf("iotoperations client is not configured")
    }
    svc := top.IoTOperations

    poller, err := svc.DataflowProfileClient.BeginDelete(ctx, rg, instance, name, nil)
    if err != nil {
        return fmt.Errorf("failed to start delete DataflowProfile: %+v", err)
    }
    _, err = poller.PollUntilDone(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to poll delete DataflowProfile: %+v", err)
    }

    return nil
}

// flatten helpers
func flattenDataflowProfileProperties(props *armiotoperations.DataflowProfileProperties) map[string]interface{} {
    if props == nil {
        return nil
    }

    m := make(map[string]interface{})

    if props.InstanceCount != nil {
        m["instance_count"] = int(*props.InstanceCount)
    }

    if props.Diagnostics != nil {
        diag := make(map[string]interface{})
        if props.Diagnostics.Logs != nil && props.Diagnostics.Logs.Level != nil {
            diag["logs"] = []interface{}{map[string]interface{}{"level": *props.Diagnostics.Logs.Level}}
        }
        if props.Diagnostics.Metrics != nil {
            metrics := make(map[string]interface{})
            if props.Diagnostics.Metrics.PrometheusPort != nil {
                metrics["prometheus_port"] = int(*props.Diagnostics.Metrics.PrometheusPort)
            }
            diag["metrics"] = []interface{}{metrics}
        }
        m["diagnostics"] = []interface{}{diag}
    }

    if props.ProvisioningState != nil {
        m["provisioning_state"] = string(*props.ProvisioningState)
    }

    return m
}

// expand helpers
func expandDataflowProfileProperties(v []interface{}) (*armiotoperations.DataflowProfileProperties, error) {
    if len(v) == 0 {
        return nil, nil
    }
    m := v[0].(map[string]interface{})
    props := &armiotoperations.DataflowProfileProperties{}

    if v, ok := m["instance_count"]; ok {
        props.InstanceCount = to.Ptr[int32](int32(v.(int)))
    }

    if v, ok := m["diagnostics"]; ok && len(v.([]interface{})) > 0 {
        el := v.([]interface{})[0].(map[string]interface{})
        props.Diagnostics = &armiotoperations.ProfileDiagnostics{}
        if v2, ok2 := el["logs"]; ok2 && len(v2.([]interface{})) > 0 {
            logMap := v2.([]interface{})[0].(map[string]interface{})
            props.Diagnostics.Logs = &armiotoperations.DiagnosticsLogs{}
            if lvl, ok3 := logMap["level"].(string); ok3 && lvl != "" {
                props.Diagnostics.Logs.Level = to.Ptr(lvl)
            }
        }
        if v2, ok2 := el["metrics"]; ok2 && len(v2.([]interface{})) > 0 {
            metMap := v2.([]interface{})[0].(map[string]interface{})
            props.Diagnostics.Metrics = &armiotoperations.Metrics{}
            if p, ok3 := metMap["prometheus_port"].(int); ok3 {
                props.Diagnostics.Metrics.PrometheusPort = to.Ptr[int32](int32(p))
            }

        }
    }
    return props, nil
}
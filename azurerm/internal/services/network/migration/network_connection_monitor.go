package migration

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
)

func NetworkConnectionMonitorV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"network_watcher_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"auto_start": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "This field has been deprecated in new api version 2020-05-01",
			},

			"destination": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"endpoint"},
				Deprecated:    "Deprecated in favor of `endpoint`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"virtual_machine_id": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"destination.0.address"},
						},

						"address": {
							Type:          schema.TypeString,
							Optional:      true,
							ConflictsWith: []string{"destination.0.virtual_machine_id"},
						},
					},
				},
			},

			"endpoint": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"source", "destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"address": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"filter": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"item": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Optional: true,
													Default:  string(network.AgentAddress),
												},

												"address": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},

									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(network.Include),
									},
								},
							},
						},

						"virtual_machine_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"interval_in_seconds": {
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Deprecated in favor of `test_frequency_sec`",
			},

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"endpoint"},
				Deprecated:    "Deprecated in favor of `endpoint`",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virtual_machine_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},

			"test_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},

						"test_frequency_sec": {
							Type:          schema.TypeInt,
							Optional:      true,
							Default:       60,
							ConflictsWith: []string{"interval_in_seconds"},
						},

						"http_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(network.Get),
									},

									"port": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									"path": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"prefer_https": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"request_header": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},

												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},

									"valid_status_code_ranges": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"icmp_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disable_trace_route": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},

						"preferred_ip_version": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"success_threshold": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"checks_failed_percent": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									"round_trip_time_ms": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},

						"tcp_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"disable_trace_route": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"test_group": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"destinations": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"sources": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"test_configurations": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"disable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"output_workspace_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				// Allow to switch workspace from specified one to default one.
				// 1. Set `output = []` in tfconfig to switch workspace from specified one to default one.
				// 2. Remove `output = []` from tfconfig to ensure no diff.
				ConfigMode: schema.SchemaConfigModeAttr,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func NetworkConnectionMonitorV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "/NetworkConnectionMonitors/", "/connectionMonitors/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}

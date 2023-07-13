// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationserviceenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/workflows"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var logicAppResourceName = "azurerm_logic_app"

func resourceLogicAppWorkflow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppWorkflowCreate,
		Read:   resourceLogicAppWorkflowRead,
		Update: resourceLogicAppWorkflowUpdate,
		Delete: resourceLogicAppWorkflowDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := workflows.ParseWorkflowID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-()_.A-Za-z0-9]{1,80}$"),
						"The Logic app name can contain only letters, numbers, periods (.), hyphens (-), brackets (()) and underscores (_), up to 80 characters",
					),
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"integration_service_environment_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: integrationserviceenvironments.ValidateIntegrationServiceEnvironmentID,
			},

			"access_control": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_caller_ip_address_range": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsCIDR,
												validation.IsIPv4Range,
											),
										},
									},
								},
							},
						},

						"content": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_caller_ip_address_range": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsCIDR,
												validation.IsIPv4Range,
											),
										},
									},
								},
							},
						},

						"trigger": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_caller_ip_address_range": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsCIDR,
												validation.IsIPv4Range,
											),
										},
									},

									"open_authentication_policy": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"name": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"claim": {
													Type:     pluginsdk.TypeSet,
													Required: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"name": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"value": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},

						"workflow_management": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_caller_ip_address_range": {
										Type:     pluginsdk.TypeSet,
										Required: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.Any(
												validation.IsCIDR,
												validation.IsIPv4Range,
											),
										},
									},
								},
							},
						},
					},
				},
			},

			"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

			"logic_app_integration_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: integrationaccounts.ValidateIntegrationAccountID,
			},

			// TODO: should Parameters be split out into their own object to allow validation on the different sub-types?
			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"workflow_schema": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "https://schema.management.azure.com/providers/Microsoft.Logic/schemas/2016-06-01/workflowdefinition.json#",
			},

			"workflow_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "1.0.0.0",
			},

			"workflow_parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"access_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"connector_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"connector_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"workflow_endpoint_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},
			"workflow_outbound_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceLogicAppWorkflowCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Logic App Workflow creation.")

	id := workflows.NewWorkflowID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Logic App Workflow %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_logic_app_workflow", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))

	workflowSchema := d.Get("workflow_schema").(string)
	workflowVersion := d.Get("workflow_version").(string)
	workflowParameters, err := expandLogicAppWorkflowWorkflowParameters(d.Get("workflow_parameters").(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `workflow_parameters`: %+v", err)
	}

	parameters, err := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]interface{}), workflowParameters)
	if err != nil {
		return err
	}
	t := d.Get("tags").(map[string]interface{})

	isEnabled := workflows.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = workflows.WorkflowStateDisabled
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// nolint gosimple
	var definition interface{}
	definition = map[string]interface{}{
		"$schema":        workflowSchema,
		"contentVersion": workflowVersion,
		"actions":        make(map[string]interface{}),
		"triggers":       make(map[string]interface{}),
		"parameters":     workflowParameters,
	}

	properties := workflows.Workflow{
		Identity: identity,
		Location: utils.String(location),
		Properties: &workflows.WorkflowProperties{
			Definition: &definition,
			Parameters: parameters,
			State:      &isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.Properties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]interface{}))
	}

	if iseID, ok := d.GetOk("integration_service_environment_id"); ok {
		properties.Properties.IntegrationServiceEnvironment = &workflows.ResourceReference{
			Id: utils.String(iseID.(string)),
		}
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.Properties.IntegrationAccount = &workflows.ResourceReference{
			Id: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
		return fmt.Errorf("[ERROR] Error creating Logic App Workflow %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogicAppWorkflowRead(d, meta)
}

func resourceLogicAppWorkflowUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	if read.Model == nil || read.Model.Properties == nil {
		return fmt.Errorf("[ERROR] Error parsing Logic App Workflow - `WorkflowProperties` is nil")
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	workflowParameters, err := expandLogicAppWorkflowWorkflowParameters(d.Get("workflow_parameters").(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `workflow_parameters`: %+v", err)
	}
	parameters, err := expandLogicAppWorkflowParameters(d.Get("parameters").(map[string]interface{}), workflowParameters)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	var definition interface{}
	if read.Model.Properties.Definition != nil {
		definitionRaw := *read.Model.Properties.Definition
		definitionMap := definitionRaw.(map[string]interface{})
		definitionMap["parameters"] = workflowParameters
		definition = definitionMap
	}

	isEnabled := workflows.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = workflows.WorkflowStateDisabled
	}

	identity, err := identity.ExpandSystemOrUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	properties := workflows.Workflow{
		Identity: identity,
		Location: utils.String(location),
		Properties: &workflows.WorkflowProperties{
			Definition: &definition,
			Parameters: parameters,
			State:      &isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.Properties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]interface{}))
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.Properties.IntegrationAccount = &workflows.ResourceReference{
			Id: utils.String(v.(string)),
		}
	}

	if iseID, ok := d.GetOk("integration_service_environment_id"); ok {
		properties.Properties.IntegrationServiceEnvironment = &workflows.ResourceReference{
			Id: utils.String(iseID.(string)),
		}
	}

	if _, err = client.CreateOrUpdate(ctx, *id, properties); err != nil {
		return fmt.Errorf("updating Logic App Workflow %s: %+v", id, err)
	}

	return resourceLogicAppWorkflowRead(d, meta)
}

func resourceLogicAppWorkflowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Logic App Workflow %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	d.Set("name", id.WorkflowName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if location := model.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return err
		}
		d.Set("identity", identity)

		if props := model.Properties; props != nil {
			d.Set("access_endpoint", props.AccessEndpoint)

			if err := d.Set("access_control", flattenLogicAppWorkflowFlowAccessControl(props.AccessControl)); err != nil {
				return fmt.Errorf("setting `access_control`: %+v", err)
			}

			if props.State != nil && *props.State != "" {
				d.Set("enabled", *props.State == workflows.WorkflowStateEnabled)
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Connector == nil {
				d.Set("connector_endpoint_ip_addresses", []interface{}{})
				d.Set("connector_outbound_ip_addresses", []interface{}{})
			} else {
				d.Set("connector_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.AccessEndpointIPAddresses))
				d.Set("connector_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Connector.OutgoingIPAddresses))
			}

			if props.EndpointsConfiguration == nil || props.EndpointsConfiguration.Workflow == nil {
				d.Set("workflow_endpoint_ip_addresses", []interface{}{})
				d.Set("workflow_outbound_ip_addresses", []interface{}{})
			} else {
				d.Set("workflow_endpoint_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.AccessEndpointIPAddresses))
				d.Set("workflow_outbound_ip_addresses", flattenIPAddresses(props.EndpointsConfiguration.Workflow.OutgoingIPAddresses))
			}
			if definition := props.Definition; definition != nil {
				definitionRaw := *props.Definition
				if v, ok := definitionRaw.(map[string]interface{}); ok {
					if v["$schema"] != nil {
						d.Set("workflow_schema", v["$schema"].(string))
					}
					if v["contentVersion"] != nil {
						d.Set("workflow_version", v["contentVersion"].(string))
					}
					if p, ok := v["parameters"]; ok {
						workflowParameters, err := flattenLogicAppWorkflowWorkflowParameters(p.(map[string]interface{}))
						if err != nil {
							return fmt.Errorf("flattening `workflow_parameters`: %+v", err)
						}
						if err := d.Set("workflow_parameters", workflowParameters); err != nil {
							return fmt.Errorf("setting `workflow_parameters`: %+v", err)
						}

						// The props.Parameters (the value of the param) is accompany with the "parameters" (the definition of the param) inside the props.Definition.
						// We will need to make use of the definition of the parameters in order to properly flatten the value of the parameters being set (for kinds of types).
						parameters, err := flattenLogicAppWorkflowParameters(d, props.Parameters, p.(map[string]interface{}))
						if err != nil {
							return fmt.Errorf("flattening `parameters`: %v", err)
						}
						if err := d.Set("parameters", parameters); err != nil {
							return fmt.Errorf("setting `parameters`: %+v", err)
						}
					}
				}
			}

			integrationServiceEnvironmentId := ""
			if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.Id != nil {
				integrationServiceEnvironmentId = *props.IntegrationServiceEnvironment.Id
			}
			d.Set("integration_service_environment_id", integrationServiceEnvironmentId)

			if props.IntegrationAccount != nil && props.IntegrationAccount.Id != nil {
				d.Set("logic_app_integration_account_id", props.IntegrationAccount.Id)
			}

			integrationAccountId := ""
			if props.IntegrationAccount != nil && props.IntegrationAccount.Id != nil {
				integrationAccountId = *props.IntegrationAccount.Id
			}
			d.Set("logic_app_integration_account_id", integrationAccountId)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceLogicAppWorkflowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := workflows.ParseWorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.WorkflowName, logicAppResourceName)
	defer locks.UnlockByName(id.WorkflowName, logicAppResourceName)

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("issuing delete request for Logic App Workflow %s: %+v", id, err)
	}

	return nil
}

func expandLogicAppWorkflowParameters(input map[string]interface{}, paramDefs map[string]interface{}) (*map[string]workflows.WorkflowParameter, error) {
	output := make(map[string]workflows.WorkflowParameter)

	for k, v := range input {
		defRaw, ok := paramDefs[k]
		if !ok {
			return nil, fmt.Errorf("no parameter definition for %s", k)
		}
		def := defRaw.(map[string]interface{})
		t := workflows.ParameterType(def["type"].(string))

		v := v.(string)

		var value interface{}
		switch t {
		case workflows.ParameterTypeBool:
			var uv bool
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to bool: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeFloat:
			var uv float64
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to float64: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeInt:
			var uv int
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to int: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeArray:
			var uv []interface{}
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to []interface{}: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeObject,
			workflows.ParameterTypeSecureObject:
			var uv map[string]interface{}
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to map[string]interface{}: %v", k, err)
			}
			value = uv
		case workflows.ParameterTypeString,
			workflows.ParameterTypeSecureString:
			value = v
		}

		output[k] = workflows.WorkflowParameter{
			Type:  &t,
			Value: &value,
		}
	}

	return &output, nil
}

func flattenLogicAppWorkflowParameters(d *pluginsdk.ResourceData, input *map[string]workflows.WorkflowParameter, paramDefs map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	if input == nil {
		return output, nil
	}

	// Read the "parameters" from state, which is used to fill in the "sensitive" properties.
	paramInState := make(map[string]interface{})
	paramsRaw := d.Get("parameters")
	if params, ok := paramsRaw.(map[string]interface{}); ok {
		paramInState = params
	}

	for k, v := range *input {
		defRaw, ok := paramDefs[k]
		if !ok {
			// This should never happen.
			log.Printf("[WARN] The parameter %s is not defined in the Logic App Workflow", k)
			continue
		}

		def := defRaw.(map[string]interface{})
		t := workflows.ParameterType(def["type"].(string))

		var value string
		switch t {
		case workflows.ParameterTypeBool:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(bool)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got %T", k, v.Value)
			}
			value = "true"
			if !tv {
				value = "false"
			}
		case workflows.ParameterTypeFloat:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := valueRaw.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.FormatFloat(tv, 'f', -1, 64)
		case workflows.ParameterTypeInt:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := valueRaw.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.Itoa(int(tv))

		case workflows.ParameterTypeArray:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.([]interface{})
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be []interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case workflows.ParameterTypeObject:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be map[string]interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case workflows.ParameterTypeString:
			if v.Value == nil {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got nil", k)
			}
			valueRaw := *v.Value
			tv, ok := valueRaw.(string)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be string, but got %T", k, v.Value)
			}
			value = tv

		case workflows.ParameterTypeSecureString,
			workflows.ParameterTypeSecureObject:
			// This is not returned from API, we will try to read them from the state instead.
			if v, ok := paramInState[k]; ok {
				value = v.(string) // The value in state here is guaranteed to be a string, so directly cast the type.
			}
		}

		output[k] = value
	}

	return output, nil
}

func expandLogicAppWorkflowWorkflowParameters(input map[string]interface{}) (map[string]interface{}, error) {
	if len(input) == 0 {
		return nil, nil
	}

	output := make(map[string]interface{})
	for k, v := range input {
		obj, err := pluginsdk.ExpandJsonFromString(v.(string))
		if err != nil {
			return nil, err
		}
		output[k] = obj
	}
	return output, nil
}

func expandLogicAppWorkflowAccessControl(input []interface{}) *workflows.FlowAccessControlConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := workflows.FlowAccessControlConfiguration{}

	if contents := v["content"].([]interface{}); len(contents) != 0 {
		result.Contents = expandLogicAppWorkflowAccessControlConfigurationPolicy(contents)
	}

	if actions := v["action"].([]interface{}); len(actions) != 0 {
		result.Actions = expandLogicAppWorkflowAccessControlConfigurationPolicy(actions)
	}

	if triggers := v["trigger"].([]interface{}); len(triggers) != 0 {
		result.Triggers = expandLogicAppWorkflowAccessControlTriggerConfigurationPolicy(triggers)
	}

	if workflowManagement := v["workflow_management"].([]interface{}); len(workflowManagement) != 0 {
		result.WorkflowManagement = expandLogicAppWorkflowAccessControlConfigurationPolicy(workflowManagement)
	}

	return &result
}

func expandLogicAppWorkflowAccessControlConfigurationPolicy(input []interface{}) *workflows.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &workflows.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}
}

func expandLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input []interface{}) *workflows.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := workflows.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}

	if openAuthenticationPolicy, ok := v["open_authentication_policy"]; ok {
		openAuthenticationPolicies := openAuthenticationPolicy.(*pluginsdk.Set).List()
		if len(openAuthenticationPolicies) != 0 {
			result.OpenAuthenticationPolicies = &workflows.OpenAuthenticationAccessPolicies{
				Policies: expandLogicAppWorkflowOpenAuthenticationPolicy(openAuthenticationPolicies),
			}
		}
	}

	return &result
}

func expandLogicAppWorkflowIPAddressRanges(input []interface{}) *[]workflows.IPAddressRange {
	results := make([]workflows.IPAddressRange, 0)

	for _, item := range input {
		results = append(results, workflows.IPAddressRange{
			AddressRange: utils.String(item.(string)),
		})
	}

	return &results
}

func expandLogicAppWorkflowOpenAuthenticationPolicy(input []interface{}) *map[string]workflows.OpenAuthenticationAccessPolicy {
	if len(input) == 0 {
		return nil
	}
	results := make(map[string]workflows.OpenAuthenticationAccessPolicy)

	for _, item := range input {
		if item == nil {
			continue
		}
		v := item.(map[string]interface{})
		policyName := v["name"].(string)

		policyType := workflows.OpenAuthenticationProviderTypeAAD
		results[policyName] = workflows.OpenAuthenticationAccessPolicy{
			Type:   &policyType,
			Claims: expandLogicAppWorkflowOpenAuthenticationPolicyClaim(v["claim"].(*pluginsdk.Set).List()),
		}
	}

	return &results
}

func expandLogicAppWorkflowOpenAuthenticationPolicyClaim(input []interface{}) *[]workflows.OpenAuthenticationPolicyClaim {
	results := make([]workflows.OpenAuthenticationPolicyClaim, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, workflows.OpenAuthenticationPolicyClaim{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		})
	}
	return &results
}

func flattenLogicAppWorkflowWorkflowParameters(input map[string]interface{}) (map[string]interface{}, error) {
	if input == nil {
		return nil, nil
	}
	output := make(map[string]interface{})
	for k, v := range input {
		objstr, err := pluginsdk.FlattenJsonToString(v.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		output[k] = objstr
	}
	return output, nil
}

func flattenIPAddresses(input *[]workflows.IPAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var addresses []interface{}
	for _, addr := range *input {
		addresses = append(addresses, *addr.Address)
	}
	return addresses
}

func flattenLogicAppWorkflowFlowAccessControl(input *workflows.FlowAccessControlConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"action":              flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.Actions),
			"content":             flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.Contents),
			"trigger":             flattenLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input.Triggers),
			"workflow_management": flattenLogicAppWorkflowAccessControlConfigurationPolicy(input.WorkflowManagement),
		},
	}
}

func flattenLogicAppWorkflowAccessControlConfigurationPolicy(input *workflows.FlowAccessControlConfigurationPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"allowed_caller_ip_address_range": flattenLogicAppWorkflowIPAddressRanges(input.AllowedCallerIPAddresses),
		},
	}
}

func flattenLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input *workflows.FlowAccessControlConfigurationPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"allowed_caller_ip_address_range": flattenLogicAppWorkflowIPAddressRanges(input.AllowedCallerIPAddresses),
			"open_authentication_policy":      flattenLogicAppWorkflowOpenAuthenticationPolicy(input.OpenAuthenticationPolicies),
		},
	}
}

func flattenLogicAppWorkflowIPAddressRanges(input *[]workflows.IPAddressRange) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var addressRange string
		if item.AddressRange != nil {
			addressRange = *item.AddressRange
		}
		results = append(results, addressRange)
	}

	return results
}

func flattenLogicAppWorkflowOpenAuthenticationPolicy(input *workflows.OpenAuthenticationAccessPolicies) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Policies == nil {
		return results
	}

	for k, v := range *input.Policies {
		results = append(results, map[string]interface{}{
			"name":  k,
			"claim": flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(v.Claims),
		})
	}

	return results
}

func flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(input *[]workflows.OpenAuthenticationPolicyClaim) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var value string
		if item.Value != nil {
			value = *item.Value
		}

		results = append(results, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}

	return results
}

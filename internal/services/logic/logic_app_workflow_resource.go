package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.WorkflowID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"integration_service_environment_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationServiceEnvironmentID,
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
				ValidateFunc: validate.IntegrationAccountID,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceLogicAppWorkflowCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	subscriptionId := meta.(*clients.Client).Logic.WorkflowClient.SubscriptionID
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Logic App Workflow creation.")

	id := parse.NewWorkflowID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Logic App Workflow %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
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

	isEnabled := logic.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = logic.WorkflowStateDisabled
	}

	identity, err := expandLogicAppWorkflowIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	properties := logic.Workflow{
		Identity: identity,
		Location: utils.String(location),
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: &map[string]interface{}{
				"$schema":        workflowSchema,
				"contentVersion": workflowVersion,
				"actions":        make(map[string]interface{}),
				"triggers":       make(map[string]interface{}),
				"parameters":     workflowParameters,
			},
			Parameters: parameters,
			State:      isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.WorkflowProperties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]interface{}))
	}

	if iseID, ok := d.GetOk("integration_service_environment_id"); ok {
		properties.WorkflowProperties.IntegrationServiceEnvironment = &logic.ResourceReference{
			ID: utils.String(iseID.(string)),
		}
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.WorkflowProperties.IntegrationAccount = &logic.ResourceReference{
			ID: utils.String(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, properties); err != nil {
		return fmt.Errorf("[ERROR] Error creating Logic App Workflow %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLogicAppWorkflowRead(d, meta)
}

func resourceLogicAppWorkflowUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.Name, logicAppResourceName)
	defer locks.UnlockByName(id.Name, logicAppResourceName)

	read, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	if read.WorkflowProperties == nil {
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

	definition := read.WorkflowProperties.Definition.(map[string]interface{})
	definition["parameters"] = workflowParameters

	isEnabled := logic.WorkflowStateEnabled
	if v := d.Get("enabled").(bool); !v {
		isEnabled = logic.WorkflowStateDisabled
	}

	identity, err := expandLogicAppWorkflowIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	properties := logic.Workflow{
		Identity: identity,
		Location: utils.String(location),
		WorkflowProperties: &logic.WorkflowProperties{
			Definition: definition,
			Parameters: parameters,
			State:      isEnabled,
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("access_control"); ok {
		properties.WorkflowProperties.AccessControl = expandLogicAppWorkflowAccessControl(v.([]interface{}))
	}

	if v, ok := d.GetOk("logic_app_integration_account_id"); ok {
		properties.WorkflowProperties.IntegrationAccount = &logic.ResourceReference{
			ID: utils.String(v.(string)),
		}
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, properties); err != nil {
		return fmt.Errorf("updating Logic App Workflow %s: %+v", id, err)
	}

	return resourceLogicAppWorkflowRead(d, meta)
}

func resourceLogicAppWorkflowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkflowID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Logic App Workflow %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error making Read request on Logic App Workflow %s: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenLogicAppWorkflowIdentity(resp.Identity)
	if err != nil {
		return err
	}
	d.Set("identity", identity)

	if props := resp.WorkflowProperties; props != nil {
		d.Set("access_endpoint", props.AccessEndpoint)

		if err := d.Set("access_control", flattenLogicAppWorkflowFlowAccessControl(props.AccessControl)); err != nil {
			return fmt.Errorf("setting `access_control`: %+v", err)
		}

		if props.State != "" {
			d.Set("enabled", props.State == logic.WorkflowStateEnabled)
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
			if v, ok := definition.(map[string]interface{}); ok {
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
		if props.IntegrationServiceEnvironment != nil && props.IntegrationServiceEnvironment.ID != nil {
			integrationServiceEnvironmentId = *props.IntegrationServiceEnvironment.ID
		}
		d.Set("integration_service_environment_id", integrationServiceEnvironmentId)

		if props.IntegrationAccount != nil && props.IntegrationAccount.ID != nil {
			d.Set("logic_app_integration_account_id", props.IntegrationAccount.ID)
		}

		integrationAccountId := ""
		if props.IntegrationAccount != nil && props.IntegrationAccount.ID != nil {
			integrationAccountId = *props.IntegrationAccount.ID
		}
		d.Set("logic_app_integration_account_id", integrationAccountId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogicAppWorkflowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.WorkflowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkflowID(d.Id())
	if err != nil {
		return err
	}

	// lock to prevent against Actions, Parameters or Triggers conflicting
	locks.ByName(id.Name, logicAppResourceName)
	defer locks.UnlockByName(id.Name, logicAppResourceName)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("issuing delete request for Logic App Workflow %s: %+v", id, err)
	}

	return nil
}

func expandLogicAppWorkflowParameters(input map[string]interface{}, paramDefs map[string]interface{}) (map[string]*logic.WorkflowParameter, error) {
	output := make(map[string]*logic.WorkflowParameter)

	for k, v := range input {
		defRaw, ok := paramDefs[k]
		if !ok {
			return nil, fmt.Errorf("no parameter definition for %s", k)
		}
		def := defRaw.(map[string]interface{})
		t := logic.ParameterType(def["type"].(string))

		v := v.(string)

		var value interface{}
		switch t {
		case logic.ParameterTypeBool:
			var uv bool
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to bool: %v", k, err)
			}
			value = uv
		case logic.ParameterTypeFloat:
			var uv float64
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to float64: %v", k, err)
			}
			value = uv
		case logic.ParameterTypeInt:
			var uv int
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to int: %v", k, err)
			}
			value = uv
		case logic.ParameterTypeArray:
			var uv []interface{}
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to []interface{}: %v", k, err)
			}
			value = uv
		case logic.ParameterTypeObject,
			logic.ParameterTypeSecureObject:
			var uv map[string]interface{}
			if err := json.Unmarshal([]byte(v), &uv); err != nil {
				return nil, fmt.Errorf("unmarshalling %s to map[string]interface{}: %v", k, err)
			}
			value = uv
		case logic.ParameterTypeString,
			logic.ParameterTypeSecureString:
			value = v
		}

		output[k] = &logic.WorkflowParameter{
			Type:  t,
			Value: value,
		}
	}

	return output, nil
}

func flattenLogicAppWorkflowParameters(d *pluginsdk.ResourceData, input map[string]*logic.WorkflowParameter, paramDefs map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{})

	// Read the "parameters" from state, which is used to fill in the "sensitive" properties.
	paramInState := make(map[string]interface{})
	paramsRaw := d.Get("parameters")
	if params, ok := paramsRaw.(map[string]interface{}); ok {
		paramInState = params
	}

	for k, v := range input {
		defRaw, ok := paramDefs[k]
		if !ok {
			// This should never happen.
			log.Printf("[WARN] The parameter %s is not defined in the Logic App Workflow", k)
			continue
		}

		if v == nil {
			log.Printf("[WARN] The value of parameter %s is nil", k)
			continue
		}

		def := defRaw.(map[string]interface{})
		t := logic.ParameterType(def["type"].(string))

		var value string
		switch t {
		case logic.ParameterTypeBool:
			tv, ok := v.Value.(bool)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be bool, but got %T", k, v.Value)
			}
			value = "true"
			if !tv {
				value = "false"
			}
		case logic.ParameterTypeFloat:
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := v.Value.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.FormatFloat(tv, 'f', -1, 64)
		case logic.ParameterTypeInt:
			// Note that the json unmarshalled response doesn't differ between float and int, as json has only type number.
			tv, ok := v.Value.(float64)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be float64, but got %T", k, v.Value)
			}
			value = strconv.Itoa(int(tv))

		case logic.ParameterTypeArray:
			tv, ok := v.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be []interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case logic.ParameterTypeObject:
			tv, ok := v.Value.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be map[string]interface{}, but got %T", k, v.Value)
			}
			obj, err := json.Marshal(tv)
			if err != nil {
				return nil, fmt.Errorf("converting %+v from json: %v", tv, err)
			}
			value = string(obj)

		case logic.ParameterTypeString:
			tv, ok := v.Value.(string)
			if !ok {
				return nil, fmt.Errorf("the value of parameter %s is expected to be string, but got %T", k, v.Value)
			}
			value = tv

		case logic.ParameterTypeSecureString,
			logic.ParameterTypeSecureObject:
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

func expandLogicAppWorkflowAccessControl(input []interface{}) *logic.FlowAccessControlConfiguration {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := logic.FlowAccessControlConfiguration{}

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

func expandLogicAppWorkflowAccessControlConfigurationPolicy(input []interface{}) *logic.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &logic.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}
}

func expandLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input []interface{}) *logic.FlowAccessControlConfigurationPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := logic.FlowAccessControlConfigurationPolicy{
		AllowedCallerIPAddresses: expandLogicAppWorkflowIPAddressRanges(v["allowed_caller_ip_address_range"].(*pluginsdk.Set).List()),
	}

	if openAuthenticationPolicy, ok := v["open_authentication_policy"]; ok {
		openAuthenticationPolicies := openAuthenticationPolicy.(*pluginsdk.Set).List()
		if len(openAuthenticationPolicies) != 0 {
			result.OpenAuthenticationPolicies = &logic.OpenAuthenticationAccessPolicies{
				Policies: expandLogicAppWorkflowOpenAuthenticationPolicy(openAuthenticationPolicies),
			}
		}
	}

	return &result
}

func expandLogicAppWorkflowIPAddressRanges(input []interface{}) *[]logic.IPAddressRange {
	results := make([]logic.IPAddressRange, 0)

	for _, item := range input {
		results = append(results, logic.IPAddressRange{
			AddressRange: utils.String(item.(string)),
		})
	}

	return &results
}

func expandLogicAppWorkflowOpenAuthenticationPolicy(input []interface{}) map[string]*logic.OpenAuthenticationAccessPolicy {
	if len(input) == 0 {
		return nil
	}
	results := make(map[string]*logic.OpenAuthenticationAccessPolicy)

	for _, item := range input {
		v := item.(map[string]interface{})
		policyName := v["name"].(string)

		results[policyName] = &logic.OpenAuthenticationAccessPolicy{
			Type:   logic.OpenAuthenticationProviderTypeAAD,
			Claims: expandLogicAppWorkflowOpenAuthenticationPolicyClaim(v["claim"].(*pluginsdk.Set).List()),
		}
	}

	return results
}

func expandLogicAppWorkflowOpenAuthenticationPolicyClaim(input []interface{}) *[]logic.OpenAuthenticationPolicyClaim {
	results := make([]logic.OpenAuthenticationPolicyClaim, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, logic.OpenAuthenticationPolicyClaim{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		})
	}
	return &results
}

func expandLogicAppWorkflowIdentity(input []interface{}) (*logic.ManagedServiceIdentity, error) {
	config, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	var identityIds map[string]*logic.UserAssignedIdentity
	if len(config.IdentityIds) != 0 {
		identityIds = map[string]*logic.UserAssignedIdentity{}
		for id := range config.IdentityIds {
			identityIds[id] = &logic.UserAssignedIdentity{}
		}
	}

	return &logic.ManagedServiceIdentity{
		Type:                   logic.ManagedServiceIdentityType(config.Type),
		UserAssignedIdentities: identityIds,
	}, nil
}

func flattenLogicAppWorkflowIdentity(input *logic.ManagedServiceIdentity) (*[]interface{}, error) {
	var config *identity.SystemOrUserAssignedMap
	if input != nil {
		identityIds := map[string]identity.UserAssignedIdentityDetails{}
		for id := range input.UserAssignedIdentities {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(id)
			if err != nil {
				return nil, err
			}
			identityIds[parsedId.ID()] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}

		principalId := ""
		if input.PrincipalID != nil {
			principalId = input.PrincipalID.String()
		}

		tenantId := ""
		if input.TenantID != nil {
			tenantId = input.TenantID.String()
		}

		config = &identity.SystemOrUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			PrincipalId: principalId,
			TenantId:    tenantId,
			IdentityIds: identityIds,
		}
	}
	return identity.FlattenSystemOrUserAssignedMap(config)
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

func flattenIPAddresses(input *[]logic.IPAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var addresses []interface{}
	for _, addr := range *input {
		addresses = append(addresses, *addr.Address)
	}
	return addresses
}

func flattenLogicAppWorkflowFlowAccessControl(input *logic.FlowAccessControlConfiguration) []interface{} {
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

func flattenLogicAppWorkflowAccessControlConfigurationPolicy(input *logic.FlowAccessControlConfigurationPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"allowed_caller_ip_address_range": flattenLogicAppWorkflowIPAddressRanges(input.AllowedCallerIPAddresses),
		},
	}
}

func flattenLogicAppWorkflowAccessControlTriggerConfigurationPolicy(input *logic.FlowAccessControlConfigurationPolicy) []interface{} {
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

func flattenLogicAppWorkflowIPAddressRanges(input *[]logic.IPAddressRange) []interface{} {
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

func flattenLogicAppWorkflowOpenAuthenticationPolicy(input *logic.OpenAuthenticationAccessPolicies) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Policies == nil {
		return results
	}

	for k, v := range input.Policies {
		results = append(results, map[string]interface{}{
			"name":  k,
			"claim": flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(v.Claims),
		})
	}

	return results
}

func flattenLogicAppWorkflowOpenAuthenticationPolicyClaim(input *[]logic.OpenAuthenticationPolicyClaim) []interface{} {
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

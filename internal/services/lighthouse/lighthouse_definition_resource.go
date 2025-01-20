// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lighthouse

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2022-10-01/registrationdefinitions"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceLighthouseDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLighthouseDefinitionCreateUpdate,
		Read:   resourceLighthouseDefinitionRead,
		Update: resourceLighthouseDefinitionCreateUpdate,
		Delete: resourceLighthouseDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := registrationdefinitions.ParseScopedRegistrationDefinitionID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managing_tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubscriptionID,
			},

			"authorization": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"role_definition_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"principal_display_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"delegated_role_definition_ids": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsUUID,
							},
						},
					},
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"eligible_authorization": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"principal_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"role_definition_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"just_in_time_access_policy": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"multi_factor_auth_provider": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(registrationdefinitions.MultiFactorAuthProviderAzure),
										}, false),
									},

									"maximum_activation_duration": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "PT8H",
										ValidateFunc: azValidate.ISO8601Duration,
									},

									"approver": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"principal_id": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validation.IsUUID,
												},

												"principal_display_name": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
											},
										},
									},
								},
							},
						},

						"principal_display_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"lighthouse_definition_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"plan": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"publisher": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"product": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},
	}
}

func resourceLighthouseDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseDefinitionID := d.Get("lighthouse_definition_id").(string)
	if lighthouseDefinitionID == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("generating UUID for Lighthouse Definition: %+v", err)
		}

		lighthouseDefinitionID = uuid
	}

	id := registrationdefinitions.NewScopedRegistrationDefinitionID(d.Get("scope").(string), lighthouseDefinitionID)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_lighthouse_definition", id.ID())
		}
	}
	authorizations := expandLighthouseDefinitionAuthorization(d.Get("authorization").(*pluginsdk.Set).List())
	parameters := registrationdefinitions.RegistrationDefinition{
		Plan: expandLighthouseDefinitionPlan(d.Get("plan").([]interface{})),
		Properties: &registrationdefinitions.RegistrationDefinitionProperties{
			Description:                utils.String(d.Get("description").(string)),
			Authorizations:             authorizations,
			RegistrationDefinitionName: utils.String(d.Get("name").(string)),
			ManagedByTenantId:          d.Get("managing_tenant_id").(string),
		},
	}

	if v, ok := d.GetOk("eligible_authorization"); ok {
		parameters.Properties.EligibleAuthorizations = expandLighthouseDefinitionEligibleAuthorization(v.(*pluginsdk.Set).List())
	}

	// NOTE: this API call uses DefinitionId then Scope - check in the future
	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLighthouseDefinitionRead(d, meta)
}

func resourceLighthouseDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registrationdefinitions.ParseScopedRegistrationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("lighthouse_definition_id", id.RegistrationDefinitionId)
	d.Set("scope", id.Scope)

	if model := resp.Model; model != nil {
		if err := d.Set("plan", flattenLighthouseDefinitionPlan(model.Plan)); err != nil {
			return fmt.Errorf("setting `plan`: %+v", err)
		}

		if props := model.Properties; props != nil {
			if err := d.Set("authorization", flattenLighthouseDefinitionAuthorization(props.Authorizations)); err != nil {
				return fmt.Errorf("setting `authorization`: %+v", err)
			}
			if err := d.Set("eligible_authorization", flattenLighthouseDefinitionEligibleAuthorization(props.EligibleAuthorizations)); err != nil {
				return fmt.Errorf("setting `eligible_authorization`: %+v", err)
			}
			d.Set("description", props.Description)
			d.Set("name", props.RegistrationDefinitionName)
			d.Set("managing_tenant_id", props.ManagedByTenantId)
		}
	}

	return nil
}

func resourceLighthouseDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := registrationdefinitions.ParseScopedRegistrationDefinitionID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func flattenLighthouseDefinitionAuthorization(input []registrationdefinitions.Authorization) []interface{} {
	results := make([]interface{}, 0)
	for _, item := range input {
		principalIDDisplayName := ""
		if item.PrincipalIdDisplayName != nil {
			principalIDDisplayName = *item.PrincipalIdDisplayName
		}

		results = append(results, map[string]interface{}{
			"role_definition_id":            item.RoleDefinitionId,
			"principal_id":                  item.PrincipalId,
			"principal_display_name":        principalIDDisplayName,
			"delegated_role_definition_ids": utils.FlattenStringSlice(item.DelegatedRoleDefinitionIds),
		})
	}

	return results
}

func expandLighthouseDefinitionAuthorization(input []interface{}) []registrationdefinitions.Authorization {
	results := make([]registrationdefinitions.Authorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		delegatedRoleDefinitionIds := utils.ExpandStringSlice(v["delegated_role_definition_ids"].(*pluginsdk.Set).List())
		result := registrationdefinitions.Authorization{
			RoleDefinitionId:           v["role_definition_id"].(string),
			PrincipalId:                v["principal_id"].(string),
			PrincipalIdDisplayName:     utils.String(v["principal_display_name"].(string)),
			DelegatedRoleDefinitionIds: delegatedRoleDefinitionIds,
		}
		results = append(results, result)
	}
	return results
}

func expandLighthouseDefinitionPlan(input []interface{}) *registrationdefinitions.Plan {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &registrationdefinitions.Plan{
		Name:      raw["name"].(string),
		Publisher: raw["publisher"].(string),
		Product:   raw["product"].(string),
		Version:   raw["version"].(string),
	}
}

func flattenLighthouseDefinitionPlan(input *registrationdefinitions.Plan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name":      input.Name,
			"publisher": input.Publisher,
			"product":   input.Product,
			"version":   input.Version,
		},
	}
}

func expandLighthouseDefinitionEligibleAuthorization(input []interface{}) *[]registrationdefinitions.EligibleAuthorization {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]registrationdefinitions.EligibleAuthorization, 0, len(input))
	for _, item := range input {
		v := item.(map[string]interface{})

		result := registrationdefinitions.EligibleAuthorization{
			PrincipalId:            v["principal_id"].(string),
			RoleDefinitionId:       v["role_definition_id"].(string),
			JustInTimeAccessPolicy: expandLighthouseDefinitionJustInTimeAccessPolicy(v["just_in_time_access_policy"].([]interface{})),
		}

		if principalDisplayName := v["principal_display_name"].(string); principalDisplayName != "" {
			result.PrincipalIdDisplayName = utils.String(principalDisplayName)
		}

		results = append(results, result)
	}

	return &results
}

func expandLighthouseDefinitionJustInTimeAccessPolicy(input []interface{}) *registrationdefinitions.JustInTimeAccessPolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	justInTimeAccessPolicy := input[0].(map[string]interface{})

	result := registrationdefinitions.JustInTimeAccessPolicy{
		MaximumActivationDuration: utils.String(justInTimeAccessPolicy["maximum_activation_duration"].(string)),
		ManagedByTenantApprovers:  expandLighthouseDefinitionApprover(justInTimeAccessPolicy["approver"].(*pluginsdk.Set).List()),
	}

	multiFactorAuthProvider := registrationdefinitions.MultiFactorAuthProviderNone
	if v := justInTimeAccessPolicy["multi_factor_auth_provider"].(string); v != "" {
		multiFactorAuthProvider = registrationdefinitions.MultiFactorAuthProvider(v)
	}
	result.MultiFactorAuthProvider = multiFactorAuthProvider

	return &result
}

func expandLighthouseDefinitionApprover(input []interface{}) *[]registrationdefinitions.EligibleApprover {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]registrationdefinitions.EligibleApprover, 0)
	for _, v := range input {
		eligibleApprover := v.(map[string]interface{})

		result := registrationdefinitions.EligibleApprover{
			PrincipalId: eligibleApprover["principal_id"].(string),
		}

		if principalDisplayName := eligibleApprover["principal_display_name"].(string); principalDisplayName != "" {
			result.PrincipalIdDisplayName = utils.String(principalDisplayName)
		}

		results = append(results, result)
	}

	return &results
}

func flattenLighthouseDefinitionEligibleAuthorization(input *[]registrationdefinitions.EligibleAuthorization) []interface{} {
	if input == nil {
		return nil
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		result := map[string]interface{}{
			"principal_id":       item.PrincipalId,
			"role_definition_id": item.RoleDefinitionId,
		}

		if item.JustInTimeAccessPolicy != nil {
			result["just_in_time_access_policy"] = flattenLighthouseDefinitionJustInTimeAccessPolicy(item.JustInTimeAccessPolicy)
		}

		if item.PrincipalIdDisplayName != nil {
			result["principal_display_name"] = *item.PrincipalIdDisplayName
		}

		results = append(results, result)
	}

	return results
}

func flattenLighthouseDefinitionJustInTimeAccessPolicy(input *registrationdefinitions.JustInTimeAccessPolicy) []interface{} {
	if input == nil {
		return nil
	}

	var results []interface{}

	result := map[string]interface{}{}

	if v := input.MultiFactorAuthProvider; v != registrationdefinitions.MultiFactorAuthProviderNone {
		result["multi_factor_auth_provider"] = string(v)
	}

	if input.ManagedByTenantApprovers != nil {
		result["approver"] = flattenLighthouseDefinitionApprover(input.ManagedByTenantApprovers)
	}

	maximumActivationDuration := "PT8H"
	if input.MaximumActivationDuration != nil {
		maximumActivationDuration = *input.MaximumActivationDuration
	}
	result["maximum_activation_duration"] = maximumActivationDuration

	return append(results, result)
}

func flattenLighthouseDefinitionApprover(input *[]registrationdefinitions.EligibleApprover) []interface{} {
	if input == nil {
		return nil
	}

	results := make([]interface{}, 0, len(*input))
	for _, item := range *input {
		result := map[string]interface{}{
			"principal_id": item.PrincipalId,
		}

		if item.PrincipalIdDisplayName != nil {
			result["principal_display_name"] = *item.PrincipalIdDisplayName
		}

		results = append(results, result)
	}

	return results
}

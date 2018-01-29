package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceArmBuiltInRoleDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBuiltInRoleDefinitionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Auditor",
					"CloudAware Collector Storage Account Keys Access",
					"DB Admin",
					"Network Security Operator",
					"API Management Service Contributor",
					"API Management Service Operator Role",
					"API Management Service Reader Role",
					"Application Insights Component Contributor",
					"Application Insights Snapshot Debugger",
					"Automation Job Operator",
					"Automation Operator",
					"Automation Runbook Operator",
					"Azure Stack Registration Owner",
					"Backup Contributor",
					"Backup Operator",
					"Backup Reader",
					"Billing Reader",
					"BizTalk Contributor",
					"CDN Endpoint Contributor",
					"CDN Endpoint Reader",
					"CDN Profile Contributor",
					"CDN Profile Reader",
					"Classic Network Contributor",
					"Classic Storage Account Contributor",
					"Classic Storage Account Key Operator Service Role",
					"Classic Virtual Machine Contributor",
					"ClearDB MySQL DB Contributor",
					"Contributor",
					"Cosmos DB Account Reader Role",
					"Data Factory Contributor",
					"Data Lake Analytics Developer",
					"DevTest Labs User",
					"DNS Zone Contributor",
					"DocumentDB Account Contributor",
					"Intelligent Systems Account Contributor",
					"Key Vault Contributor",
					"Lab Accounts User",
					"Log Analytics Contributor",
					"Log Analytics Reader",
					"Logic App Contributor",
					"Logic App Operator",
					"Managed Identity Contributor",
					"Managed Identity Operator",
					"Monitoring Contributor",
					"Monitoring Reader",
					"Network Contributor",
					"New Relic APM Account Contributor",
					"Owner",
					"Reader",
					"Redis Cache Contributor",
					"Scheduler Job Collections Contributor",
					"Search Service Contributor",
					"Security Admin",
					"Security Manager",
					"Security Reader",
					"Site Recovery Contributor",
					"Site Recovery Operator",
					"Site Recovery Reader",
					"SQL DB Contributor",
					"SQL Security Manager",
					"SQL Server Contributor",
					"Storage Account Contributor",
					"Storage Account Key Operator Service Role",
					"Support Request Contributor",
					"Traffic Manager Contributor",
					"User Access Administrator",
					// TODO: make this `Virtual Machine Contributor` and handle deprecation
					"VirtualMachineContributor",
					"Web Plan Contributor",
					"Website Contributor",
				}, false),
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"not_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"assignable_scopes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmBuiltInRoleDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).roleDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	filter := fmt.Sprintf("roleName eq '%s'", name)
	roleDefinitions, err := client.List(ctx, "", filter)
	if err != nil {
		return fmt.Errorf("Error loading Role Definition List: %+v", err)
	}
	if len(roleDefinitions.Values()) != 1 {
		return fmt.Errorf("Error loading Role Definition List: could not find role '%s'", name)
	}

	roleDefinitionId := *roleDefinitions.Values()[0].ID

	d.SetId(roleDefinitionId)

	role, err := client.GetByID(ctx, roleDefinitionId)
	if err != nil {
		return fmt.Errorf("Error loading Role Definition: %+v", err)
	}

	if props := role.Properties; props != nil {
		d.Set("name", props.RoleName)
		d.Set("description", props.Description)
		d.Set("type", props.Type)

		permissions := flattenRoleDefinitionPermissions(props.Permissions)
		if err := d.Set("permissions", permissions); err != nil {
			return err
		}

		assignableScopes := flattenRoleDefinitionAssignableScopes(props.AssignableScopes)
		if err := d.Set("assignable_scopes", assignableScopes); err != nil {
			return err
		}
	}

	return nil
}

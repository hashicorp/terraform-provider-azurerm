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
	roleDefinitionIds := map[string]string{
		"Auditor": "/providers/Microsoft.Authorization/roleDefinitions/8da36ee2-a5e5-4c9d-ad4a-6d281636c8b1",
		"CloudAware Collector Storage Account Keys Access": "/providers/Microsoft.Authorization/roleDefinitions/caea9d35-7d6c-4bd5-ac15-1e92cb428f43",
		"DB Admin":                                          "/providers/Microsoft.Authorization/roleDefinitions/d80f878a-b665-46fd-a179-5b231a7a126e",
		"Network Security Operator":                         "/providers/Microsoft.Authorization/roleDefinitions/ff1e2e01-fceb-4dc1-a91a-ffb9ca766f7d",
		"API Management Service Contributor":                "/providers/Microsoft.Authorization/roleDefinitions/312a565d-c81f-4fd8-895a-4e21e48d571c",
		"API Management Service Operator Role":              "/providers/Microsoft.Authorization/roleDefinitions/e022efe7-f5ba-4159-bbe4-b44f577e9b61",
		"API Management Service Reader Role":                "/providers/Microsoft.Authorization/roleDefinitions/71522526-b88f-4d52-b57f-d31fc3546d0d",
		"Application Insights Component Contributor":        "/providers/Microsoft.Authorization/roleDefinitions/ae349356-3a1b-4a5e-921d-050484c6347e",
		"Application Insights Snapshot Debugger":            "/providers/Microsoft.Authorization/roleDefinitions/08954f03-6346-4c2e-81c0-ec3a5cfae23b",
		"Automation Job Operator":                           "/providers/Microsoft.Authorization/roleDefinitions/4fe576fe-1146-4730-92eb-48519fa6bf9f",
		"Automation Operator":                               "/providers/Microsoft.Authorization/roleDefinitions/d3881f73-407a-4167-8283-e981cbba0404",
		"Automation Runbook Operator":                       "/providers/Microsoft.Authorization/roleDefinitions/5fb5aef8-1081-4b8e-bb16-9d5d0385bab5",
		"Azure Stack Registration Owner":                    "/providers/Microsoft.Authorization/roleDefinitions/6f12a6df-dd06-4f3e-bcb1-ce8be600526a",
		"Backup Contributor":                                "/providers/Microsoft.Authorization/roleDefinitions/5e467623-bb1f-42f4-a55d-6e525e11384b",
		"Backup Operator":                                   "/providers/Microsoft.Authorization/roleDefinitions/00c29273-979b-4161-815c-10b084fb9324",
		"Backup Reader":                                     "/providers/Microsoft.Authorization/roleDefinitions/a795c7a0-d4a2-40c1-ae25-d81f01202912",
		"Billing Reader":                                    "/providers/Microsoft.Authorization/roleDefinitions/fa23ad8b-c56e-40d8-ac0c-ce449e1d2c64",
		"BizTalk Contributor":                               "/providers/Microsoft.Authorization/roleDefinitions/5e3c6656-6cfa-4708-81fe-0de47ac73342",
		"CDN Endpoint Contributor":                          "/providers/Microsoft.Authorization/roleDefinitions/426e0c7f-0c7e-4658-b36f-ff54d6c29b45",
		"CDN Endpoint Reader":                               "/providers/Microsoft.Authorization/roleDefinitions/871e35f6-b5c1-49cc-a043-bde969a0f2cd",
		"CDN Profile Contributor":                           "/providers/Microsoft.Authorization/roleDefinitions/ec156ff8-a8d1-4d15-830c-5b80698ca432",
		"CDN Profile Reader":                                "/providers/Microsoft.Authorization/roleDefinitions/8f96442b-4075-438f-813d-ad51ab4019af",
		"Classic Network Contributor":                       "/providers/Microsoft.Authorization/roleDefinitions/b34d265f-36f7-4a0d-a4d4-e158ca92e90f",
		"Classic Storage Account Contributor":               "/providers/Microsoft.Authorization/roleDefinitions/86e8f5dc-a6e9-4c67-9d15-de283e8eac25",
		"Classic Storage Account Key Operator Service Role": "/providers/Microsoft.Authorization/roleDefinitions/985d6b00-f706-48f5-a6fe-d0ca12fb668d",
		"Classic Virtual Machine Contributor":               "/providers/Microsoft.Authorization/roleDefinitions/d73bb868-a0df-4d4d-bd69-98a00b01fccb",
		"ClearDB MySQL DB Contributor":                      "/providers/Microsoft.Authorization/roleDefinitions/9106cda0-8a86-4e81-b686-29a22c54effe",
		"Contributor":                                       "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c",
		"Cosmos DB Account Reader Role":                     "/providers/Microsoft.Authorization/roleDefinitions/fbdf93bf-df7d-467e-a4d2-9458aa1360c8",
		"Data Factory Contributor":                          "/providers/Microsoft.Authorization/roleDefinitions/673868aa-7521-48a0-acc6-0f60742d39f5",
		"Data Lake Analytics Developer":                     "/providers/Microsoft.Authorization/roleDefinitions/47b7735b-770e-4598-a7da-8b91488b4c88",
		"DevTest Labs User":                                 "/providers/Microsoft.Authorization/roleDefinitions/76283e04-6283-4c54-8f91-bcf1374a3c64",
		"DNS Zone Contributor":                              "/providers/Microsoft.Authorization/roleDefinitions/befefa01-2a29-4197-83a8-272ff33ce314",
		"DocumentDB Account Contributor":                    "/providers/Microsoft.Authorization/roleDefinitions/5bd9cd88-fe45-4216-938b-f97437e15450",
		"Intelligent Systems Account Contributor":           "/providers/Microsoft.Authorization/roleDefinitions/03a6d094-3444-4b3d-88af-7477090a9e5e",
		"Key Vault Contributor":                             "/providers/Microsoft.Authorization/roleDefinitions/f25e0fa2-a7c8-4377-a976-54943a77a395",
		"Lab Accounts User":                                 "/providers/Microsoft.Authorization/roleDefinitions/b97fb8bc-a8b2-4522-a38b-dd33c7e65ead",
		"Log Analytics Contributor":                         "/providers/Microsoft.Authorization/roleDefinitions/92aaf0da-9dab-42b6-94a3-d43ce8d16293",
		"Log Analytics Reader":                              "/providers/Microsoft.Authorization/roleDefinitions/73c42c96-874c-492b-b04d-ab87d138a893",
		"Logic App Contributor":                             "/providers/Microsoft.Authorization/roleDefinitions/87a39d53-fc1b-424a-814c-f7e04687dc9e",
		"Logic App Operator":                                "/providers/Microsoft.Authorization/roleDefinitions/515c2055-d9d4-4321-b1b9-bd0c9a0f79fe",
		"Managed Identity Contributor":                      "/providers/Microsoft.Authorization/roleDefinitions/e40ec5ca-96e0-45a2-b4ff-59039f2c2b59",
		"Managed Identity Operator":                         "/providers/Microsoft.Authorization/roleDefinitions/f1a07417-d97a-45cb-824c-7a7467783830",
		"Monitoring Contributor":                            "/providers/Microsoft.Authorization/roleDefinitions/749f88d5-cbae-40b8-bcfc-e573ddc772fa",
		"Monitoring Reader":                                 "/providers/Microsoft.Authorization/roleDefinitions/43d0d8ad-25c7-4714-9337-8ba259a9fe05",
		"Network Contributor":                               "/providers/Microsoft.Authorization/roleDefinitions/4d97b98b-1d4f-4787-a291-c67834d212e7",
		"New Relic APM Account Contributor":                 "/providers/Microsoft.Authorization/roleDefinitions/5d28c62d-5b37-4476-8438-e587778df237",
		"Owner":                                     "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
		"Reader":                                    "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
		"Redis Cache Contributor":                   "/providers/Microsoft.Authorization/roleDefinitions/e0f68234-74aa-48ed-b826-c38b57376e17",
		"Scheduler Job Collections Contributor":     "/providers/Microsoft.Authorization/roleDefinitions/188a0f2f-5c9e-469b-ae67-2aa5ce574b94",
		"Search Service Contributor":                "/providers/Microsoft.Authorization/roleDefinitions/7ca78c08-252a-4471-8644-bb5ff32d4ba0",
		"Security Admin":                            "/providers/Microsoft.Authorization/roleDefinitions/fb1c8493-542b-48eb-b624-b4c8fea62acd",
		"Security Manager":                          "/providers/Microsoft.Authorization/roleDefinitions/e3d13bf0-dd5a-482e-ba6b-9b8433878d10",
		"Security Reader":                           "/providers/Microsoft.Authorization/roleDefinitions/39bc4728-0917-49c7-9d2c-d95423bc2eb4",
		"Site Recovery Contributor":                 "/providers/Microsoft.Authorization/roleDefinitions/6670b86e-a3f7-4917-ac9b-5d6ab1be4567",
		"Site Recovery Operator":                    "/providers/Microsoft.Authorization/roleDefinitions/494ae006-db33-4328-bf46-533a6560a3ca",
		"Site Recovery Reader":                      "/providers/Microsoft.Authorization/roleDefinitions/dbaa88c4-0c30-4179-9fb3-46319faa6149",
		"SQL DB Contributor":                        "/providers/Microsoft.Authorization/roleDefinitions/9b7fa17d-e63e-47b0-bb0a-15c516ac86ec",
		"SQL Security Manager":                      "/providers/Microsoft.Authorization/roleDefinitions/056cd41c-7e88-42e1-933e-88ba6a50c9c3",
		"SQL Server Contributor":                    "/providers/Microsoft.Authorization/roleDefinitions/6d8ee4ec-f05a-4a1d-8b00-a9b17e38b437",
		"Storage Account Contributor":               "/providers/Microsoft.Authorization/roleDefinitions/17d1049b-9a84-46fb-8f53-869881c3d3ab",
		"Storage Account Key Operator Service Role": "/providers/Microsoft.Authorization/roleDefinitions/81a9662b-bebf-436f-a333-f67b29880f12",
		"Support Request Contributor":               "/providers/Microsoft.Authorization/roleDefinitions/cfd33db0-3dd1-45e3-aa9d-cdbdf3b6f24e",
		"Traffic Manager Contributor":               "/providers/Microsoft.Authorization/roleDefinitions/a4b10055-b0c7-44c2-b00f-c7b5b3550cf7",
		"User Access Administrator":                 "/providers/Microsoft.Authorization/roleDefinitions/18d7d88d-d35e-4fb5-a5c3-7773c20a72d9",
		"VirtualMachineContributor":                 "/providers/Microsoft.Authorization/roleDefinitions/9980e02c-c2be-4d73-94e8-173b1dc7cf3c",
		"Web Plan Contributor":                      "/providers/Microsoft.Authorization/roleDefinitions/2cc479cb-7b4d-49a8-b449-8c00fd0f0a4b",
		"Website Contributor":                       "/providers/Microsoft.Authorization/roleDefinitions/de139f84-1756-47ae-9be6-808fbbe84772",
	}
	roleDefinitionId := roleDefinitionIds[name]

	d.SetId(roleDefinitionId)

	role, err := client.GetByID(ctx, roleDefinitionId)
	if err != nil {
		return fmt.Errorf("Error loadng Role Definition: %+v", err)
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

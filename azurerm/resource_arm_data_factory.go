package azurerm

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmDataFactory() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryCreateOrUpdate,
		Read:   resourceArmDataFactoryRead,
		Update: resourceArmDataFactoryCreateOrUpdate,
		Delete: resourceArmDataFactoryDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Some error message`,
				),
			},

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"github_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name":         {},
						"collaboration_branch": {},
						"host_name":            {},
						"repository_name":      {},
						"root_folder":          {},
					},
				},
				ConflictsWith: []string{"vsts_configuration"},
			},

			"vsts_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_name":         {},
						"collaboration_branch": {},
						"project_name":         {},
						"repository_name":      {},
						"root_folder":          {},
						"tenant_id":            {},
					},
				},
				ConflictsWith: []string{"github_configuration"},
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

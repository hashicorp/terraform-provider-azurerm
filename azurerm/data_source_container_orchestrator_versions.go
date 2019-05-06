package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmContainerOrchestratorVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataArmSourceContainerOrchestratorVersionsRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"orchestrator_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Kubernetes",
			},
			"version_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataArmSourceContainerOrchestratorVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerServicesClient
	ctx := meta.(*ArmClient).StopContext

	location := d.Get("location").(string)
	ot := d.Get("orchestrator_type").(string)

	listResp, err := client.ListOrchestrators(ctx, location, "managedClusters")
	if err != nil {
		if utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("Error: Container Orchestrators not found for location %q", location)
		}
		return fmt.Errorf("Error making Read request on AzureRM Container Service %q: %+v", location, err)
	}

	bv, err := version.NewVersion("0.0.0")
	if err != nil {
		return fmt.Errorf("Cannot set version baseline (likely an issue in go-version): %+v", err)
	}
	lv := bv
	versions := []string{}
	for _, rawV := range *listResp.Orchestrators {
		if rawV.OrchestratorType != nil && rawV.OrchestratorVersion != nil {
			if !strings.HasPrefix(*rawV.OrchestratorVersion, d.Get("version_prefix").(string)) {
				continue
			}
			if *rawV.OrchestratorType == ot {
				versions = append(versions, *rawV.OrchestratorVersion)
				v, err := version.NewVersion(*rawV.OrchestratorVersion)
				if err != nil {
					log.Printf("[WARN] Cannot parse orchestrator version '%s'. Skipping because %s", *rawV.OrchestratorVersion, err)
					continue
				}
				if v.GreaterThan(lv) {
					lv = v
				}
			}
		}
	}
	d.SetId(*listResp.ID)
	d.Set("versions", versions)
	if lv.GreaterThan(bv) {
		d.Set("latest_version", lv.Original())
	}

	return nil
}

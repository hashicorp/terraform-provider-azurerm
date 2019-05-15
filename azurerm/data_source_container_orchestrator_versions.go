package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmKubernetesServiceVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmKubernetesServiceVersionsRead,

		Schema: map[string]*schema.Schema{
			"location": locationSchema(),
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

func dataSourceArmKubernetesServiceVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).containerServicesClient
	ctx := meta.(*ArmClient).StopContext

	location := d.Get("location").(string)

	listResp, err := client.ListOrchestrators(ctx, location, "managedClusters")
	if err != nil {
		if utils.ResponseWasNotFound(listResp.Response) {
			return fmt.Errorf("Error: Container Orchestrators not found for location %q", location)
		}
		return fmt.Errorf("Error retrieving Kubernetes Versions in %q: %+v", location, err)
	}

	lv, err := version.NewVersion("0.0.0")
	if err != nil {
		return fmt.Errorf("Cannot set version baseline (likely an issue in go-version): %+v", err)
	}
	versions := []string{}
	for _, rawV := range *listResp.Orchestrators {
		if rawV.OrchestratorType != nil && rawV.OrchestratorVersion != nil {
			if *rawV.OrchestratorType == "Kubernetes" {
				if !strings.HasPrefix(*rawV.OrchestratorVersion, d.Get("version_prefix").(string)) {
					continue
				}
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
	d.Set("latest_version", lv.Original())

	return nil
}

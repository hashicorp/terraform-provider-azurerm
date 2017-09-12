package azurerm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/arm/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerCreate,
		Read:   resourceArmPostgreSQLServerRead,
		Update: resourceArmPostgreSQLServerUpdate,
		Delete: resourceArmPostgreSQLServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"PGSQLB50",
								"PGSQLB100",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validateIntInSlice([]int{
								50,
								100,
							}),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.Basic),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.NineFullStopFive),
					string(postgresql.NineFullStopSix),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ForceNew:         true,
			},

			"storage_mb": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateIntInSlice([]int{
					51200,
					102400,
				}),
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.SslEnforcementEnumDisabled),
					string(postgresql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPostgreSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)

	adminLogin := d.Get("administrator_login").(string)
	adminLoginPassword := d.Get("administrator_login_password").(string)
	sslEnforcement := d.Get("ssl_enforcement").(string)
	version := d.Get("version").(string)
	storageMB := d.Get("storage_mb").(int)

	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmPostgreSQLServerSku(d, storageMB)

	properties := postgresql.ServerForCreate{
		Location: &location,
		Sku:      sku,
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			Version:                    postgresql.ServerVersion(version),
			StorageMB:                  utils.Int64(int64(storageMB)),
			SslEnforcement:             postgresql.SslEnforcementEnum(sslEnforcement),
			AdministratorLogin:         utils.String(adminLogin),
			AdministratorLoginPassword: utils.String(adminLoginPassword),
			CreateMode:                 postgresql.CreateModeDefault,
		},
		Tags: expandTags(tags),
	}

	_, error := client.CreateOrUpdate(resGroup, name, properties, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server update.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	adminLoginPassword := d.Get("administrator_login_password").(string)
	sslEnforcement := d.Get("ssl_enforcement").(string)
	version := d.Get("version").(string)
	storageMB := d.Get("storage_mb").(int)
	sku := expandAzureRmPostgreSQLServerSku(d, storageMB)

	tags := d.Get("tags").(map[string]interface{})

	properties := postgresql.ServerUpdateParameters{
		Sku: sku,
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			SslEnforcement:             postgresql.SslEnforcementEnum(sslEnforcement),
			StorageMB:                  utils.Int64(int64(storageMB)),
			Version:                    postgresql.ServerVersion(version),
			AdministratorLoginPassword: utils.String(adminLoginPassword),
		},
		Tags: expandTags(tags),
	}

	_, error := client.Update(resGroup, name, properties, make(chan struct{}))
	err := <-error
	if err != nil {
		return err
	}

	read, err := client.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Server '%s' was not found (resource group '%s')", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Server %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("storage_mb", int(*resp.StorageMB))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))
	d.Set("sku", flattenPostgreSQLServerSku(resp.Sku))

	flattenAndSetTags(d, resp.Tags)

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

	return nil
}

func resourceArmPostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["servers"]

	_, deleteErr := client.Delete(resGroup, name, make(chan struct{}))
	err = <-deleteErr

	return err
}

func expandAzureRmPostgreSQLServerSku(d *schema.ResourceData, storageMB int) *postgresql.Sku {
	skus := d.Get("sku").(*schema.Set).List()
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)

	return &postgresql.Sku{
		Name:     utils.String(name),
		Capacity: utils.Int32(int32(capacity)),
		Tier:     postgresql.SkuTier(tier),
		Size:     utils.String(strconv.Itoa(storageMB)),
	}
}

func flattenPostgreSQLServerSku(resp *postgresql.Sku) []interface{} {
	values := map[string]interface{}{}

	values["name"] = *resp.Name
	values["capacity"] = int(*resp.Capacity)
	values["tier"] = string(resp.Tier)

	sku := []interface{}{values}
	return sku
}

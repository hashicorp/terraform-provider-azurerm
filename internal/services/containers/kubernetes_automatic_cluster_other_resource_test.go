// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

func TestAccKubernetesAutomaticCluster_basicVMSS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVMSSConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicVMSSConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImportConfig(data),
			ExpectError: acceptance.RequiresImportError("azurerm_kubernetes_automatic_cluster"),
		},
	})
}

func TestAccKubernetesAutomaticCluster_linuxProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxProfileConfig(data, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kube_config.0.cluster_ca_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.host").Exists(),
				check.That(data.ResourceName).Key("linux_profile.0.admin_username").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_nodeResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	nodeResourceGroupName := fmt.Sprintf("acctestRGAKS-%d", data.RandomInteger)
	nodeResourceGroupId := commonids.NewResourceGroupID(data.Subscriptions.Primary, nodeResourceGroupName)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nodeResourceGroupConfig(data, nodeResourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("node_resource_group_id").HasValue(nodeResourceGroupId.ID()),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeSkuTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.skuConfigFree(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.skuConfigStandard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.skuConfigFree(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_upgrade(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeConfig(data, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
				check.That(data.ResourceName).Key("current_kubernetes_version").Exists(),
			),
		},
		{
			Config: r.upgradeConfig(data, currentKubernetesVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(currentKubernetesVersion),
				check.That(data.ResourceName).Key("current_kubernetes_version").Exists(),
			),
		},
	})
}

func TestAccKubernetesAutomaticCluster_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.tagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.tagsUpdatedConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_privateClusterPublicFqdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateClusterPublicFqdn(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.privateClusterPublicFqdn(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_windowsProfile(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsProfileConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kube_config.0.cluster_ca_certificate").Exists(),
				check.That(data.ResourceName).Key("kube_config.0.host").Exists(),
				check.That(data.ResourceName).Key("linux_profile.0.admin_username").Exists(),
				check.That(data.ResourceName).Key("windows_profile.0.admin_username").Exists(),
			),
		},
		data.ImportStep("windows_profile.0.admin_password"),
	})
}

func TestAccKubernetesAutomaticCluster_windowsProfileGMSA(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsProfileGMSAEmptyPropertyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("windows_profile.0.admin_password"),
	})
}

func TestAccKubernetesAutomaticCluster_windowsProfileLicense(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsProfileLicense(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("windows_profile.0.admin_password"),
	})
}

func TestAccKubernetesAutomaticCluster_diskEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.diskEncryptionConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("disk_encryption_set_id").Exists(),
			),
		},
		data.ImportStep(
			"windows_profile.0.admin_password",
		),
	})
}

func TestAccKubernetesAutomaticCluster_upgradeChannel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.upgradeChannelConfig(data, olderKubernetesAutomaticVersion),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("kubernetes_version").HasValue(olderKubernetesAutomaticVersion),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_microsoftDefender(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.microsoftDefender(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.microsoftDefenderDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingWithMultipleDnsZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingWithMultipleDnsZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingWithEmptyDnsZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingWitEmptyDnsZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_webAppRoutingWithNginxControllerType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.webAppRoutingWithNginxController(data, "None"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.webAppRoutingWithNginxController(data, "Internal"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.webAppRoutingWithNginxController(data, "External"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.webAppRoutingWithNginxController(data, "AnnotationControlled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.webAppRoutingWithNginxController(data, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("web_app_routing.0.web_app_routing_identity.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_azureMonitorKubernetesMetrics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureMonitorKubernetesMetricsEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.azureMonitorKubernetesMetricsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.azureMonitorKubernetesMetricsDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_supportPlanKubernetesOfficial(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.supportPlanKubernetesOfficial(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_costAnalysis(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.costAnalysisEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.costAnalysisEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.costAnalysisEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_customCaTrustCerts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	fakeCertList := []string{
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjRENDQWxpZ0F3SUJBZ0lFU1QwSUhEQU5CZ2txaGtpRzl3MEJBUXNGQURCUk1Rc3dDUVlEVlFRR0V3SlEKVERFTk1Bc0dBMVVFQXd3RVZHVnpkREVWTUJNR0ExVUVCd3dNUkdWbVlYVnNkQ0JEYVhSNU1Sd3dHZ1lEVlFRSwpEQk5FWldaaGRXeDBJRU52YlhCaGJua2dUSFJrTUI0WERUSXpNRFV5T0RFeE1qY3dNMW9YRFRNek1EVXlOVEV4Ck1qY3dNMW93VVRFTE1Ba0dBMVVFQmhNQ1VFd3hEVEFMQmdOVkJBTU1CRlJsYzNReEZUQVRCZ05WQkFjTURFUmwKWm1GMWJIUWdRMmwwZVRFY01Cb0dBMVVFQ2d3VFJHVm1ZWFZzZENCRGIyMXdZVzU1SUV4MFpEQ0NBU0l3RFFZSgpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFLN2JIYWtxSkdRMWVBOUFHUmlhNGl2anNDRXlGMDhDCjNpSzJZeWthNkREeldmTk1tRWpOUjJiQVZOMEhlLy9pWTd1VjJ2dXl6V1UxMzZGVkdMZkdyeTZGOHNQQUZaSzYKSE4vcWk1QVp6MUpoOGdWSTRwS1pjZEFxQS81clF3VVlvWVN3Q245dGVOYytsbU1ZUk5OcTVwdlV2NjcrNEM3MgpPc3BOSUxSclhBbWNUb1YveVRZVzFKWDBOeEJJSHZZaFZXUE9LQXpRZDQ5UEpSeFpqMUgydCszMEFsazgzTDFwClFzTGx2SzV3MjJpeXdkYVpRN1lmV0xXd1hPQzVPWXdRTUw1R3BHUFNQaEdxdjhqSUhpcHBVeTdrRDlNWFFZOFoKdDl2QkczMzVWSEdlUjI2QnNQQXRFbTJjR05ocjA5cmRvdWJGd2tDR05OYXNVamFoVW9CKzhPY0NBd0VBQWFOUQpNRTR3SFFZRFZSME9CQllFRk9CNmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1COEdBMVVkSXdRWU1CYUFGT0I2CmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1Bd0dBMVVkRXdRRk1BTUJBZjh3RFFZSktvWklodmNOQVFFTEJRQUQKZ2dFQkFKTklHdHJpeFlCRUc1Yy9iQWdOMHlMOEJvOW9nN29ha0hVMUc5TjBxOUNWWXhjOVhma2ZUaEhYOVBUeApMbVNGcHJEQlAyYnVGTzVIUDFpbnNFT1E2N1lGanAvRjVJWGdaQ2twZUpGdDBTL0R3N2ZRbFJJN2RCNGQzNmIzCmE1R2txU0M4aFlZemxLUm9DRGNhalp4QmdoVUFxK0tnTnV4RmNsM1Fnd1Uyam1QbkU4a1A4TmgyM3hlVUJ3WEkKL3pqbU1rdjV4SFhKdHBpdlpzTlpSSUttQW56RU9TWGlRK2JMTStTdlhtSkhYd29YYTZyTXg4YmkySzV4WkhIRwpkUHA1TnQ3L2dxOUdXcm95SkVjSFpEclBiSnR2WGFibTZYUXpxTTFYUzA3SDlaSFBXc0dENGlBM1k0T3JUUlRCClZ5blRPUDl5U3cwbklaVEk4YjZuR2RHTzBOOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==",
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURmakNDQW1hZ0F3SUJBZ0lFZnlWdk56QU5CZ2txaGtpRzl3MEJBUXNGQURCWU1Rc3dDUVlEVlFRR0V3SlEKVERFVU1CSUdBMVVFQXd3TFJtRnJaU0JEWlhKMElESXhGVEFUQmdOVkJBY01ERVJsWm1GMWJIUWdRMmwwZVRFYwpNQm9HQTFVRUNnd1RSR1ZtWVhWc2RDQkRiMjF3WVc1NUlFeDBaREFlRncweU16QTJNRFF3TnpJME1qZGFGdzB5Ck5UQTJNRE13TnpJME1qZGFNRmd4Q3pBSkJnTlZCQVlUQWxCTU1SUXdFZ1lEVlFRRERBdEdZV3RsSUVObGNuUWcKTWpFVk1CTUdBMVVFQnd3TVJHVm1ZWFZzZENCRGFYUjVNUnd3R2dZRFZRUUtEQk5FWldaaGRXeDBJRU52YlhCaApibmtnVEhSa01JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMENTdVdUaGNjSG5MCkhFdjk4SUVNc2JLY3h4YVh4YTZiRXl1Yy9sUjRackpVN2p6eVlWNGVscTV5WTgwdDFCM0MyV3E2SXFoajErSGYKYW0xaStsU1FTejM1eWNnTWlwSWp2cUxKOVIzMVF0Wi9TRURkdGV2b2JqbytEa1dCOE55cG9Ia0pVbEIyQnR6ZgpOK09KeVFSdXU1b1cya2c5OE5Bd3JuTGpmQ0lremVWcFh5d0l4Tkx2ZmFrVGxpNWpYdG9WWG5pOTU5bmtINWVwClkrRnVoSEQwaU5CS25XYVkxR2QwVGhhSHNwTERmNFUycmo2WE5SZHd6QVZoVkdhUm02cndvSHRZeDVrYys1ZWMKQ0F4UEdRWFRzTzJUTHVrQzJ2YXI0M3RUM0ZjSC9taDRST2JaaThZS2xSQ3Fldm1QU1RmZ293RUFkTjlvSmxyRApXN2lzN2NnQjhRSURBUUFCbzFBd1RqQWRCZ05WSFE0RUZnUVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293Ckh3WURWUjBqQkJnd0ZvQVVuRkRqN0pBQW9WZ2NzQkgyNzdMOHZlM0Q4U293REFZRFZSMFRCQVV3QXdFQi96QU4KQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBT0diT0Zyek4rN2YxbzhJSDNtMXZxT3IyTUtvNEZMWExGRjBVbEhkNApwZXRhL05aQjArUmQ3TnUrOCtnUnlUbEJWZU9EZjN5SXU0TlFCUU92MlNqdS9Jakd0MUtmaUF3WkUwT1RUQXc3CnhIWStsMVBJWEFFVWNqNk00cjFKQzc4ZVZrc2pycTZoV1RPZ0RrSVZuRjY3bXlReXduR25EY1k0d0Fqc2pUajgKKzR4NTIrRi9QaVNQVGtjUFNuN0s2UjQzaEt5QUs2Z0poOHE5cVNhME5RQ2U2czhwTGU2SVY5SElWVVFFVERVOQpsM1VWWHNBMGx4dlB0blU1TXo2QWQ5cDA5L2w4d3o0cUdBdGFCUEd3K0R2cTNlaHdTd2VZZ3VHSktDQjhjb01JCjJRVUo0Zi9mNkFNVWtMeWxYZ3RSUEt1QjA3d3YwTmk1eWI5MjlFY1FJQ0l2dFE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
	}

	fakeCertList2 := []string{
		"LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURjRENDQWxpZ0F3SUJBZ0lFU1QwSUhEQU5CZ2txaGtpRzl3MEJBUXNGQURCUk1Rc3dDUVlEVlFRR0V3SlEKVERFTk1Bc0dBMVVFQXd3RVZHVnpkREVWTUJNR0ExVUVCd3dNUkdWbVlYVnNkQ0JEYVhSNU1Sd3dHZ1lEVlFRSwpEQk5FWldaaGRXeDBJRU52YlhCaGJua2dUSFJrTUI0WERUSXpNRFV5T0RFeE1qY3dNMW9YRFRNek1EVXlOVEV4Ck1qY3dNMW93VVRFTE1Ba0dBMVVFQmhNQ1VFd3hEVEFMQmdOVkJBTU1CRlJsYzNReEZUQVRCZ05WQkFjTURFUmwKWm1GMWJIUWdRMmwwZVRFY01Cb0dBMVVFQ2d3VFJHVm1ZWFZzZENCRGIyMXdZVzU1SUV4MFpEQ0NBU0l3RFFZSgpLb1pJaHZjTkFRRUJCUUFEZ2dFUEFEQ0NBUW9DZ2dFQkFLN2JIYWtxSkdRMWVBOUFHUmlhNGl2anNDRXlGMDhDCjNpSzJZeWthNkREeldmTk1tRWpOUjJiQVZOMEhlLy9pWTd1VjJ2dXl6V1UxMzZGVkdMZkdyeTZGOHNQQUZaSzcKSE4vcWk1QVp6MUpoOGdWSTRwS1pjZEFxQS81clF3VVlvWVN3Q245dGVOYytsbU1ZUk5OcTVwdlV2NjcrNEM3MgpPc3BOSUxSclhBbWNUb1YveVRZVzFKWDBOeEJJSHZZaFZXUE9LQXpRZDQ5UEpSeFpqMUgydCszMEFsazgzTDFwClFzTGx2SzV3MjJpeXdkYVpRN1lmV0xXd1hPQzVPWXdRTUw1R3BHUFNQaEdxdjhqSUhpcHBVeTdrRDlNWFFZOFoKdDl2QkczMzVWSEdlUjI2QnNQQXRFbTJjR05ocjA5cmRvdWJGd2tDR05OYXNVamFoVW9CKzhPY0NBd0VBQWFOUQpNRTR3SFFZRFZSME9CQllFRk9CNmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1COEdBMVVkSXdRWU1CYUFGT0I2CmNpTGtUL21Cc2xXSm5Na2phQzZqbjd4ek1Bd0dBMVVkRXdRRk1BTUJBZjh3RFFZSktvWklodmNOQVFFTEJRQUQKZ2dFQkFKTklHdHJpeFlCRUc1Yy9iQWdOMHlMOEJvOW9nN29ha0hVMUc5TjBxOUNWWXhjOVhma2ZUaEhYOVBUeApMbVNGcHJEQlAyYnVGTzVIUDFpbnNFT1E2N1lGanAvRjVJWGdaQ2twZUpGdDBTL0R3N2ZRbFJJN2RCNGQzNmIzCmE1R2txU0M4aFlZemxLUm9DRGNhalp4QmdoVUFxK0tnTnV4RmNsM1Fnd1Uyam1QbkU4a1A4TmgyM3hlVUJ3WEkKL3pqbU1rdjV4SFhKdHBpdlpzTlpSSUttQW56RU9TWGlRK2JMTStTdlhtSkhYd29YYTZyTXg4YmkySzV4WkhIRwpkUHA1TnQ3L2dxOUdXcm95SkVjSFpEclBiSnR2WGFibTZYUXpxTTFYUzA3SDlaSFBXc0dENGlBM1k0T3JUUlRCClZ5blRPUDl5U3cwbklaVEk4YjZuR2RHTzBOOD0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQ==",
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.customCATrustCertificates(data, fakeCertList),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.0").Exists(),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.1").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.customCATrustCertificates(data, fakeCertList2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.0").Exists(),
				check.That(data.ResourceName).Key("custom_ca_trust_certificates_base64.1").DoesNotExist(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesAutomaticCluster_aiToolchainOperatorProfileToggle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_automatic_cluster", "test")
	r := KubernetesAutomaticClusterResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.aiToolchainOperatorProfile(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.aiToolchainOperatorProfile(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.aiToolchainOperatorProfile(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KubernetesAutomaticClusterResource) basicVMSSConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r KubernetesAutomaticClusterResource) requiresImportConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_automatic_cluster" "import" {
  name                = azurerm_kubernetes_automatic_cluster.test.name
  location            = azurerm_kubernetes_automatic_cluster.test.location
  resource_group_name = azurerm_kubernetes_automatic_cluster.test.resource_group_name
  dns_prefix          = azurerm_kubernetes_automatic_cluster.test.dns_prefix

  identity {
    type = "SystemAssigned"
  }
}
`, r.basicVMSSConfig(data))
}

func (KubernetesAutomaticClusterResource) linuxProfileConfig(data acceptance.TestData, keyData string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key_data = "%s"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, keyData)
}

func (KubernetesAutomaticClusterResource) nodeResourceGroupConfig(data acceptance.TestData, nodeResourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                     = "acctestaks%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  dns_prefix               = "acctestaks%d"
  node_resource_group_name = "%s"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, nodeResourceGroupName)
}

func (KubernetesAutomaticClusterResource) skuConfigStandard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) skuConfigFree(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) tagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    dimension = "C-137"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) tagsUpdatedConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    dimension = "D-99"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) upgradeConfig(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = "%s"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
  }

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, version, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) windowsProfileConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!h@2h1C0rP"
    group_managed_service_accounts {
      dns_server  = "10.10.0.10/2"
      root_domain = "contoso.com"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network {
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) windowsProfileGMSAEmptyPropertyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  linux_profile {
    admin_username = "acctestuser%[1]d"

    ssh_key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!h@2h1C0rP"
    group_managed_service_accounts {
      dns_server  = "vnet"
      root_domain = "vnet"
    }
  }

  identity {
    type = "SystemAssigned"
  }

  network {
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (KubernetesAutomaticClusterResource) windowsProfileLicense(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
  }

  windows_profile {
    admin_username = "azureuser"
    admin_password = "P@55W0rd1234!h@2h1C0rP"
    license        = "Windows_Server"
  }

  identity {
    type = "SystemAssigned"
  }

  network {
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) diskEncryptionConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                        = "acctestkeyvault%s"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "standard"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
  soft_delete_retention_days  = 7
}

resource "azurerm_key_vault_access_policy" "acctest" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "Create",
    "Delete",
    "Purge",
    "GetRotationPolicy",
  ]
}

resource "azurerm_key_vault_key" "test" {
  name         = "destestkey"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

  depends_on = [azurerm_key_vault_access_policy.acctest]
}

resource "azurerm_disk_encryption_set" "test" {
  name                = "acctestDES-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_key_id    = azurerm_key_vault_key.test.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "disk-encryption-perm" {
  key_vault_id = azurerm_key_vault.test.id

  tenant_id = azurerm_disk_encryption_set.test.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.test.identity.0.principal_id

  key_permissions = [
    "Get",
    "WrapKey",
    "UnwrapKey",
  ]
}

resource "azurerm_role_assignment" "disk-encryption-read-keyvault" {
  scope                = azurerm_key_vault.test.id
  role_definition_name = "Reader"
  principal_id         = azurerm_disk_encryption_set.test.identity.0.principal_id
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                   = "acctestaks%d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  dns_prefix             = "acctestaks%d"
  disk_encryption_set_id = azurerm_disk_encryption_set.test.id

  linux_profile {
    admin_username = "acctestuser%d"

    ssh_key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
  }

  identity {
    type = "SystemAssigned"
  }

  network {
    dns_service_ip = "10.10.0.10"
    service_cidr   = "10.10.0.0/16"
  }

  depends_on = [
    azurerm_key_vault_access_policy.disk-encryption-perm,
    azurerm_role_assignment.disk-encryption-read-keyvault
  ]

}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) upgradeChannelConfig(data acceptance.TestData, controlPlaneVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  kubernetes_version  = %q

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, controlPlaneVersion)
}

func (KubernetesAutomaticClusterResource) privateClusterPublicFqdn(data acceptance.TestData, privateClusterPublicFqdnEnabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}
resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"

  private_cluster {
    public_fully_qualified_domain_name_enabled = %t
  }

  identity {
    type = "SystemAssigned"
  }
  network {
    load_balancer_sku = "standard"
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, privateClusterPublicFqdnEnabled)
}

func (KubernetesAutomaticClusterResource) microsoftDefender(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}
resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  identity {
    type = "SystemAssigned"
  }
  microsoft_defender {
    log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) microsoftDefenderDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}
resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}
resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%d"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) webAppRoutingWithMultipleDnsZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[2]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dns_zone" "test2" {
  name                = "acctestzone2%[2]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }

  web_app_routing {
    dns_zone_ids = [azurerm_dns_zone.test.id, azurerm_dns_zone.test2.id]
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) webAppRoutingWitEmptyDnsZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }

  web_app_routing {
    dns_zone_ids = []
  }
}
 `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) webAppRoutingWithNginxController(data acceptance.TestData, controllerType string) string {
	defaultNginxController := ""
	if controllerType != "" {
		defaultNginxController = fmt.Sprintf(`default_nginx_controller = "%s"`, controllerType)
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }

  web_app_routing {
    dns_zone_ids = []
    %s
  }
}
 `, data.Locations.Primary, data.RandomInteger, defaultNginxController)
}

func (KubernetesAutomaticClusterResource) azureMonitorKubernetesMetricsEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }

  monitor_metrics {
    enabled = true
  }
}
  `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) azureMonitorKubernetesMetricsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }

  monitor_metrics {
    enabled             = true
    annotations_allowed = "pods=[k8s-annotation-1,k8s-annotation-n]"
    labels_allowed      = "namespaces=[k8s-label-1,k8s-label-n]"
  }
}
  `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) azureMonitorKubernetesMetricsDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }
}
  `, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) supportPlanKubernetesOfficial(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"

  identity {
    type = "SystemAssigned"
  }
  support_plan = "KubernetesOfficial"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (KubernetesAutomaticClusterResource) costAnalysisEnabled(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  identity {
    type = "SystemAssigned"
  }
  cost_analysis_enabled = %[3]t

}
`, data.Locations.Primary, data.RandomInteger, enabled)
}

func (KubernetesAutomaticClusterResource) customCATrustCertificates(data acceptance.TestData, certsList []string) string {
	certsString := ""
	if certsList != nil {
		certsString = "\"" + strings.Join(certsList, "\" ,\"") + "\""
	}
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[2]d"
  location = "%[1]s"
}
resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[2]d"
  identity {
    type = "SystemAssigned"
  }
  custom_ca_trust_certificates_base64 = [%[3]s]
}
`, data.Locations.Primary, data.RandomInteger, certsString)
}

func (KubernetesAutomaticClusterResource) aiToolchainOperatorProfile(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-aks-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_automatic_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"
  kubernetes_version  = "1.35.0"

  ai_toolchain_operator_enabled = %[3]t

  identity {
    type = "SystemAssigned"
  }
}
  `, data.RandomInteger, data.Locations.Primary, enabled)
}

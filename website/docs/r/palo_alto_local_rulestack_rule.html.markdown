---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rulestack_rule"
description: |-
  Manages a Palo Alto Local Rulestack Rule.
---

# azurerm_palo_alto_local_rulestack_rule

Manages a Palo Alto Local Rulestack Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "lrs-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_palo_alto_local_rulestack_rule" "example" {
  name         = "example-rule"
  rulestack_id = azurerm_palo_alto_local_rulestack.example.id
  priority     = 1000
  action       = "Allow"
  protocol     = "application-default"

  applications = ["any"]

  source {
    cidrs = ["10.0.0.0/8"]
  }

  destination {
    cidrs = ["192.168.16.0/24"]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `applications` - (Required) Specifies a list of Applications.

* `rulestack_id` - (Required) The ID of the Local Rulestack in which to create this Rule. Changing this forces a new Palo Alto Local Rulestack Rule to be created.

* `priority` - (Required) The Priority of this rule. Rules are executed in numerical order. Changing this forces a new Palo Alto Local Rulestack Rule to be created.

~> **Note:** This is the primary identifier of a rule, as such it is not possible to change the Priority of a rule once created.

* `action` - (Required) The action to take on the rule being triggered. Possible values are `Allow`, `DenyResetBoth`, `DenyResetServer` and `DenySilent`.

* `name` - (Required) The name which should be used for this Palo Alto Local Rulestack Rule. 

* `destination` - (Required) One or more `destination` blocks as defined below.

* `source` - (Required) One or more `source` blocks as defined below.

---


* `audit_comment` - (Optional) The comment for Audit purposes.

* `category` - (Optional) A `category` block as defined below.

* `decryption_rule_type` - (Optional) The type of Decryption to perform on the rule. Possible values include `SSLInboundInspection`, `SSLOutboundInspection`, and `None`. Defaults to `None`.

* `description` - (Optional) The description for the rule.

* `enabled` - (Optional) Should this Rule be enabled? Defaults to `true`.

* `inspection_certificate_id` - (Optional) The ID of the certificate for inbound inspection. Only valid when `decryption_rule_type` is set to `SSLInboundInspection`.

* `logging_enabled` - (Optional) Should Logging be enabled? Defaults to `false`.

* `negate_destination` - (Optional) Should the inverse of the Destination configuration be used. Defaults to `false`.

* `negate_source` - (Optional) Should the inverse of the Source configuration be used. Defaults to `false`.

* `protocol` - (Optional) The Protocol and port to use in the form `[protocol]:[port_number]` e.g. `TCP:8080` or `UDP:53`. Conflicts with `protocol_ports`. Defaults to `application-default`.

~> **Note:** In 4.0 or later versions, the default of `protocol` will no longer be set by provider, exactly one of `protocol` and `protocol_ports` must be specified. You need to explicitly specify `protocol="application-default"` to keep the the current default of the `protocol`.
 
* `protocol_ports` - (Optional) Specifies a list of Protocol:Port entries. E.g. `[ "TCP:80", "UDP:5431" ]`. Conflicts with `protocol`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Local Rulestack Rule.

---

A `category` block supports the following:

* `feeds` - (Optional) Specifies a list of feeds to match.

* `custom_urls` - (Required) Specifies a list of URL categories to match. Possible values include `abortion`, `abused-drugs`, `adult`, `alcohol-and-tobacco`, `auctions`, `business-and-economy`, `command-and-control`, `computer-and-internet-info`, `content-delivery-networks`, `copyright-infringement`, `cryptocurrency`, `dating`, `dynamic-dns`, `educational-institutions`, `entertainment-and-arts`, `extremism`, `financial-services`, `gambling`, `games`, `government`, `grayware`, `hacking`, `health-and-medicine`, `high-risk`, `home-and-garden`, `hunting-and-fishing`, `insufficient-content`, `internet-communications-and-telephony`, `internet-portals`, `job-search`, `legal`, `low-risk`, `malware`, `medium-risk`, `military`, `motor-vehicles`, `music`, `newly-registered-domain`, `news`, `not-resolved`, `nudity`, `online-storage-and-backup`, `parked`, `peer-to-peer`, `personal-sites-and-blogs`, `philosophy-and-political-advocacy`, `phishing`, `private-ip-addresses`, `proxy-avoidance-and-anonymizers`, `questionable`, `real-estate`, `real-time-detection`, `recreation-and-hobbies`, `reference-and-research`, `religion`, `search-engines`, `sex-education`, `shareware-and-freeware`, `shopping`, `social-networking`, `society`, `sports`, `stock-advice-and-tools`, `streaming-media`, `swimsuits-and-intimate-apparel`, `training-and-tools`, `translation`, `travel`, `unknown`, `weapons`, `web-advertisements`, `web-based-email`, and `web-hosting`. 

---

A `destination` block supports the following:

~> **Note:** At least one of the following properties must be specified.

* `cidrs` - (Optional) Specifies a list of CIDR's.

* `countries` - (Optional) Specifies a list of ISO3361-1 Alpha-2 Country codes. Possible values include `AF`, `AX`, `AL`, `DZ`, `AS`, `AD`, `AO`, `AI`, `AQ`, `AG`, `AR`, `AM`, `AW`, `AU`, `AT`, `AZ`, `BS`, `BH`, `BD`, `BB`, `BY`, `BE`, `BZ`, `BJ`, `BM`, `BT`, `BO`, `BQ`, `BA`, `BW`, `BV`, `BR`, `IO`, `BN`, `BG`, `BF`, `BI`, `KH`, `CM`, `CA`, `CV`, `KY`, `CF`, `TD`, `CL`, `CN`, `CX`, `CC`, `CO`, `KM`, `CG`, `CD`, `CK`, `CR`, `CI`, `HR`, `CU`, `CW`, `CY`, `CZ`, `DK`, `DJ`, `DM`, `DO`, `EC`, `EG`, `SV`, `GQ`, `ER`, `EE`, `ET`, `FK`, `FO`, `FJ`, `FI`, `FR`, `GF`, `PF`, `TF`, `GA`, `GM`, `GE`, `DE`, `GH`, `GI`, `GR`, `GL`, `GD`, `GP`, `GU`, `GT`, `GG`, `GN`, `GW`, `GY`, `HT`, `HM`, `VA`, `HN`, `HK`, `HU`, `IS`, `IN`, `ID`, `IR`, `IQ`, `IE`, `IM`, `IL`, `IT`, `JM`, `JP`, `JE`, `JO`, `KZ`, `KE`, `KI`, `KP`, `KR`, `KW`, `KG`, `LA`, `LV`, `LB`, `LS`, `LR`, `LY`, `LI`, `LT`, `LU`, `MO`, `MK`, `MG`, `MW`, `MY`, `MV`, `ML`, `MT`, `MH`, `MQ`, `MR`, `MU`, `YT`, `MX`, `FM`, `MD`, `MC`, `MN`, `ME`, `MS`, `MA`, `MZ`, `MM`, `NA`, `NR`, `NP`, `NL`, `NC`, `NZ`, `NI`, `NE`, `NG`, `NU`, `NF`, `MP`, `NO`, `OM`, `PK`, `PW`, `PS`, `PA`, `PG`, `PY`, `PE`, `PH`, `PN`, `PL`, `PT`, `PR`, `QA`, `RE`, `RO`, `RU`, `RW`, `BL`, `SH`, `KN`, `LC`, `MF`, `PM`, `VC`, `WS`, `SM`, `ST`, `SA`, `SN`, `RS`, `SC`, `SL`, `SG`, `SX`, `SK`, `SI`, `SB`, `SO`, `ZA`, `GS`, `SS`, `ES`, `LK`, `SD`, `SR`, `SJ`, `SZ`, `SE`, `CH`, `SY`, `TW`, `TJ`, `TZ`, `TH`, `TL`, `TG`, `TK`, `TO`, `TT`, `TN`, `TR`, `TM`, `TC`, `TV`, `UG`, `UA`, `AE`, `GB`, `US`, `UM`, `UY`, `UZ`, `VU`, `VE`, `VN`, `VG`, `VI`, `WF`, `EH`, `YE`, `ZM`, `ZW` 

* `feeds` - (Optional) Specifies a list of Feeds.

* `local_rulestack_fqdn_list_ids` - (Optional) Specifies a list of FQDN lists.

~> **Note:** This is a list of names of FQDN Lists configured on the same Local Rulestack as this Rule is being created.

* `local_rulestack_prefix_list_ids` - (Optional) Specifies a list of Prefix Lists.

~> **Note:** This is a list of names of Prefix Lists configured on the same Local Rulestack as this Rule is being created.

---

A `source` block supports the following:

~> **Note:** At least one of the following properties must be specified.

* `cidrs` - (Optional) Specifies a list of CIDRs.

* `countries` - (Optional) Specifies a list of ISO3361-1 Alpha-2 Country codes. Possible values include `AF`, `AX`, `AL`, `DZ`, `AS`, `AD`, `AO`, `AI`, `AQ`, `AG`, `AR`, `AM`, `AW`, `AU`, `AT`, `AZ`, `BS`, `BH`, `BD`, `BB`, `BY`, `BE`, `BZ`, `BJ`, `BM`, `BT`, `BO`, `BQ`, `BA`, `BW`, `BV`, `BR`, `IO`, `BN`, `BG`, `BF`, `BI`, `KH`, `CM`, `CA`, `CV`, `KY`, `CF`, `TD`, `CL`, `CN`, `CX`, `CC`, `CO`, `KM`, `CG`, `CD`, `CK`, `CR`, `CI`, `HR`, `CU`, `CW`, `CY`, `CZ`, `DK`, `DJ`, `DM`, `DO`, `EC`, `EG`, `SV`, `GQ`, `ER`, `EE`, `ET`, `FK`, `FO`, `FJ`, `FI`, `FR`, `GF`, `PF`, `TF`, `GA`, `GM`, `GE`, `DE`, `GH`, `GI`, `GR`, `GL`, `GD`, `GP`, `GU`, `GT`, `GG`, `GN`, `GW`, `GY`, `HT`, `HM`, `VA`, `HN`, `HK`, `HU`, `IS`, `IN`, `ID`, `IR`, `IQ`, `IE`, `IM`, `IL`, `IT`, `JM`, `JP`, `JE`, `JO`, `KZ`, `KE`, `KI`, `KP`, `KR`, `KW`, `KG`, `LA`, `LV`, `LB`, `LS`, `LR`, `LY`, `LI`, `LT`, `LU`, `MO`, `MK`, `MG`, `MW`, `MY`, `MV`, `ML`, `MT`, `MH`, `MQ`, `MR`, `MU`, `YT`, `MX`, `FM`, `MD`, `MC`, `MN`, `ME`, `MS`, `MA`, `MZ`, `MM`, `NA`, `NR`, `NP`, `NL`, `NC`, `NZ`, `NI`, `NE`, `NG`, `NU`, `NF`, `MP`, `NO`, `OM`, `PK`, `PW`, `PS`, `PA`, `PG`, `PY`, `PE`, `PH`, `PN`, `PL`, `PT`, `PR`, `QA`, `RE`, `RO`, `RU`, `RW`, `BL`, `SH`, `KN`, `LC`, `MF`, `PM`, `VC`, `WS`, `SM`, `ST`, `SA`, `SN`, `RS`, `SC`, `SL`, `SG`, `SX`, `SK`, `SI`, `SB`, `SO`, `ZA`, `GS`, `SS`, `ES`, `LK`, `SD`, `SR`, `SJ`, `SZ`, `SE`, `CH`, `SY`, `TW`, `TJ`, `TZ`, `TH`, `TL`, `TG`, `TK`, `TO`, `TT`, `TN`, `TR`, `TM`, `TC`, `TV`, `UG`, `UA`, `AE`, `GB`, `US`, `UM`, `UY`, `UZ`, `VU`, `VE`, `VN`, `VG`, `VI`, `WF`, `EH`, `YE`, `ZM`, `ZW`

* `feeds` - (Optional) Specifies a list of Feeds.

* `local_rulestack_prefix_list_ids` - (Optional) Specifies a list of Prefix Lists.

~> **Note:** This is a list of names of Prefix Lists configured on the same Local Rulestack as this Rule is being created.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Local Rulestack Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Local Rulestack Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Local Rulestack Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Local Rulestack Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Local Rulestack Rule.

## Import

Palo Alto Local Rulestack Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rulestack_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack/localRules/myRule1
```

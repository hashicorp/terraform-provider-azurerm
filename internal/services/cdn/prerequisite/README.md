# CDN Front Door acceptance test prerequisite: public DNS parent zone

These CDN Front Door custom-domain acceptance tests require a **real, publicly delegated** DNS parent zone.

Terraform can create the Azure DNS zone, but you must still configure your registrar to delegate the subdomain to the Azure DNS name servers.

This folder contains a tiny Terraform config you can apply once to create that Azure DNS zone and print the Azure name servers you need to paste into your registrar.

## Defaults (override as needed)

This module includes sensible defaults so you only need to supply the DNS zone name:

- `resource_group_name`: `acctest-afdx-dns-prereq`
- `location`: `westeurope`
- `tags`: `{}`

You can override these either via `-var ...` flags, or via the standard `TF_VAR_*` environment variables.

## What this creates

- A resource group (optional; configurable)
- An **Azure DNS Public Zone** for a subdomain like `acctest.example.com`
- Outputs of `name_servers` (the Azure DNS NS values)

## Usage

From this folder:

1) Set variables (example for `acctest.example.com`):

**Note:** `dns_zone_name` must be a real DNS zone name with two or more labels (e.g., `acctest.example.com`).

```powershell
$env:ARM_SUBSCRIPTION_ID = "..."
$env:ARM_TENANT_ID = "..."
$env:ARM_CLIENT_ID = "..."
$env:ARM_CLIENT_SECRET = "..."

terraform init
terraform apply -auto-approve `
  -var "dns_zone_name=acctest.example.com"

# Optional: override defaults if needed
# terraform apply -auto-approve `
#   -var "dns_zone_name=acctest.example.com" `
#   -var "resource_group_name=acctest-afdx-dns-prereq" `
#   -var "location=westeurope"
```

```bash
export ARM_SUBSCRIPTION_ID="..."
export ARM_TENANT_ID="..."
export ARM_CLIENT_ID="..."
export ARM_CLIENT_SECRET="..."

terraform init
terraform apply -auto-approve \
  -var "dns_zone_name=acctest.example.com"

# Optional: override defaults if needed
# terraform apply -auto-approve \
#   -var "dns_zone_name=acctest.example.com" \
#   -var "resource_group_name=acctest-afdx-dns-prereq" \
#   -var "location=westeurope"
```

2) Copy the output `name_servers`.

3) In your registrar (for `example.com`), add **NS** records:

- Type: `NS`
- Host: `acctest`
- Value: each Azure name server from the output (add all 4)

Do **not** change the domain-level nameservers for the whole domain.

You can verify delegation once it propagates:

```powershell
nslookup -type=ns acctest.example.com
```

```bash
dig NS acctest.example.com +short
```

4) After delegation propagates, set acceptance test env vars:

```powershell
$env:ARM_TEST_DNS_ZONE_NAME = "acctest.example.com"
$env:ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME = "acctest-afdx-dns-prereq"
```

```bash
export ARM_TEST_DNS_ZONE_NAME="acctest.example.com"
export ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME="acctest-afdx-dns-prereq"
```

## Notes

- DNS propagation can take minutes to hours.
- If you already have a delegated parent zone you can use, you do not need to apply this step.

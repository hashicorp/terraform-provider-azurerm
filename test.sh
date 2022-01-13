#!/bin/bash -e

export ARM_ENVIRONMENT="public"
export ARM_TENANT_ID="95cdfee8-f785-457d-a5ac-17346563394a"
export ARM_SUBSCRIPTION_ID="3f164b83-a88c-4170-8745-dfbc6df07a79"
export ARM_CLIENT_SECRET='yaU7Q~1egf9~rRSIDUiQvla5DHeThy_KMMDhX'
export ARM_CLIENT_ID="dc9509e0-c401-4661-bfe1-f9e236bbcc73"
export ARM_TEST_LOCATION="westeurope"
export ARM_TEST_LOCATION_ALT="northeurope"
export ARM_TEST_LOCATION_ALT2="uksouth"

make acctests SERVICE='storage' TESTARGS='-run=TestAccStorageAccount' TESTTIMEOUT='60m'


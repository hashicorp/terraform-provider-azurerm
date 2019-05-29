---
layout: "azurerm"
page_title: "Terraform Registry Contribution Guide"
sidebar_current: "docs-azurerm-guide-registry-contribute"
description: |-
    Terraform Registry: Contribution Guide

---
# Contribute to Terraform Modules Registry
This document describes the basic process of authoring and contributing re-usable, verifiable Terraform templates to Terraform Registry (https://registry.terraform.io). Terraform Registry serves as a centralized location to discover, learn and use Terraform artifacts to provision infrastructure into Azure, as well as into other cloud platforms.
## Basic Steps
Below are the basic steps to create and publish a verified Azure module into the Terraform Registry:

1. Following best practices for Terraform template authoring, create a template deploying infrastructure components into Azure.
2. Create a set of tests verifying correctness of Terraform template created in Step 1.
3. Set up a basic Azure DevOps pipeline for Continuous Integration and Deployment, including running the tests created in Step 2 above.
4. Using your GitHub credentials, sign into Terraform Registry and publish your module.

While Steps 2 and 3 are not a requirement for Terraform Registry, we strongly encourage all modules published as "Azure Verified" to include those steps.

### Step 1 - Creating Terraform template 
While Terraform offers an easy-to-create, easy-to-understand template authoring language in the form HCL, creating components that could be re-usable is often a time-consuming art form. With Terraform Registry, you are forced to think about various customer scenarios a template could be deployed with, and the permutation of various use cases could be daunting. While it is impossible to anticipate all the scenarios, we recommend that you familiarize yourself with a [set of Terraform Best Practices] (https://www.terraform.io/docs/enterprise/guides/recommended-practices/index.html) published by HashiCorp. Then, think through the most common use cases your module would likely serve, and try to author it in such a way that it will accommodate a set of those scenarios (as opposed to being strongly coupled to a single use case).

We have [created and documented a tool](https://docs.microsoft.com/en-us/azure/terraform/terraform-vscode-module-generator) that builds basic scaffolding for an Azure module, including a sample test. You are also welcome to clone one of the [existing modules] (https://github.com/Azure/terraform-azurerm-postgresql), providing you a set of Terraform files that you can customize for your module. 

### Step 2 - Create a set of tests 
After experimenting with several approaches for testing Terraform modules, we have settled on using [Terratest](https://github.com/gruntwork-io/terratest). The benefits of Terratest are its ease of integration with CI/CD tooling, as well its similarity to native [Terraform acceptance testing](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#writing-an-acceptance-test) framework. You will, however, need to be familiar with Go to create tests using Terratest.

If you used [Module Generator](https://docs.microsoft.com/en-us/azure/terraform/terraform-vscode-module-generator) in Step 1, a default test using Terratest is created for you in the test folder. You will need to think through additional tests that you can create that verify correctness of your template.

### Step 3 - Integrate with Azure DevOps 
As you introduce changes to the module, or approve pull requests, you'd like to be certain that new changes did not break the functionality the users have become dependent on. This is where CI/CD integration comes in. Integration with Azure DevOps is straightforward to setup - you will need to make sure you have an [azure-pipelines.yml](https://github.com/Azure/terraform-azurerm-postgresql/blob/master/azure-pipelines.yml) in your repo, as well as add the [test.sh](https://github.com/Azure/terraform-azurerm-postgresql/blob/master/test.sh) file to the root of the repo. 

The pipelines file is very generic - it builds a docker container and then calls into test.sh for Terratest execution. The only thing you will change in that file is the imageName variable - this will be the docker image name that the pipeline will create. The test.sh file is generic and stays as is, calling all of the test methods that were authored in Step 2.

After your pipelines succeed, please send a request to Terraform on Azure distribution list to pull this pipeline into Terraform ADO environment.

### Step 4 - Publish Module to Terraform Registry
This is the easiest step of all. Simply sign into [Terraform Registry](https://registry.terraform.io) using your GitHub credentials, then click on Publish in the top right corner. Select the repo to publish, then click publish to complete the process.

After your module has been published, please send a request to Terraform on Azure distribution list to enusre we add an "Azure Verified" check box next to the module you've created.

Summarize this Terraform AzureRM provider repository for me. Please provide an overview that covers:

## Repository Overview
- What this Terraform provider does and its primary purpose
- Key Azure services and resources it supports
- The provider's position in the HashiCorp ecosystem

## Architecture & Structure
- Main directory structure and organization
- How resources are organized by Azure service categories
- Key components: resources, data sources, client configurations
- Relationship between internal packages and their purposes

## Resource Implementation Patterns
- How Terraform resources are structured and implemented
- Common patterns for Azure resource lifecycle management (CRUD operations)
- Schema definitions and validation approaches
- State management and drift detection strategies

## Azure Integration
- How the provider integrates with Azure APIs
- Azure SDK usage and client management
- Authentication and authorization mechanisms
- API versioning and compatibility handling

## Development Workflow
- Build system and tooling (Make targets, scripts)
- Testing framework and acceptance tests
- Code generation and automation tools
- Documentation and example management

## Key Dependencies
- Major Go dependencies and their roles
- Azure SDK components used
- HashiCorp Plugin SDK integration
- Testing and development tools

Please assume the reader is new to both Terraform provider development and this specific repository. Focus on providing a clear mental model of how this large codebase is organized and how different components work together to enable Azure resource management through Terraform.
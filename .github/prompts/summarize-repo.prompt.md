Summarize this Terraform AzureRM provider repository for me. Please provide a comprehensive technical overview that covers:

## Repository Overview
- What this Terraform provider does and its primary purpose in the Infrastructure as Code ecosystem
- Comprehensive overview of Azure services and resources it supports (with specific examples from major categories)
- The provider's position in the HashiCorp ecosystem and relationship to other providers
- Scale and scope: number of resources, data sources, and maintainer information
- Target audience and use cases (from individual developers to enterprise deployments)

## Architecture & Structure
- Detailed main directory structure and organization philosophy
- How resources are organized by Azure service categories with specific examples
- Key components breakdown:
  - Resources vs data sources: implementation differences and use cases
  - Client configurations and Azure SDK integration patterns
  - Provider registration and initialization flow
- Relationship between internal packages and their specific purposes:
  - `/internal/services/*` - service-specific implementations
  - `/internal/clients/` - Azure API client management
  - `/internal/sdk/` - modern implementation framework
  - `/internal/acceptance/` - testing infrastructure
  - `/internal/common/` - shared utilities and helpers
- Code organization patterns and conventions used throughout the codebase

## Resource Implementation Patterns
- Modern vs Legacy implementation approaches:
  - Modern SDK-based implementation (preferred for new resources)
  - Legacy Plugin SDK implementation (maintenance mode)
  - Migration strategies and compatibility considerations
- Detailed breakdown of Terraform resource lifecycle:
  - Create: Resource provisioning and initial configuration
  - Read: State refresh and drift detection mechanisms
  - Update: In-place updates vs resource recreation (ForceNew)
  - Delete: Resource cleanup and dependency management
  - Import: Existing resource adoption into Terraform state
- Schema definitions and validation approaches:
  - Type system usage (TypeString, TypeBool, TypeList, TypeSet, TypeMap)
  - Validation functions and custom validators
  - Required vs Optional vs Computed fields
  - Nested block structures and complex data types
- State management strategies:
  - State file structure and resource tracking
  - Drift detection and reconciliation
  - Handling Azure resource changes outside Terraform
- Error handling patterns and user experience considerations

## Azure Integration Deep Dive
- Azure API integration architecture:
  - REST API consumption patterns
  - Azure Resource Manager (ARM) API usage
  - Service-specific API endpoints and versions
- Azure SDK usage and evolution:
  - HashiCorp Go Azure SDK (preferred)
  - Official Azure SDK for Go (selective usage)
  - Custom wrapper implementations and abstractions
- Client management and lifecycle:
  - Authentication flow and credential management
  - Connection pooling and retry mechanisms
  - Rate limiting and throttling handling
- Authentication and authorization mechanisms:
  - Service Principal authentication
  - Managed Identity integration
  - Azure CLI credential chain
  - Certificate-based authentication
  - Multi-tenant scenarios
- API versioning and compatibility:
  - Azure API version management
  - Preview API integration strategies
  - Backward compatibility maintenance
  - Feature flag management for new capabilities

## Development Workflow & Tooling
- Build system comprehensive breakdown:
  - GNUmakefile targets and their purposes
  - Cross-platform build considerations
  - Release automation and versioning
- Testing framework architecture:
  - Unit testing patterns and coverage
  - Acceptance testing against live Azure APIs
  - Test data management and cleanup strategies
  - Parallel test execution and resource isolation
  - CI/CD pipeline integration
- Code quality and maintenance tools:
  - Linting rules and code formatting (golangci-lint, gofmt)
  - Static analysis and security scanning
  - Dependency management and updates
- Documentation generation and validation:
  - Resource documentation standards
  - Example validation and testing
  - Website generation and deployment
- Development environment setup:
  - Required tools and dependencies
  - Local testing setup and Azure subscription requirements
  - Debugging and troubleshooting workflows

## Implementation Best Practices & Patterns
- Resource naming conventions and consistency
- Error handling and user messaging standards
- Azure-specific patterns:
  - Resource tagging implementation
  - Location/region handling
  - Resource group management
  - Subscription and tenant scoping
- Performance considerations:
  - API call optimization
  - Bulk operations where supported
  - Caching strategies
- Security implementation:
  - Sensitive data handling
  - Secret management integration
  - Principle of least privilege

## Key Dependencies & Ecosystem
- Core Go dependencies detailed analysis:
  - github.com/hashicorp/terraform-plugin-sdk/v2: Plugin framework
  - github.com/hashicorp/go-azure-sdk: Primary Azure integration
  - github.com/hashicorp/go-azure-helpers: Azure utilities
  - github.com/Azure/azure-sdk-for-go: Official SDK components
- Development and testing dependencies:
  - Testing frameworks and assertion libraries
  - Mock generation and test utilities
  - Build and automation tools
- Runtime dependencies and their roles in provider operation
- Version compatibility matrix and update strategies

## Community & Contribution Model
- Contribution guidelines and development process
- Issue triage and feature request handling
- Release cycle and versioning strategy
- Relationship with Azure product teams and feedback loops
- Community resources and support channels

Please assume the reader is new to both Terraform provider development and this specific repository. Focus on providing a clear mental model of how this large codebase is organized, how different components work together to enable Azure resource management through Terraform, and what makes this provider unique in the Terraform ecosystem. Include specific code examples, file paths, and architectural diagrams where helpful for understanding.

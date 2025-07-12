Summarize this Terraform AzureRM provider repository for me. Please provide a comprehensive technical overview that covers:

## Repository Overview
- Primary purpose in the Infrastructure as Code ecosystem
- Major Azure service categories supported (compute, network, storage, database, security, etc.)
- Position in HashiCorp ecosystem and scale (resource/data source counts)
- Target audience from individual developers to enterprise deployments

## Architecture & Structure
- Main directory organization philosophy and key directories:
  - `/internal/services/*` - service implementations
  - `/internal/clients/` - Azure client management  
  - `/internal/sdk/` - typed implementation framework
  - `/internal/acceptance/` - testing infrastructure
- Resource vs data source implementation patterns
- Provider registration and initialization flow
- Code organization conventions across the codebase

## Implementation Approaches
- **Typed vs Untyped Resource Patterns:**
  - Typed resource implementation (preferred): Uses internal/sdk framework
  - Untyped resource implementation (maintenance): Traditional Plugin SDK
  - Migration strategies and when to use each approach
- **Resource Lifecycle (CRUD + Import):**
  - Create, Read, Update, Delete operation patterns
  - State management and drift detection strategies
  - Import functionality for existing Azure resources
- **Schema Design:**
  - Type system usage and validation approaches
  - Complex nested structures and data handling
  - ForceNew vs in-place update patterns

## Azure Integration
- **API Integration Architecture:**
  - Azure Resource Manager (ARM) API consumption
  - HashiCorp Go Azure SDK (preferred) vs Official Azure SDK usage
  - Service-specific endpoint management
- **Authentication & Client Management:**
  - Service Principal, Managed Identity, Azure CLI credential chains
  - Connection pooling, retry mechanisms, rate limiting
  - Multi-tenant and subscription scoping
- **API Versioning:**
  - Azure API version management and preview integration
  - Backward compatibility and feature flag strategies

## Development Workflow
- **Build System:** GNUmakefile targets, cross-platform builds, release automation
- **Testing Framework:** Unit tests, acceptance tests against live APIs, parallel execution
- **Code Quality:** Linting, formatting, static analysis, dependency management
- **Documentation:** Resource docs generation, example validation, website deployment

## Key Dependencies & Ecosystem
- **Core Dependencies:**
  - terraform-plugin-sdk/v2: Plugin framework foundation
  - go-azure-sdk: Primary Azure integration layer
  - go-azure-helpers: Azure-specific utilities
  - azure-sdk-for-go: Official SDK for selective usage
- **Development Tools:** Testing, mocking, build automation
- **Community Model:** Contribution process, release cycles, Azure team collaboration

## Best Practices & Patterns
- Resource naming conventions and Azure-specific patterns
- Error handling standards and user experience design
- Performance optimization (API efficiency, caching, bulk operations)
- Security implementation (sensitive data, secrets, least privilege)

Please focus on providing a clear mental model of how this large codebase enables Azure resource management through Terraform, what makes it unique in the ecosystem, and how the different architectural components work together. Include specific examples and file paths where helpful for understanding the technical implementation.

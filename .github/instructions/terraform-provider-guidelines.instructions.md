---
applyTo: "internal/**/*.go"
description: This document outlines the Azure-specific guidelines for Go files in the Terraform Azure Provider repository. It includes best practices for Azure Resource Manager integration, Terraform provider patterns, and resource implementation.
---

## Azure Terraform Provider Guidelines
Given below are the Azure-specific guidelines for this Terraform Provider project which **MUST** be followed.

### Azure Resource Manager (ARM) Integration
- Use the HashiCorp Go Azure SDK as the primary SDK for Azure integrations
- Implement proper error handling for Azure API responses
- Use appropriate polling for long-running operations (LROs)
- Implement proper retry logic with exponential backoff
- Handle Azure API rate limits and throttling gracefully
- Use managed identity authentication when possible
- Always validate resource IDs using the proper parsing utilities

### Terraform Resource Implementation Patterns
- Follow the standard CRUD lifecycle: Create, Read, Update, Delete
- Use proper Terraform Plugin SDK v2 patterns and conventions
- Implement proper state management and drift detection
- Use `ForceNew` for properties that require resource recreation
- Implement proper timeout configurations for all operations
- Use appropriate validation functions for resource properties
- Handle nested resource configurations properly using TypeSet, TypeList, and TypeMap

### Azure-Specific Resource Patterns
- Use standardized resource naming patterns with resourceToken
- Implement proper location/region handling across all resources
- Follow Azure resource tagging conventions consistently
- Handle Azure API versioning correctly
- Implement proper subscription and resource group scoping
- Use Azure resource IDs consistently across all resources
- Handle cascading deletes and dependencies properly

### Security and Authentication
- Never hardcode sensitive values in code, tests, or examples
- Use Azure Key Vault references where appropriate
- Implement proper authentication handling with service principals
- Validate all inputs to prevent injection attacks
- Follow principle of least privilege for service principals
- Use context with appropriate timeouts for all Azure API calls

### Testing Requirements
- Write comprehensive acceptance tests for all resources
- Use the standard acceptance test framework in `/internal/acceptance`
- Test both success and failure scenarios thoroughly
- Ensure tests are idempotent and can run in parallel
- Mock external dependencies appropriately
- Test resource import functionality for all resources
- Use environment variables for test configuration (avoid hardcoding)

### Error Handling and Logging
- Use structured logging with appropriate log levels
- Include resource IDs and operation context in all logs
- Implement proper error wrapping and context propagation
- Use Terraform's diagnostic system for user-facing errors
- Avoid logging sensitive information (passwords, keys, tokens)
- Handle Azure API errors gracefully with meaningful messages
- Field names in error messages should be wrapped in backticks for clarity
- Field values in error messages should be wrapped in backticks for clarity
- Error messages must follow Go standards (lowercase, no punctuation, descriptive)
- Do not use contractions in error messages. Always use the full form of words. For example, write 'cannot' instead of 'can't' and 'is not' instead of 'isn't'
- Error messages must use '%w' instead of '%+v' for proper Go error wrapping standards
- Error messages must be clear, concise, and provide actionable guidance

### Performance Considerations
- Implement efficient resource queries to minimize API calls
- Use bulk operations where supported by Azure APIs
- Implement proper caching where appropriate
- Monitor and optimize API call patterns
- Use context with appropriate timeouts for all operations
- Consider pagination for large result sets

# TRex Core Library

[![Go Reference](https://pkg.go.dev/badge/github.com/openshift-online/rh-trex-core.svg)](https://pkg.go.dev/github.com/openshift-online/rh-trex-core)
[![Go Report Card](https://goreportcard.com/badge/github.com/openshift-online/rh-trex-core)](https://goreportcard.com/report/github.com/openshift-online/rh-trex-core)

The TRex Core Library provides reusable patterns and frameworks for building REST API microservices in Go. This library extracts the proven patterns from the TRex template project into a reusable package that can be imported by new microservices.

## Features

- **Generic CRUD Services**: Type-safe CRUD operations using Go generics
- **Event-Driven Controllers**: Automatic event handling with PostgreSQL LISTEN/NOTIFY
- **Generic DAOs**: Database access patterns with GORM integration
- **Project Templates**: Code generation for new microservices
- **Standardized APIs**: Common patterns for REST API development
- **Built-in Testing**: Comprehensive test utilities and mocks

## Installation

```bash
go get github.com/openshift-online/rh-trex-core
```

## Quick Start

### 1. Define Your Resource Type

```go
package main

import (
    "github.com/openshift-online/rh-trex-core/api"
)

type User struct {
    api.Meta
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"unique;not null"`
}
```

### 2. Create Service with Generic CRUD

```go
import (
    "github.com/openshift-online/rh-trex-core/dao"
    "github.com/openshift-online/rh-trex-core/services"
)

func NewUserService(db *gorm.DB, eventEmitter events.EventEmitter) services.CRUDService[User] {
    userDAO := dao.NewBaseDAO[User](db)
    return services.NewBaseCRUDService[User](userDAO, eventEmitter, "Users")
}
```

### 3. Auto-Register Controllers

```go
import (
    "github.com/openshift-online/rh-trex-core/controllers"
)

func main() {
    userService := NewUserService(db, eventEmitter)
    
    // Automatically handles CREATE/UPDATE/DELETE events
    controllers.AutoRegisterCRUDController(controllerManager, userService, "Users")
}
```

## Architecture

The core library provides these main components:

- **`api/`** - Base API types and metadata patterns
- **`services/`** - Generic CRUD service implementations
- **`controllers/`** - Event-driven controller framework
- **`dao/`** - Generic data access objects with GORM
- **`generator/`** - Code generation utilities
- **`template/`** - Project scaffolding tools

## Usage Patterns

### Generic Services

Create type-safe services for any resource:

```go
type Product struct {
    api.Meta
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

productDAO := dao.NewBaseDAO[Product](db)
productService := services.NewBaseCRUDService[Product](productDAO, eventEmitter, "Products")

// Automatically provides: Get, Create, Replace, Delete, List, FindByIDs
```

### Event-Driven Architecture

Services automatically emit events for database changes:

```go
// Implement custom event handlers
func (s *CustomProductService) OnUpsert(ctx context.Context, id string) error {
    // Custom business logic for CREATE/UPDATE events
    return nil
}

func (s *CustomProductService) OnDelete(ctx context.Context, id string) error {
    // Custom business logic for DELETE events
    return nil
}
```

### Project Generation

Use the template system to create new microservices:

```go
import "github.com/openshift-online/rh-trex-core/template"

config := template.ProjectConfig{
    Name:        "my-service",
    Module:      "github.com/myorg/my-service",
    Resources:   []template.ResourceConfig{
        {Name: "User", Fields: []template.FieldConfig{...}},
        {Name: "Product", Fields: []template.FieldConfig{...}},
    },
}

generator := template.NewProjectTemplate(config)
err := generator.Generate("./my-service")
```

## Benefits

- **Shared Evolution**: All projects get improvements automatically via `go get -u`
- **Smaller Codebases**: Projects contain only business logic, not infrastructure
- **Consistency**: All projects use the same proven patterns
- **Better Testing**: Framework is tested once, thoroughly
- **Faster Development**: Focus on business logic, not boilerplate

## Migration from TRex Template

If you have an existing project based on the TRex template, see our [Migration Guide](MIGRATION.md) for step-by-step instructions.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `go test ./...`
6. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [rh-trex](https://github.com/openshift-online/rh-trex) - The original TRex template project
- [Examples](https://github.com/openshift-online/rh-trex-core-examples) - Example projects using the core library
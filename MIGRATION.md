# Migration Guide: TRex Template to Core Library

This guide helps you migrate an existing project based on the TRex template to use the rh-trex-core library.

## Overview

**Before Migration:**
- Project contains 200+ files copied from TRex
- Full maintenance burden for infrastructure code
- Manual updates required when TRex improves

**After Migration:**
- ~20 business logic files
- Infrastructure provided by `github.com/openshift-online/rh-trex-core`
- Automatic updates via `go get -u`

## Step-by-Step Migration

### Phase 1: Add Core Library Dependency

1. **Update go.mod**:
```bash
go mod edit -require github.com/openshift-online/rh-trex-core@latest
go mod tidy
```

2. **Update imports** in your service files:
```go
// Before
import "yourproject/pkg/dao"
import "yourproject/pkg/services"

// After  
import "github.com/openshift-online/rh-trex-core/dao"
import "github.com/openshift-online/rh-trex-core/services"
```

### Phase 2: Migrate Resource Types

For each of your business domain types (keep your custom types, remove template dinosaurs):

1. **Update type definitions** to use core API base:
```go
// Before
type YourResource struct {
    ID        string    `json:"id" gorm:"primary_key"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // ... other fields
}

// After
import "github.com/openshift-online/rh-trex-core/api"

type YourResource struct {
    api.Meta
    // ... your business fields only
}
```

2. **Update services** to use generic CRUD:
```go
// Before: Custom service implementation
type yourResourceService struct {
    dao dao.YourResourceDAO
}

// After: Use generic service
import "github.com/openshift-online/rh-trex-core/services"

func NewYourResourceService(db *gorm.DB, eventEmitter events.EventEmitter) services.CRUDService[YourResource] {
    dao := dao.NewBaseDAO[YourResource](db)
    return services.NewBaseCRUDService[YourResource](dao, eventEmitter, "YourResources")
}
```

### Phase 3: Remove Infrastructure Files

Delete these copied TRex infrastructure files (keep only your business logic):

```bash
# Remove copied infrastructure
rm -rf pkg/auth/
rm -rf pkg/config/
rm -rf pkg/db/
rm -rf pkg/errors/
rm -rf pkg/logger/
rm -rf pkg/controllers/framework.go
rm -rf pkg/handlers/framework.go
rm -rf pkg/handlers/errors.go
rm -rf pkg/handlers/helpers.go
rm -rf pkg/handlers/rest.go
rm -rf pkg/handlers/validation.go
rm -rf pkg/services/generic.go
rm -rf pkg/dao/generic.go

# Keep only your business domain files
# pkg/api/yourresource.go
# pkg/services/yourresource.go  
# pkg/dao/yourresource.go
# pkg/handlers/yourresource.go
```

### Phase 4: Update Main Application

Update your main application to use core library patterns:

```go
// cmd/yourapp/main.go
package main

import (
    "github.com/openshift-online/rh-trex-core/generator"
    "github.com/openshift-online/rh-trex-core/controllers"
    "yourproject/pkg/api"
)

func main() {
    // Use core library for framework setup
    factory := generator.NewResourceFactory(db, controllerMgr, eventEmitter)
    
    // Register your business resources
    userService := generator.RegisterResourceType[api.User](factory, "Users")
    productService := generator.RegisterResourceType[api.Product](factory, "Products")
    
    // Core library handles all CRUD operations and event processing
    server.Start()
}
```

### Phase 5: Update Environment Configuration

Simplify your environment configuration:

```go
// Before: Complex environment framework
type Services struct {
    Users     UserServiceLocator
    Products  ProductServiceLocator
    Generic   GenericServiceLocator
    Events    EventServiceLocator
}

// After: Use core library patterns
import "github.com/openshift-online/rh-trex-core/services"

type Services struct {
    Users    services.CRUDService[api.User]
    Products services.CRUDService[api.Product] 
}
```

## Example: Complete Migration

Here's a complete before/after example:

### Before (200+ files)
```
yourproject/
├── cmd/yourapp/
│   ├── main.go
│   ├── environments/     # 10+ files
│   ├── server/          # 15+ files
│   └── migrate/         # 5+ files
├── pkg/
│   ├── api/            # 20+ files
│   ├── auth/           # 10+ files  
│   ├── config/         # 10+ files
│   ├── controllers/    # 5+ files
│   ├── dao/            # 15+ files
│   ├── db/             # 15+ files
│   ├── errors/         # 5+ files
│   ├── handlers/       # 15+ files
│   ├── logger/         # 5+ files
│   └── services/       # 15+ files
└── test/               # 25+ files
```

### After (~20 files)
```
yourproject/
├── cmd/yourapp/
│   ├── main.go          # Uses core library
│   └── migrate/         # Simple migration script
├── pkg/
│   ├── api/
│   │   ├── user.go      # Your business type
│   │   └── product.go   # Your business type
│   ├── services/
│   │   ├── user.go      # Your business logic
│   │   └── product.go   # Your business logic
│   └── handlers/
│       ├── user.go      # Your API handlers
│       └── product.go   # Your API handlers
└── test/
    ├── user_test.go     # Your business logic tests
    └── product_test.go  # Your business logic tests
```

## Testing Migration

After migration, verify everything works:

```bash
# Test compilation
go build ./...

# Run tests
go test ./...

# Test integration
make test-integration

# Test API endpoints
curl localhost:8000/api/yourapp/v1/users
```

## Benefits Realized

After migration, your project will have:

- **90% fewer files** - Focus on business logic only
- **Automatic updates** - `go get -u github.com/openshift-online/rh-trex-core`
- **Shared improvements** - All projects benefit from core library enhancements
- **Consistent patterns** - Standard approach across all microservices
- **Better testing** - Core library is thoroughly tested

## Troubleshooting

### Common Issues

1. **Import path errors**: Update all imports to use core library paths
2. **Type compatibility**: Ensure your types embed `api.Meta`
3. **Service interfaces**: Use generic `CRUDService[T]` interface
4. **Event handling**: Implement `OnUpsert`/`OnDelete` if needed

### Getting Help

- Check the [examples repository](https://github.com/openshift-online/rh-trex-core-examples)
- Review the [API documentation](https://pkg.go.dev/github.com/openshift-online/rh-trex-core)
- Open an issue on the [core library repository](https://github.com/openshift-online/rh-trex-core/issues)

## Next Steps

Once migration is complete:
1. Remove unused template code (like dinosaurs)
2. Implement your business-specific event handlers
3. Add custom validation logic
4. Deploy and test in your environment
5. Set up automatic updates for the core library
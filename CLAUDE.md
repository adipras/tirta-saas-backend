# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a multi-tenant water billing SaaS backend built with Go and Gin framework. The system manages water usage, billing, and payments for different tenants (water utility companies).

## Development Commands

### Running the Application
```bash
go run main.go
```
The server runs on port 8080 by default (configurable via PORT env var).

### Database Setup
- Requires MySQL database connection
- Database configuration via environment variables (see .env file)
- Auto-migration runs on startup via `config.Migrate()`

### Environment Variables
Required environment variables:
- `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_PORT`, `DB_NAME` - Database connection
- `JWT_SECRET` - JWT token signing secret
- `PORT` - Server port (optional, defaults to 8080)

## Architecture

### Multi-tenant Architecture
- **Tenant Isolation**: Every model includes `tenant_id` for data isolation
- **JWT Claims**: Include `tenant_id`, `user_id`, and `role` for authorization
- **Middleware**: Automatic tenant context extraction from JWT tokens

### Database Models Structure
All models extend `BaseModel` which provides:
- UUID primary keys (`char(36)`)
- Timestamps (`created_at`, `updated_at`)
- Soft deletes (`deleted_at`)
- Auto-UUID generation via `BeforeCreate` hook

Key entity relationships:
- `Tenant` → `User` (admin users per tenant)
- `Customer` → `SubscriptionType` (water service plans)
- `Customer` → `WaterUsage` → `Invoice` → `Payment` (billing flow)

### Controller Patterns
Controllers follow consistent patterns:
- Tenant ID extraction from JWT context: `c.MustGet("tenant_id").(uuid.UUID)`
- Request validation via struct binding with `binding` tags
- Database queries always filter by `tenant_id` for isolation
- Transaction usage for multi-table operations (e.g., customer registration)

### Request/Response Structure
- **Requests**: Defined in `requests/` package with validation tags
- **Responses**: Defined in `responses/` package for consistent API output
- **Validation**: Uses Gin's binding with struct tags (`required`, `email`, etc.)

### Authentication & Authorization
- **JWT Authentication**: `middleware.JWTAuthMiddleware()` extracts user context
- **Role-based Access**: `middleware.AdminOnly()` restricts admin endpoints
- **Multi-tenant Security**: All queries automatically filtered by tenant_id

### API Routing Structure
- `/auth` - Registration and login (no auth required)
- `/api/*` - Protected endpoints requiring JWT + admin role
- Swagger documentation available at `/swagger/*any`

Routes are organized by domain:
- `/api/customers` - Customer management
- `/api/invoices` - Invoice generation and management
- `/api/payments` - Payment processing
- `/api/water-usage` - Usage tracking
- `/api/water-rates` - Tariff management

## Business Logic

### Customer Registration Flow
1. Create customer with subscription type
2. Auto-generate registration invoice with subscription fee
3. Customer remains inactive until registration fee is paid

### Billing Process
1. Record water usage (`WaterUsage`)
2. Generate monthly invoice based on usage and tariff
3. Process payments against invoices
4. Track payment status and remaining balances

## Database Migration
Database schema is managed via GORM's AutoMigrate feature in `config.Migrate()`. Migration runs automatically on application startup and includes all models in dependency order.
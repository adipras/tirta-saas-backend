# Tirta SaaS Backend

## 🌊 Multi-Tenant Water Utility Billing System

A production-ready, enterprise-grade SaaS backend for water utility companies (Paguyuban Air Bersih) built with Go and Gin framework. This system provides complete billing lifecycle management from customer registration to payment processing with robust multi-tenancy support.

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Technology Stack](#️-technology-stack)
- [Project Status](#-project-status)
- [Getting Started](#-getting-started)
- [Development Roadmap](#-development-roadmap)
- [Documentation](#-documentation)
- [API Documentation](#-api-documentation)
- [Contributing](#-contributing)
- [License](#-license)

---

## 🎯 Overview

**Tirta SaaS Backend** is a B2B2C multi-tenant platform designed specifically for water utility cooperatives (Paguyuban Air Bersih) in Indonesia. Each tenant operates independently with complete data isolation, managing their own customers, billing, and payment operations.

### Key Capabilities

- **Multi-tenant Architecture**: Unlimited tenants with complete data isolation
- **Complete Billing Lifecycle**: Registration → Usage Tracking → Invoice Generation → Payment Processing
- **Flexible Pricing Models**: Multiple subscription types with complex fee structures
- **Customer Self-Service**: Portal for customers to view bills, make payments, and track usage
- **Enterprise Security**: JWT authentication, role-based access control, rate limiting
- **Production-Grade Performance**: Optimized queries, connection pooling, comprehensive monitoring

---

## ✨ Features

### 🏢 Multi-Tenant Management
- Tenant registration with admin user creation
- Complete tenant data isolation with UUID identification
- Tenant-specific security policies and rate limiting

### 🔐 Authentication & Authorization
- **Admin Authentication**: JWT-based with role-based access control
- **Customer Authentication**: Email/password with activation workflow
- **Secure Password Management**: bcrypt hashing (cost 14)
- **Multi-layer Rate Limiting**: IP-based and role-based protection

### 👥 Customer Management
- Admin-managed customer registration and profile management
- Customer activation workflow via registration payment
- Customer self-service portal with full account access
- Password management and profile updates

### 📊 Subscription Management
- Complete CRUD operations for subscription types
- Complex fee structures (registration, monthly, maintenance, late fees)
- Tenant-specific subscription types with full isolation

### 💧 Water Usage & Rate Management
- Monthly meter reading entry with validation
- Automatic usage calculation with continuity checks
- Dynamic water rate management with effective dates
- Subscription-type specific pricing models
- Rate versioning and historical tracking

### 🧾 Invoice Management
- Bulk monthly invoice generation with duplicate prevention
- Multi-component billing (usage + subscription + maintenance)
- Registration vs monthly invoice types
- Real-time payment status tracking
- Customer self-service invoice viewing

### 💳 Payment Processing
- Complete payment lifecycle management
- Partial and full payment support
- Automatic customer activation on registration payment
- Payment history tracking and audit trails
- Overpayment prevention with business rule validation

### 🛡️ Enterprise Security
- Multi-layer rate limiting (global, endpoint-specific, authentication)
- CORS protection with configurable policies
- Security headers (XSS, clickjacking, MIME sniffing protection)
- SQL injection protection via GORM
- Input validation and sanitization
- Audit trails for all critical operations

### 📊 Monitoring & Observability
- Structured JSON logging with request tracking
- Health check endpoints (liveness, readiness, database)
- Metrics collection for monitoring
- Performance tracking and optimization
- Error tracking and alerting

---

## 🛠️ Technology Stack

- **Language**: Go 1.24.2
- **Web Framework**: Gin Web Framework
- **Database**: MySQL 8.0+ with GORM ORM
- **Authentication**: JWT tokens with bcrypt password hashing
- **Security**: Rate limiting, CORS, Security headers
- **Documentation**: Swagger/OpenAPI integration
- **Architecture**: Clean Architecture with separated concerns

### Production Features
- Connection pooling for optimal database performance
- 35+ performance indexes for query optimization
- Comprehensive middleware stack
- Transaction safety for data integrity
- Business rule validation throughout

---

## 🎉 Project Status

### **Current Version: 1.0 (Core System Complete)**

**Development Progress: 95% Complete**

| Phase | Status | Completion |
|-------|--------|------------|
| Phase 1-4 (Core System) | ✅ Complete | 100% |
| Phase 5 (Platform Management) | ⏳ Planned | 0% |
| Phase 6 (Tenant Admin Enhancements) | ⏳ Planned | 0% |
| Phase 7 (Operational Management) | ⏳ Planned | 0% |

### ✅ Production-Ready Features
- Multi-tenant water billing system (core operations)
- Enterprise security and authentication
- Real-time performance monitoring
- Complete API documentation
- Deployment guides and procedures

### 🚀 Ready for Deployment
The core system is **production-ready** and can be deployed immediately for:
- Customer registration and management
- Water usage tracking and billing
- Invoice generation and payment processing
- Customer self-service portal
- Basic operational reporting

---

## 🚀 Getting Started

### Prerequisites
- Go 1.24.2 or higher
- MySQL 8.0 or higher
- Git

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/tirta-saas-backend.git
cd tirta-saas-backend
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
Create a `.env` file in the root directory:
```bash
# Database Configuration
DB_USER=your_db_user
DB_PASS=your_db_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=tirta_saas

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Configuration
PORT=8080
GIN_MODE=release

# Rate Limiting (optional)
RATE_LIMIT_ENABLED=true
RATE_LIMIT_RPS=100
```

4. **Run database migrations**
The application will automatically run migrations on startup via GORM's AutoMigrate.

5. **Start the application**
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Quick Test

Test the health check endpoint:
```bash
curl http://localhost:8080/health
```

### API Documentation

Access Swagger UI documentation:
```
http://localhost:8080/swagger/index.html
```

---

## 📈 Development Roadmap

### Phase 5: Platform Management (4 weeks)
**Status**: Not Started

See [`PLATFORM_MANAGEMENT_TODO.md`](./PLATFORM_MANAGEMENT_TODO.md) for details.

**Key Features**:
- Platform owner dashboard for multi-tenant monitoring
- Tenant subscription plans (Basic, Premium, Enterprise)
- Platform analytics and revenue tracking
- System monitoring with audit logs
- Enhanced notification system (Email, SMS, In-app)
- Bulk operations and data import/export
- Advanced reporting and business intelligence

### Phase 6: Tenant Admin Enhancements (4 weeks)
**Status**: Not Started

See [`TENANT_ADMIN_TODO.md`](./TENANT_ADMIN_TODO.md) for details.

**Key Features**:
- Enhanced user and role management
- Master data settings (Paguyuban profile, tariff categories)
- Service area management (RT/RW/Blok zones)
- Automated meter reading schedules
- Invoice template customization
- Business rules engine
- Operational dashboards and analytics

### Phase 7: Operational Management (8 weeks)
**Status**: Not Started

See [`OPERATIONAL_MANAGEMENT_TODO.md`](./OPERATIONAL_MANAGEMENT_TODO.md) for details.

**Key Features**:
- Mobile-friendly meter reading interface with offline sync
- Automated invoice distribution (Email, WhatsApp)
- Advanced payment collection and reconciliation
- Service request management system
- Complaint and ticketing system
- Inventory and asset management
- Field staff mobile app with GPS tracking
- Customer self-service mobile app
- Operational analytics and KPIs
- External integrations (Payment gateways, SMS, Banking)
- IoT and smart meter integration

---

## 📚 Documentation

Comprehensive documentation is available in the following files:

### Core Documentation
- **[PROJECT_ANALYSIS.md](./PROJECT_ANALYSIS.md)** - Complete project analysis, architecture, and development status
- **[PRODUCTION_DEPLOYMENT_GUIDE.md](./PRODUCTION_DEPLOYMENT_GUIDE.md)** - Step-by-step deployment procedures
- **[POSTMAN_TESTING_GUIDE.md](./POSTMAN_TESTING_GUIDE.md)** - API testing guide with Postman collection
- **[CLAUDE.md](./CLAUDE.md)** - Development guidance and code patterns

### Feature Roadmap
- **[PLATFORM_MANAGEMENT_TODO.md](./PLATFORM_MANAGEMENT_TODO.md)** - Platform owner features
- **[TENANT_ADMIN_TODO.md](./TENANT_ADMIN_TODO.md)** - Enhanced tenant administration
- **[OPERATIONAL_MANAGEMENT_TODO.md](./OPERATIONAL_MANAGEMENT_TODO.md)** - Day-to-day operations

---

## 📖 API Documentation

### Authentication Endpoints
```
POST /register         - Register new tenant with admin user
POST /login           - Admin/Operator login
POST /customer/login  - Customer login
```

### Customer Management (Admin)
```
POST   /api/customers              - Register new customer
GET    /api/customers              - List all customers
GET    /api/customers/:id          - Get customer details
PUT    /api/customers/:id          - Update customer
DELETE /api/customers/:id          - Delete customer
POST   /api/customers/:id/activate - Activate customer
```

### Customer Self-Service
```
GET /api/customer/profile          - View own profile
PUT /api/customer/profile          - Update own profile
PUT /api/customer/password         - Change password
GET /api/customer/invoices         - View own invoices
GET /api/customer/payments         - View payment history
GET /api/customer/water-usage      - View usage history
```

### Subscription Types
```
POST   /api/subscription-types     - Create subscription type
GET    /api/subscription-types     - List subscription types
GET    /api/subscription-types/:id - Get subscription details
PUT    /api/subscription-types/:id - Update subscription type
DELETE /api/subscription-types/:id - Delete subscription type
```

### Water Usage & Rates
```
POST /api/water-usage              - Record meter reading
GET  /api/water-usage              - List water usage records
PUT  /api/water-usage/:id          - Update usage record
POST /api/water-rates              - Create water rate
GET  /api/water-rates/active       - Get active rates
PUT  /api/water-rates/:id          - Update rate
```

### Invoices
```
POST /api/invoices/generate-monthly - Generate monthly invoices
GET  /api/invoices                  - List invoices
GET  /api/invoices/:id              - Get invoice details
PUT  /api/invoices/:id              - Update invoice
```

### Payments
```
POST /api/payments                  - Record payment
GET  /api/payments                  - List payments
GET  /api/payments/:id              - Get payment details
PUT  /api/payments/:id              - Update payment
```

### Health & Monitoring
```
GET /health         - Basic health check
GET /health/live    - Liveness probe
GET /health/ready   - Readiness probe
GET /metrics        - System metrics
```

For complete API documentation with request/response examples, visit the Swagger UI at `/swagger/index.html` when the server is running.

---

## 🧪 Testing

### Postman Collection

A complete Postman collection is available in the repository:
- **Collection**: `Tirta-SaaS-Backend.postman_collection.json`
- **Environment**: `Tirta-SaaS-Development.postman_environment.json`

See [`POSTMAN_TESTING_GUIDE.md`](./POSTMAN_TESTING_GUIDE.md) for detailed testing instructions.

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

---

## 🚀 Deployment

### Docker Deployment

```bash
# Build Docker image
docker build -t tirta-saas-backend .

# Run container
docker run -p 8080:8080 --env-file .env tirta-saas-backend
```

### Server Deployment

See [`PRODUCTION_DEPLOYMENT_GUIDE.md`](./PRODUCTION_DEPLOYMENT_GUIDE.md) for complete deployment instructions including:
- Environment setup
- Database configuration
- Security hardening
- Monitoring setup
- Backup procedures
- Troubleshooting guide

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write unit tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting PR
- Follow the existing code structure and patterns

---

## 📄 License

This project is proprietary software. All rights reserved.

---

## 🙏 Acknowledgments

Built with:
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Validator](https://github.com/go-playground/validator)

---

## 📞 Support

For questions, issues, or feature requests:
- Open an issue on GitHub
- Contact the development team
- Check the documentation files

---

## 🏆 Project Status Summary

**✅ Production-Ready Core System**
- Multi-tenant water billing operations
- Enterprise security and performance
- Complete documentation and deployment guides
- Customer self-service portal

**⏳ Planned Enhancements (16 weeks)**
- Platform management dashboard
- Advanced operational features
- Mobile applications
- External integrations
- IoT and smart meter support

**🎉 Ready for immediate deployment with comprehensive roadmap for future enhancements!**

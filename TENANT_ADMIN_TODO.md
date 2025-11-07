# Tenant Admin Management TODO

## 🎯 Overview
This document outlines the tasks needed for Tenant Admin to manage their Paguyuban (water utility cooperative) operations, including user management, master data settings, and operational configurations.

## 📋 User Management

### 1. User Account Management
- [ ] Enhanced user creation with role assignment
  - [ ] POST `/api/tenant-users/create-with-profile` - Create user with complete profile
  - [ ] POST `/api/tenant-users/bulk-create` - Create multiple users from CSV
  - [ ] GET `/api/tenant-users/templates/csv` - Download CSV template for bulk user creation
- [ ] User profile management
  - [ ] GET `/api/tenant-users/:id/profile` - Get detailed user profile
  - [ ] PUT `/api/tenant-users/:id/profile` - Update user profile
  - [ ] POST `/api/tenant-users/:id/avatar` - Upload user avatar
  - [ ] GET `/api/tenant-users/:id/activity` - View user activity history
- [ ] User access control
  - [ ] POST `/api/tenant-users/:id/suspend` - Temporarily suspend user
  - [ ] POST `/api/tenant-users/:id/activate` - Reactivate suspended user
  - [ ] PUT `/api/tenant-users/:id/permissions` - Customize user permissions
  - [ ] GET `/api/tenant-users/:id/sessions` - View active sessions
  - [ ] POST `/api/tenant-users/:id/logout-all` - Force logout from all devices

### 2. Role & Permission Management
- [ ] Custom role creation for Paguyuban
  - [ ] GET `/api/roles` - List all available roles
  - [ ] POST `/api/roles/custom` - Create custom role (future enhancement)
  - [ ] GET `/api/permissions` - List all available permissions
  - [ ] GET `/api/roles/:role/permissions` - View permissions for a role

## 📋 Master Data Settings

### 1. Paguyuban Profile & Configuration
- [ ] Create comprehensive tenant profile model
  - [ ] Paguyuban name and legal information
  - [ ] Address and contact details
  - [ ] Bank account information for payments
  - [ ] Operating hours
  - [ ] Service area/coverage
- [ ] Tenant profile endpoints
  - [ ] GET `/api/tenant/profile` - Get Paguyuban profile
  - [ ] PUT `/api/tenant/profile` - Update Paguyuban profile
  - [ ] POST `/api/tenant/logo` - Upload Paguyuban logo
  - [ ] POST `/api/tenant/documents` - Upload legal documents

### 2. Water Tariff & Pricing Management
- [ ] Enhanced water rate configuration
  - [ ] GET `/api/water-rates/categories` - List rate categories
  - [ ] POST `/api/water-rates/categories` - Create rate category (residential, commercial, etc.)
  - [ ] PUT `/api/water-rates/categories/:id` - Update rate category
  - [ ] DELETE `/api/water-rates/categories/:id` - Delete rate category
- [ ] Progressive rate setup
  - [ ] POST `/api/water-rates/progressive` - Set progressive rates (0-10m³, 11-20m³, etc.)
  - [ ] GET `/api/water-rates/simulation` - Simulate bill calculation
  - [ ] POST `/api/water-rates/bulk-update` - Bulk update rates
  - [ ] GET `/api/water-rates/history` - View rate change history

### 3. Subscription Type Management
- [ ] Enhanced subscription configuration
  - [ ] GET `/api/subscription-types/templates` - Get subscription templates
  - [ ] POST `/api/subscription-types/from-template` - Create from template
  - [ ] PUT `/api/subscription-types/:id/features` - Configure subscription features
  - [ ] POST `/api/subscription-types/:id/clone` - Clone subscription type
- [ ] Subscription pricing
  - [ ] PUT `/api/subscription-types/:id/pricing` - Update registration fees and monthly fees
  - [ ] POST `/api/subscription-types/:id/discounts` - Add discount rules
  - [ ] GET `/api/subscription-types/:id/customers` - View customers by subscription

### 4. Payment Method Configuration
- [ ] Payment channel setup
  - [ ] GET `/api/payment-methods` - List available payment methods
  - [ ] POST `/api/payment-methods` - Add payment method (cash, bank transfer, e-wallet)
  - [ ] PUT `/api/payment-methods/:id` - Update payment method details
  - [ ] POST `/api/payment-methods/:id/toggle` - Enable/disable payment method
- [ ] Bank account management
  - [ ] GET `/api/bank-accounts` - List Paguyuban bank accounts
  - [ ] POST `/api/bank-accounts` - Add bank account
  - [ ] PUT `/api/bank-accounts/:id` - Update bank account
  - [ ] POST `/api/bank-accounts/:id/set-primary` - Set as primary account

### 5. Invoice & Billing Configuration
- [ ] Invoice settings
  - [ ] GET `/api/invoice-settings` - Get invoice configuration
  - [ ] PUT `/api/invoice-settings` - Update invoice settings
    - [ ] Invoice prefix/numbering format
    - [ ] Due date calculation (days after generation)
    - [ ] Late payment penalty percentage
    - [ ] Grace period days
    - [ ] Minimum bill amount
- [ ] Invoice template customization
  - [ ] GET `/api/invoice-templates` - List invoice templates
  - [ ] PUT `/api/invoice-templates/:id` - Customize invoice template
  - [ ] POST `/api/invoice-templates/:id/preview` - Preview invoice template

## 📋 Operational Settings

### 1. Service Area Management
- [ ] Area/zone configuration
  - [ ] GET `/api/service-areas` - List service areas
  - [ ] POST `/api/service-areas` - Create service area (RT/RW/Blok)
  - [ ] PUT `/api/service-areas/:id` - Update service area
  - [ ] POST `/api/service-areas/:id/assign-reader` - Assign meter reader to area
  - [ ] GET `/api/service-areas/:id/customers` - List customers in area

### 2. Meter Reading Configuration
- [ ] Reading schedule setup
  - [ ] GET `/api/reading-schedules` - List reading schedules
  - [ ] POST `/api/reading-schedules` - Create reading schedule
  - [ ] PUT `/api/reading-schedules/:id` - Update schedule
  - [ ] POST `/api/reading-schedules/:id/assign-areas` - Assign areas to schedule
- [ ] Reading validation rules
  - [ ] GET `/api/reading-rules` - Get validation rules
  - [ ] PUT `/api/reading-rules` - Update validation rules
    - [ ] Maximum usage threshold
    - [ ] Minimum usage warning
    - [ ] Abnormal usage detection

### 3. Customer Categories
- [ ] Category management
  - [ ] GET `/api/customer-categories` - List customer categories
  - [ ] POST `/api/customer-categories` - Create category (household, business, mosque, etc.)
  - [ ] PUT `/api/customer-categories/:id` - Update category
  - [ ] POST `/api/customer-categories/:id/rates` - Assign special rates to category

### 4. Notification Templates
- [ ] Template management
  - [ ] GET `/api/notification-templates` - List notification templates
  - [ ] PUT `/api/notification-templates/:id` - Customize template content
  - [ ] POST `/api/notification-templates/:id/preview` - Preview notification
  - [ ] PUT `/api/notification-templates/:id/toggle` - Enable/disable template
- [ ] Notification types
  - [ ] Bill ready notification
  - [ ] Payment reminder
  - [ ] Payment received confirmation
  - [ ] Service interruption notice
  - [ ] Meter reading schedule

### 5. Business Rules Configuration
- [ ] Operational rules
  - [ ] GET `/api/business-rules` - List business rules
  - [ ] PUT `/api/business-rules` - Update rules
    - [ ] Minimum payment amount
    - [ ] Partial payment allowed (yes/no)
    - [ ] Auto-suspend after X days unpaid
    - [ ] Reconnection fee amount
    - [ ] New connection requirements

## 📋 Data Management

### 1. Import/Export Tools
- [ ] Data import
  - [ ] POST `/api/import/customers` - Import customers from Excel/CSV
  - [ ] POST `/api/import/meter-readings` - Import historical meter readings
  - [ ] POST `/api/import/payments` - Import payment history
  - [ ] GET `/api/import/status/:jobId` - Check import job status
- [ ] Data export
  - [ ] POST `/api/export/customers` - Export customer data
  - [ ] POST `/api/export/financial-report` - Export financial reports
  - [ ] POST `/api/export/usage-report` - Export usage reports
  - [ ] GET `/api/export/download/:fileId` - Download exported file

### 2. Data Backup & Archive
- [ ] Backup configuration
  - [ ] GET `/api/backup/settings` - Get backup settings
  - [ ] PUT `/api/backup/settings` - Configure automatic backup
  - [ ] POST `/api/backup/manual` - Trigger manual backup
  - [ ] GET `/api/backup/history` - View backup history

## 📋 Dashboard & Analytics

### 1. Operational Dashboard
- [ ] Dashboard widgets
  - [ ] GET `/api/dashboard/summary` - Get dashboard summary
    - [ ] Total active customers
    - [ ] Monthly collection rate
    - [ ] Outstanding payments
    - [ ] Today's collections
  - [ ] GET `/api/dashboard/charts/collection-trend` - Collection trend chart
  - [ ] GET `/api/dashboard/charts/usage-pattern` - Usage pattern analysis
  - [ ] GET `/api/dashboard/alerts` - Operational alerts

### 2. Reports Configuration
- [ ] Report templates
  - [ ] GET `/api/reports/templates` - List report templates
  - [ ] POST `/api/reports/schedule` - Schedule automatic reports
  - [ ] PUT `/api/reports/schedule/:id` - Update report schedule
  - [ ] GET `/api/reports/generated` - List generated reports

## 🛠️ Technical Requirements

### Database Models
- [ ] Create `tenant_settings` table with JSON column for flexible configuration
- [ ] Create `service_areas` table for area management
- [ ] Create `customer_categories` table
- [ ] Create `reading_schedules` table
- [ ] Create `business_rules` table
- [ ] Create `import_jobs` table for tracking imports
- [ ] Create `export_jobs` table for tracking exports

### API Patterns
- [ ] All endpoints should be under `/api/` prefix
- [ ] Use consistent REST patterns
- [ ] Implement pagination for list endpoints
- [ ] Add filtering and search capabilities
- [ ] Include validation for all inputs
- [ ] Return consistent error messages

### Security Considerations
- [ ] Validate tenant context for all operations
- [ ] Implement audit logging for all changes
- [ ] Ensure proper permission checks
- [ ] Sanitize file uploads
- [ ] Implement rate limiting for exports/imports

## 🚀 Implementation Priority

### Phase 1 (Week 1) - Core Settings
1. Paguyuban profile configuration
2. User management enhancements
3. Water tariff configuration
4. Basic payment method setup

### Phase 2 (Week 2) - Operational Setup
1. Service area management
2. Customer category configuration
3. Invoice settings
4. Meter reading configuration

### Phase 3 (Week 3) - Advanced Features
1. Notification template customization
2. Business rules configuration
3. Import/Export tools
4. Dashboard implementation

### Phase 4 (Week 4) - Polish & Optimize
1. Report generation
2. Backup configuration
3. Advanced analytics
4. Performance optimization

## ✅ Definition of Done

- [ ] All CRUD operations implemented
- [ ] Validation rules applied
- [ ] Permission checks in place
- [ ] Audit logging implemented
- [ ] API documentation updated
- [ ] Unit tests written
- [ ] Integration tests passed
- [ ] UI mockups approved
- [ ] Performance benchmarks met

## 📝 Notes

- All settings should have sensible defaults
- Configuration changes should be audited
- Consider adding configuration templates for quick setup
- Allow configuration export/import between environments
- Implement configuration versioning for rollback capability
- Add configuration validation before saving
- Consider multi-language support for templates
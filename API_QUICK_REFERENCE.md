# Phase 1 API Quick Reference

## 🚀 Platform Owner Endpoints

### Authentication Required
All endpoints require:
- JWT token in `Authorization: Bearer <token>` header
- Platform Owner or Admin role

---

## Tenant Management

### List All Tenants
```http
GET /api/platform/tenants?page=1&page_size=20&search=&status=&sort_by=created_at&sort_order=desc
```

**Query Parameters:**
- `page` (int) - Page number (default: 1)
- `page_size` (int) - Items per page (default: 20)
- `search` (string) - Search by name, village code, email
- `status` (string) - Filter by status: ACTIVE, SUSPENDED, INACTIVE
- `subscription_plan` (string) - Filter by plan: FREE, BASIC, PREMIUM, ENTERPRISE
- `sort_by` (string) - Sort field: created_at, name, total_customers
- `sort_order` (string) - Sort order: asc, desc

### Get Tenant Details
```http
GET /api/platform/tenants/:id
```

### Update Tenant
```http
PUT /api/platform/tenants/:id
Content-Type: application/json

{
  "name": "Updated Tenant Name",
  "email": "newemail@example.com",
  "phone": "081234567890",
  "address": "New Address",
  "notes": "Updated notes"
}
```

### Suspend Tenant
```http
POST /api/platform/tenants/:id/suspend
Content-Type: application/json

{
  "reason": "Payment overdue"
}
```

### Activate Tenant
```http
POST /api/platform/tenants/:id/activate
```

### Delete Tenant (Soft Delete)
```http
DELETE /api/platform/tenants/:id
```

### Get Tenant Statistics
```http
GET /api/platform/tenants/:id/statistics
```

---

## Platform Analytics

### Overview Statistics
```http
GET /api/platform/analytics/overview
```

**Returns:**
- Tenant counts (total, active, suspended, trial)
- Revenue metrics (total, monthly, outstanding)
- Growth metrics (new tenants, churn rate)
- Usage metrics (users, customers, storage)

### Tenant Growth Analytics
```http
GET /api/platform/analytics/tenants?months=6
```

**Query Parameters:**
- `months` (int) - Number of months to analyze (1-24, default: 6)

**Returns:**
- Monthly breakdown of tenant acquisition/churn
- Growth rates
- Distribution by plan and status

### Revenue Analytics
```http
GET /api/platform/analytics/revenue?months=6
```

**Query Parameters:**
- `months` (int) - Number of months to analyze (1-24, default: 6)

**Returns:**
- Total revenue and MRR
- Average revenue per tenant
- Monthly revenue breakdown
- Revenue by subscription plan
- Payment method statistics

### Usage Analytics
```http
GET /api/platform/analytics/usage?months=6
```

**Query Parameters:**
- `months` (int) - Number of months to analyze (1-24, default: 6)

**Returns:**
- Total users, customers, invoices, payments
- Water usage statistics
- Storage consumption
- Monthly usage trends
- Top tenants by usage

---

## 🏢 Tenant Admin Endpoints

### Authentication Required
All endpoints require:
- JWT token in `Authorization: Bearer <token>` header
- Tenant Admin role
- Valid tenant context

---

## Tenant Settings

### Get Settings
```http
GET /api/tenant/settings
```

### Update Settings
```http
PUT /api/tenant/settings
Content-Type: application/json

{
  "company_name": "Water Utility Company",
  "address": "123 Main Street",
  "phone": "081234567890",
  "email": "info@company.com",
  "website": "https://company.com",
  "primary_color": "#1976D2",
  "secondary_color": "#424242",
  "invoice_prefix": "INV",
  "invoice_due_days": 30,
  "invoice_footer_text": "Thank you for your business",
  "late_penalty_percent": 2.0,
  "late_penalty_max_cap": 100000,
  "grace_period_days": 7,
  "minimum_bill_amount": 10000,
  "bank_name": "Bank ABC",
  "bank_account_name": "Company Account",
  "bank_account_no": "1234567890",
  "operating_hours": "Mon-Fri 08:00-17:00",
  "service_area": "City Area",
  "timezone": "Asia/Jakarta",
  "language": "id"
}
```

### Upload Logo
```http
POST /api/tenant/settings/logo
Content-Type: multipart/form-data

logo=@/path/to/image.png
```

**File Requirements:**
- **Allowed types:** JPEG, JPG, PNG, GIF, WebP
- **Max size:** 5MB
- **Field name:** `logo`

**Returns:**
```json
{
  "status": "success",
  "message": "Logo uploaded successfully",
  "data": {
    "logo_url": "uploads/tenants/{tenant-id}/logos/{uuid}_{timestamp}.png"
  }
}
```

---

## 📊 Response Format

### Success Response
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description",
  "error": "Detailed error message"
}
```

### List Response with Pagination
```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": {
    "items": [...],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

---

## 🔐 Authentication

### Get JWT Token
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "role": "ADMIN"
    }
  }
}
```

### Use Token in Requests
```http
GET /api/platform/tenants
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 🧪 Testing with cURL

### Example: List Tenants
```bash
curl -X GET "http://localhost:8080/api/platform/tenants?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Example: Upload Logo
```bash
curl -X POST "http://localhost:8080/api/tenant/settings/logo" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "logo=@/path/to/logo.png"
```

### Example: Get Revenue Analytics
```bash
curl -X GET "http://localhost:8080/api/platform/analytics/revenue?months=12" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📝 Notes

### File Upload
- Files are stored in `uploads/tenants/{tenant-id}/logos/`
- Old logo is automatically deleted when uploading new one
- Files are named with UUID for security
- Maximum file size: 5MB

### Analytics
- All analytics support `months` parameter (1-24)
- Default period is 6 months
- Statistics are calculated in real-time
- Monthly breakdowns show historical trends

### Pagination
- Default page size: 20 items
- Maximum page size: 100 items
- Total count included in response
- Zero-indexed pages (page 1 is first page)

### Filtering
- Multiple filters can be combined
- Search is case-insensitive
- Partial matching supported for search
- Status and plan filters use exact match

---

## 🔗 API Documentation

Interactive API documentation available at:
```
http://localhost:8080/docs
```

---

**Version:** 1.0.0  
**Last Updated:** November 30, 2025

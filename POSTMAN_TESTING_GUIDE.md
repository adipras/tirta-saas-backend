# Complete Postman Testing Guide for Tirta-SaaS Backend

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Testing Sequence](#testing-sequence)
4. [Error Testing Scenarios](#error-testing-scenarios)
5. [Tips and Troubleshooting](#tips-and-troubleshooting)

## Prerequisites

Before starting the tests, ensure:
- The Tirta-SaaS backend is running on `http://localhost:8080`
- MySQL database is running and properly configured
- You have Postman installed (version 9.0 or higher recommended)
- Import the provided Postman collection: `Tirta-SaaS-Backend.postman_collection.json`

## Environment Setup

### 1. Create a New Environment in Postman

Create a new environment called "Tirta-SaaS Development" with these variables:

| Variable Name | Initial Value | Current Value | Description |
|--------------|---------------|---------------|-------------|
| base_url | http://localhost:8080 | http://localhost:8080 | Base URL of the API (REQUIRED) |
| admin_token | | | JWT token for admin authentication |
| customer_token | | | JWT token for customer authentication |
| tenant_id | | | UUID of the created tenant |
| tenant_name | PT Air Bersih Jakarta | PT Air Bersih Jakarta | Name of the tenant company |
| village_code | JKT001 | JKT001 | Village/area code for the tenant |
| admin_email | admin@airbersih.co.id | admin@airbersih.co.id | Admin email for login |
| admin_password | SecurePass123! | SecurePass123! | Admin password |
| customer_id | | | UUID of the created customer |
| customer_email | john.doe@email.com | john.doe@email.com | Customer email for login |
| customer_password | CustomerPass123! | CustomerPass123! | Customer password |
| subscription_type_id | | | UUID of the subscription type |
| water_rate_id | | | UUID of the water rate |
| invoice_id | | | UUID of the invoice |
| payment_id | | | UUID of the payment |
| usage_id | | | UUID of water usage record |

### 2. Pre-request Scripts

The collection includes pre-request scripts that automatically:
- Set Content-Type headers
- Add Authorization headers when tokens are available
- Generate timestamps for date fields

### 3. Tests Scripts

Each request includes test scripts that:
- Validate response status codes
- Save response data to environment variables
- Verify response structure
- Log important information to console

## Testing Sequence

Follow this sequence for a complete end-to-end test flow:

### Phase 1: Tenant Setup and Admin Authentication

#### Test 1: Register New Tenant
- **Endpoint**: `POST /auth/register`
- **Purpose**: Create a new tenant with admin user
- **Expected**: 201 Created, saves tenant_id
- **Validates**: Tenant multi-tenancy setup

#### Test 2: Admin Login
- **Endpoint**: `POST /auth/login`
- **Purpose**: Authenticate as admin user
- **Expected**: 200 OK, saves admin_token
- **Validates**: JWT authentication system

### Phase 2: Master Data Setup

#### Test 3: Create Subscription Type
- **Endpoint**: `POST /api/subscription-types`
- **Purpose**: Define customer subscription plans
- **Expected**: 201 Created, saves subscription_type_id
- **Validates**: Subscription management

#### Test 4: List Subscription Types
- **Endpoint**: `GET /api/subscription-types`
- **Purpose**: Verify subscription creation
- **Expected**: 200 OK with pagination
- **Validates**: List endpoints and pagination

#### Test 5: Create Water Rate
- **Endpoint**: `POST /api/water-rates`
- **Purpose**: Set pricing for water usage
- **Expected**: 201 Created, saves water_rate_id
- **Validates**: Water rate management

### Phase 3: Customer Management

#### Test 6: Register Customer
- **Endpoint**: `POST /api/customers`
- **Purpose**: Create new customer account
- **Expected**: 201 Created with registration invoice
- **Validates**: Customer creation and auto-invoice generation

#### Test 7: List Customers
- **Endpoint**: `GET /api/customers`
- **Purpose**: View all customers (paginated)
- **Expected**: 200 OK with customer list
- **Validates**: Customer listing and filtering

#### Test 8: Get Registration Invoice
- **Endpoint**: `GET /api/invoices`
- **Purpose**: Find customer's registration invoice
- **Expected**: 200 OK, saves invoice_id
- **Validates**: Invoice filtering and retrieval

### Phase 4: Payment Processing

#### Test 9: Process Registration Payment
- **Endpoint**: `POST /api/payments`
- **Purpose**: Pay registration fee to activate customer
- **Expected**: 201 Created, customer activated
- **Validates**: Payment processing and customer activation

### Phase 5: Customer Portal

#### Test 10: Customer Login
- **Endpoint**: `POST /customer/auth/login`
- **Purpose**: Authenticate as customer
- **Expected**: 200 OK, saves customer_token
- **Validates**: Customer authentication system

#### Test 11: View Customer Profile
- **Endpoint**: `GET /customer/profile`
- **Purpose**: Customer views own profile
- **Expected**: 200 OK with customer details
- **Validates**: Customer self-service access

#### Test 12: Customer Change Password
- **Endpoint**: `PUT /customer/change-password`
- **Purpose**: Update customer password
- **Expected**: 200 OK
- **Validates**: Password management

### Phase 6: Water Usage and Billing

#### Test 13: Record Water Usage
- **Endpoint**: `POST /api/water-usage`
- **Purpose**: Record monthly meter reading
- **Expected**: 201 Created
- **Validates**: Usage tracking

#### Test 14: Generate Monthly Invoices
- **Endpoint**: `POST /api/invoices/generate-monthly`
- **Purpose**: Create invoices for all customers
- **Expected**: 200 OK with invoice count
- **Validates**: Bulk invoice generation

#### Test 15: Customer View Invoices
- **Endpoint**: `GET /customer/invoices`
- **Purpose**: Customer views their invoices
- **Expected**: 200 OK with invoice list
- **Validates**: Customer invoice access

#### Test 16: Process Monthly Payment (Customer)
- **Endpoint**: `POST /customer/payments`
- **Purpose**: Customer pays monthly bill
- **Expected**: 201 Created
- **Validates**: Customer payment processing

### Phase 7: Admin Operations

#### Test 17: View All Payments
- **Endpoint**: `GET /api/payments`
- **Purpose**: Admin views payment history
- **Expected**: 200 OK with paginated payments
- **Validates**: Payment tracking and filtering

#### Test 18: Update Subscription Type
- **Endpoint**: `PUT /api/subscription-types/{id}`
- **Purpose**: Modify subscription details
- **Expected**: 200 OK
- **Validates**: Update operations

#### Test 19: Delete Water Rate
- **Endpoint**: `DELETE /api/water-rates/{id}`
- **Purpose**: Remove old water rate
- **Expected**: 200 OK
- **Validates**: Soft delete functionality

### Phase 8: System Health and Security

#### Test 20: Health Check
- **Endpoint**: `GET /health`
- **Purpose**: Check system health
- **Expected**: 200 OK with component status
- **Validates**: Health monitoring

#### Test 21: Readiness Check
- **Endpoint**: `GET /ready`
- **Purpose**: Check if system is ready
- **Expected**: 200 OK
- **Validates**: Kubernetes readiness

#### Test 22: System Metrics
- **Endpoint**: `GET /metrics`
- **Purpose**: View system metrics
- **Expected**: 200 OK with detailed metrics
- **Validates**: Performance monitoring

#### Test 23: Rate Limiting Test
- **Endpoint**: Multiple rapid requests to any endpoint
- **Purpose**: Verify rate limiting
- **Expected**: 429 Too Many Requests after limit
- **Validates**: API security

## Error Testing Scenarios

### Authentication Errors

#### Test E1: Invalid Login Credentials
- **Endpoint**: `POST /auth/login`
- **Body**: Wrong password
- **Expected**: 401 Unauthorized

#### Test E2: Expired Token
- **Endpoint**: Any protected endpoint
- **Header**: Expired or invalid token
- **Expected**: 401 Unauthorized

#### Test E3: Unauthorized Access
- **Endpoint**: `GET /api/customers` 
- **Header**: Customer token (not admin)
- **Expected**: 403 Forbidden

### Validation Errors

#### Test E4: Invalid Email Format
- **Endpoint**: `POST /auth/register`
- **Body**: Invalid email format
- **Expected**: 400 Bad Request with validation errors

#### Test E5: Missing Required Fields
- **Endpoint**: `POST /api/customers`
- **Body**: Missing name or email
- **Expected**: 400 Bad Request

#### Test E6: Duplicate Customer Email
- **Endpoint**: `POST /api/customers`
- **Body**: Existing email
- **Expected**: 409 Conflict

#### Test E7: Invalid Meter Number Length
- **Endpoint**: `POST /auth/customers`
- **Body**: Meter number less than 3 or more than 20 characters
- **Expected**: 400 Bad Request

#### Test E8: Invalid Customer Name Length
- **Endpoint**: `POST /auth/customers`
- **Body**: Name less than 2 or more than 100 characters
- **Expected**: 400 Bad Request

### Business Logic Errors

#### Test E9: Invalid Payment Amount
- **Endpoint**: `POST /api/payments`
- **Body**: Negative or zero amount
- **Expected**: 400 Bad Request

#### Test E10: Payment Exceeds Maximum Limit
- **Endpoint**: `POST /api/payments`
- **Body**: Amount exceeding 999,999
- **Expected**: 400 Bad Request

#### Test E11: Overpayment
- **Endpoint**: `POST /api/payments`
- **Body**: Amount exceeding invoice remaining balance
- **Expected**: 400 Bad Request

#### Test E12: Invalid Meter Reading (Negative)
- **Endpoint**: `POST /api/water-usage`
- **Body**: Negative meter_end value
- **Expected**: 400 Bad Request

#### Test E13: Invalid Meter Reading (Too High)
- **Endpoint**: `POST /api/water-usage`
- **Body**: meter_end value exceeding 99,999,999
- **Expected**: 400 Bad Request

#### Test E14: Meter Reading Going Backwards
- **Endpoint**: `POST /api/water-usage`
- **Body**: meter_end less than previous month's meter_end
- **Expected**: 400 Bad Request

#### Test E15: Excessive Water Usage
- **Endpoint**: `POST /api/water-usage`
- **Body**: Usage calculating to more than 1000 m3/month
- **Expected**: 400 Bad Request

#### Test E16: Duplicate Invoice Generation
- **Endpoint**: `POST /api/invoices/generate-monthly`
- **Body**: Same month twice
- **Expected**: 200 OK with skipped count

#### Test E17: Inactive Customer Login
- **Endpoint**: `POST /customer/auth/login`
- **Body**: Unpaid customer credentials
- **Expected**: 401 Unauthorized

## Tips and Troubleshooting

### Common Issues and Solutions

1. **Connection Refused Error**
   - Ensure the backend is running: `go run main.go`
   - Check if port 8080 is available
   - Verify database connection in `.env` file

2. **401 Unauthorized Errors**
   - Check if token is properly saved in environment
   - Verify token hasn't expired (24-hour validity)
   - Re-run login endpoint to get fresh token

3. **Tenant ID Not Found**
   - Ensure you're using the correct tenant_id
   - Check if previous requests saved variables correctly
   - Look at Postman Console for saved values

4. **Rate Limiting Hit**
   - Wait 1 minute before retrying
   - Check rate limit headers in response
   - Different endpoints have different limits

### Best Practices

1. **Sequential Testing**
   - Always run tests in order first time
   - Some tests depend on data from previous tests
   - Use collection runner for automated testing

2. **Environment Management**
   - Create separate environments for dev/staging/prod
   - Don't commit sensitive data like passwords
   - Use initial values for defaults

3. **Debugging**
   - Check Postman Console for detailed logs
   - Verify environment variables are set
   - Look at response headers for debugging info
   - Check application logs for server-side errors

4. **Data Cleanup**
   - The system uses soft deletes
   - Create new test data for each test run
   - Use unique emails/customer IDs

### Performance Testing

1. **Load Testing**
   - Use Postman Collection Runner
   - Set iterations for stress testing
   - Monitor response times

2. **Concurrent Users**
   - Create multiple environments with different users
   - Test simultaneous operations
   - Verify data isolation between tenants

### Security Testing

1. **SQL Injection**
   - Try SQL in input fields
   - System should reject with 400 Bad Request

2. **XSS Prevention**
   - Try script tags in text fields
   - Check sanitization in responses

3. **Authorization**
   - Try accessing other tenant's data
   - Verify proper isolation

## Appendix: Quick Reference

### HTTP Status Codes Used
- **200 OK**: Successful GET, PUT
- **201 Created**: Successful POST
- **400 Bad Request**: Validation error, business rule violation
- **401 Unauthorized**: Invalid/missing token, inactive customer
- **403 Forbidden**: Insufficient permissions
- **404 Not Found**: Resource doesn't exist
- **409 Conflict**: Duplicate resource (email, meter number)
- **429 Too Many Requests**: Rate limit exceeded
- **500 Internal Server Error**: Server error

### Required Headers
- **Content-Type**: `application/json` (for POST/PUT)
- **Authorization**: `Bearer {{token}}` (for protected endpoints)

### Response Structure

All API responses now follow standardized formats:

#### Success Response Format
```json
{
  "id": "uuid",
  "meter_number": "MTR001",
  "name": "John Doe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2024-01-31T10:00:00Z"
}
```

#### List Response Format
```json
{
  "customers": [
    {
      "id": "uuid",
      "meter_number": "MTR001",
      "name": "John Doe",
      "email": "john@example.com",
      "is_active": true
    }
  ],
  "total": 25
}
```

#### Error Response Format
```json
{
  "error": "Detailed error message"
}
```

#### Business Validation Limits
- **Meter Numbers**: 3-20 characters (alphanumeric)
- **Customer Names**: 2-100 characters
- **Water Usage**: Maximum 1000 m3 per month
- **Meter Readings**: 0 to 99,999,999 range
- **Payment Amounts**: Positive values up to 999,999
- **Invoice Totals**: Must be positive and reasonable

### Date Formats
- **Dates**: `YYYY-MM-DD` (e.g., "2024-01-31")
- **Month**: `YYYY-MM` (e.g., "2024-01")

### Pagination Parameters
- **page**: Page number (default: 1)
- **per_page**: Items per page (default: 10, max: 100)
- **search**: Search term (optional)
- **sort**: Sort field (optional)
- **order**: Sort order - asc/desc (optional)
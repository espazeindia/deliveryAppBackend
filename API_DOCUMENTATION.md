# Delivery App Backend - API Documentation

Complete API reference for the Espaze Delivery Partner Backend.

## Table of Contents
1. [Authentication](#authentication)
2. [Order Management](#order-management)
3. [Profile Management](#profile-management)
4. [Earnings](#earnings)
5. [Error Responses](#error-responses)

## Base URL
```
Production: https://api.espaze.com/api/v1
Development: http://localhost:8081/api/v1
```

## Authentication

### Headers
All protected endpoints require JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 1. Authentication Endpoints

### 1.1 Login with PIN

Login using phone number and PIN.

**Endpoint:** `POST /delivery/login`

**Request Body:**
```json
{
  "phoneNumber": "9876543210",
  "pin": 123456
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phoneNumber": "9876543210",
    "isAvailable": false
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "message": "Invalid PIN"
}
```

---

### 1.2 Request OTP

Request OTP for phone number verification.

**Endpoint:** `POST /delivery/request-otp`

**Request Body:**
```json
{
  "phoneNumber": "9876543210"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP sent successfully"
}
```

---

### 1.3 Verify OTP

Verify OTP and login.

**Endpoint:** `POST /delivery/verify-otp`

**Request Body:**
```json
{
  "phoneNumber": "9876543210",
  "otp": 123456
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP verified successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phoneNumber": "9876543210",
    "isAvailable": false
  }
}
```

---

## 2. Order Management

All order endpoints are protected and require authentication.

### 2.1 Get Active Orders

Get all active orders assigned to the delivery partner.

**Endpoint:** `GET /delivery/orders/active`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "orders": [
    {
      "id": "507f1f77bcf86cd799439012",
      "orderId": "ORD123456",
      "status": "in_transit",
      "address": "123 Main St, Bangalore",
      "amount": 1500,
      "deliveryFee": 50,
      "itemsCount": 3,
      "distance": 5.2,
      "createdAt": "2025-10-26T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

### 2.2 Get Order History

Get completed order history with pagination.

**Endpoint:** `GET /delivery/orders/history`

**Query Parameters:**
- `limit` (optional): Number of records per page (default: 20)
- `offset` (optional): Number of records to skip (default: 0)

**Example:** `GET /delivery/orders/history?limit=20&offset=0`

**Success Response (200 OK):**
```json
{
  "success": true,
  "orders": [
    {
      "id": "507f1f77bcf86cd799439013",
      "orderId": "ORD123457",
      "status": "delivered",
      "address": "456 Park Ave, Bangalore",
      "amount": 2500,
      "deliveryFee": 75,
      "itemsCount": 5,
      "distance": 8.5,
      "createdAt": "2025-10-25T15:20:00Z"
    }
  ],
  "total": 50,
  "limit": 20,
  "offset": 0,
  "hasNext": true,
  "hasPrevious": false
}
```

---

### 2.3 Get Order Details

Get detailed information about a specific order.

**Endpoint:** `GET /delivery/orders/:id`

**Path Parameters:**
- `id`: Delivery ID

**Success Response (200 OK):**
```json
{
  "success": true,
  "order": {
    "id": "507f1f77bcf86cd799439012",
    "orderId": "ORD123456",
    "partnerId": "507f1f77bcf86cd799439011",
    "customerId": "507f1f77bcf86cd799439020",
    "customerName": "Jane Smith",
    "customerPhone": "9876543211",
    "status": "in_transit",
    "pickupAddress": "Warehouse A, Electronics City",
    "deliveryAddress": "123 Main St, Bangalore",
    "pickupLatitude": 12.8456,
    "pickupLongitude": 77.6632,
    "deliveryLatitude": 12.9716,
    "deliveryLongitude": 77.5946,
    "distance": 5.2,
    "orderAmount": 1500,
    "deliveryFee": 50,
    "itemsCount": 3,
    "items": [
      {
        "productId": "PROD001",
        "name": "Product 1",
        "quantity": 2,
        "price": 500,
        "imageUrl": "https://example.com/image1.jpg"
      }
    ],
    "assignedAt": "2025-10-26T10:00:00Z",
    "pickedUpAt": "2025-10-26T10:15:00Z",
    "inTransitAt": "2025-10-26T10:20:00Z",
    "createdAt": "2025-10-26T10:00:00Z",
    "updatedAt": "2025-10-26T10:20:00Z",
    "paymentMethod": "online",
    "notes": ""
  }
}
```

---

### 2.4 Accept Order

Accept a pending order.

**Endpoint:** `POST /delivery/orders/:id/accept`

**Path Parameters:**
- `id`: Delivery ID

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Order accepted successfully"
}
```

---

### 2.5 Update Order Status

Update the status of an order.

**Endpoint:** `POST /delivery/orders/:id/status`

**Request Body:**
```json
{
  "status": "in_transit",
  "latitude": 12.9716,
  "longitude": 77.5946
}
```

**Valid Status Values:**
- `picked_up`: Order picked up from warehouse
- `in_transit`: On the way to delivery location
- `delivered`: Order delivered (use complete delivery endpoint instead)

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Status updated successfully"
}
```

---

### 2.6 Complete Delivery

Mark an order as delivered.

**Endpoint:** `POST /delivery/orders/:id/complete`

**Request Body:**
```json
{
  "latitude": 12.9716,
  "longitude": 77.5946,
  "signature": "base64_encoded_signature_image",
  "notes": "Delivered to customer directly"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Delivery completed successfully"
}
```

---

## 3. Profile Management

### 3.1 Get Profile

Get delivery partner profile information.

**Endpoint:** `GET /delivery/profile`

**Success Response (200 OK):**
```json
{
  "success": true,
  "profile": {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phoneNumber": "9876543210",
    "email": "john@example.com",
    "isAvailable": true,
    "isVerified": true,
    "rating": 4.8,
    "totalDeliveries": 342,
    "aadharNumber": "123456789012",
    "panNumber": "ABCDE1234F",
    "drivingLicense": "DL1234567890",
    "vehicleNumber": "KA01AB1234",
    "vehicleType": "bike",
    "bankAccountNumber": "1234567890",
    "ifsc": "ABCD0001234",
    "currentLatitude": 12.9716,
    "currentLongitude": 77.5946,
    "lastLocationAt": "2025-10-26T10:30:00Z",
    "createdAt": "2025-01-01T00:00:00Z"
  }
}
```

---

### 3.2 Update Profile

Update delivery partner profile.

**Endpoint:** `PUT /delivery/profile`

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "aadharNumber": "123456789012",
  "panNumber": "ABCDE1234F",
  "drivingLicense": "DL1234567890",
  "vehicleNumber": "KA01AB1234",
  "vehicleType": "bike",
  "bankAccountNumber": "1234567890",
  "ifsc": "ABCD0001234"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully"
}
```

---

### 3.3 Update Location

Update current location of delivery partner.

**Endpoint:** `POST /delivery/location`

**Request Body:**
```json
{
  "latitude": 12.9716,
  "longitude": 77.5946
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Location updated successfully"
}
```

---

### 3.4 Toggle Availability

Update availability status.

**Endpoint:** `POST /delivery/availability`

**Request Body:**
```json
{
  "isAvailable": true
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "You are now online"
}
```

---

## 4. Earnings

### 4.1 Get Earnings

Get earnings summary for a specific period.

**Endpoint:** `GET /delivery/earnings`

**Query Parameters:**
- `period` (optional): `today`, `week`, or `month` (default: `week`)

**Example:** `GET /delivery/earnings?period=week`

**Success Response (200 OK):**
```json
{
  "success": true,
  "totalEarnings": 12540,
  "deliveriesCount": 47,
  "avgPerDelivery": 266,
  "bonusEarnings": 1881,
  "weeklyCount": 47,
  "period": "week"
}
```

---

### 4.2 Get Earnings History

Get detailed earnings history with pagination.

**Endpoint:** `GET /delivery/earnings/history`

**Query Parameters:**
- `limit` (optional): Number of records per page (default: 50)
- `offset` (optional): Number of records to skip (default: 0)

**Example:** `GET /delivery/earnings/history?limit=50&offset=0`

**Success Response (200 OK):**
```json
{
  "success": true,
  "history": [
    {
      "orderId": "ORD123456",
      "amount": 75,
      "completedAt": "2025-10-26T10:45:00Z"
    },
    {
      "orderId": "ORD123457",
      "amount": 50,
      "completedAt": "2025-10-26T09:30:00Z"
    }
  ],
  "total": 342,
  "limit": 50,
  "offset": 0,
  "hasNext": true,
  "hasPrevious": false
}
```

---

## Error Responses

### Standard Error Format
```json
{
  "success": false,
  "error": "Error message description"
}
```

### HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Invalid or missing token |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error |

### Common Error Messages

#### Authentication Errors
```json
{
  "success": false,
  "error": "Authorization header required"
}
```

```json
{
  "success": false,
  "error": "Invalid or expired token"
}
```

#### Validation Errors
```json
{
  "success": false,
  "error": "Key: 'DeliveryPartnerLoginRequest.PhoneNumber' Error:Field validation for 'PhoneNumber' failed on the 'required' tag"
}
```

#### Resource Not Found
```json
{
  "success": false,
  "error": "Order not found"
}
```

---

## Rate Limiting

- Authentication endpoints: 5 requests per minute per IP
- Other endpoints: 100 requests per minute per user

## Versioning

The API uses URL versioning. Current version: `v1`

Example: `/api/v1/delivery/orders/active`

## Support

For API issues:
- Email: api-support@espaze.com
- Slack: #delivery-api-support


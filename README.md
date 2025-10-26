# Espaze Delivery App Backend

A robust Go backend service for the Espaze delivery partner mobile application. Built with clean architecture principles using Gin framework and MongoDB.

## Features

- 🔐 **Authentication** - JWT-based authentication with PIN and OTP support
- 📦 **Order Management** - Real-time order tracking and status updates
- 💰 **Earnings Tracking** - Comprehensive earnings and transaction history
- 📍 **Location Services** - Real-time location updates for delivery partners
- 👤 **Profile Management** - Partner profile and document management
- 🔄 **Availability Toggle** - Dynamic availability status management

## Tech Stack

- **Go** (1.24.4) - Programming language
- **Gin** - HTTP web framework
- **MongoDB** - NoSQL database
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing

## Architecture

The project follows Clean Architecture principles with clear separation of concerns:

```
deliveryAppBackend/
├── config/                 # Configuration and database setup
│   └── db.go              # MongoDB connection
├── domain/                # Business logic layer
│   ├── entities/          # Domain entities and DTOs
│   └── repositories/      # Repository interfaces
├── infrastructure/        # External dependencies
│   └── mongodb/           # MongoDB implementations
├── usecase/              # Application business logic
├── handlers/             # HTTP request handlers
├── routes/               # Route definitions
├── middlewares/          # HTTP middlewares
├── utils/                # Utility functions
└── main.go              # Application entry point
```

## Prerequisites

- Go 1.24.4 or higher
- MongoDB 4.4 or higher
- Git

## Installation

1. Clone the repository:
```bash
cd deliveryAppBackend
```

2. Install dependencies:
```bash
go mod download
```

3. Create environment file:
```bash
cp .env.example .env
```

4. Update `.env` file with your configuration:
```env
PORT=8081
MONGO_URI=mongodb://localhost:27017/espaze_delivery
JWT_SECRET=your-secure-secret-key
```

## Running the Application

### Development Mode

```bash
go run main.go
```

The server will start on `http://localhost:8081`

### Production Build

```bash
go build -o bin/delivery-backend main.go
./bin/delivery-backend
```

## API Documentation

### Base URL
```
http://localhost:8081/api/v1
```

### Authentication Endpoints

#### 1. Login with PIN
```http
POST /delivery/login
Content-Type: application/json

{
  "phoneNumber": "9876543210",
  "pin": 123456
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "name": "John Doe",
    "phoneNumber": "9876543210",
    "isAvailable": false
  }
}
```

#### 2. Request OTP
```http
POST /delivery/request-otp
Content-Type: application/json

{
  "phoneNumber": "9876543210"
}
```

#### 3. Verify OTP
```http
POST /delivery/verify-otp
Content-Type: application/json

{
  "phoneNumber": "9876543210",
  "otp": 123456
}
```

### Order Management Endpoints (Protected)

All order endpoints require authentication header:
```
Authorization: Bearer <token>
```

#### 1. Get Active Orders
```http
GET /delivery/orders/active
```

#### 2. Get Order History
```http
GET /delivery/orders/history?limit=20&offset=0
```

#### 3. Get Order Details
```http
GET /delivery/orders/:id
```

#### 4. Accept Order
```http
POST /delivery/orders/:id/accept
```

#### 5. Update Order Status
```http
POST /delivery/orders/:id/status
Content-Type: application/json

{
  "status": "in_transit",
  "latitude": 12.9716,
  "longitude": 77.5946
}
```

**Valid Status Transitions:**
- `pending` → `picked_up`
- `picked_up` → `in_transit`
- `in_transit` → `delivered`

#### 6. Complete Delivery
```http
POST /delivery/orders/:id/complete
Content-Type: application/json

{
  "latitude": 12.9716,
  "longitude": 77.5946,
  "signature": "base64_signature_data",
  "notes": "Delivered successfully"
}
```

### Profile Endpoints (Protected)

#### 1. Get Profile
```http
GET /delivery/profile
```

#### 2. Update Profile
```http
PUT /delivery/profile
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "aadharNumber": "123456789012",
  "panNumber": "ABCDE1234F",
  "drivingLicense": "DL1234567890",
  "vehicleNumber": "KA01AB1234",
  "vehicleType": "bike",
  "bankAccountNumber": "1234567890",
  "ifsc": "ABCD0001234"
}
```

#### 3. Update Location
```http
POST /delivery/location
Content-Type: application/json

{
  "latitude": 12.9716,
  "longitude": 77.5946
}
```

#### 4. Toggle Availability
```http
POST /delivery/availability
Content-Type: application/json

{
  "isAvailable": true
}
```

### Earnings Endpoints (Protected)

#### 1. Get Earnings
```http
GET /delivery/earnings?period=week
```

**Query Parameters:**
- `period`: `today`, `week`, or `month`

#### 2. Get Earnings History
```http
GET /delivery/earnings/history?limit=50&offset=0
```

## Database Schema

### Collections

#### delivery_partners
```javascript
{
  _id: ObjectId,
  name: String,
  phoneNumber: String,
  email: String,
  otp: Number,
  numberOfRetriesOTP: Number,
  otpGeneratedAt: Date,
  pin: Number,
  numberOfRetriesPIN: Number,
  isAvailable: Boolean,
  isVerified: Boolean,
  rating: Number,
  totalDeliveries: Number,
  lastLoginAt: Date,
  createdAt: Date,
  updatedAt: Date,
  aadharNumber: String,
  panNumber: String,
  drivingLicense: String,
  vehicleNumber: String,
  vehicleType: String,
  bankAccountNumber: String,
  ifsc: String,
  currentLatitude: Number,
  currentLongitude: Number,
  lastLocationAt: Date
}
```

#### deliveries
```javascript
{
  _id: ObjectId,
  orderId: String,
  partnerId: String,
  customerId: String,
  customerName: String,
  customerPhone: String,
  warehouseId: String,
  status: String,
  pickupAddress: String,
  deliveryAddress: String,
  pickupLatitude: Number,
  pickupLongitude: Number,
  deliveryLatitude: Number,
  deliveryLongitude: Number,
  distance: Number,
  orderAmount: Number,
  deliveryFee: Number,
  itemsCount: Number,
  items: Array,
  assignedAt: Date,
  pickedUpAt: Date,
  inTransitAt: Date,
  deliveredAt: Date,
  createdAt: Date,
  updatedAt: Date,
  paymentMethod: String,
  notes: String
}
```

#### earnings
```javascript
{
  _id: ObjectId,
  partnerId: String,
  deliveryId: String,
  orderId: String,
  amount: Number,
  deliveryFee: Number,
  bonus: Number,
  totalEarning: Number,
  earnedAt: Date,
  createdAt: Date
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8081` |
| `MONGO_URI` | MongoDB connection string | Required |
| `JWT_SECRET` | Secret key for JWT tokens | Required |

## Testing

### Run Tests
```bash
go test ./...
```

### Run Tests with Coverage
```bash
go test -cover ./...
```

### Run Specific Test
```bash
go test -run TestFunctionName ./path/to/package
```

## Deployment

### Using Docker

1. Build Docker image:
```bash
docker build -t delivery-backend .
```

2. Run container:
```bash
docker run -p 8081:8081 \
  -e MONGO_URI=mongodb://mongo:27017/espaze_delivery \
  -e JWT_SECRET=your-secret \
  delivery-backend
```

### Using Docker Compose

```bash
docker-compose up -d
```

## Security Considerations

1. **JWT Secret**: Use a strong, random secret key in production
2. **PIN Storage**: Currently stored as integer; consider hashing in production
3. **CORS**: Configure CORS properly for production environment
4. **Rate Limiting**: Implement rate limiting for authentication endpoints
5. **Input Validation**: All inputs are validated using Gin's binding
6. **HTTPS**: Always use HTTPS in production

## Performance Optimization

1. **Database Indexes**: Create indexes on frequently queried fields
```javascript
db.delivery_partners.createIndex({ phoneNumber: 1 })
db.deliveries.createIndex({ partnerId: 1, status: 1 })
db.earnings.createIndex({ partnerId: 1, earnedAt: -1 })
```

2. **Connection Pooling**: MongoDB driver uses connection pooling by default
3. **Caching**: Consider implementing Redis for session management

## Troubleshooting

### MongoDB Connection Issues
```bash
# Check MongoDB is running
mongosh --eval "db.adminCommand('ping')"

# Check connection string
echo $MONGO_URI
```

### Port Already in Use
```bash
# Find process using port 8081
lsof -i :8081

# Kill the process
kill -9 <PID>
```

## Contributing

1. Create a feature branch
2. Make your changes
3. Write tests
4. Submit a pull request

## License

Proprietary - Espaze

## Support

For issues and support:
- Create an issue in the repository
- Contact the development team
- Check internal documentation

## Changelog

### Version 1.0.0
- Initial release
- Authentication with PIN and OTP
- Order management
- Earnings tracking
- Profile management
- Location services


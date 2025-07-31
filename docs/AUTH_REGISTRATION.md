# Google OAuth Registration API

This document describes the Google OAuth registration endpoint implementation for the Sprintaza backend.

## Overview

The registration feature allows users to create accounts using Google OAuth. When a user completes Google OAuth on the frontend, the email and name data are sent to this backend endpoint to create or authenticate the user.

## API Endpoint

### Register with Google OAuth

**Endpoint:** `POST /api/v1/auth/register`

**Description:** Creates a new user account or authenticates an existing user using Google OAuth data.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "user@example.com",
  "name": "John Doe"
}
```

**Parameters:**
- `email` (string, required): User's email from Google OAuth
- `name` (string, required): User's name from Google OAuth (2-100 characters)

#### Response

**Success Response (New User):**
- **Status Code:** `201 Created`
```json
{
  "success": true,
  "data": {
    "message": "User registered successfully",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com",
      "avatar_url": ""
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "is_new_user": true
  }
}
```

**Success Response (Existing User):**
- **Status Code:** `200 OK`
```json
{
  "success": true,
  "data": {
    "message": "User already exists, logged in successfully",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "user@example.com",
      "avatar_url": ""
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "is_new_user": false
  }
}
```

**Error Response (Validation Error):**
- **Status Code:** `400 Bad Request`
```json
{
  "error": "Validation failed",
  "message": "Email is required",
  "field": "email"
}
```

**Error Response (Server Error):**
- **Status Code:** `500 Internal Server Error`
```json
{
  "error": "Registration failed",
  "message": "Unable to process registration at this time"
}
```

## JWT Token

The API returns a JWT token that contains:
- `user_id`: User's unique identifier
- `email`: User's email address
- `exp`: Token expiration time (7 days from issue)
- `iat`: Token issued at time
- `jti`: Unique token identifier

**Token Usage:**
Include the token in subsequent API requests:
```
Authorization: Bearer <token>
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(100),
    user_xp_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### User XP Table
```sql
CREATE TABLE user_xp (
    user_xp_id SERIAL PRIMARY KEY,
    total_xp INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## Environment Variables

Required environment variables:

```env
# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here-change-this-in-production
JWT_EXPIRES_IN=24h

# Database Configuration
DATABASE_URI_DEV=postgres://username:password@localhost:5432/sprintaza_dev?sslmode=disable
DATABASE_URI_PROD=postgres://username:password@localhost:5432/sprintaza_prod?sslmode=disable

# Application Configuration
GO_PORT=8080
GO_ENV=development
```

## Architecture

The implementation follows a clean architecture pattern:

```
handlers/auth-handlers/register/
├── register.go          # HTTP handler layer
controllers/auth-controllers/register/
├── input.go            # Input validation structures
├── service.go          # Business logic layer
└── repository.go       # Data access layer
```

### Components:

1. **Handler (`register.go`)**: Handles HTTP requests and responses
2. **Service (`service.go`)**: Contains business logic for user registration
3. **Repository (`repository.go`)**: Manages database operations
4. **Input (`input.go`)**: Defines and validates input structures

## Security Features

1. **Password Security**: Random secure passwords are generated for OAuth users
2. **JWT Security**: Tokens include unique identifiers and expiration times
3. **Input Validation**: All inputs are validated before processing
4. **Database Transactions**: User creation uses database transactions for consistency
5. **Logging**: Comprehensive logging for security monitoring

## Testing

Run the tests with:
```bash
go test -v ./...
```

Test coverage includes:
- Valid registration scenarios
- Input validation
- Existing user handling
- Error cases

## Usage Example

### Frontend Integration

```javascript
// After Google OAuth success
const googleOAuthData = {
  email: googleUser.email,
  name: googleUser.name
};

const response = await fetch('/api/v1/auth/register', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(googleOAuthData)
});

const result = await response.json();

if (result.success) {
  // Store JWT token
  localStorage.setItem('token', result.data.token);
  
  // Check if new user for onboarding
  if (result.data.is_new_user) {
    // Redirect to onboarding
  } else {
    // Redirect to dashboard
  }
}
```

## Health Check

**Endpoint:** `GET /api/v1/auth/health`

**Response:**
```json
{
  "status": "OK",
  "service": "Authentication Service",
  "timestamp": {},
  "version": "1.0.0"
}
```

## Error Handling

The API provides detailed error responses with appropriate HTTP status codes:

- `400 Bad Request`: Invalid input data
- `500 Internal Server Error`: Server-side errors
- `201 Created`: New user successfully created
- `200 OK`: Existing user successfully authenticated

All errors are logged with relevant context for debugging and monitoring.
# AppleShop e-commerce frontend API contract

The current Go backend exposes product endpoints only. The account-related routes below are frontend integration contracts for future backend implementation. All requests use same-origin cookies (`credentials: "same-origin"`); the frontend does not store passwords, session IDs, JWTs, or other credentials in `localStorage`.

## Existing endpoint

### Products

- `GET /api/v1/products?search={value}` — secure product listing/search.
- `GET /api/v1/vulnerable/products?search={value}` — vulnerable product listing/search.

These routes are unchanged and remain mode-aware through the existing API client.

## Authentication

Authentication routes use the normal `/api/v1` prefix and are not routed through the product vulnerable prefix.

### `POST /api/v1/auth/register`

Request:

```json
{
  "email": "user@example.test",
  "password": "test-only-password",
  "password_confirmation": "test-only-password"
}
```

Success: `201` or `200`, set an HttpOnly session cookie, and optionally return:

```json
{
  "user": {
    "id": 1,
    "email": "user@example.test"
  }
}
```

Errors: `400` validation error, `409` email already exists, `429` rate limited, `500` server error.

### `POST /api/v1/auth/login`

Request:

```json
{
  "email": "user@example.test",
  "password": "test-only-password"
}
```

Success: `200`, set or rotate an HttpOnly session cookie, and optionally return the user object. Errors: `400`, `401`, `429`, `500`.

### `POST /api/v1/auth/logout`

Request: no body. Success: `204` or `200`, invalidate the session and clear the cookie. Errors: `401`, `500`.

### `GET /api/v1/auth/me`

Request: no body; authentication comes from the HttpOnly cookie. Success: `200` with `{ "user": { "id": 1, "email": "user@example.test" } }`. Errors: `401` when unauthenticated, `500` on server failure.

## Account

### `GET /api/v1/account`

Returns the authenticated user’s account summary:

```json
{
  "user": { "id": 1, "email": "user@example.test" },
  "balance": 125000,
  "orders_count": 3,
  "recent_orders": [
    {
      "id": 42,
      "created_at": "2026-11-03T10:00:00Z",
      "status": "paid",
      "total": 99900
    }
  ]
}
```

Errors: `401`, `500`.

### `GET /api/v1/account/balance`

Success: `200` with `{ "balance": 125000 }`. This is display-only in the current frontend. No deposit, withdrawal, or payment operation is defined. Errors: `401`, `500`.

## Orders

### `GET /api/v1/orders`

Success: `200` with an array or `{ "orders": [] }`:

```json
[
  {
    "id": 42,
    "created_at": "2026-11-03T10:00:00Z",
    "status": "paid",
    "total": 99900
  }
]
```

Errors: `401`, `500`.

### `GET /api/v1/orders/{id}`

Success: `200` with an order or `{ "order": { ... } }`:

```json
{
  "id": 42,
  "created_at": "2026-11-03T10:00:00Z",
  "status": "paid",
  "total": 99900,
  "items": [
    {
      "product_id": 1,
      "name": "iPhone",
      "quantity": 1,
      "unit_price": 99900
    }
  ]
}
```

Errors: `401`, `403`, `404`, `500`.

## Frontend behavior when routes are unavailable

- Login and registration show the backend error or an integration message; they never claim success without a successful response.
- Account, balance, orders, and order detail show loading, empty, error, or `Backend integration pending` states.
- No account, order, balance, or authentication data is fabricated in the frontend.
- Password inputs are sent only over the request and are never persisted by the frontend.

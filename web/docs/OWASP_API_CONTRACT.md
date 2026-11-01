# OWASP Labs API contract

This document separates the endpoint already implemented by the Go backend from the endpoints reserved for future laboratory work. The frontend derives the mode prefix from the existing global mode switch; it does not independently decide whether a request is secure.

## Existing endpoint

### A05 — Injection: product search / SQL Injection

| Mode | Method | Path |
| --- | --- | --- |
| Secure | `GET` | `/api/v1/products?search={value}` |
| Vulnerable | `GET` | `/api/v1/vulnerable/products?search={value}` |

Request payload: the `search` query parameter, URL encoded. An empty value returns all products.

Success response (`200`): JSON array of product objects:

```json
[
  {
    "id": 1,
    "name": "iPhone",
    "description": "...",
    "price": 99900,
    "stock_quantity": 10,
    "category_id": 1,
    "image_url": "/images/iphone.jpg"
  }
]
```

Secure behavior: the repository uses a parameterized `ILIKE $1` query. Vulnerable behavior: the repository interpolates the search value into the SQL string. The frontend displays the actual HTTP status and body and makes no claim about exploit success.

Expected errors: non-2xx response from the API, with the backend’s existing JSON error format where available. The frontend must preserve and display the status/body rather than inventing a result.

## Proposed endpoints

The following contracts are placeholders for backend implementation. Until the routes exist, the corresponding UI remains disabled and reports **Backend integration pending**. For every proposed route, the `{modePrefix}` is empty in Secure Mode and `/vulnerable` in Vulnerable Mode.

Common success envelope:

```json
{
  "result": "protected|vulnerable|blocked|allowed|error",
  "message": "Human-readable educational result",
  "details": {},
  "request_id": "optional-request-id"
}
```

Common error responses: `400` for malformed input, `401` for unauthenticated requests, `403` for denied access, `404` for an unknown resource, `409` for an invalid state, and `500` for an unhandled server error. Error bodies should be JSON with `message`, optional `code`, and optional `request_id`. Secure and vulnerable implementations may intentionally differ in result and details, but must keep the response shape stable.

### A01 — Broken Access Control

`POST /api/v1{modePrefix}/labs/access-control/check`

Request: `{ "resource_id": "1001", "subject": "demo-user", "action": "read" }`.

Response details: `{ "allowed": false, "resource_owner": "optional", "authorization_reason": "..." }`. Secure behavior enforces the subject’s server-side authorization; vulnerable behavior may demonstrate an object-level authorization flaw. Errors: `400`, `401`, `403`, `404`, `500`.

### A02 — Security Misconfiguration

`GET /api/v1{modePrefix}/labs/misconfiguration/headers?target=/api/v1/products`

Response details: `{ "target": "...", "headers": { "header-name": "value" }, "findings": [{ "name": "...", "severity": "low|medium|high" }] }`. Secure behavior returns the hardened configuration; vulnerable behavior exposes the intentionally configured weakness. Errors: `400`, `404`, `500`.

### A03 — Software Supply Chain Failures

`GET /api/v1{modePrefix}/labs/supply-chain/dependencies?package=example-package`

Response details: `{ "package": "...", "version": "...", "source": "...", "integrity": { "algorithm": "sha256", "expected": "...", "actual": "...", "valid": true } }`. Secure behavior validates approved provenance/integrity; vulnerable behavior demonstrates missing or bypassable validation. Errors: `400`, `404`, `500`.

### A04 — Cryptographic Failures

`POST /api/v1{modePrefix}/labs/cryptography/inspect`

Request: `{ "value": "test-only-value", "operation": "hash|encrypt|inspect" }`.

Response details: `{ "algorithm": "...", "protected": true, "exposed_value": false, "key_management": "..." }`. Secure behavior uses approved handling without returning secrets; vulnerable behavior demonstrates weak or exposed handling. Errors: `400`, `413`, `500`.

### A06 — Insecure Design

`POST /api/v1{modePrefix}/labs/design/cart-check`

Request: `{ "product_id": 1, "quantity": 1, "discount": "TEST10", "total": 100 }`.

Response details: `{ "accepted": false, "server_total": 100, "client_total": 100, "violations": [] }`. Secure behavior enforces server-side business invariants; vulnerable behavior demonstrates the specified design flaw. Errors: `400`, `409`, `500`.

### A07 — Authentication Failures

`POST /api/v1{modePrefix}/labs/authentication/check`

Request: `{ "username": "demo-user", "password": "test-only-password" }`.

Response details: `{ "authenticated": false, "session": { "cookie_flags": {}, "token_type": "..." }, "attempts_remaining": 0 }`. Never return plaintext credentials or real tokens. Secure behavior applies authentication/session controls; vulnerable behavior demonstrates only synthetic lab behavior. Errors: `400`, `401`, `429`, `500`.

### A08 — Software or Data Integrity Failures

`POST /api/v1{modePrefix}/labs/integrity/verify`

Request: `{ "artifact": "demo-release.tar", "digest": "sha256:...", "signature": "optional-test-signature" }`.

Response details: `{ "valid": false, "algorithm": "sha256", "expected": "...", "actual": "...", "signature_verified": false }`. Secure behavior rejects altered or untrusted data; vulnerable behavior demonstrates the missing verification. Errors: `400`, `404`, `422`, `500`.

### A09 — Security Logging and Alerting Failures

`GET /api/v1{modePrefix}/labs/logging/events?event=failed_login`

Response details: `{ "events": [{ "timestamp": "...", "event": "...", "severity": "...", "alerted": false, "redacted": true }] }`. Secure behavior emits useful, protected events and alerts; vulnerable behavior demonstrates missing, excessive, or ineffective detection. Errors: `400`, `500`.

### A10 — Mishandling of Exceptional Conditions

`POST /api/v1{modePrefix}/labs/exceptions/trigger`

Request: `{ "scenario": "timeout|dependency_error|invalid_state" }`.

Response details: `{ "handled": true, "fail_mode": "closed|open", "status": 503, "error_id": "...", "sensitive_details_exposed": false }`. Secure behavior fails closed and returns bounded diagnostics; vulnerable behavior demonstrates unsafe exception handling. Errors: `400`, `500`, with no stack traces or secrets in the response.

## PortSwigger topic contracts

These are proposed frontend integration contracts for the PortSwigger topic pages. They are not currently implemented routes. The pages display them as `Backend Integration Pending` and never send requests until the project owner adds the corresponding Go handlers. `{modePrefix}` is empty in Secure Mode and `/vulnerable` in Vulnerable Mode.

| Topic | Method and path | Example request |
| --- | --- | --- |
| Authentication | `POST /api/v1{modePrefix}/labs/authentication/check` | `{ "username": "demo-user", "password": "test-only-password" }` |
| Path traversal | `GET /api/v1{modePrefix}/labs/path-traversal/read?path=images/iphone.jpg` | `path=images/iphone.jpg` |
| Command injection | `POST /api/v1{modePrefix}/labs/command-injection/execute` | `{ "host": "example.test" }` |
| Business logic vulnerabilities | `POST /api/v1{modePrefix}/labs/design/cart-check` | `{ "product_id": 1, "quantity": 1, "discount": "TEST10", "total": 100 }` |
| Information disclosure | `GET /api/v1{modePrefix}/labs/information-disclosure/inspect?resource=product-details` | `resource=product-details` |
| Access control | `POST /api/v1{modePrefix}/labs/access-control/check` | `{ "resource_id": "1001", "subject": "demo-user", "action": "read" }` |
| File upload vulnerabilities | `POST /api/v1{modePrefix}/labs/file-upload/inspect` | `{ "filename": "avatar.png", "content_type": "image/png" }` |
| Race conditions | `POST /api/v1{modePrefix}/labs/race-conditions/check` | `{ "operation": "redeem", "token": "TEST10" }` |
| Server-side request forgery (SSRF) | `POST /api/v1{modePrefix}/labs/ssrf/fetch` | `{ "url": "https://example.test/health" }` |
| XXE injection | `POST /api/v1{modePrefix}/labs/xxe/parse` | `<order><item>demo</item></order>` |
| NoSQL injection | `POST /api/v1{modePrefix}/labs/nosql-injection/query` | `{ "username": "demo-user", "password": "test-only-password" }` |
| API testing | `POST /api/v1{modePrefix}/labs/api-testing/inspect` | `{ "path": "/api/v1/products", "method": "GET" }` |

For these proposed routes, the response should use the common envelope above and return the real HTTP status/body from the backend. Secure behavior should enforce the relevant control; vulnerable behavior may demonstrate the intentionally unsafe implementation. Expected errors are the common `400`, `401`, `403`, `404`, `409`, `422`, `429`, and `500` responses appropriate to the operation. No response is fabricated by the frontend.

## Frontend integration rules

The shared API module owns the mode prefix. Individual lab views must use the central contract and must not hardcode secure/vulnerable decisions. A missing route is not a successful exploit: the UI must retain the pending state until the backend is available.

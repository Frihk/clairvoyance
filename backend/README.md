# Backend — ProofPass API

Go / Gin REST API. Handles auth, credential issuance, blockchain writes, and verification. Deployed on Railway.

---

## Requirements

- Go 1.25 or later
- PostgreSQL database (a Supabase project works directly)
- Access to the Polygon Amoy RPC endpoint (see `.env` below)
- The deployed ProofPass smart contract address

---

## Setup

```bash
cd backend
cp .env.example .env   # fill in the values described below
go mod download
go run cmd/server/main.go
```

The server listens on the port defined by `PORT` (default `8080`). A `GET /health` endpoint confirms the server is running.

---

## Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | yes | Full Postgres connection string, e.g. `postgres://user:pass@host:5432/db` |
| `JWT_SECRET` | yes | Random string used to sign HS256 tokens. At least 32 characters. |
| `PORT` | no | Defaults to `8080` |
| `APP_ENV` | no | `development` (debug logging) or `production` (warn-only). Defaults to `development`. |
| `JWT_EXPIRY_MINUTES` | no | Access token lifetime in minutes. Defaults to `10`. |
| `RPC_URL` | yes (chain) | Polygon Amoy JSON-RPC endpoint |
| `CONTRACT_ADDRESS` | yes (chain) | Deployed ProofPass contract address |
| `ISSUER_EMAIL` | yes (demo) | Hardcoded issuer account email |
| `ISSUER_PASS` | yes (demo) | Hardcoded issuer account password (stored as bcrypt hash in DB at seed time) |

`DATABASE_URL` is a required field — the server will not start without it. The `.env` file ships with the key named `DB_URL`; rename it to `DATABASE_URL` to match `config.go`.

---

## Project structure

```
backend/
├── cmd/server/main.go              Entry point. Loads config, starts Gin.
├── go.mod
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── auth.go             Register, Login, Refresh, Logout, Profile
│   │   ├── middleware/
│   │   │   ├── auth.go             AuthRequired — validates access tokens, sets user_id in context
│   │   │   └── cors.go             CORS headers, OPTIONS short-circuit
│   │   └── router/
│   │       └── router.go           Route registration
│   │
│   ├── auth/
│   │   ├── jwt.go                  GenerateTokenPair, ValidateAccessToken, ValidateRefreshToken
│   │   └── password.go             bcrypt hash and check
│   │
│   ├── blockchain/                 go-ethereum client (stub — Person 01)
│   │
│   ├── config/
│   │   └── config.go               Reads env vars into a typed Config struct
│   │
│   ├── db/                         Database connection (stub)
│   │
│   ├── models/
│   │   ├── users.go                User struct + BeforeCreate hook
│   │   └── credentials.go          Credential struct
│   │
│   ├── repositories/
│   │   ├── auth.go                 UserRepository — CRUD for users
│   │   └── credentials.go          CredentialRepository — CRUD + hash existence check
│   │
│   ├── services/
│   │   ├── auth.go                 RegisterUser, LoginUser, RefreshTokens, GetAuthedUser
│   │   └── credentials.go          Issue, Verify, Mine, ByIssuer + ChainClient interface + MockChainClient
│   │
│   └── utils/
│       ├── errors.go               Sentinel errors (ErrNotFound, ErrUnauthorized, etc.)
│       ├── logger.go               slog initialisation
│       └── responses.go            SuccessResponse, ErrorResponse, PaginatedResponse
```

---

## API routes

All routes are prefixed `/api`.

### Auth — `/api/auth`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/wallet/connect` | Public | Verify MetaMask signature, upsert user, return token pair |
| `POST` | `/email/login` | Public | Hardcoded issuer check, return token pair |
| `POST` | `/refresh` | Public | Accept refresh token, return new token pair |
| `POST` | `/logout` | JWT | Client discards tokens |
| `GET` | `/me` | JWT | Return current user from token |

### Credentials — `/api/credentials`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/issue` | Issuer JWT | Hash fields, write to chain, persist, return credential + QR URL |
| `GET` | `/verify/:id` | Public | Re-hash DB data, compare on-chain, return `verified` / `tampered` / `not_found` |
| `GET` | `/my` | JWT | All credentials for the authenticated user |

### Teams — `/api/teams`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/:slug` | Public | Team + members + contribution percentages (seeded data) |
| `POST` | `/create` | JWT | Create team, caller becomes leader |
| `POST` | `/:id/members` | Leader JWT | Add member with role and contribution percentage |

### Utility — `/api`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/stats` | Public | Total credentials, verifications, institutions |
| `GET` | `/health` | Public | Server up, DB connected, chain RPC responding |

---

## Authentication design

The API uses short-lived access tokens (10 minutes by default) and long-lived refresh tokens (7 days). Both are HS256 JWTs. The `token_type` claim (`access` or `refresh`) is validated in middleware — a refresh token cannot be used as a bearer token on protected routes.

The frontend should call `POST /api/auth/refresh` proactively around the 9-minute mark while the user is active. If the user is idle for longer than the access token lifetime, the natural consequence is sign-out on the next request.

---

## Credential hash

The canonical hash used for both issuance and verification is:

```
SHA-256( recipient_ref | title | credential_type | issue_date (RFC3339 UTC) | skills (comma-joined) )
```

Field order is fixed. Changing it would invalidate all existing on-chain records. The implementation lives in `services/credentials.go` as `HashCredential`.

---

## Blockchain interface

The `ChainClient` interface in `services/credentials.go` decouples the credential service from the go-ethereum implementation:

```go
type ChainClient interface {
    IssueCredential(credentialID uuid.UUID, dataHash string) (txHash string, blockNumber int64, err error)
    GetCredentialHash(credentialID uuid.UUID) (string, error)
}
```


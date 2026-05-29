# ProofPass

ProofPass lets institutions issue tamper-proof credentials on Polygon Amoy and lets anyone — no account, no app — verify them instantly by scanning a QR code or opening a link.

---

## Why it exists

Academic certificates and skill badges are easy to fake. Verifying them today means emailing registrars, waiting days, and hoping someone replies. ProofPass removes that friction by anchoring a SHA-256 hash of every credential onto a public blockchain at the moment of issuance. The full credential data lives in a database; the hash lives on-chain. At verification time the backend re-computes the hash from stored data and compares it against what the chain recorded. If they match, the credential is genuine. If they diverge, the credential has been altered. The result is always one of three unambiguous states: **verified**, **tampered**, or **not found** — returned in under a second, from any device.

The QR code on a certificate points to `proofpass.io/verify/:id`. Scanning it requires nothing more than a camera. No wallet, no login, no extension.

---

## Who uses it

| Role | What they do |
|------|-------------|
| Issuer | Institution staff (university registrar, bootcamp admin). Logs in with email, fills the issue form, clicks sign. One action — credential is on-chain. |
| Holder | Graduate or participant. Views credentials on their dashboard, downloads QR codes, shares a portfolio link covering all their credentials at once. |
| Verifier | Recruiter, judge, anyone. Opens a URL or scans a QR. No account required. |

---

## Repository layout

```
clairvoyance/
├── backend/          Go / Gin API server
└── frontend/         React single-page application
```

Each directory contains its own README with full setup instructions.

---

## How the core flow works

```
Issuer submits form
        |
        v
Backend computes SHA-256(recipient | title | type | issueDate | skills)
        |
        v
issueCredential(id, hash) written to Polygon Amoy smart contract
        |
        v
Full credential record + tx_hash + block_number stored in Postgres
        |
        v
QR code URL returned: proofpass.io/verify/:id
```

At verify time the same hash computation runs over the stored database record and is compared against what `verifyCredential(id)` returns from the chain. The comparison is the entire verification logic — no trusted intermediary involved.

---

## Tech stack

| Layer | Technology |
|-------|-----------|
| API | Go 1.25, Gin, GORM |
| Database | PostgreSQL (Supabase) |
| Blockchain | Polygon Amoy testnet, go-ethereum |
| Auth | JWT (HS256), access + refresh token pair, MetaMask signature verification |
| Frontend | React, Vite, wagmi |
| Deploy | Railway (backend), Vercel (frontend) |

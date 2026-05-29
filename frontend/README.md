# Frontend — ProofPass

React single-page application. Four pages: landing, dashboard, issue credential, and verify. Deployed on Vercel.

---

## Requirements

- Node.js 18 or later
- A running instance of the ProofPass backend API (see `backend/README.md`)
- MetaMask browser extension (optional — email login is the fallback for demo)

---

## Setup

```bash
cd frontend
cp .env.example .env.local   # set VITE_API_URL
npm install
npm run dev
```

The dev server runs on `http://localhost:5173` by default.

---

## Environment variables

| Variable | Description |
|----------|-------------|
| `VITE_API_URL` | Base URL of the backend API, e.g. `http://localhost:8080/api` |

---

## Pages and routes

| Route | Auth required | Description |
|-------|--------------|-------------|
| `/` | No | Landing page. Hero, stats from `GET /api/stats`, how-it-works steps, quick verify input. |
| `/dashboard` | JWT | Credential list for the authenticated user. QR code modal per credential. Portfolio share link. |
| `/issue` | Issuer JWT | Issue form. Five-step signing animation while the backend writes to chain. Success screen with QR. |
| `/verify/:id` | No | Public verification result. Three states: verified (green), tampered (red), not found (grey). |
| `/profile/:slug` | No | Public portfolio. All credentials for one person, aggregated skills panel, single shareable QR. |

---

## Authentication flow

On load, the app checks `localStorage` for a JWT. If none is present, the user is redirected to `/` from any protected route.

Two login paths exist:

**Wallet** — uses wagmi's `useConnect()`. If no injected wallet is detected, the MetaMask button is hidden and email login is shown with an explanatory note.

**Email** — sends `POST /api/auth/email/login` with the hardcoded issuer credentials. Returns an access token and refresh token. The access token expires in 10 minutes; the app calls `POST /api/auth/refresh` proactively at the 9-minute mark to keep the session alive.

Role is read from the decoded JWT. The Issue page and the Issue nav link are only rendered when `user.role === 'issuer'`.

---

## Key components

**VerifyPage** — the most important page from a demo perspective. Called by judges who scan a QR code on a phone. It must work without any authentication and must always resolve to one of three states. Never show a spinner that does not resolve.

**QR code** — rendered with `qrcode.react`. Use `QRCodeSVG` for display and `QRCodeCanvas` (hidden) for download. The download triggers `canvas.toDataURL('image/png')` and a programmatic anchor click.

**Portfolio page** — a single URL, `proofpass.io/profile/:slug`, that shows all of a person's credentials. The slug is derived from the holder's first name in lowercase. Credentials are expandable; each has a link to the individual verify page.

---

## Wallet detection

```js
// Hide MetaMask button when no injected provider is found
const { connectors } = useConnect()
const hasInjected = connectors.some(c => c.type === 'injected')
```

If `hasInjected` is false, render only the email login button with a note: "No MetaMask? Use Email login."

---

## Design tokens

The UI matches the spec's design system. All values live in a central token object in the application. The key ones are:

| Token | Value | Use |
|-------|-------|-----|
| `bg` | `#07090F` | Page background |
| `bg2` | `#111620` | Card background |
| `acc` | `#5B8DEF` | Primary interactive colour |
| `teal` | `#00C9A7` | Confirmed / success accent |
| `amber` | `#F5A623` | Issuer badge, warnings |
| `red` | `#FF4D6D` | Tampered state, destructive |
| `green` | `#22C55E` | Verified state |
| `mono` | IBM Plex Mono | Hashes, labels, tags |
| `display` | Bricolage Grotesque | Headings |
| `body` | DM Sans | Body copy |

Fonts are loaded from Google Fonts. All three families must be present or the layout will shift.

---

## Verify page behaviour

The verify page is a public entry point. Judges will open it on a phone by scanning a QR. It must:

- require no login
- always resolve (never hang on a spinner)
- show one of exactly three outcomes: verified, tampered, not found
- be readable on a narrow viewport (max-width 480px centered)

The tampered state must show both the recomputed hash and the on-chain hash so the divergence is visible.
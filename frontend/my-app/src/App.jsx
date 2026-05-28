import { useState, useEffect, useCallback } from "react";
import { useAccount, useConnect, useDisconnect } from "wagmi";

// ── Design Tokens (from spec) ──────────────────────────────────────────
const C = {
  bg: "#07090F", bg1: "#0C1018", bg2: "#111620", bg3: "#182030",
  border: "#1C2740", border2: "#243050",
  acc: "#5B8DEF", teal: "#00C9A7", amber: "#F5A623",
  red: "#FF4D6D", green: "#22C55E", purple: "#A78BFA",
  white: "#F0F4FF", text: "#B8C4DE", dim: "#5A6A8A",
};

const API_BASE_URL = (import.meta.env.VITE_P2_API_URL || import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

// ── Seed Credentials ───────────────────────────────────────────────────
const SEED_CREDENTIALS = [
  {
    id: "cred_alice_001",
    title: "BSc. Computer Science",
    type: "degree",
    recipient: "Alice Wanjiku",
    recipientEmail: "alice@kenyatta.edu",
    issuer: "Kenyatta University",
    issuerWallet: "0x7f3a...c91e",
    issueDate: "2024-03-15",
    description: "4-year programme, graduated with first class honours.",
    skills: ["Solidity", "Go", "Python"],
    txHash: "0x7f3a8b2c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a",
    blockNumber: "4,823,901",
    network: "Polygon Amoy",
    dataHash: "0x9a8b7c6d5e4f3a2b1c0d9e8f7a6b5c4d3e2f1a0b9c8d7e6f5a4b3c2d1e0f9a8b",
    status: "verified",
    emoji: "🎓",
    color: C.acc,
  },
  {
    id: "cred_brian_002",
    title: "Solidity Developer — BuildKE 2025",
    type: "certificate",
    recipient: "Brian Omondi",
    recipientEmail: "brian@buildke.io",
    issuer: "BuildKE",
    issuerWallet: "0x4d8e...f21c",
    issueDate: "2025-05-20",
    description: "Completed 12-week intensive Solidity & smart contract development bootcamp.",
    skills: ["Solidity", "Smart Contracts", "Hardhat"],
    txHash: "0x4d8ef21c9b0a7e6d5c4b3a2f1e0d9c8b7a6f5e4d3c2b1a0f9e8d7c6b5a4f3e2d",
    blockNumber: "4,891,234",
    network: "Polygon Amoy",
    dataHash: "0x3c2b1a0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2f1e0d9c8b7a6f5e4d3c2b",
    status: "verified",
    emoji: "⛓",
    color: C.teal,
  },
  {
    id: "cred_carol_003",
    title: "Smart Contract Auditor — Moringa School",
    type: "certificate",
    recipient: "Carol Njeri",
    recipientEmail: "carol@moringa.school",
    issuer: "Moringa School",
    issuerWallet: "0xa1b2...d4e5",
    issueDate: "2024-11-01",
    description: "Advanced security auditing course covering reentrancy, overflow, and access control vulnerabilities.",
    skills: ["Auditing", "Security", "EVM"],
    txHash: "0xa1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1d2e3f4a5b6c7d8e9f0a1b2",
    blockNumber: "4,756,012",
    network: "Polygon Amoy",
    dataHash: "0x5e4d3c2b1a0f9e8d7c6b5a4f3e2d1c0b9a8f7e6d5c4b3a2f1e0d9c8b7a6f5e4d",
    status: "verified",
    emoji: "💼",
    color: C.purple,
  },
];

// Portfolio slug → person name mapping (for public profile links)
const PROFILE_SLUGS = {
  "alice": "Alice Wanjiku",
  "brian": "Brian Omondi",
  "carol": "Carol Njeri",
};

const TEAM_MEMBERS = [
  { slug: "alice", name: "Alice Wanjiku", role: "Protocol Engineer", focus: "Smart contracts", contribution: 86, color: C.acc },
  { slug: "brian", name: "Brian Omondi", role: "Full-stack Builder", focus: "Credential workflows", contribution: 72, color: C.teal },
  { slug: "carol", name: "Carol Njeri", role: "Security Reviewer", focus: "Verification quality", contribution: 64, color: C.purple },
];

const GITHUB_SYNC = {
  alice: {
    handle: "alicewanjiku",
    lastSynced: "2026-05-28",
    repos: [
      { name: "proofpass-contracts", commits: 34, prs: 8, contribution: 92 },
      { name: "credential-registry", commits: 18, prs: 5, contribution: 74 },
    ],
  },
  brian: {
    handle: "brianomondi",
    lastSynced: "2026-05-28",
    repos: [
      { name: "proofpass-web", commits: 42, prs: 11, contribution: 88 },
      { name: "issuer-console", commits: 25, prs: 6, contribution: 69 },
    ],
  },
  carol: {
    handle: "carolnjeri",
    lastSynced: "2026-05-27",
    repos: [
      { name: "verification-audits", commits: 21, prs: 7, contribution: 81 },
      { name: "proofpass-contracts", commits: 12, prs: 4, contribution: 58 },
    ],
  },
};

// Hardcoded issuer accounts (from spec)
const ISSUER_ACCOUNTS = [
  { email: "issuer@kenyatta.edu", password: "demo123", name: "Kenyatta University", wallet: "0x7f3a...c91e", role: "issuer" },
  { email: "admin@buildke.io", password: "demo123", name: "BuildKE Registrar", wallet: "0x4d8e...f21c", role: "issuer" },
];

// Simulated SHA-256
function fakeHash(str) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) - hash) + str.charCodeAt(i);
    hash |= 0;
  }
  const hex = Math.abs(hash).toString(16).padStart(8, "0");
  return "0x" + hex.repeat(8).slice(0, 64);
}

function shortAddress(address) {
  if (!address) return "";
  return `${address.slice(0, 6)}...${address.slice(-4)}`;
}

function credentialVerifyUrl(id) {
  return `proofpass.io/verify/${id}`;
}

function statusBadge(status) {
  const normalized = String(status || "pending").toLowerCase();
  const badgeMap = {
    verified: { label: "✓ VERIFIED", color: C.green, bg: "rgba(34,197,94,0.1)", border: "rgba(34,197,94,0.3)" },
    pending: { label: "PENDING", color: C.amber, bg: "rgba(245,166,35,0.1)", border: "rgba(245,166,35,0.3)" },
    revoked: { label: "REVOKED", color: C.red, bg: "rgba(255,77,109,0.1)", border: "rgba(255,77,109,0.3)" },
    invalid: { label: "INVALID", color: C.red, bg: "rgba(255,77,109,0.1)", border: "rgba(255,77,109,0.3)" },
  };
  return badgeMap[normalized] || { label: normalized.toUpperCase(), color: C.dim, bg: "rgba(90,106,138,0.1)", border: "rgba(90,106,138,0.3)" };
}

function personSlug(name) {
  const entry = Object.entries(PROFILE_SLUGS).find(([, personName]) => personName === name);
  return entry?.[0] || name?.split(" ")[0]?.toLowerCase() || "";
}

function githubTotals(sync) {
  const repos = sync?.repos || [];
  return repos.reduce((totals, repo) => ({
    commits: totals.commits + repo.commits,
    prs: totals.prs + repo.prs,
    contribution: totals.contribution + repo.contribution,
  }), { commits: 0, prs: 0, contribution: 0 });
}

function GitHubRepoRows({ sync, color = C.teal }) {
  if (!sync) {
    return (
      <div style={{ fontSize: 12, color: C.dim, lineHeight: 1.6 }}>
        No GitHub account linked for this profile yet.
      </div>
    );
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
      {sync.repos.map(repo => (
        <div key={repo.name}>
          <div style={{ display: "flex", justifyContent: "space-between", gap: 10, marginBottom: 5 }}>
            <div style={{ ...styles.mono, fontSize: 11, color: C.white }}>{repo.name}</div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim }}>{repo.commits} commits · {repo.prs} PRs</div>
          </div>
          <div style={{ height: 6, background: C.bg3, border: `1px solid ${C.border}`, borderRadius: 20, overflow: "hidden" }}>
            <div style={{ height: "100%", width: `${repo.contribution}%`, background: color, borderRadius: 20 }} />
          </div>
        </div>
      ))}
    </div>
  );
}

async function issueCredential(payload) {
  const response = await fetch(`${API_BASE_URL}/credentials/issue`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  const text = await response.text();
  let body;
  try {
    body = text ? JSON.parse(text) : {};
  } catch {
    body = { message: text };
  }

  if (!response.ok) {
    throw new Error(body.message || body.error || `Issue request failed with status ${response.status}`);
  }

  return body;
}

function normalizeIssuedCredential(payload, responseBody) {
  const source = responseBody.credential || responseBody.data || responseBody;

  return {
    id: source.id || source.credentialId || payload.id,
    title: source.title || payload.title,
    type: source.type || payload.type,
    recipient: source.recipient || payload.recipient,
    recipientEmail: source.recipientEmail || source.recipient_email || payload.recipientEmail,
    issuer: source.issuer || payload.issuer,
    issuerWallet: source.issuerWallet || source.issuer_wallet || payload.issuerWallet,
    issueDate: source.issueDate || source.issue_date || payload.issueDate,
    description: source.description || payload.description,
    skills: source.skills || payload.skills,
    txHash: source.txHash || source.tx_hash || "",
    blockNumber: source.blockNumber || source.block_number || "",
    network: source.network || "Polygon Amoy",
    dataHash: source.dataHash || source.data_hash || payload.dataHash,
    status: source.status || "verified",
    emoji: payload.type === "degree" ? "🎓" : payload.type === "participation" ? "🏅" : "⛓",
    color: C.teal,
  };
}

// ── SVG QR Mock ───────────────────────────────────────────────────────
function QRCode({ value, size = 120 }) {
  const cells = [];
  const grid = 10;
  const cell = size / grid;
  for (let r = 0; r < grid; r++) {
    for (let c = 0; c < grid; c++) {
      const seed = (value.charCodeAt((r * grid + c) % value.length) + r * 7 + c * 13) % 2;
      const corner = (r < 3 && c < 3) || (r < 3 && c >= grid - 3) || (r >= grid - 3 && c < 3);
      const fill = corner ? true : seed === 0;
      cells.push(
        <rect key={`${r}-${c}`} x={c * cell + 1} y={r * cell + 1}
          width={cell - 2} height={cell - 2} rx={1}
          fill={fill ? "#000" : "#fff"} />
      );
    }
  }
  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} style={{ background: "#fff", borderRadius: 8 }}>
      {cells}
    </svg>
  );
}

// ── Base Styles ────────────────────────────────────────────────────────
const styles = {
  app: { background: C.bg, minHeight: "100vh", fontFamily: "'DM Sans', sans-serif", color: C.text, fontSize: 14 },
  nav: { background: C.bg1, borderBottom: `1px solid ${C.border}`, padding: "13px 28px", display: "flex", alignItems: "center", justifyContent: "space-between", position: "sticky", top: 0, zIndex: 50, backdropFilter: "blur(16px)" },
  brand: { fontFamily: "'Bricolage Grotesque', sans-serif", fontWeight: 800, fontSize: 16, color: C.white },
  card: { background: C.bg2, border: `1px solid ${C.border}`, borderRadius: 12, padding: 20 },
  input: { width: "100%", background: C.bg3, border: `1px solid ${C.border2}`, borderRadius: 8, padding: "9px 13px", color: C.white, fontSize: 13, fontFamily: "inherit", outline: "none", boxSizing: "border-box" },
  label: { display: "block", fontFamily: "'IBM Plex Mono', monospace", fontSize: 10, color: C.dim, letterSpacing: "1.5px", textTransform: "uppercase", marginBottom: 5 },
  btnBlue: { background: C.acc, color: "#fff", border: "none", borderRadius: 8, padding: "9px 18px", fontSize: 13, fontWeight: 600, cursor: "pointer", display: "inline-flex", alignItems: "center", gap: 6 },
  btnGhost: { background: "transparent", color: C.dim, border: `1px solid ${C.border2}`, borderRadius: 8, padding: "8px 16px", fontSize: 12, fontWeight: 500, cursor: "pointer", display: "inline-flex", alignItems: "center", gap: 6 },
  btnTeal: { background: C.teal, color: "#000", border: "none", borderRadius: 8, padding: "9px 18px", fontSize: 13, fontWeight: 700, cursor: "pointer", display: "inline-flex", alignItems: "center", gap: 6 },
  mono: { fontFamily: "'IBM Plex Mono', monospace" },
  display: { fontFamily: "'Bricolage Grotesque', sans-serif" },
};

function Tag({ color, bg, border, children }) {
  return (
    <span style={{ fontFamily: "'IBM Plex Mono', monospace", fontSize: 9, padding: "2px 8px", borderRadius: 20, letterSpacing: 1, color, background: bg, border: `1px solid ${border}` }}>
      {children}
    </span>
  );
}

// ── LANDING PAGE ────────────────────────────────────────────────────────
function Landing({ onLogin, onVerify, onPortfolio }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [err, setErr] = useState("");
  const [verifyId, setVerifyId] = useState("");
  const { address, chain, isConnected } = useAccount();
  const { connect, connectors, error: connectError, isPending } = useConnect();

  useEffect(() => {
    if (!isConnected || !address) return;

    onLogin({
      email: `${shortAddress(address)}@wallet.local`,
      name: "Connected Wallet",
      wallet: shortAddress(address),
      walletAddress: address,
      network: chain?.name || "Unknown Network",
      role: "holder",
    });
  }, [address, chain?.name, isConnected, onLogin]);

  function handleEmailLogin() {
    setErr("");
    const acct = ISSUER_ACCOUNTS.find(a => a.email === email && a.password === password);
    if (acct) { onLogin(acct); } else { setErr("Invalid credentials. Try issuer@kenyatta.edu / demo123"); }
  }

  function handleWalletLogin() {
    setErr("");

    if (isConnected && address) {
      onLogin({
        email: `${shortAddress(address)}@wallet.local`,
        name: "Connected Wallet",
        wallet: shortAddress(address),
        walletAddress: address,
        network: chain?.name || "Unknown Network",
        role: "holder",
      });
      return;
    }

    const connector = connectors.find(c => c.id === "metaMask") || connectors[0];
    if (!connector) {
      setErr("No wallet connector found. Install MetaMask or enable a browser wallet.");
      return;
    }

    connect({ connector });
  }

  const stats = [
    { label: "Credentials", value: SEED_CREDENTIALS.length + " issued", color: C.teal },
    { label: "Verifications", value: "8,903", color: C.acc },
    { label: "Institutions", value: "34", color: C.purple },
    { label: "Network", value: "Polygon Amoy", color: C.amber },
  ];

  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />

      {/* NAV */}
      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ display: "flex", gap: 24, alignItems: "center" }}>
          <span style={{ fontSize: 13, color: C.dim }}>How it Works</span>
          <button style={styles.btnGhost} onClick={() => onVerify(null)}>🔍 Verify</button>
          <button style={{ ...styles.btnBlue, fontSize: 12 }} onClick={() => handleWalletLogin()} disabled={isPending}>
            {isPending ? "Connecting..." : "🦊 Connect Wallet"}
          </button>
        </div>
      </nav>

      {/* HERO */}
      <div style={{ maxWidth: 980, margin: "0 auto", padding: "64px 32px 48px", display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(280px, 1fr))", gap: 32, alignItems: "center" }}>
        <div>
          <div style={{ ...styles.mono, fontSize: 10, color: C.teal, letterSpacing: 3, marginBottom: 16, textTransform: "uppercase" }}>
            Blockchain Identity · Polygon Amoy
          </div>
          <h1 style={{ ...styles.display, fontSize: 44, fontWeight: 800, color: C.white, lineHeight: 1.1, marginBottom: 14 }}>
            Verifiable Skills.<br />
            <span style={{ color: C.teal }}>Trusted Contributions.</span>
          </h1>
          <p style={{ color: C.dim, fontSize: 15, maxWidth: 460, margin: "0 0 28px", lineHeight: 1.7 }}>
            Blockchain-signed credentials anyone can verify instantly — no account, no app, no trust required.
          </p>
          <div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
            <button style={{ ...styles.btnTeal, padding: "11px 24px" }} onClick={() => onVerify(null)}>🔍 Verify a Credential</button>
            <button style={{ ...styles.btnGhost, padding: "11px 24px" }} onClick={() => onPortfolio("alice")}>View Sample Portfolio</button>
          </div>
        </div>

        <div style={{ ...styles.card, padding: 28 }}>
          <div style={{ ...styles.display, fontSize: 20, fontWeight: 700, color: C.white, marginBottom: 6 }}>Sign in</div>
          <div style={{ fontSize: 12, color: C.dim, marginBottom: 18 }}>Use MetaMask as a holder, or email/password as an issuer.</div>

          <button style={{ ...styles.btnBlue, width: "100%", justifyContent: "center", padding: "12px 18px", marginBottom: 12 }} onClick={handleWalletLogin} disabled={isPending}>
            {isPending ? "Connecting..." : "🦊 Connect MetaMask"}
          </button>

          {connectError && (
            <div style={{ ...styles.mono, fontSize: 11, color: C.red, marginBottom: 12 }}>
              {connectError.shortMessage || connectError.message}
            </div>
          )}

          <div style={{ display: "flex", alignItems: "center", gap: 10, margin: "16px 0" }}>
            <div style={{ height: 1, background: C.border, flex: 1 }} />
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, letterSpacing: 2 }}>OR</div>
            <div style={{ height: 1, background: C.border, flex: 1 }} />
          </div>

          <label style={styles.label}>Email</label>
          <input style={{ ...styles.input, marginBottom: 12 }} value={email} onChange={e => setEmail(e.target.value)} placeholder="issuer@kenyatta.edu" />
          <label style={styles.label}>Password</label>
          <input style={{ ...styles.input, marginBottom: 14 }} type="password" value={password} onChange={e => setPassword(e.target.value)} placeholder="demo123" onKeyDown={e => e.key === "Enter" && handleEmailLogin()} />
          {err && <div style={{ fontSize: 12, color: C.red, marginBottom: 12 }}>{err}</div>}
          <button style={{ ...styles.btnGhost, width: "100%", justifyContent: "center", color: C.white }} onClick={handleEmailLogin}>
            Sign in with Email
          </button>
          <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginTop: 12 }}>
            Demo issuer: issuer@kenyatta.edu / demo123
          </div>
        </div>
      </div>

      {/* STATS */}
      <div style={{ maxWidth: 900, margin: "0 auto", padding: "0 32px 40px" }}>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(4,1fr)", gap: 14, marginBottom: 40 }}>
          {stats.map(s => (
            <div key={s.label} style={{ ...styles.card, textAlign: "center" }}>
              <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, textTransform: "uppercase", marginBottom: 6 }}>{s.label}</div>
              <div style={{ ...styles.mono, fontSize: 22, fontWeight: 700, color: s.color }}>{s.value}</div>
            </div>
          ))}
        </div>

        {/* HOW IT WORKS */}
        <div style={{ ...styles.display, fontWeight: 700, fontSize: 15, color: C.white, marginBottom: 16 }}>How ProofPass works</div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3,1fr)", gap: 14 }}>
          {[
            { step: "01", color: C.acc, title: "Institution Issues", desc: "Issuer logs in, fills the form, signs. Hash written to Polygon Amoy." },
            { step: "02", color: C.teal, title: "Credential Stored", desc: "Full data stored in database. SHA-256 hash anchored on-chain." },
            { step: "03", color: C.green, title: "Anyone Verifies", desc: "Scan QR or open URL. Instant result. No login. No app. Any browser." },
          ].map(s => (
            <div key={s.step} style={styles.card}>
              <div style={{ ...styles.mono, fontSize: 10, color: s.color, marginBottom: 6 }}>STEP {s.step}</div>
              <div style={{ fontSize: 14, color: C.white, fontWeight: 600, marginBottom: 4 }}>{s.title}</div>
              <div style={{ fontSize: 12, color: C.dim, lineHeight: 1.6 }}>{s.desc}</div>
            </div>
          ))}
        </div>

        {/* Inline verify search */}
        <div style={{ ...styles.card, marginTop: 32, padding: "24px 28px" }}>
          <div style={{ ...styles.display, fontWeight: 700, fontSize: 15, color: C.white, marginBottom: 12 }}>🔍 Quick Verify a Credential</div>
          <div style={{ display: "flex", gap: 10 }}>
            <input style={{ ...styles.input, flex: 1 }} value={verifyId}
              onChange={e => setVerifyId(e.target.value)}
              placeholder="Enter credential ID, e.g. cred_alice_001" />
            <button style={styles.btnTeal} onClick={() => onVerify(verifyId || "cred_alice_001")}>Verify</button>
          </div>
          <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginTop: 8 }}>
            Try: cred_alice_001 · cred_brian_002 · cred_carol_003 · cred_fake_999
          </div>
          <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginTop: 6 }}>
            Or view a full portfolio:{" "}
            {["alice","brian","carol"].map((s,i) => (
              <span key={s}>{i > 0 && " · "}<span style={{ color: C.teal, cursor: "pointer", textDecoration: "underline" }} onClick={() => onPortfolio(s)}>{s}</span></span>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

// ── DASHBOARD ─────────────────────────────────────────────────────────
function Dashboard({ user, credentials, onVerify, onIssue, onPortfolio, onTeam, onLogout }) {
  const [qrModal, setQrModal] = useState(null);
  const [portfolioQr, setPortfolioQr] = useState(false);
  const [copiedId, setCopiedId] = useState(null);
  const myCredentials = user.role === "holder"
    ? credentials.filter(c => c.recipientEmail === user.email || c.recipient === user.name)
    : credentials;
  const dashboardMode = user.role === "issuer" ? "issued" : "received";

  // derive slug from name for portfolio URL
  const slug = personSlug(user.name);
  const githubSync = GITHUB_SYNC[slug];
  const githubStats = githubTotals(githubSync);

  async function copyCredentialQrLink(cred) {
    const link = credentialVerifyUrl(cred.id);
    await navigator.clipboard?.writeText(link);
    setCopiedId(cred.id);
    setTimeout(() => setCopiedId(current => current === cred.id ? null : current), 1400);
  }

  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />

      {/* NAV */}
      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ display: "flex", gap: 20, alignItems: "center" }}>
          <span style={{ fontSize: 13, color: C.white, fontWeight: 500, borderBottom: `2px solid ${C.acc}`, paddingBottom: 2 }}>Dashboard</span>
          <span style={{ fontSize: 13, color: C.dim, cursor: "pointer" }} onClick={onTeam}>Team</span>
          <span style={{ fontSize: 13, color: C.dim, cursor: "pointer" }} onClick={() => onPortfolio(slug)}>Portfolio</span>
          {user.role === "issuer" && (
            <button style={{ ...styles.btnBlue, fontSize: 12 }} onClick={onIssue}>⛓ Issue</button>
          )}
        </div>
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <Tag color={C.green} bg="rgba(34,197,94,0.1)" border="rgba(34,197,94,0.3)">● Connected</Tag>
          <span style={{ ...styles.mono, fontSize: 11, color: C.dim }}>{user.wallet}</span>
          {user.role === "issuer" && (
            <Tag color={C.amber} bg="rgba(245,166,35,0.1)" border="rgba(245,166,35,0.3)">ISSUER</Tag>
          )}
          <button style={{ ...styles.btnGhost, fontSize: 11 }} onClick={onLogout}>Sign Out</button>
        </div>
      </nav>

      <div style={{ maxWidth: 860, margin: "0 auto", padding: "36px 28px" }}>
        {/* Header */}
        <div style={{ display: "flex", alignItems: "flex-start", justifyContent: "space-between", marginBottom: 28 }}>
          <div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, letterSpacing: 2, marginBottom: 5, textTransform: "uppercase" }}>Welcome Back</div>
            <div style={{ ...styles.display, fontSize: 26, fontWeight: 800, color: C.white }}>{user.name}</div>
            <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginTop: 3 }}>{user.wallet}</div>
          </div>
          <div style={{ display: "flex", gap: 12 }}>
            {[
              { label: "Credentials", value: myCredentials.length, color: C.teal },
              { label: "Verified", value: myCredentials.filter(c => c.status === "verified").length, color: C.green },
            ].map(s => (
              <div key={s.label} style={{ ...styles.card, textAlign: "center", minWidth: 90, padding: "14px 20px" }}>
                <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginBottom: 4, textTransform: "uppercase" }}>{s.label}</div>
                <div style={{ ...styles.mono, fontSize: 28, fontWeight: 700, color: s.color }}>{s.value}</div>
              </div>
            ))}
          </div>
        </div>

        <div style={{ ...styles.card, marginBottom: 24, display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(260px, 1fr))", gap: 22, alignItems: "start" }}>
          <div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 2, marginBottom: 8, textTransform: "uppercase" }}>GitHub Sync</div>
            <div style={{ display: "flex", alignItems: "center", gap: 10, flexWrap: "wrap", marginBottom: 8 }}>
              <div style={{ ...styles.display, fontSize: 18, fontWeight: 700, color: C.white }}>
                {githubSync ? `@${githubSync.handle}` : "No handle linked"}
              </div>
              <Tag color={githubSync ? C.green : C.amber} bg={githubSync ? "rgba(34,197,94,0.1)" : "rgba(245,166,35,0.1)"} border={githubSync ? "rgba(34,197,94,0.3)" : "rgba(245,166,35,0.3)"}>
                {githubSync ? "SYNCED" : "PENDING"}
              </Tag>
            </div>
            <div style={{ fontSize: 12, color: C.dim, lineHeight: 1.6 }}>
              {githubSync
                ? `${githubStats.commits} commits and ${githubStats.prs} pull requests synced from ${githubSync.repos.length} repositories.`
                : "Connect GitHub to sync repository contributions into credential evidence."}
            </div>
            {githubSync && (
              <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginTop: 8 }}>
                Last synced {githubSync.lastSynced}
              </div>
            )}
          </div>
          <GitHubRepoRows sync={githubSync} color={C.acc} />
        </div>

        {/* Portfolio share banner */}
        {user.role !== "issuer" && myCredentials.length > 0 && (
          <div style={{ ...styles.card, borderColor: "rgba(0,201,167,0.3)", background: "rgba(0,201,167,0.04)", display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 24, padding: "14px 20px" }}>
            <div>
              <div style={{ fontSize: 13, color: C.white, fontWeight: 600, marginBottom: 2 }}>📋 Your Public Portfolio</div>
              <div style={{ ...styles.mono, fontSize: 11, color: C.dim }}>proofpass.io/profile/{slug} — all {myCredentials.length} credentials in one link</div>
            </div>
            <div style={{ display: "flex", gap: 8, flexShrink: 0 }}>
              <button style={{ ...styles.btnGhost, fontSize: 11 }} onClick={() => setPortfolioQr(true)}>📱 Portfolio QR</button>
              <button style={{ ...styles.btnTeal, fontSize: 11 }} onClick={() => onPortfolio(slug)}>View Portfolio</button>
            </div>
          </div>
        )}

        <div style={{ ...styles.display, fontWeight: 700, fontSize: 15, color: C.white, marginBottom: 14 }}>
          {dashboardMode === "issued" ? "Issued Credentials" : "Received Credentials"}
        </div>

        {myCredentials.length === 0 ? (
          <div style={{ ...styles.card, textAlign: "center", padding: 48, color: C.dim }}>
            No credentials yet.{" "}
            {user.role === "issuer" && <button style={{ ...styles.btnBlue, marginLeft: 8 }} onClick={onIssue}>Issue First Credential</button>}
          </div>
        ) : (
          <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
            {myCredentials.map(cred => (
              <div key={cred.id} style={{ ...styles.card, display: "flex", alignItems: "center", gap: 16, flexWrap: "wrap" }}>
                <div style={{ width: 46, height: 46, background: `${cred.color}22`, border: `1px solid ${cred.color}44`, borderRadius: 10, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 22, flexShrink: 0 }}>
                  {cred.emoji}
                </div>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div style={{ fontSize: 14, color: C.white, fontWeight: 600, marginBottom: 2 }}>{cred.title}</div>
                  <div style={{ fontSize: 11, color: C.dim }}>
                    {dashboardMode === "issued" ? "Issued to" : "Received from"}{" "}
                    <span style={{ color: C.teal }}>{dashboardMode === "issued" ? cred.recipient : cred.issuer}</span>
                    {" "}· {cred.issueDate} · <span style={styles.mono}>{cred.id}</span>
                  </div>
                  <div style={{ display: "flex", gap: 5, marginTop: 7, flexWrap: "wrap" }}>
                    {cred.skills.map(s => (
                      <Tag key={s} color={cred.color} bg={`${cred.color}18`} border={`${cred.color}33`}>{s}</Tag>
                    ))}
                  </div>
                </div>
                <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: 7, flexShrink: 0, marginLeft: "auto" }}>
                  {(() => {
                    const badge = statusBadge(cred.status);
                    return <Tag color={badge.color} bg={badge.bg} border={badge.border}>{badge.label}</Tag>;
                  })()}
                  <div style={{ display: "flex", gap: 6, flexWrap: "wrap", justifyContent: "flex-end" }}>
                    <button style={{ ...styles.btnGhost, fontSize: 11, padding: "4px 10px" }} onClick={() => setQrModal(cred)}>📱 QR Code</button>
                    <button style={{ ...styles.btnGhost, fontSize: 11, padding: "4px 10px" }} onClick={() => copyCredentialQrLink(cred)}>
                      {copiedId === cred.id ? "✓ Copied" : "Copy QR link"}
                    </button>
                    <button style={{ ...styles.btnGhost, fontSize: 11, padding: "4px 10px" }} onClick={() => onVerify(cred.id)}>🔍 Verify</button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {user.role === "issuer" && (
          <div style={{ marginTop: 24 }}>
            <button style={{ ...styles.btnBlue, width: "100%", justifyContent: "center", padding: 14 }} onClick={onIssue}>
              ⛓ Issue New Credential
            </button>
          </div>
        )}
      </div>

      {/* PORTFOLIO QR MODAL */}
      {portfolioQr && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.8)", zIndex: 200, display: "flex", alignItems: "center", justifyContent: "center" }}
          onClick={e => e.target === e.currentTarget && setPortfolioQr(false)}>
          <div style={{ ...styles.card, width: 340, padding: 32, textAlign: "center" }}>
            <div style={{ ...styles.mono, fontSize: 10, color: C.teal, letterSpacing: 3, marginBottom: 8, textTransform: "uppercase" }}>Public Portfolio</div>
            <div style={{ ...styles.display, fontSize: 18, fontWeight: 700, color: C.white, marginBottom: 2 }}>{user.name}</div>
            <div style={{ fontSize: 12, color: C.dim, marginBottom: 18 }}>{myCredentials.length} verified credentials</div>
            <div style={{ display: "flex", justifyContent: "center", marginBottom: 14 }}>
              <QRCode value={`proofpass.io/profile/${slug}`} size={148} />
            </div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginBottom: 16 }}>
              proofpass.io/profile/{slug}
            </div>
            <div style={{ display: "flex", gap: 8, justifyContent: "center" }}>
              <button style={styles.btnGhost} onClick={() => navigator.clipboard?.writeText(`proofpass.io/profile/${slug}`)}>Copy Link</button>
              <button style={styles.btnTeal} onClick={() => { setPortfolioQr(false); onPortfolio(slug); }}>Open Portfolio</button>
            </div>
          </div>
        </div>
      )}

      {/* QR MODAL */}
      {qrModal && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.8)", zIndex: 200, display: "flex", alignItems: "center", justifyContent: "center" }}
          onClick={e => e.target === e.currentTarget && setQrModal(null)}>
          <div style={{ ...styles.card, width: 320, padding: 32, textAlign: "center" }}>
            <div style={{ ...styles.display, fontSize: 16, fontWeight: 700, color: C.white, marginBottom: 4 }}>{qrModal.title}</div>
            <div style={{ fontSize: 12, color: C.dim, marginBottom: 20 }}>{qrModal.issuer}</div>
            <div style={{ display: "flex", justifyContent: "center", marginBottom: 14 }}>
              <QRCode value={`proofpass.io/verify/${qrModal.id}`} size={140} />
            </div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginBottom: 14 }}>
              proofpass.io/verify/{qrModal.id}
            </div>
            <div style={{ display: "flex", gap: 8, justifyContent: "center" }}>
              <button style={styles.btnGhost} onClick={() => navigator.clipboard?.writeText(credentialVerifyUrl(qrModal.id))}>Copy Link</button>
              <button style={styles.btnBlue} onClick={() => setQrModal(null)}>Close</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// ── TEAM PAGE ─────────────────────────────────────────────────────────
function TeamPage({ user, credentials, onBack, onPortfolio, onLogout }) {
  const members = TEAM_MEMBERS.map(member => {
    const memberCredentials = credentials.filter(c => c.recipient === member.name);
    const sync = GITHUB_SYNC[member.slug];
    const syncTotals = githubTotals(sync);
    return {
      ...member,
      github: sync,
      githubCommits: syncTotals.commits,
      githubPrs: syncTotals.prs,
      credentialCount: memberCredentials.length,
      verifiedCount: memberCredentials.filter(c => c.status === "verified").length,
    };
  });
  const totalCredentials = members.reduce((sum, member) => sum + member.credentialCount, 0);
  const totalCommits = members.reduce((sum, member) => sum + member.githubCommits, 0);
  const avgContribution = Math.round(members.reduce((sum, member) => sum + member.contribution, 0) / members.length);

  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />

      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ display: "flex", gap: 20, alignItems: "center" }}>
          <span style={{ fontSize: 13, color: C.dim, cursor: "pointer" }} onClick={onBack}>Dashboard</span>
          <span style={{ fontSize: 13, color: C.white, fontWeight: 500, borderBottom: `2px solid ${C.acc}`, paddingBottom: 2 }}>Team</span>
        </div>
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <Tag color={C.green} bg="rgba(34,197,94,0.1)" border="rgba(34,197,94,0.3)">● Connected</Tag>
          <span style={{ ...styles.mono, fontSize: 11, color: C.dim }}>{user.wallet}</span>
          <button style={{ ...styles.btnGhost, fontSize: 11 }} onClick={onLogout}>Sign Out</button>
        </div>
      </nav>

      <div style={{ maxWidth: 900, margin: "0 auto", padding: "36px 28px" }}>
        <div style={{ display: "flex", alignItems: "flex-start", justifyContent: "space-between", gap: 20, marginBottom: 28, flexWrap: "wrap" }}>
          <div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 3, marginBottom: 6, textTransform: "uppercase" }}>Team Console</div>
            <div style={{ ...styles.display, fontSize: 26, fontWeight: 800, color: C.white, marginBottom: 4 }}>Credential Contributors</div>
            <div style={{ fontSize: 13, color: C.dim }}>Read-only demo view of team contribution progress and credential coverage.</div>
          </div>
          <div style={{ display: "flex", gap: 12 }}>
            {[
              { label: "Members", value: members.length, color: C.teal },
              { label: "Credentials", value: totalCredentials, color: C.green },
              { label: "Commits", value: totalCommits, color: C.acc },
              { label: "Avg Progress", value: `${avgContribution}%`, color: C.amber },
            ].map(s => (
              <div key={s.label} style={{ ...styles.card, textAlign: "center", minWidth: 92, padding: "14px 18px" }}>
                <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginBottom: 4, textTransform: "uppercase" }}>{s.label}</div>
                <div style={{ ...styles.mono, fontSize: 24, fontWeight: 700, color: s.color }}>{s.value}</div>
              </div>
            ))}
          </div>
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
          {members.map(member => (
            <div key={member.slug} style={{ ...styles.card, display: "grid", gridTemplateColumns: "52px minmax(0,1fr)", gap: 16, alignItems: "start" }}>
              <div style={{ width: 52, height: 52, borderRadius: "50%", background: `${member.color}22`, border: `1px solid ${member.color}55`, display: "flex", alignItems: "center", justifyContent: "center", color: member.color, fontWeight: 800, ...styles.display }}>
                {member.name.split(" ").map(w => w[0]).join("")}
              </div>
              <div style={{ minWidth: 0 }}>
                <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", gap: 12, marginBottom: 4 }}>
                  <div>
                    <div style={{ fontSize: 14, color: C.white, fontWeight: 700 }}>{member.name}</div>
                    <div style={{ fontSize: 11, color: C.dim }}>{member.role} · {member.focus}</div>
                  </div>
                  <Tag color={C.teal} bg="rgba(0,201,167,0.1)" border="rgba(0,201,167,0.3)">
                    {member.credentialCount} Credential{member.credentialCount === 1 ? "" : "s"}
                  </Tag>
                </div>
                <div style={{ background: C.bg3, border: `1px solid ${C.border}`, borderRadius: 8, padding: "10px 12px", marginTop: 10 }}>
                  <div style={{ display: "flex", justifyContent: "space-between", gap: 10, marginBottom: 8, flexWrap: "wrap" }}>
                    <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 1.5, textTransform: "uppercase" }}>
                      GitHub @{member.github?.handle || "unlinked"}
                    </div>
                    <div style={{ ...styles.mono, fontSize: 10, color: C.dim }}>
                      {member.githubCommits} commits · {member.githubPrs} PRs
                    </div>
                  </div>
                  <GitHubRepoRows sync={member.github} color={member.color} />
                </div>
                <div style={{ display: "flex", alignItems: "center", gap: 10, marginTop: 12 }}>
                  <div style={{ height: 8, background: C.bg3, border: `1px solid ${C.border}`, borderRadius: 20, flex: 1, overflow: "hidden" }}>
                    <div style={{ height: "100%", width: `${member.contribution}%`, background: member.color, borderRadius: 20 }} />
                  </div>
                  <div style={{ ...styles.mono, fontSize: 11, color: member.color, width: 38, textAlign: "right" }}>{member.contribution}%</div>
                </div>
                <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginTop: 6 }}>
                  {member.verifiedCount} verified · Demo contribution progress
                </div>
              </div>
              <button style={{ ...styles.btnGhost, justifyContent: "center", gridColumn: "2", width: 120 }} onClick={() => onPortfolio(member.slug)}>
                Portfolio
              </button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// ── ISSUE CREDENTIAL ────────────────────────────────────────────────────
function IssueCredential({ user, onIssued, onBack }) {
  const [form, setForm] = useState({ recipient: "", recipientEmail: "", type: "certificate", title: "", description: "", issueDate: new Date().toISOString().slice(0, 10), skills: "" });
  const [step, setStep] = useState("form"); // form | signing | success
  const [sigStep, setSigStep] = useState(0);
  const [result, setResult] = useState(null);
  const [submitError, setSubmitError] = useState("");

  const sigSteps = [
    { label: "Computing SHA-256 hash...", color: C.acc },
    { label: "Submitting POST /credentials/issue...", color: C.amber },
    { label: "Waiting for P2 confirmation...", color: C.teal },
    { label: "Saving issued credential...", color: C.purple },
    { label: "Generating QR URL...", color: C.green },
  ];

  async function handleIssue() {
    if (!form.recipient || !form.title || !form.issueDate) return;
    setStep("signing");
    setSigStep(0);
    setSubmitError("");

    const skillsArr = form.skills.split(",").map(s => s.trim()).filter(Boolean);
    const credId = "cred_" + Math.random().toString(36).slice(2, 10);
    const dataHash = fakeHash(form.title + form.recipient + form.issueDate + skillsArr.join());

    const payload = {
      id: credId,
      title: form.title,
      type: form.type,
      recipient: form.recipient,
      recipientEmail: form.recipientEmail,
      issuer: user.name,
      issuerWallet: user.walletAddress || user.wallet,
      issueDate: form.issueDate,
      description: form.description,
      skills: skillsArr,
      dataHash,
    };

    try {
      setSigStep(1);
      const responseBody = await issueCredential(payload);
      setSigStep(sigSteps.length);
      const issuedCredential = normalizeIssuedCredential(payload, responseBody);

      setResult(issuedCredential);
      setStep("success");
      onIssued(issuedCredential);
    } catch (error) {
      setSubmitError(error.message || "Unable to issue credential.");
      setStep("form");
    }
  }

  function setField(k, v) { setForm(f => ({ ...f, [k]: v })); }

  function handleSubmit(e) {
    e.preventDefault();
    handleIssue();
  }

  if (step === "signing") {
    return (
      <div style={{ ...styles.app, display: "flex", alignItems: "center", justifyContent: "center", minHeight: "100vh" }}>
        <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />
        <div style={{ ...styles.card, width: 420, padding: 40, textAlign: "center" }}>
          <div style={{ fontSize: 32, marginBottom: 16 }}>⛓</div>
          <div style={{ ...styles.display, fontSize: 20, fontWeight: 700, color: C.white, marginBottom: 24 }}>Writing to Polygon Amoy...</div>
          <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
            {sigSteps.map((s, i) => (
              <div key={i} style={{ display: "flex", alignItems: "center", gap: 12, padding: "10px 14px", borderRadius: 8, background: i < sigStep ? `${s.color}18` : C.bg3, border: `1px solid ${i < sigStep ? s.color + "44" : C.border}`, transition: "all 0.4s" }}>
                <div style={{ width: 20, height: 20, borderRadius: "50%", background: i < sigStep ? s.color : C.bg3, border: `1px solid ${i < sigStep ? s.color : C.border2}`, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 10, color: i < sigStep ? "#000" : C.dim, flexShrink: 0, fontFamily: "'IBM Plex Mono', monospace", fontWeight: 700 }}>
                  {i < sigStep ? "✓" : i + 1}
                </div>
                <span style={{ fontSize: 12, color: i < sigStep ? C.white : C.dim }}>{s.label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (step === "success" && result) {
    return (
      <div style={{ ...styles.app, minHeight: "100vh" }}>
        <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />
        <nav style={styles.nav}>
          <div style={styles.brand}>⬡ ProofPass</div>
          <button style={styles.btnGhost} onClick={onBack}>← Back to Dashboard</button>
        </nav>
        <div style={{ maxWidth: 680, margin: "0 auto", padding: "40px 28px", textAlign: "center" }}>
          <div style={{ ...styles.display, fontSize: 28, fontWeight: 800, color: C.green, marginBottom: 6 }}>✓ Credential Issued!</div>
          <div style={{ fontSize: 14, color: C.dim, marginBottom: 32 }}>Successfully anchored on Polygon Amoy.</div>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 24, textAlign: "left" }}>
            <div style={styles.card}>
              <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginBottom: 10, textTransform: "uppercase" }}>Credential Details</div>
              {[
                { l: "Title", v: result.title, c: C.white },
                { l: "Recipient", v: result.recipient, c: C.teal },
                { l: "Issuer", v: result.issuer, c: C.acc },
                { l: "Date", v: result.issueDate, c: C.text },
              ].map(r => (
                <div key={r.l} style={{ marginBottom: 10 }}>
                  <div style={{ ...styles.mono, fontSize: 9, color: C.dim, textTransform: "uppercase", letterSpacing: 1.5 }}>{r.l}</div>
                  <div style={{ fontSize: 13, color: r.c, fontWeight: 500, marginTop: 2 }}>{r.v}</div>
                </div>
              ))}
              <div style={{ ...styles.mono, fontSize: 9, color: C.dim, textTransform: "uppercase", letterSpacing: 1.5, marginTop: 4 }}>Skills</div>
              <div style={{ display: "flex", gap: 5, flexWrap: "wrap", marginTop: 5 }}>
                {result.skills.map(s => <Tag key={s} color={C.teal} bg="rgba(0,201,167,0.1)" border="rgba(0,201,167,0.3)">{s}</Tag>)}
              </div>
            </div>
            <div style={{ ...styles.card, borderColor: "rgba(0,201,167,0.3)" }}>
              <div style={{ ...styles.mono, fontSize: 9, color: C.teal, letterSpacing: 2, marginBottom: 10, textTransform: "uppercase" }}>On-Chain Proof</div>
              {[
                { l: "Transaction", v: result.txHash.slice(0, 20) + "...", c: C.white },
                { l: "Block", v: "#" + result.blockNumber, c: C.teal },
                { l: "Network", v: result.network, c: C.amber },
              ].map(r => (
                <div key={r.l} style={{ marginBottom: 10 }}>
                  <div style={{ ...styles.mono, fontSize: 9, color: C.dim, textTransform: "uppercase", letterSpacing: 1.5 }}>{r.l}</div>
                  <div style={{ ...styles.mono, fontSize: 12, color: r.c, marginTop: 2 }}>{r.v}</div>
                </div>
              ))}
            </div>
          </div>

          {/* QR */}
          <div style={{ ...styles.card, padding: 28 }}>
            <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginBottom: 14, textTransform: "uppercase" }}>Share This Credential</div>
            <div style={{ display: "flex", justifyContent: "center", marginBottom: 12 }}>
              <QRCode value={`proofpass.io/verify/${result.id}`} size={140} />
            </div>
            <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginBottom: 14 }}>proofpass.io/verify/{result.id}</div>
            <div style={{ display: "flex", gap: 10, justifyContent: "center" }}>
              <button style={styles.btnGhost} onClick={() => navigator.clipboard?.writeText(`proofpass.io/verify/${result.id}`)}>Copy Link</button>
              <button style={styles.btnTeal} onClick={onBack}>← Dashboard</button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // FORM
  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />
      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
          <Tag color={C.amber} bg="rgba(245,166,35,0.1)" border="rgba(245,166,35,0.3)">ISSUER</Tag>
          <span style={{ fontSize: 12, color: C.dim }}>{user.name}</span>
        </div>
        <button style={styles.btnGhost} onClick={onBack}>← Dashboard</button>
      </nav>

      <div style={{ maxWidth: 860, margin: "0 auto", padding: "36px 28px" }}>
        <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 3, marginBottom: 6, textTransform: "uppercase" }}>Issuer Console</div>
        <div style={{ ...styles.display, fontSize: 26, fontWeight: 800, color: C.white, marginBottom: 4 }}>Issue New Credential</div>
        <div style={{ fontSize: 13, color: C.dim, marginBottom: 32 }}>Fill the form · backend hashes fields · writes to Polygon · stores in DB · generates QR URL.</div>

        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 28 }}>
          {/* FORM */}
          <form onSubmit={handleSubmit}>
            {[
              { key: "recipientEmail", label: "Recipient Email or Wallet", placeholder: "alice@example.com or 0x..." },
              { key: "recipient", label: "Recipient Name", placeholder: "Alice Wanjiku" },
            ].map(f => (
              <div key={f.key} style={{ marginBottom: 14 }}>
                <label style={styles.label}>{f.label}</label>
                <input style={styles.input} value={form[f.key]} onChange={e => setField(f.key, e.target.value)} placeholder={f.placeholder} />
              </div>
            ))}

            <div style={{ marginBottom: 14 }}>
              <label style={styles.label}>Credential Type</label>
              <select style={{ ...styles.input }} value={form.type} onChange={e => setField("type", e.target.value)}>
                <option value="degree">Degree</option>
                <option value="certificate">Certificate</option>
                <option value="participation">Participation</option>
                <option value="skill">Skill Badge</option>
              </select>
            </div>

            <div style={{ marginBottom: 14 }}>
              <label style={styles.label}>Title</label>
              <input style={styles.input} value={form.title} onChange={e => setField("title", e.target.value)} placeholder="BSc. Computer Science" />
            </div>

            <div style={{ marginBottom: 14 }}>
              <label style={styles.label}>Description</label>
              <textarea style={{ ...styles.input, height: 72, resize: "vertical" }} value={form.description} onChange={e => setField("description", e.target.value)} placeholder="4-year programme, graduated with honours..." />
            </div>

            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12, marginBottom: 14 }}>
              <div>
                <label style={styles.label}>Issue Date</label>
                <input style={styles.input} type="date" value={form.issueDate} onChange={e => setField("issueDate", e.target.value)} />
              </div>
              <div>
                <label style={styles.label}>Skills (comma-separated)</label>
                <input style={styles.input} value={form.skills} onChange={e => setField("skills", e.target.value)} placeholder="Solidity, Go, Blockchain" />
              </div>
            </div>

            {submitError && (
              <div style={{ background: "rgba(255,77,109,0.08)", border: "1px solid rgba(255,77,109,0.25)", borderRadius: 8, padding: "10px 12px", color: C.red, fontSize: 12, marginBottom: 12 }}>
                {submitError}
              </div>
            )}

            <button type="submit" style={{ ...styles.btnBlue, marginTop: 6, padding: "11px 22px" }}
              disabled={!form.recipient || !form.title || !form.issueDate}>
              ⛓ Sign &amp; Issue on Blockchain
            </button>
            <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginTop: 8 }}>Signing with: {user.wallet} ({user.name})</div>
          </form>

          {/* WHAT HAPPENS */}
          <div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, letterSpacing: 2, marginBottom: 14, textTransform: "uppercase" }}>What happens when you click issue</div>
            <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
              {[
                { n: 1, c: C.acc, title: "Hash credential data", desc: "SHA-256 of all fields combined" },
                { n: 2, c: C.amber, title: "Write hash to Polygon", desc: "issueCredential(id, hash) on-chain" },
                { n: 3, c: C.teal, title: "Store in database", desc: "Credential + tx_hash + block_number saved" },
                { n: 4, c: C.green, title: "Return QR URL", desc: "proofpass.io/verify/:id — ready to share", highlight: true },
              ].map(s => (
                <div key={s.n} style={{ ...styles.card, display: "flex", gap: 12, alignItems: "flex-start", borderColor: s.highlight ? "rgba(34,197,94,0.3)" : C.border }}>
                  <div style={{ width: 22, height: 22, background: s.c, borderRadius: "50%", display: "flex", alignItems: "center", justifyContent: "center", fontFamily: "'IBM Plex Mono', monospace", fontSize: 10, color: s.n === 1 ? "#fff" : "#000", flexShrink: 0, marginTop: 1 }}>{s.n}</div>
                  <div>
                    <div style={{ fontSize: 12, color: C.white, fontWeight: 500 }}>{s.title}</div>
                    <div style={{ fontSize: 11, color: C.dim }}>{s.desc}</div>
                  </div>
                </div>
              ))}
            </div>
            <div style={{ background: "rgba(245,166,35,0.06)", border: "1px solid rgba(245,166,35,0.2)", borderRadius: 8, padding: "12px 16px", marginTop: 14, fontSize: 12, color: C.amber, lineHeight: 1.6 }}>
              <strong style={{ color: C.amber }}>On success:</strong> A success screen shows the QR code and verify URL. Download the QR immediately.
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ── VERIFY PAGE ──────────────────────────────────────────────────────
function VerifyPage({ credentialId, allCredentials, onBack }) {
  const [inputId, setInputId] = useState(credentialId || "");
  const [searching, setSearching] = useState(false);
  const [state, setState] = useState(null); // null | 'verified' | 'tampered' | 'notfound'
  const [found, setFound] = useState(null);

  const doVerify = useCallback(async (id) => {
    if (!id) return;
    setSearching(true);
    setState(null);
    await new Promise(r => setTimeout(r, 1000));

    if (id === "cred_fake_999" || id.includes("fake") || id.includes("tampered")) {
      setState("tampered");
      setFound(null);
    } else {
      const cred = allCredentials.find(c => c.id === id);
      if (cred) { setState("verified"); setFound(cred); }
      else { setState("notfound"); setFound(null); }
    }
    setSearching(false);
  }, [allCredentials]);

  useEffect(() => {
    if (!credentialId) return undefined;
    const timer = setTimeout(() => {
      doVerify(credentialId);
    }, 0);
    return () => clearTimeout(timer);
  }, [credentialId, doVerify]);

  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />

      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ fontSize: 12, color: C.dim }}>Public Verification · No login required</div>
        <button style={styles.btnGhost} onClick={onBack}>← Home</button>
      </nav>

      <div style={{ maxWidth: 520, margin: "0 auto", padding: "48px 24px" }}>
        {/* Search bar */}
        <div style={{ display: "flex", gap: 10, marginBottom: 32 }}>
          <input style={{ ...styles.input, flex: 1 }} value={inputId}
            onChange={e => setInputId(e.target.value)}
            placeholder="Enter credential ID..."
            onKeyDown={e => e.key === "Enter" && doVerify(inputId)} />
          <button style={styles.btnTeal} onClick={() => doVerify(inputId)}>🔍 Verify</button>
        </div>

        {/* Loading */}
        {searching && (
          <div style={{ textAlign: "center", padding: "48px 0" }}>
            <div style={{ ...styles.mono, fontSize: 12, color: C.acc, letterSpacing: 2 }}>QUERYING POLYGON AMOY...</div>
            <div style={{ marginTop: 16, display: "flex", justifyContent: "center", gap: 8 }}>
              {[0, 1, 2].map(i => (
                <div key={i} style={{ width: 8, height: 8, borderRadius: "50%", background: C.acc, opacity: 0.3 + i * 0.3 }} />
              ))}
            </div>
          </div>
        )}

        {/* VERIFIED */}
        {state === "verified" && found && (
          <div style={{ textAlign: "center" }}>
            <div style={{ display: "inline-flex", alignItems: "center", gap: 8, padding: "8px 22px", borderRadius: 100, background: "rgba(34,197,94,0.12)", color: C.green, border: `2px solid ${C.green}`, ...styles.mono, fontSize: 13, fontWeight: 600, letterSpacing: 1, marginBottom: 16 }}>
              ✓ CREDENTIAL VERIFIED
            </div>
            <div style={{ ...styles.display, fontSize: 24, fontWeight: 800, color: C.white, marginBottom: 4 }}>{found.title}</div>
            <div style={{ fontSize: 13, color: C.dim, marginBottom: 24 }}>{found.issuer} · {found.issueDate}</div>

            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12, textAlign: "left", marginBottom: 20 }}>
              <div style={styles.card}>
                {[
                  { l: "Holder", v: found.recipient, c: C.white },
                  { l: "Issued By", v: found.issuer, c: C.acc },
                ].map(r => (
                  <div key={r.l} style={{ marginBottom: 12 }}>
                    <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, textTransform: "uppercase" }}>{r.l}</div>
                    <div style={{ fontSize: 13, color: r.c, fontWeight: 500, marginTop: 3 }}>{r.v}</div>
                  </div>
                ))}
                <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginTop: 4, textTransform: "uppercase" }}>Skills</div>
                <div style={{ display: "flex", gap: 5, flexWrap: "wrap", marginTop: 5 }}>
                  {found.skills.map(s => <Tag key={s} color={found.color} bg={`${found.color}18`} border={`${found.color}33`}>{s}</Tag>)}
                </div>
              </div>
              <div style={{ ...styles.card, borderColor: "rgba(0,201,167,0.3)" }}>
                <div style={{ ...styles.mono, fontSize: 9, color: C.teal, letterSpacing: 2, marginBottom: 10, textTransform: "uppercase" }}>On-Chain Proof</div>
                {[
                  { l: "Transaction", v: found.txHash.slice(0, 18) + "...", c: C.white },
                  { l: "Block", v: "#" + found.blockNumber, c: C.teal },
                  { l: "Network", v: found.network, c: C.amber },
                ].map(r => (
                  <div key={r.l} style={{ marginBottom: 9 }}>
                    <div style={{ ...styles.mono, fontSize: 9, color: C.dim, textTransform: "uppercase", letterSpacing: 1.5 }}>{r.l}</div>
                    <div style={{ ...styles.mono, fontSize: 11, color: r.c, marginTop: 2 }}>{r.v}</div>
                  </div>
                ))}
              </div>
            </div>

            <div style={{ ...styles.card, padding: 24 }}>
              <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 2, marginBottom: 12, textTransform: "uppercase" }}>Share</div>
              <div style={{ display: "flex", justifyContent: "center", marginBottom: 12 }}>
                <QRCode value={`proofpass.io/verify/${found.id}`} size={110} />
              </div>
              <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginBottom: 10 }}>proofpass.io/verify/{found.id}</div>
              <button style={styles.btnGhost} onClick={() => navigator.clipboard?.writeText(`proofpass.io/verify/${found.id}`)}>Copy Link</button>
            </div>
          </div>
        )}

        {/* TAMPERED */}
        {state === "tampered" && (
          <div style={{ textAlign: "center" }}>
            <div style={{ display: "inline-flex", alignItems: "center", gap: 8, padding: "8px 22px", borderRadius: 100, background: "rgba(255,77,109,0.12)", color: C.red, border: `2px solid ${C.red}`, ...styles.mono, fontSize: 13, fontWeight: 600, letterSpacing: 1, marginBottom: 16 }}>
              ✗ CREDENTIAL INVALID
            </div>
            <div style={{ ...styles.display, fontSize: 22, fontWeight: 800, color: C.white, marginBottom: 8 }}>Hash Mismatch Detected</div>
            <div style={{ fontSize: 13, color: C.dim, maxWidth: 320, margin: "0 auto 20px", lineHeight: 1.7 }}>
              The credential data does not match what was recorded on the blockchain. This credential may have been tampered with.
            </div>
            <div style={{ ...styles.card, borderColor: "rgba(255,77,109,0.3)", textAlign: "left", maxWidth: 340, margin: "0 auto 20px" }}>
              <div style={{ ...styles.mono, fontSize: 9, color: C.red, letterSpacing: 2, marginBottom: 10, textTransform: "uppercase" }}>Mismatch Details</div>
              <div style={{ fontSize: 11, color: C.dim, marginBottom: 4 }}>Expected (on-chain)</div>
              <div style={{ ...styles.mono, fontSize: 11, color: C.red, wordBreak: "break-all" }}>0x7f3a8b2c4d5e6f7a8b9c0d1e2f3a4b5c...</div>
              <div style={{ fontSize: 11, color: C.dim, marginTop: 8, marginBottom: 4 }}>Found (computed)</div>
              <div style={{ ...styles.mono, fontSize: 11, color: C.red, wordBreak: "break-all" }}>0xb9c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6...</div>
            </div>
            <button style={{ background: "rgba(255,77,109,0.15)", color: C.red, border: `1px solid rgba(255,77,109,0.3)`, borderRadius: 8, padding: "8px 18px", fontSize: 12, cursor: "pointer" }}>
              Report This Credential
            </button>
          </div>
        )}

        {/* NOT FOUND */}
        {state === "notfound" && (
          <div style={{ textAlign: "center", padding: "40px 0" }}>
            <div style={{ fontSize: 40, marginBottom: 16 }}>🔍</div>
            <div style={{ ...styles.display, fontSize: 22, fontWeight: 700, color: C.white, marginBottom: 8 }}>Credential Not Found</div>
            <div style={{ fontSize: 13, color: C.dim, marginBottom: 24, lineHeight: 1.7 }}>
              No credential with this ID exists on ProofPass. Check the ID and try again.
            </div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim }}>
              Try: cred_alice_001 · cred_brian_002 · cred_carol_003
            </div>
          </div>
        )}

        {/* Idle hint */}
        {!state && !searching && (
          <div style={{ textAlign: "center", padding: "32px 0", color: C.dim }}>
            <div style={{ fontSize: 32, marginBottom: 14 }}>⬡</div>
            <div style={{ fontSize: 13, marginBottom: 8 }}>Enter a credential ID above to verify</div>
            <div style={{ ...styles.mono, fontSize: 11 }}>
              Try: cred_alice_001 · cred_brian_002 · cred_carol_003 · cred_fake_999
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

// ── PORTFOLIO PAGE ────────────────────────────────────────────────────
function PortfolioPage({ slug, allCredentials, onVerify, onBack }) {
  const personName = PROFILE_SLUGS[slug] || slug;
  const creds = allCredentials.filter(c => c.recipient === personName);
  const [shareQr, setShareQr] = useState(false);
  const [expandedId, setExpandedId] = useState(null);
  const githubSync = GITHUB_SYNC[slug];
  const githubStats = githubTotals(githubSync);
  const githubAvgContribution = githubSync?.repos.length
    ? Math.round(githubStats.contribution / githubSync.repos.length)
    : 0;

  // Aggregate all unique skills across credentials
  const allSkills = [...new Set(creds.flatMap(c => c.skills))];

  return (
    <div style={styles.app}>
      <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500;600;700&family=Bricolage+Grotesque:wght@400;600;700;800&family=DM+Sans:wght@300;400;500;600&display=swap" rel="stylesheet" />

      {/* NAV */}
      <nav style={styles.nav}>
        <div style={styles.brand}>⬡ ProofPass</div>
        <div style={{ fontSize: 12, color: C.dim }}>Public Portfolio · No login required</div>
        <button style={styles.btnGhost} onClick={onBack}>← Back</button>
      </nav>

      <div style={{ maxWidth: 700, margin: "0 auto", padding: "48px 24px" }}>

        {creds.length === 0 ? (
          <div style={{ textAlign: "center", padding: "64px 0" }}>
            <div style={{ fontSize: 40, marginBottom: 16 }}>🔍</div>
            <div style={{ ...styles.display, fontSize: 22, fontWeight: 700, color: C.white, marginBottom: 8 }}>Portfolio Not Found</div>
            <div style={{ fontSize: 13, color: C.dim }}>No credentials found for "{slug}".<br />Try: alice · brian · carol</div>
          </div>
        ) : (
          <>
            {/* Profile header */}
            <div style={{ textAlign: "center", marginBottom: 36 }}>
              {/* Avatar */}
              <div style={{ width: 72, height: 72, borderRadius: "50%", background: `linear-gradient(135deg, ${C.acc}44, ${C.teal}44)`, border: `2px solid ${C.teal}55`, display: "flex", alignItems: "center", justifyContent: "center", margin: "0 auto 16px", fontSize: 28 }}>
                {personName.split(" ").map(w => w[0]).join("")}
              </div>

              <div style={{ ...styles.display, fontSize: 28, fontWeight: 800, color: C.white, marginBottom: 4 }}>{personName}</div>
              <div style={{ ...styles.mono, fontSize: 11, color: C.dim, marginBottom: 16 }}>
                proofpass.io/profile/{slug}
              </div>

              {/* Summary badges */}
              <div style={{ display: "flex", gap: 10, justifyContent: "center", flexWrap: "wrap", marginBottom: 20 }}>
                <span style={{ background: "rgba(0,201,167,0.1)", border: "1px solid rgba(0,201,167,0.3)", color: C.teal, borderRadius: 20, padding: "4px 14px", ...styles.mono, fontSize: 11 }}>
                  {creds.length} Credential{creds.length !== 1 ? "s" : ""}
                </span>
                <span style={{ background: "rgba(34,197,94,0.1)", border: "1px solid rgba(34,197,94,0.3)", color: C.green, borderRadius: 20, padding: "4px 14px", ...styles.mono, fontSize: 11 }}>
                  ✓ All Verified
                </span>
                <span style={{ background: "rgba(245,166,35,0.1)", border: "1px solid rgba(245,166,35,0.3)", color: C.amber, borderRadius: 20, padding: "4px 14px", ...styles.mono, fontSize: 11 }}>
                  Polygon Amoy
                </span>
                <span style={{ background: githubSync ? "rgba(91,141,239,0.1)" : "rgba(245,166,35,0.1)", border: githubSync ? "1px solid rgba(91,141,239,0.3)" : "1px solid rgba(245,166,35,0.3)", color: githubSync ? C.acc : C.amber, borderRadius: 20, padding: "4px 14px", ...styles.mono, fontSize: 11 }}>
                  {githubSync ? `GitHub @${githubSync.handle}` : "GitHub Unlinked"}
                </span>
              </div>

              {/* Share buttons */}
              <div style={{ display: "flex", gap: 8, justifyContent: "center" }}>
                <button style={{ ...styles.btnGhost, fontSize: 12 }} onClick={() => setShareQr(true)}>📱 Share Portfolio QR</button>
                <button style={{ ...styles.btnGhost, fontSize: 12 }} onClick={() => navigator.clipboard?.writeText(`proofpass.io/profile/${slug}`)}>🔗 Copy Link</button>
              </div>
            </div>

            {/* GitHub contribution evidence */}
            <div style={{ ...styles.card, marginBottom: 24, borderColor: githubSync ? "rgba(91,141,239,0.25)" : C.border }}>
              <div style={{ display: "flex", justifyContent: "space-between", gap: 16, alignItems: "flex-start", flexWrap: "wrap", marginBottom: 16 }}>
                <div>
                  <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 2, marginBottom: 8, textTransform: "uppercase" }}>GitHub Sync</div>
                  <div style={{ ...styles.display, fontSize: 18, fontWeight: 700, color: C.white, marginBottom: 4 }}>
                    {githubSync ? `@${githubSync.handle}` : "No handle linked"}
                  </div>
                  <div style={{ fontSize: 12, color: C.dim, lineHeight: 1.6 }}>
                    {githubSync
                      ? `${githubStats.commits} commits and ${githubStats.prs} pull requests synced from ${githubSync.repos.length} repositories.`
                      : "Repository contribution evidence has not been synced for this portfolio."}
                  </div>
                  {githubSync && (
                    <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginTop: 8 }}>
                      Last synced {githubSync.lastSynced}
                    </div>
                  )}
                </div>
                <div style={{ display: "flex", gap: 10, flexWrap: "wrap" }}>
                  <div style={{ minWidth: 84, background: C.bg3, border: `1px solid ${C.border}`, borderRadius: 8, padding: "10px 12px", textAlign: "center" }}>
                    <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 1.5, textTransform: "uppercase", marginBottom: 4 }}>Repos</div>
                    <div style={{ ...styles.mono, fontSize: 20, color: C.acc, fontWeight: 700 }}>{githubSync?.repos.length || 0}</div>
                  </div>
                  <div style={{ minWidth: 84, background: C.bg3, border: `1px solid ${C.border}`, borderRadius: 8, padding: "10px 12px", textAlign: "center" }}>
                    <div style={{ ...styles.mono, fontSize: 9, color: C.dim, letterSpacing: 1.5, textTransform: "uppercase", marginBottom: 4 }}>Avg</div>
                    <div style={{ ...styles.mono, fontSize: 20, color: C.teal, fontWeight: 700 }}>{githubAvgContribution}%</div>
                  </div>
                </div>
              </div>
              <GitHubRepoRows sync={githubSync} color={C.acc} />
            </div>

            {/* Skills overview */}
            {allSkills.length > 0 && (
              <div style={{ ...styles.card, marginBottom: 24, borderColor: "rgba(91,141,239,0.25)" }}>
                <div style={{ ...styles.mono, fontSize: 10, color: C.acc, letterSpacing: 2, marginBottom: 12, textTransform: "uppercase" }}>Verified Skills</div>
                <div style={{ display: "flex", gap: 7, flexWrap: "wrap" }}>
                  {allSkills.map(s => (
                    <span key={s} style={{ background: "rgba(91,141,239,0.1)", border: "1px solid rgba(91,141,239,0.25)", color: C.acc, borderRadius: 20, padding: "4px 12px", ...styles.mono, fontSize: 10 }}>{s}</span>
                  ))}
                </div>
              </div>
            )}

            {/* Credentials list */}
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, letterSpacing: 2, marginBottom: 12, textTransform: "uppercase" }}>
              {creds.length} Verified Credential{creds.length !== 1 ? "s"  : ""}
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
              {creds.map(cred => {
                const expanded = expandedId === cred.id;
                return (
                  <div key={cred.id} style={{ ...styles.card, borderColor: expanded ? `${cred.color}55` : C.border, transition: "border-color 0.2s" }}>
                    {/* Credential row */}
                    <div style={{ display: "flex", alignItems: "center", gap: 14, cursor: "pointer" }}
                      onClick={() => setExpandedId(expanded ? null : cred.id)}>
                      <div style={{ width: 44, height: 44, background: `${cred.color}22`, border: `1px solid ${cred.color}44`, borderRadius: 10, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 20, flexShrink: 0 }}>
                        {cred.emoji}
                      </div>
                      <div style={{ flex: 1, minWidth: 0 }}>
                        <div style={{ fontSize: 14, color: C.white, fontWeight: 600, marginBottom: 2 }}>{cred.title}</div>
                        <div style={{ fontSize: 11, color: C.dim }}>{cred.issuer} · {cred.issueDate}</div>
                        <div style={{ display: "flex", gap: 5, marginTop: 6, flexWrap: "wrap" }}>
                          {cred.skills.map(s => <Tag key={s} color={cred.color} bg={`${cred.color}18`} border={`${cred.color}33`}>{s}</Tag>)}
                        </div>
                      </div>
                      <div style={{ display: "flex", flexDirection: "column", alignItems: "flex-end", gap: 6, flexShrink: 0 }}>
                        <Tag color={C.green} bg="rgba(34,197,94,0.1)" border="rgba(34,197,94,0.3)">✓ VERIFIED</Tag>
                        <span style={{ ...styles.mono, fontSize: 10, color: C.dim }}>{expanded ? "▲ less" : "▼ details"}</span>
                      </div>
                    </div>

                    {/* Expanded detail */}
                    {expanded && (
                      <div style={{ marginTop: 16, paddingTop: 16, borderTop: `1px solid ${C.border}` }}>
                        {cred.description && (
                          <div style={{ fontSize: 12, color: C.text, lineHeight: 1.7, marginBottom: 14 }}>{cred.description}</div>
                        )}
                        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 10, marginBottom: 14 }}>
                          {[
                            { l: "Issuer", v: cred.issuer, c: C.acc },
                            { l: "Block", v: "#" + cred.blockNumber, c: C.teal },
                            { l: "Tx Hash", v: cred.txHash.slice(0, 18) + "...", c: C.white },
                            { l: "Network", v: cred.network, c: C.amber },
                          ].map(r => (
                            <div key={r.l}>
                              <div style={{ ...styles.mono, fontSize: 9, color: C.dim, textTransform: "uppercase", letterSpacing: 1.5 }}>{r.l}</div>
                              <div style={{ ...styles.mono, fontSize: 11, color: r.c, marginTop: 2 }}>{r.v}</div>
                            </div>
                          ))}
                        </div>
                        <button style={{ ...styles.btnGhost, fontSize: 11 }} onClick={() => onVerify(cred.id)}>
                          🔍 Verify this credential
                        </button>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>

            {/* Footer verification note */}
            <div style={{ background: "rgba(34,197,94,0.05)", border: "1px solid rgba(34,197,94,0.15)", borderRadius: 10, padding: "14px 18px", marginTop: 28, display: "flex", gap: 10, alignItems: "flex-start" }}>
              <span style={{ color: C.green, fontSize: 16, flexShrink: 0 }}>✓</span>
              <div style={{ fontSize: 12, color: C.dim, lineHeight: 1.7 }}>
                All credentials on this portfolio are independently verifiable on <span style={{ color: C.amber }}>Polygon Amoy</span>. No login required. Click any credential above and tap "Verify this credential" for full blockchain proof.
              </div>
            </div>
          </>
        )}
      </div>

      {/* Share QR modal */}
      {shareQr && (
        <div style={{ position: "fixed", inset: 0, background: "rgba(0,0,0,0.8)", zIndex: 200, display: "flex", alignItems: "center", justifyContent: "center" }}
          onClick={e => e.target === e.currentTarget && setShareQr(false)}>
          <div style={{ ...styles.card, width: 320, padding: 32, textAlign: "center" }}>
            <div style={{ ...styles.mono, fontSize: 10, color: C.teal, letterSpacing: 3, marginBottom: 10, textTransform: "uppercase" }}>Share Full Portfolio</div>
            <div style={{ ...styles.display, fontSize: 18, fontWeight: 700, color: C.white, marginBottom: 2 }}>{personName}</div>
            <div style={{ fontSize: 12, color: C.dim, marginBottom: 18 }}>{creds.length} credentials · scan to view all</div>
            <div style={{ display: "flex", justifyContent: "center", marginBottom: 12 }}>
              <QRCode value={`proofpass.io/profile/${slug}`} size={148} />
            </div>
            <div style={{ ...styles.mono, fontSize: 10, color: C.dim, marginBottom: 16 }}>proofpass.io/profile/{slug}</div>
            <div style={{ display: "flex", gap: 8, justifyContent: "center" }}>
              <button style={styles.btnGhost} onClick={() => navigator.clipboard?.writeText(`proofpass.io/profile/${slug}`)}>Copy Link</button>
              <button style={styles.btnBlue} onClick={() => setShareQr(false)}>Close</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// ── ROOT APP ────────────────────────────────────────────────────────────
export default function App() {
  const [page, setPage] = useState("landing"); // landing | dashboard | issue | verify | portfolio
  const [user, setUser] = useState(null);
  const [credentials, setCredentials] = useState(SEED_CREDENTIALS);
  const [verifyId, setVerifyId] = useState(null);
  const [portfolioSlug, setPortfolioSlug] = useState(null);
  const { disconnect } = useDisconnect();

  function handleLogin(u) { setUser(u); setPage("dashboard"); }
  function handleLogout() {
    disconnect();
    setUser(null);
    setPage("landing");
  }
  function handleVerify(id) { setVerifyId(id); setPage("verify"); }
  function handlePortfolio(slug) { setPortfolioSlug(slug); setPage("portfolio"); }
  function handleIssued(cred) { setCredentials(cs => [...cs, cred]); }

  if (page === "portfolio") return <PortfolioPage slug={portfolioSlug} allCredentials={credentials} onVerify={handleVerify} onBack={() => setPage(user ? "dashboard" : "landing")} />;
  if (page === "verify") return <VerifyPage credentialId={verifyId} allCredentials={credentials} onBack={() => setPage(user ? "dashboard" : "landing")} />;
  if (page === "team" && user) return <TeamPage user={user} credentials={credentials} onBack={() => setPage("dashboard")} onPortfolio={handlePortfolio} onLogout={handleLogout} />;
  if (page === "issue" && user) return <IssueCredential user={user} onIssued={handleIssued} onBack={() => setPage("dashboard")} />;
  if (page === "dashboard" && user) return <Dashboard user={user} credentials={credentials} onVerify={handleVerify} onIssue={() => setPage("issue")} onPortfolio={handlePortfolio} onTeam={() => setPage("team")} onLogout={handleLogout} />;

  return <Landing onLogin={handleLogin} onVerify={handleVerify} onPortfolio={handlePortfolio} />;
}

# Running ProofPass Locally

Use this order when running the local blockchain, backend, and frontend.

## 1. Start Hardhat Node

Open terminal 1:

```bash
cd proofpass/proofpass-contract
npx hardhat node
```

Keep this terminal running. If you stop or restart it, the local blockchain resets and you must deploy the contract again.

## 2. Deploy ProofPass Contract

Open terminal 2:

```bash
cd proofpass/proofpass-contract
npx hardhat run scripts/deploy.ts --network localhost
```

Copy the deployed contract address from:

```txt
ProofPass deployed to: 0x...
```

Update `backend/.env`:

```env
RPC_URL=http://127.0.0.1:8545
CONTRACT_ADDRESS=0xPASTE_DEPLOYED_ADDRESS_HERE
PRIVATE_KEY=0xYOUR_BACKEND_WALLET_PRIVATE_KEY
```

For local Hardhat, the first account private key from the `npx hardhat node` output is usually fine for development.

## 3. Authorize Backend Wallet

Still in terminal 2:

```bash
cd proofpass/proofpass-contract
npx hardhat run scripts/addIssuer.ts --network localhost
```

This allows the backend wallet to call `issueCredential`.

## 4. Start Backend

Open terminal 3:

```bash
cd backend
go run cmd/server/main.go
```

You should see:

```txt
Listening and serving HTTP on :8080
```

If port `8080` is already in use:

```bash
lsof -i :8080
kill <PID>
go run cmd/server/main.go
```

Test the backend:

```txt
http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

## 5. Start Frontend

Open terminal 4:

```bash
cd frontend/my-app
npm run dev
```

Open the URL printed by Vite, usually:

```txt
http://localhost:5173
```

## Full Sequence

```txt
1. npx hardhat node
2. npx hardhat run scripts/deploy.ts --network localhost
3. Update backend/.env CONTRACT_ADDRESS
4. npx hardhat run scripts/addIssuer.ts --network localhost
5. go run cmd/server/main.go
6. npm run dev
```

## Common Errors

### `Failed to fetch`

The frontend cannot reach the backend.

Check:

```txt
http://localhost:8080/health
```

Also make sure the frontend API URL points to:

```txt
http://localhost:8080/api
```

### `no contract code at given address`

`CONTRACT_ADDRESS` does not match the contract deployed on the current Hardhat node.

Fix:

```txt
Redeploy the contract and update backend/.env with the new address.
```

### `Not issuer`

The backend wallet has not been authorized.

Fix:

```bash
cd proofpass/proofpass-contract
npx hardhat run scripts/addIssuer.ts --network localhost
```

### `invalid chainId. The expected chainId is 31337`

Restart the backend so it uses the latest code and connects to the running Hardhat node.

```bash
lsof -i :8080
kill <PID>
cd backend
go run cmd/server/main.go
```

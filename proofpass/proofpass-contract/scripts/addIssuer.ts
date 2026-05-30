import { network } from "hardhat";
import { getAddress, isAddress } from "viem";
import { privateKeyToAccount } from "viem/accounts";

const proofPassAddress = process.env.PROOFPASS_ADDRESS ?? process.env.CONTRACT_ADDRESS;
const backendWalletAddress =
  process.env.BACKEND_WALLET_ADDRESS ??
  (process.env.PRIVATE_KEY
    ? privateKeyToAccount(normalizePrivateKey(process.env.PRIVATE_KEY)).address
    : undefined);

async function main() {
  if (!proofPassAddress || !isAddress(proofPassAddress)) {
    throw new Error(`Invalid PROOFPASS_ADDRESS: ${proofPassAddress}`);
  }

  if (!backendWalletAddress || !isAddress(backendWalletAddress)) {
    throw new Error(`Invalid BACKEND_WALLET_ADDRESS: ${backendWalletAddress}`);
  }

  const proofPass = getAddress(proofPassAddress);
  const backendWallet = getAddress(backendWalletAddress);

  const { viem } = await network.create();
  const [owner] = await viem.getWalletClients();
  const publicClient = await viem.getPublicClient();
  const contract = await viem.getContractAt("ProofPass", proofPass);

  const contractOwner = await contract.read.owner();
  const alreadyIssuer = await contract.read.authorizedIssuers([backendWallet]);

  console.log("ProofPass:", proofPass);
  console.log("Signer:", owner.account.address);
  console.log("Contract owner:", contractOwner);
  console.log("Backend wallet:", backendWallet);

  if (getAddress(owner.account.address) !== getAddress(contractOwner)) {
    throw new Error("Signer is not the ProofPass owner and cannot add issuers");
  }

  if (alreadyIssuer) {
    console.log("Issuer is already authorized.");
    return;
  }

  console.log("Adding issuer...");
  const hash = await contract.write.addIssuer([backendWallet]);
  console.log("Transaction:", hash);

  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  console.log("Issuer added in block:", receipt.blockNumber.toString());
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});

function normalizePrivateKey(privateKey: string): `0x${string}` {
  return privateKey.startsWith("0x")
    ? (privateKey as `0x${string}`)
    : (`0x${privateKey}` as `0x${string}`);
}

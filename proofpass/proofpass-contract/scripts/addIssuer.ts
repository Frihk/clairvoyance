import { network } from "hardhat";
import { getAddress, isAddress } from "viem";

const proofPassAddress =
  process.env.PROOFPASS_ADDRESS ?? "0x8a4dcdc1db7502948b43713e8eb67817d6b642a4";
const backendWalletAddress =
  process.env.BACKEND_WALLET_ADDRESS ?? "0x06b1c4327d2cf5e3aa3cb0c84f0158fd89bb32cd";

async function main() {
  if (!isAddress(proofPassAddress)) {
    throw new Error(`Invalid PROOFPASS_ADDRESS: ${proofPassAddress}`);
  }

  if (!isAddress(backendWalletAddress)) {
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

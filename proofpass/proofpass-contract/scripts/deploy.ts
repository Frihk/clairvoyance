import { artifacts, network } from "hardhat";

async function main() {
  const { viem } = await network.create();
  const [deployer] = await viem.getWalletClients();
  const publicClient = await viem.getPublicClient();
  const gasLimit = BigInt(process.env.DEPLOY_GAS_LIMIT ?? "1000000");

  console.log("Deploying with:", deployer.account.address);
  console.log("Using gas limit:", gasLimit.toString());

  const balance = await publicClient.getBalance({
    address: deployer.account.address,
  });
  console.log("Deployer balance:", balance.toString(), "wei");
  console.log("Deployer balance:", Number(balance) / 1e18, "POL");

  const contract = await viem.deployContract("ProofPass", [], {
    gas: gasLimit,
  });
  const artifact = await artifacts.readArtifact("ProofPass");

  console.log("ProofPass deployed to:", contract.address);
  console.log("ABI:");
  console.log(JSON.stringify(artifact.abi, null, 2));
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});

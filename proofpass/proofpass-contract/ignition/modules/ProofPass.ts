import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

export default buildModule("ProofPassModule", (m) => {
  const proofPass = m.contract("ProofPass");

  return { proofPass };
});

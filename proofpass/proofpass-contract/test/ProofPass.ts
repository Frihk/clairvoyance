import { expect } from "chai";
import hre from "hardhat";
import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { zeroAddress } from "viem";

import "@nomicfoundation/hardhat-toolbox-viem";

describe("ProofPass Smart Contract", function () {
  async function deployFixture() {
    const { viem } = await hre.network.create();
    const [owner, issuer, unauthorizedUser] = await viem.getWalletClients();
    const proofPass = await viem.deployContract("ProofPass");

    return { proofPass, owner, issuer, unauthorizedUser };
  }

  describe("Deployment & Access Control", function () {
    it("Should set the right owner", async function () {
      const { proofPass, owner } = await deployFixture();
      const contractOwner = await proofPass.read.owner();
      
      expect(contractOwner.toLowerCase()).to.equal(owner.account.address.toLowerCase());
    });

    it("Should allow the owner to add an issuer", async function () {
      const { proofPass, issuer } = await deployFixture();
      await proofPass.write.addIssuer([issuer.account.address]);
      
      const isIssuer = await proofPass.read.authorizedIssuers([issuer.account.address]);
      expect(isIssuer).to.be.true;
    });

    it("Should revert if a non-owner tries to add an issuer", async function () {
      const { proofPass, issuer, unauthorizedUser } = await deployFixture();

      await assert.rejects(
        proofPass.write.addIssuer([issuer.account.address], {
          account: unauthorizedUser.account,
        }),
        /Not owner/
      );
    });
  });

  describe("Issuing Credentials", function () {
    const credId = "cred_12345";
    const dataHash = "0xabcdef1234567890";

    it("Should allow an authorized issuer to issue a credential", async function () {
      const { proofPass, issuer } = await deployFixture();
      await proofPass.write.addIssuer([issuer.account.address]);
      
      await proofPass.write.issueCredential([credId, dataHash], {
        account: issuer.account,
      });

      const result = await proofPass.read.verifyCredential([credId]);
      expect(result[0]).to.equal(dataHash);
      expect(result[1].toLowerCase()).to.equal(issuer.account.address.toLowerCase());
    });

    it("Should revert if an unauthorized user tries to issue", async function () {
      const { proofPass, unauthorizedUser } = await deployFixture();

      await assert.rejects(
        proofPass.write.issueCredential([credId, dataHash], {
          account: unauthorizedUser.account,
        }),
        /Not issuer/
      );
    });

    it("Should revert if the credential ID has already been issued", async function () {
      const { proofPass, issuer } = await deployFixture();
      await proofPass.write.addIssuer([issuer.account.address]);
      
      await proofPass.write.issueCredential([credId, dataHash], {
        account: issuer.account,
      });

      await assert.rejects(
        proofPass.write.issueCredential([credId, "0xnewdatahash"], {
          account: issuer.account,
        }),
        /Already issued/
      );
    });
  });

  describe("Verifying Credentials", function () {
    const credId = "cred_alice_001";
    const dataHash = "0xdeadbeef";

    it("Should return correct data for a valid, existing credential", async function () {
      const { proofPass, issuer } = await deployFixture();
      await proofPass.write.addIssuer([issuer.account.address]);
      await proofPass.write.issueCredential([credId, dataHash], {
        account: issuer.account,
      });

      const [returnedHash, returnedIssuer, issuedAt] = await proofPass.read.verifyCredential([credId]);

      expect(returnedHash).to.equal(dataHash);
      expect(returnedIssuer.toLowerCase()).to.equal(issuer.account.address.toLowerCase());
      
      // Fixes the BigInt vs Number type mismatch error
      expect(issuedAt > 0n).to.be.true; 
    });

    it("Should return default empty values for a non-existent credential", async function () {
      const { proofPass } = await deployFixture();
      
      const [returnedHash, returnedIssuer, issuedAt] = await proofPass.read.verifyCredential(["fake_cred_id"]);

      expect(returnedHash).to.equal("");
      expect(returnedIssuer).to.equal(zeroAddress);
      expect(issuedAt).to.equal(0n);
    });
  });
});

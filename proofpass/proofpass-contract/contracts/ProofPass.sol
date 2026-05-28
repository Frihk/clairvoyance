// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract ProofPass {

    struct CredentialRecord {
        string  dataHash;
        address issuer;
        uint256 issuedAt;
    }

    mapping(string => CredentialRecord) private records;
    mapping(address => bool) public  authorizedIssuers;
    address public owner;

    event CredentialIssued(string credentialId, string dataHash, address issuer);

    constructor() { owner = msg.sender; }

    modifier onlyOwner()  { require(msg.sender == owner, "Not owner"); _; }
    modifier onlyIssuer() { require(authorizedIssuers[msg.sender], "Not issuer"); _; }

    function addIssuer(address issuer) external onlyOwner {
        authorizedIssuers[issuer] = true;
    }

    function issueCredential(
        string calldata credentialId,
        string calldata dataHash
    ) external onlyIssuer {
        require(records[credentialId].issuedAt == 0, "Already issued");
        records[credentialId] = CredentialRecord(dataHash, msg.sender, block.timestamp);
        emit CredentialIssued(credentialId, dataHash, msg.sender);
    }

    function verifyCredential(string calldata credentialId)
        external view returns (string memory, address, uint256)
    {
        CredentialRecord memory r = records[credentialId];
        return (r.dataHash, r.issuer, r.issuedAt);
    }
}

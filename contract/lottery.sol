// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.0;

import '@openzeppelin/contracts/token/ERC721/ERC721.sol';
import '@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol';
import '@openzeppelin/contracts/access/AccessControl.sol';
import '@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol';

contract Modoo is ERC721URIStorage, AccessControl {

    // Overriden contract functions configuration
    function supportsInterface(bytes4 interfaceId) public view virtual override(ERC721, AccessControl) returns (bool) {
        return super.supportsInterface(interfaceId);
    }

    // Token attributes & Admin role
    constructor(string memory _unopenedTokenURI, string memory _openedTokenURI) ERC721("Modoo", "MODO") {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        unopenedTokenURI = _unopenedTokenURI;
        openedTokenURI = _openedTokenURI;
    }

    mapping (uint256 => Attributes) public tokenMap;
    struct Attributes {
        uint256     tokenId;
        string      tokenURI;
        uint256     timestamp;
        uint256     randomNumber;
        bool        isOpened;
    }
    uint256 tokenId;
    string unopenedTokenURI;
    string openedTokenURI;

    // Minting token
    function mintToken(address to) public onlyRole(DEFAULT_ADMIN_ROLE) {
        tokenId++;
        string memory tokenURI = unopenedTokenURI;
        uint randomNumber;
        _mint(to, tokenId);
        _setTokenURI(tokenId, tokenURI);
        Attributes memory newAttributes = Attributes({
            tokenId:        tokenId,
            tokenURI:       tokenURI,
            timestamp:      block.timestamp,
            randomNumber:   randomNumber,
            isOpened:       false
        });
        tokenMap[tokenId] = newAttributes;
    }

    // Purhcase & transfer of token ownership
    function purchaseToken(uint256 _tokenId) public payable {
        require(msg.value == 100, "100 wei required");
        require(_ownerOf(_tokenId) == address(this), "Cannot purchase tokens that are already purchased by others");
        _approve(msg.sender, _tokenId);
        safeTransferFrom(address(this), msg.sender, _tokenId);
    }

    function transferToken(uint256 _tokenId, address to) public payable {
        safeTransferFrom(msg.sender, to, _tokenId);
    }

    // Reveal token
    function revealToken(uint _tokenId) public onlyRole(DEFAULT_ADMIN_ROLE) returns(uint256) {
        tokenMap[_tokenId].isOpened = true;
        tokenMap[_tokenId].tokenURI = openedTokenURI;
        _setTokenURI(_tokenId, openedTokenURI);
        tokenMap[_tokenId].randomNumber = randMod(3); // A random number among 0, 1, 2
        return tokenMap[_tokenId].randomNumber;
    }
    // Random number generation
    uint randNonce = 0;
    function randMod(uint modulus) internal returns(uint) {
        randNonce++;
        return uint(keccak256(abi.encodePacked(block.timestamp, msg.sender, randNonce))) % modulus;
        }

    // Read token information
    function getToken(uint _tokenId) public view returns(uint256, string memory, uint256, uint256, bool) {
        require(tokenMap[_tokenId].tokenId != 0, "Not existing");
        return(
            tokenMap[_tokenId].tokenId,
            tokenMap[_tokenId].tokenURI,
            tokenMap[_tokenId].timestamp,
            tokenMap[_tokenId].randomNumber,
            tokenMap[_tokenId].isOpened
        );
    }

    // Contract receiver
    receive() external payable {}

    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external virtual returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }

}

// SPDX-License-Identifier: MIT

pragma solidity ^0.8.17;

interface IERC20 {
    function totalSupply() external view returns (uint);
    function balanceOf(address account) external view returns (uint);
    function transfer(address recipient, uint amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint);
    function approve(address spender, uint amount) external returns (bool);
    function transferFrom(
        address sender,
        address recipient,
        uint amount
    ) external returns (bool);
    event Transfer(address indexed from, address indexed to, uint value);
    event Approval(address indexed owner, address indexed spender, uint value);
}

contract TotoroToken is IERC20 {
    uint public totalSupply;
    mapping(address => uint) public balanceOf;
    mapping(address => mapping(address => uint)) public allowance;
    string public name = "TOTORO Token";
    string public symbol = "TTR";
    uint8 public decimals = 18;
    address owner;
    mapping(address => bool) welcomeUser; 

    constructor() {
        owner = msg.sender;
        totalSupply += 1000000000e18;
        balanceOf[owner] += 1000000000e18;
    }

    function welcomeToken() external returns (bool) {
        require(welcomeUser[msg.sender] == false, "Aready Welcomed");
        _transferFromOwner(msg.sender, 100e18);
        welcomeUser[msg.sender] = true;
        return true;
    }

    function getSymbol() external view returns (string memory) {
        return symbol;
    }

    function transfer(address recipient, uint amount) external returns (bool) {
        return _transfer(recipient, amount);
    }

    function _transfer(address recipient, uint amount) internal returns (bool) {
        balanceOf[msg.sender] -= amount;
        balanceOf[recipient] += amount;
        emit Transfer(msg.sender, recipient, amount);
        return true;
    }

    function approve(address spender, uint amount) external returns (bool) {
        allowance[msg.sender][spender] = amount;
        emit Approval(msg.sender, spender, amount);
        return true;
    }

    function transferFrom(address sender, address recipient, uint amount) external returns (bool) {
        return _transferFrom(sender, recipient, amount);
    }

    function _transferFrom(address sender, address recipient, uint amount) internal returns (bool) {
        allowance[sender][msg.sender] -= amount;
        balanceOf[sender] -= amount;
        balanceOf[recipient] += amount;
        emit Transfer(sender, recipient, amount);
        return true;
    }

    function _transferFromOwner(address recipient, uint amount) internal returns (bool) {
        balanceOf[owner] -= amount;
        balanceOf[recipient] += amount;
        emit Transfer(owner, recipient, amount);
        return true;
    }

    function _transferToOwner(address sender, uint amount) internal returns (bool) {
        balanceOf[sender] -= amount;
        balanceOf[owner] += amount;
        emit Transfer(sender, owner, amount);
        return true;
    }

    function mint(uint amount) external {
        balanceOf[msg.sender] += amount;
        totalSupply += amount;
        emit Transfer(address(0), msg.sender, amount);
    }

    function burn(uint amount) external {
        balanceOf[msg.sender] -= amount;
        totalSupply -= amount;
        emit Transfer(msg.sender, address(0), amount);
    }
}
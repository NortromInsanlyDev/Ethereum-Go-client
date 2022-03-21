// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract Vault {
  
  address token;

  mapping (address=>uint) userBalance;

  event Deposit(address from, uint256 amount);
  event Withdrawal(address to, uint256 amount);

  constructor (address _token) {
    token = _token;
  }

  function deposit(uint amount) public {
    bool success = IERC20(token).transferFrom(msg.sender, address(this), amount);
    require(success, "Deposit failed!");
    userBalance[msg.sender] += amount;
    emit Deposit(msg.sender, amount);
  }

  function withdraw(uint amount) public {
    require(amount <= userBalance[msg.sender], "Insufficient fund!");
    userBalance[msg.sender] -= amount;
    bool success = IERC20(token).transfer(msg.sender, amount);
    require(success, "Withdrawal failed!");
    emit Withdrawal(msg.sender, amount);
  }
}
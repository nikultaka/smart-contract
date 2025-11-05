// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DataAnchor {
    event Anchored(bytes data);
    event Received(address sender, uint256 amount);

    function anchor(bytes calldata data) external payable {
        emit Anchored(data);
    }

    // Optional: allow receiving ETH
    receive() external payable {
        emit Received(msg.sender, msg.value);
    }

    // Optional fallback
    fallback() external payable {
        emit Received(msg.sender, msg.value);
    }
}

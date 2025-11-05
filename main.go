package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// 1️⃣ Connect to Ethereum node
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/cebeebdc01b34aedafb9cacdbe45637b")
	if err != nil {
		log.Fatal("Failed to connect to Ethereum node:", err)
	}
	defer client.Close()

	// 2️⃣ Load your wallet private key
	privateKey, err := crypto.HexToECDSA("103b65b2af9630593b0e7755e1dc627736d0a85618a52153c4f99eac951e5ae9")
	if err != nil {
		log.Fatal("Invalid private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert public key type")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 3️⃣ Get nonce, gas price, etc.
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 4️⃣ Load contract ABI
	abiData, err := os.ReadFile("DataAnchor.abi")
	if err != nil {
		log.Fatal(err)
	}
	contractABI, err := abi.JSON(stringReader(string(abiData)))
	if err != nil {
		log.Fatal(err)
	}

	// 5️⃣ Contract address (after deploying it)
	contractAddress := common.HexToAddress("0xf8e81D47203A594245E36C48e151709F0C19fBe8")

	// 6️⃣ Prepare custom data (convert string to hex bytes)
	message := "Hello Blockchain!"
	dataHex := hex.EncodeToString([]byte(message))
	dataBytes, _ := hex.DecodeString(dataHex)

	// 7️⃣ Pack the ABI call for `anchor(bytes)`
	input, err := contractABI.Pack("anchor", dataBytes)
	if err != nil {
		log.Fatal("ABI pack error:", err)
	}

	// 8️⃣ Build transaction
	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), 500000, gasPrice, input)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 9️⃣ Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 🔟 Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal("Failed to send tx:", err)
	}

	fmt.Printf("✅ Transaction sent! Hash: %s\n", signedTx.Hash().Hex())

	address := common.HexToAddress("0x554Fd2f5a43F9D292ADF555a83b0246bC7BaEFF3")
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Balance:", new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18)), "ETH")
}

// helper to convert string to io.Reader
func stringReader(s string) *os.File {
	f, _ := os.CreateTemp("", "abi")
	f.WriteString(s)
	f.Seek(0, 0)
	return f
}

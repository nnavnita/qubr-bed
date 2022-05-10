package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	// "errors"
	"fmt"
	"io"
	"os"

	// "strings"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/future"

	// "github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/algorand/go-algorand-sdk/crypto"

	transaction "github.com/algorand/go-algorand-sdk/future"
)

// Utility function that takes a file and returns the sha256 hash value
func hashFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

// TODO: insert aditional utility functions here

func main() {
	// Create account
	account := crypto.GenerateAccount()
	myAddress := account.Address.String()

	fmt.Printf("Pikachu's address: %s\n", myAddress)
	// Fund account
	fmt.Println("Fund Pikachu's new account using testnet faucet:\n--> https://dispenser.testnet.aws.algodev.network?account=" + myAddress)
	fmt.Println("--> Once funded, press ENTER key to continue...")
	fmt.Scanln()

	// Hash the source image file
	fmt.Println("Hashing the source file...")
	imgHash := hashFile("pikachu-nft.png")
	print(string(imgHash))
	imgSRI := "sha256-" + base64.StdEncoding.EncodeToString(imgHash)
	fmt.Printf("--> The SRI of pikachu-nft.png is: '%s'\n\n", imgSRI)

	// Add data to template file
	fmt.Println("Creating metadata.json with Pikachu's asset data..")
	// see metadata.json

	// Hash the metadata.json file
	fmt.Println("Hashing the metadata file...")
	metadataHash := hashFile("metadata.json")
	fmt.Printf("--> The metaDataHash value for metadata.json is: '%s'\n\n", metadataHash)

	// Pin the file to storage platform
	fmt.Println("Pinning files to storage platform...")
	fmt.Println("--> pikachu-nft.png")
	fmt.Println("--> metadata.json")

	// Instantiate algod client
	const algodAddress = "https://academy-algod.dev.aws.algodev.network"
	const algodToken = "2f3203f21e738a1de6110eba6984f9d03e5a95d7a577b34616854064cf2c0e7b"

	algodClient, err := algod.MakeClient(algodAddress, algodToken)
	if err != nil {
		fmt.Printf("Issue with creating algod client: %s\n", err)
		return
	}

	// Create asset
	fmt.Println("Making the assetCreate transaction...")
	txParams, err := algodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		fmt.Printf("Error getting suggested tx params: %s\n", err)
		return
	}
	creator := account.Address.String()
	assetName := "pikachu@arc3"
	unitName := "pikachu"
	assetURL := "https://assets.pokemon.com/assets/cms2/img/pokedex/full/025.png"
	assetMetadataHash := string(metadataHash)
	totalIssuance := uint64(1) // NFTs set totalIssuance to exactly 1
	decimals := uint32(0)      // NFTs set decimals to 0 (not divisible)
	manager := ""
	reserve := ""
	clawback := ""
	freeze := ""
	defaultFrozen := false
	note := []byte(nil)
	txn, err := transaction.MakeAssetCreateTxn(
		creator, note, txParams, totalIssuance, decimals,
		defaultFrozen, manager, reserve, freeze, clawback,
		unitName, assetName, assetURL, assetMetadataHash)
	if err != nil {
		fmt.Printf("Failed to make asset: %s\n", err)
		return
	}

	// sign the transaction
	txid, stx, err := crypto.SignTransaction(account.PrivateKey, txn)
	if err != nil {
		fmt.Printf("Failed to sign transaction: %s\n", err)
		return
	}
	fmt.Printf("Signing transaction ID: %s\n", txid)
	// Broadcast the transaction to the network
	txID, err := algodClient.SendRawTransaction(stx).Do(context.Background())
	if err != nil {
		fmt.Printf("failed to send transaction: %s\n", err)
		return
	}
	fmt.Println("Submitting transaction...")
	fmt.Printf("waiting for confirmation\n")
	// Wait for confirmation
	confirmedTxn, err := future.WaitForConfirmation(algodClient, txID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txID)
		return
	}
	fmt.Printf("Confirmed Transaction: %s in Round %d\n", txID, confirmedTxn.ConfirmedRound)

	assetId := confirmedTxn.AssetIndex
	println("Created assetID:", assetId)

	// Destroy asset
	println("Destroying asset...")
	txn, err = transaction.MakeAssetDestroyTxn(creator, note, txParams, assetId)
	if err != nil {
		fmt.Printf("Failed to destroy asset: %s\n", err)
		return
	}
	txid, stx, err = crypto.SignTransaction(account.PrivateKey, txn)
	txID, err = algodClient.SendRawTransaction(stx).Do(context.Background())

	// Closeout account to dispenser
	println("Closing creator account to dispenser...")
	dispenser := "HZ57J3K46JIJXILONBBZOHX6BKPXEM2VVXNRFSUED6DKFD5ZD24PMJ3MVA"
	txn, err = transaction.MakePaymentTxn(creator, dispenser, 0, nil, dispenser, txParams)
	if err != nil {
		fmt.Printf("Failed to close account: %s\n", err)
		return
	}
	txid, stx, err = crypto.SignTransaction(account.PrivateKey, txn)
	txID, err = algodClient.SendRawTransaction(stx).Do(context.Background())

	// TODO: insert additional codeblocks here
}

package transactiontest

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime/debug"
	"testing"
	"time"

	"github.com/babylonchain-io/bbld/btcjson"
	"github.com/babylonchain-io/bbld/btcutil"
	"github.com/babylonchain-io/bbld/chaincfg"
	"github.com/babylonchain-io/bbld/chaincfg/chainhash"
	"github.com/babylonchain-io/bbld/integration/rpctest"
	"github.com/babylonchain-io/bbld/txscript"
	"github.com/babylonchain-io/bbld/wire"
)

func contains(hashes []*chainhash.Hash, hash *chainhash.Hash) bool {
	for _, h := range hashes {
		if h.IsEqual(hash) {
			return true
		}
	}
	return false
}

func generateRadomBytes(n int) []byte {
	var b = make([]byte, n)
	rand.Read(b)
	return b
}

func generateRandomBytesString(n int) string {
	return hex.EncodeToString(generateRadomBytes(n))
}

func getCoinbaseHash(r *rpctest.Harness, bNum int) (*chainhash.Hash, error) {
	bh, err := r.Client.GetBlockHash(int64(bNum))
	if err != nil {
		return nil, err
	}

	block, err := r.Client.GetBlock(bh)

	if err != nil {
		return nil, err
	}

	// it will panic if block does not have coinbase, but it is fine as every block
	// should have coinbase, if it does not we are doing something seriously wrong
	hash := block.Transactions[0].TxHash()

	return &hash, nil
}

func getBestBlock(r *rpctest.Harness) (*btcutil.Block, error) {
	bh, _, err := r.Client.GetBestBlock()
	if err != nil {
		return nil, err
	}

	block, err := r.Client.GetBlock(bh)

	if err != nil {
		return nil, err
	}

	return btcutil.NewBlock(block), nil
}

func assertInputAgainstCommitment(t *testing.T, input *btcjson.PosDataInput, comm *wire.Commitmment) {
	if input == nil && comm == nil {
		return
	}

	if input == nil && comm != nil {
		t.Fatal("Received commitment for nil input")
	}

	if input != nil && comm == nil {
		t.Fatal("Expected non nil commitment")
	}

	tagBytes, _ := hex.DecodeString(input.Tag)
	if !bytes.Equal(tagBytes, comm.Tag[:]) {
		t.Fatal("Tag bytes of commitment do not match input")
	}

	if input.ProtectionLevel != comm.ProtectionLevel() {
		t.Fatal("Commitment protection level do not match input")
	}

	sigBytes, _ := hex.DecodeString(input.PosSig)
	if !bytes.Equal(sigBytes, comm.PosSig) {
		t.Fatal("Pos signature of commitment do not match input")
	}

	if input.ProtectionLevel == 0 {
		if comm.DataSize > 0 {
			t.Fatal("Commitment data size should be 0 with protection level = 0")
		}

		hashBytes, _ := hex.DecodeString(input.HashOrData)

		if !bytes.Equal(hashBytes, comm.HashCommitment[:]) {
			t.Fatal("Protection level 0, Commitment hash different than provided")
		}

	} else if input.ProtectionLevel == 1 {
		inputBytes, _ := hex.DecodeString(input.HashOrData)
		if int(comm.DataSize) != len(inputBytes) {
			t.Fatal("Commitment data size do not match provided data")
		}

		dataHash := sha256.Sum256(inputBytes)

		if !bytes.Equal(dataHash[:], comm.HashCommitment[:]) {
			t.Fatal("Hash of commitment should be sha256 hash")
		}

	} else {
		t.Fatal("Unexpected input protection level")
	}
}

func testCreateTransaction(r *rpctest.Harness, t *testing.T) {
	// Grab a fresh address from the wallet.
	addr, err := r.NewAddress()
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	coinbaseHash, err := getCoinbaseHash(r, 1)

	if err != nil {
		t.Fatalf("Unable to get coinbase transaction: %v", err)
	}

	tin := btcjson.TransactionInput{Txid: coinbaseHash.String(), Vout: 0}
	am := map[btcutil.Address]btcutil.Amount{addr: 1}
	inps := []btcjson.TransactionInput{tin}

	posDataInput := btcjson.PosDataInput{
		Tag:             generateRandomBytesString(32),
		HashOrData:      generateRandomBytesString(64),
		ProtectionLevel: 1,
		Nonce:           0,
		PosSig:          generateRandomBytesString(64),
	}

	rawTx, err := r.Client.CreateRawTransaction(inps, am, nil, &posDataInput)

	if err != nil {
		t.Fatalf("Unable to create rawTx: %v", err)
	}

	assertInputAgainstCommitment(t, &posDataInput, rawTx.PosCommitment)

	posDataInput1 := btcjson.PosDataInput{
		Tag:             generateRandomBytesString(32),
		HashOrData:      generateRandomBytesString(32),
		ProtectionLevel: 0,
		Nonce:           0,
		PosSig:          generateRandomBytesString(64),
	}

	rawTx1, err := r.Client.CreateRawTransaction(inps, am, nil, &posDataInput1)

	if err != nil {
		t.Fatalf("Unable to create rawTx: %v", err)
	}

	assertInputAgainstCommitment(t, &posDataInput1, rawTx1.PosCommitment)

	// TODO Add error paths
}

func testSendRawTransactionWithCommitment(r *rpctest.Harness, t *testing.T) {
	// Grab a fresh address from the wallet.
	addr, err := r.NewAddress()
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	addrScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	coinbaseHash, err := getCoinbaseHash(r, 1)

	if err != nil {
		t.Fatalf("Unable to get coinbase transaction: %v", err)
	}

	tin := btcjson.TransactionInput{Txid: coinbaseHash.String(), Vout: 0}
	am := map[btcutil.Address]btcutil.Amount{addr: 1e8}
	inps := []btcjson.TransactionInput{tin}

	commData := generateRandomBytesString(64)
	dataBytes, _ := hex.DecodeString(commData)

	posDataInput := btcjson.PosDataInput{
		Tag:             generateRandomBytesString(32),
		HashOrData:      commData,
		ProtectionLevel: 1,
		Nonce:           0,
		PosSig:          generateRandomBytesString(64),
	}

	rawTx, err := r.Client.CreateRawTransaction(inps, am, nil, &posDataInput)

	if err != nil {
		t.Fatalf("Unable to create rawTx: %v", err)
	}

	output := wire.NewTxOut(5e8, addrScript)

	// TODO-BABYLON It is a little bit hacky as we are re-using commitment created by
	// createRawTx endpoint in creating transasction by memory wallet. Other alternative
	// would be to pass coinbase output directly to signing code which also is not pretty.
	testTx, err := r.CreateTransactionWithCommitment([]*wire.TxOut{output}, 10, true, rawTx.PosCommitment)
	if err != nil {
		t.Fatalf("Cannot create tx with commitment: %v", err)
	}

	txHash, err := r.Client.SendRawTransactionWithData(testTx, false, commData)

	if err != nil {
		t.Fatalf("Unable to send raw transaction with data: %v", err)
	}

	mempoolTxs, err := r.Client.GetRawMempool()

	if err != nil {
		t.Fatalf("Unable to get mempool transactions: %v", err)
	}

	if !contains(mempoolTxs, txHash) {
		t.Fatalf("Mempool transactions did not contain sent transaction")
	}

	txns := make([]*btcutil.Tx, 0, 1)
	txns = append(txns, btcutil.NewTx(testTx))

	data := make([][]byte, 0, 1)
	data = append(data, dataBytes)

	// TODO-BABYLON we should probably modify getBlockTemplate endpoint and use it here like
	// t := r.getBlockTemplate()
	// block := r.submitBlock(t.transactions, t.data)
	block, e := r.GenerateAndSubmitBlock(txns, data, -1, time.Time{})

	if e != nil {
		t.Fatalf("Unable to generate block: %v", e)
	}

	// we should have at least two transactions, coinbase + transaction from mempool
	if len(block.Transactions()) != 2 {
		t.Fatal("Block should contain at leas 2 transactions")
	}

	if len(block.MsgBlock().PosData) != 1 {
		t.Fatal("Block should contang at least one piece of data")
	}

	blockTransactionHash := block.Transactions()[1].Hash()

	if !blockTransactionHash.IsEqual(txHash) {
		t.Fatal("Unexpected transaction in mined block")
	}

	mempoolTxs, err = r.Client.GetRawMempool()

	if err != nil {
		t.Fatalf("Unable to get mempool transactions: %v", err)
	}

	// After importing block, transaciton in it should be removed from the mempool
	if len(mempoolTxs) > 0 {
		t.Fatalf("Mempoool should be empty")
	}

	// double check that block was really imported in block chain
	hash, height, err := r.Client.GetBestBlock()

	if err != nil {
		t.Fatalf("Unable to get best block: %v", err)
	}

	if !block.Hash().IsEqual(hash) || block.Height() != height {
		t.Fatal("Best block is not the last imported block")
	}
}

// Check that transactionw with commitment and data is properly propagated
// and then introduced into block
func testPropagateTxWithData(r *rpctest.Harness, t *testing.T) {
	// Create a second harness with only the genesis block so it is behind
	// the main harness.
	harness, err := rpctest.New(&chaincfg.SimNetParams, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := harness.SetUp(false, 0); err != nil {
		t.Fatalf("unable to complete rpctest setup: %v", err)
	}
	defer harness.TearDown()

	nodeSlice := []*rpctest.Harness{r, harness}

	// Establish an outbound connection from the local harness to the main
	// harness and wait for the chains to be synced.
	if err := rpctest.ConnectNode(harness, r); err != nil {
		t.Fatalf("unable to connect harnesses: %v", err)
	}

	if err := rpctest.JoinNodes(nodeSlice, rpctest.Blocks); err != nil {
		t.Fatalf("unable to join node on blocks: %v", err)
	}

	mainHarnessBestHash, _, e := r.Client.GetBestBlock()

	if e != nil {
		t.Fatal("Cannot get main node best block")
	}

	secondHarnessBestHash, _, e := harness.Client.GetBestBlock()

	if e != nil {
		t.Fatal("Cannot get main node best block")
	}

	if !mainHarnessBestHash.IsEqual(secondHarnessBestHash) {
		t.Fatal("Nodes are not synchronized")
	}

	// Grab a fresh address from the wallet.
	addr, err := r.NewAddress()
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	addrScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		t.Fatalf("unable to get new address: %v", err)
	}

	output := wire.NewTxOut(5e8, addrScript)

	data := generateRadomBytes(100)
	dataHash := sha256.Sum256(data)
	tagBytes := generateRadomBytes(32)
	var tag [wire.TagSize]byte
	copy(tag[:], tagBytes)

	sig := generateRadomBytes(64)
	comm := wire.NewTxCommitment(tag, 0, 1, uint32(len(data)), dataHash, 0, sig)

	testTx, err := r.CreateTransactionWithCommitment([]*wire.TxOut{output}, 10, true, comm)
	if err != nil {
		t.Fatalf("Cannot create tx with commitment: %v", err)
	}

	txHash, err := r.Client.SendRawTransactionWithData(testTx, false, hex.EncodeToString(data))

	if err != nil {
		t.Fatalf("Unable to send raw transaction with data: %v", err)
	}

	mempool, err := r.Client.GetRawMempool()

	if err != nil {
		t.Fatalf("unable to get main node mempoool: %v", err)
	}

	if !contains(mempool, txHash) {
		t.Fatal("Main node mempool should contain sent transction")
	}

	if err := rpctest.JoinNodes(nodeSlice, rpctest.Mempools); err != nil {
		t.Fatalf("unable to join node on mempool: %v", err)
	}

	clientNodeMempool, err := harness.Client.GetRawMempool()

	if err != nil {
		t.Fatalf("unable to get main node mempoool: %v", err)
	}

	if !contains(clientNodeMempool, txHash) {
		t.Fatal("Client node mempool should contain sent transction")
	}

	txns := make([]*btcutil.Tx, 0, 1)
	txns = append(txns, btcutil.NewTx(testTx))

	blockData := make([][]byte, 0, 1)
	blockData = append(blockData, data)

	block, err := r.GenerateAndSubmitBlock(txns, blockData, -1, time.Time{})

	if err != nil {
		t.Fatalf("Unable to sumbmit block to main node: %v", err)
	}

	bestBlockHash, _, err := r.Client.GetBestBlock()

	if err != nil {
		t.Fatalf("Unable to get best block of main node: %v", err)
	}

	if !block.Hash().IsEqual(bestBlockHash) {
		t.Fatal("Submitted block is not best block")
	}

	if len(block.Transactions()) != 2 || len(block.MsgBlock().PosData) != 1 {
		t.Fatal("Best block does not have expected number of data or transactions")
	}

	// Now new best block should be propagated to client node
	if err := rpctest.JoinNodes(nodeSlice, rpctest.Blocks); err != nil {
		t.Fatalf("unable to join node on blocks: %v", err)
	}

	clientBestBlock, err := getBestBlock(harness)

	if err != nil {
		t.Fatal("Cannot get client node best block")
	}

	if !bestBlockHash.IsEqual(clientBestBlock.Hash()) {
		t.Fatal("Client best block hash is not equal to main node best block hash")
	}

	clientBestBlockData := clientBestBlock.MsgBlock().PosData[0]
	clientBestBlockCommitment := clientBestBlock.Transactions()[1].MsgTx().PosCommitment

	if !bytes.Equal(clientBestBlockData, data) {
		t.Fatal("Client best block does not have expected data")
	}

	if !bytes.Equal(clientBestBlockCommitment.HashCommitment[:], comm.HashCommitment[:]) {
		t.Fatal("Client best block does not have expected commitment")
	}

}

var rpcTestCases = []rpctest.HarnessTestCase{
	testCreateTransaction,
	testSendRawTransactionWithCommitment,
	testPropagateTxWithData,
}

var primaryHarness *rpctest.Harness

func TestMain(m *testing.M) {
	var err error

	// In order to properly test scenarios on as if we were on mainnet,
	// ensure that non-standard transactions aren't accepted into the
	// mempool or relayed.
	btcdCfg := []string{"--rejectnonstd"}
	primaryHarness, err = rpctest.New(
		&chaincfg.SimNetParams, nil, btcdCfg, "",
	)
	if err != nil {
		fmt.Println("unable to create primary harness: ", err)
		os.Exit(1)
	}

	// Initialize the primary mining node with a chain of length 125,
	// providing 25 mature coinbases to allow spending from for testing
	// purposes.
	if err := primaryHarness.SetUp(true, 25); err != nil {
		fmt.Println("unable to setup test chain: ", err)

		// Even though the harness was not fully setup, it still needs
		// to be torn down to ensure all resources such as temp
		// directories are cleaned up.  The error is intentionally
		// ignored since this is already an error path and nothing else
		// could be done about it anyways.
		_ = primaryHarness.TearDown()
		os.Exit(1)
	}

	exitCode := m.Run()

	// Clean up any active harnesses that are still currently running.This
	// includes removing all temporary directories, and shutting down any
	// created processes.
	if err := rpctest.TearDownAll(); err != nil {
		fmt.Println("unable to tear down all harnesses: ", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func TestRpcTransaction(t *testing.T) {
	var currentTestNum int
	defer func() {
		// If one of the integration tests caused a panic within the main
		// goroutine, then tear down all the harnesses in order to avoid
		// any leaked btcd processes.
		if r := recover(); r != nil {
			fmt.Println("recovering from test panic: ", r)
			if err := rpctest.TearDownAll(); err != nil {
				fmt.Println("unable to tear down all harnesses: ", err)
			}
			t.Fatalf("test #%v panicked: %s", currentTestNum, debug.Stack())
		}
	}()

	for _, testCase := range rpcTestCases {
		testCase(primaryHarness, t)

		currentTestNum++
	}
}

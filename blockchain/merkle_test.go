// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"testing"

	"github.com/babylonchain-io/bbld/btcutil"
)

// TestMerkle tests the BuildMerkleTreeStore API.
func TestMerkle(t *testing.T) {
	t.Skip("Need to mine proper babylon block with transactions")

	block := btcutil.NewBlock(&Block100000)
	merkles := BuildMerkleTreeStore(block.Transactions(), false)
	calculatedMerkleRoot := merkles[len(merkles)-1]
	wantMerkle := &Block100000.Header.MerkleRoot
	if !wantMerkle.IsEqual(calculatedMerkleRoot) {
		t.Errorf("BuildMerkleTreeStore: merkle root mismatch - "+
			"got %v, want %v", calculatedMerkleRoot, wantMerkle)
	}
}

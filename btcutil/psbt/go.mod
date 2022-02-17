module github.com/babylonchain-io/bbld/btcutil/psbt

go 1.17

require (
	github.com/babylonchain-io/bbld v0.0.0-00010101000000-000000000000
	github.com/babylonchain-io/bbld/btcec v0.0.0-00010101000000-000000000000
	github.com/babylonchain-io/bbld/btcutil v0.0.0-00010101000000-000000000000
	github.com/davecgh/go-spew v1.1.1
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
)

replace github.com/babylonchain-io/bbld/btcec => ../../btcec

replace github.com/babylonchain-io/bbld/btcutil => ../

replace github.com/babylonchain-io/bbld => ../..

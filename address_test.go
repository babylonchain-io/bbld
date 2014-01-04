// Copyright (c) 2013 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package btcscript_test

import (
	"bytes"
	"github.com/conformal/btcscript"
	"testing"
)

type addressTest struct {
	script     []byte
	addrhash   []byte
	shouldFail error
	class      btcscript.ScriptType
}

var addressTests = []addressTest{
	{script: []byte{btcscript.OP_DATA_65,
		0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
		0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc,
		0x1e, 0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48,
		0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a,
		0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
		0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b,
		0xfa, 0x9b, 0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f,
		0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56, 0xb4, 0x12,
		0xa3, btcscript.OP_CHECKSIG},
		addrhash: []byte{0x11, 0xb3, 0x66, 0xed, 0xfc, 0x0a,
			0x8b, 0x66, 0xfe, 0xeb, 0xae, 0x5c, 0x2e,
			0x25, 0xa7, 0xb6, 0xa5, 0xd1, 0xcf, 0x31},
		class: btcscript.ScriptPubKey,
	},
	{script: []byte{btcscript.OP_DATA_65,
		0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
		0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc,
		0x1e, 0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48,
		0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a,
		0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
		0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b,
		0xfa, 0x9b, 0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f,
		0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56, 0xb4, 0x12,
		0xa3},
		shouldFail: btcscript.StackErrUnknownAddress,
	},
	{script: []byte{btcscript.OP_DATA_71,
		0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
		0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3,
		0xa1, 0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32,
		0xe9, 0xd6, 0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab,
		0x5f, 0xb8, 0xcd, 0x41, 0x02, 0x20, 0x18, 0x15,
		0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60,
		0xa4, 0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c,
		0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22, 0x08, 0x22,
		0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09, 0x01},
		addrhash: nil,
		class:    btcscript.ScriptPubKey,
	},
	{script: []byte{btcscript.OP_DUP, btcscript.OP_HASH160,
		btcscript.OP_DATA_20,
		0xad, 0x06, 0xdd, 0x6d, 0xde, 0xe5, 0x5c, 0xbc,
		0xa9, 0xa9, 0xe3, 0x71, 0x3b, 0xd7, 0x58, 0x75,
		0x09, 0xa3, 0x05, 0x64,
		btcscript.OP_EQUALVERIFY, btcscript.OP_CHECKSIG,
	},
		addrhash: []byte{0xad, 0x06, 0xdd, 0x6d, 0xde, 0xe5,
			0x5c, 0xbc, 0xa9, 0xa9, 0xe3, 0x71, 0x3b,
			0xd7, 0x58, 0x75, 0x09, 0xa3, 0x05, 0x64},
		class: btcscript.ScriptAddr,
	},
	{script: []byte{btcscript.OP_DATA_73,
		0x30, 0x46, 0x02, 0x21, 0x00, 0xdd, 0xc6, 0x97,
		0x38, 0xbf, 0x23, 0x36, 0x31, 0x8e, 0x4e, 0x04,
		0x1a, 0x5a, 0x77, 0xf3, 0x05, 0xda, 0x87, 0x42,
		0x8a, 0xb1, 0x60, 0x6f, 0x02, 0x32, 0x60, 0x01,
		0x78, 0x54, 0x35, 0x0d, 0xdc, 0x02, 0x21, 0x00,
		0x81, 0x7a, 0xf0, 0x9d, 0x2e, 0xec, 0x36, 0x86,
		0x2d, 0x16, 0x00, 0x98, 0x52, 0xb7, 0xe3, 0xa0,
		0xf6, 0xdd, 0x76, 0x59, 0x82, 0x90, 0xb7, 0x83,
		0x4e, 0x14, 0x53, 0x66, 0x03, 0x67, 0xe0, 0x7a,
		0x01,
		btcscript.OP_DATA_65,
		0x04, 0xcd, 0x42, 0x40, 0xc1, 0x98, 0xe1, 0x25,
		0x23, 0xb6, 0xf9, 0xcb, 0x9f, 0x5b, 0xed, 0x06,
		0xde, 0x1b, 0xa3, 0x7e, 0x96, 0xa1, 0xbb, 0xd1,
		0x37, 0x45, 0xfc, 0xf9, 0xd1, 0x1c, 0x25, 0xb1,
		0xdf, 0xf9, 0xa5, 0x19, 0x67, 0x5d, 0x19, 0x88,
		0x04, 0xba, 0x99, 0x62, 0xd3, 0xec, 0xa2, 0xd5,
		0x93, 0x7d, 0x58, 0xe5, 0xa7, 0x5a, 0x71, 0x04,
		0x2d, 0x40, 0x38, 0x8a, 0x4d, 0x30, 0x7f, 0x88,
		0x7d},
		addrhash: []byte{0x40, 0x92, 0x08, 0xf3, 0x87, 0xf4,
			0x7f, 0xd2, 0x3a, 0x9f, 0x44, 0x5e, 0x14,
			0xdc, 0x1f, 0x99, 0xbb, 0xb8, 0x0d, 0xaa},
		class: btcscript.ScriptAddr,
	},
	{script: []byte{btcscript.OP_DATA_73,
		0x30, 0x46, 0x02, 0x21, 0x00, 0xac, 0x7e, 0x4e,
		0x27, 0xf2, 0xb1, 0x1c, 0xb8, 0x6f, 0xb5, 0xaa,
		0x87, 0x2a, 0xb9, 0xd3, 0x2c, 0xdc, 0x08, 0x33,
		0x80, 0x73, 0x3e, 0x3e, 0x98, 0x47, 0xff, 0x77,
		0xa0, 0x69, 0xcd, 0xdf, 0xab, 0x02, 0x21, 0x00,
		0xc0, 0x4c, 0x3e, 0x6f, 0xfe, 0x88, 0xa1, 0x5b,
		0xc5, 0x07, 0xb8, 0xe5, 0x71, 0xaa, 0x35, 0x92,
		0x8a, 0xcf, 0xe1, 0x5a, 0x4a, 0x23, 0x20, 0x1b,
		0x08, 0xfe, 0x3c, 0x7b, 0x3c, 0x97, 0xc8, 0x8f,
		0x01,
		btcscript.OP_DATA_33,
		0x02, 0x40, 0x05, 0xc9, 0x45, 0xd8, 0x6a, 0xc6,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xa7, 0xa8, 0x45, 0xbd, 0x25, 0x68, 0x9e, 0xdb,
		0x72, 0x3d, 0x5a, 0xd4, 0x06, 0x8d, 0xdd, 0x30,
		0x36,
	},
		addrhash: []byte{0x0c, 0x1b, 0x83, 0xd0, 0x1d, 0x0f,
			0xfb, 0x2b, 0xcc, 0xae, 0x60, 0x69, 0x63,
			0x37, 0x6c, 0xca, 0x38, 0x63, 0xa7, 0xce},
		class: btcscript.ScriptAddr,
	},
	{script: []byte{btcscript.OP_DATA_32,
		0x30, 0x46, 0x02, 0x21, 0x00, 0xac, 0x7e, 0x4e,
		0x27, 0xf2, 0xb1, 0x1c, 0xb8, 0x6f, 0xb5, 0xaa,
		0x87, 0x2a, 0xb9, 0xd3, 0x2c, 0xdc, 0x08, 0x33,
		0x80, 0x73, 0x3e, 0x3e, 0x98, 0x47, 0xff, 0x77,
	},
		addrhash: nil,
		class:    btcscript.ScriptStrange,
	},
	{script: []byte{btcscript.OP_DATA_33,
		0x02, 0x40, 0x05, 0xc9, 0x45, 0xd8, 0x6a, 0xc6,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xa7, 0xa8, 0x45, 0xbd, 0x25, 0x68, 0x9e, 0xdb,
		0x72, 0x3d, 0x5a, 0xd4, 0x06, 0x8d, 0xdd, 0x30,
		0x36,
		btcscript.OP_CHECKSIG,
	},
		addrhash: []byte{0x0c, 0x1b, 0x83, 0xd0, 0x1d, 0x0f,
			0xfb, 0x2b, 0xcc, 0xae, 0x60, 0x69, 0x63,
			0x37, 0x6c, 0xca, 0x38, 0x63, 0xa7, 0xce},
		class: btcscript.ScriptPubKey,
	},
	{script: []byte{btcscript.OP_DATA_33,
		0x02, 0x40, 0x05, 0xc9, 0x45, 0xd8, 0x6a, 0xc6,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xa7, 0xa8, 0x45, 0xbd, 0x25, 0x68, 0x9e, 0xdb,
		0x72, 0x3d, 0x5a, 0xd4, 0x06, 0x8d, 0xdd, 0x30,
		0x36,
		btcscript.OP_CHECKMULTISIG, // note this isn't a real tx
	},
		addrhash: nil,
		class:    btcscript.ScriptStrange,
	},
	{script: []byte{btcscript.OP_0, btcscript.OP_DATA_33,
		0x02, 0x40, 0x05, 0xc9, 0x45, 0xd8, 0x6a, 0xc6,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xa7, 0xa8, 0x45, 0xbd, 0x25, 0x68, 0x9e, 0xdb,
		0x72, 0x3d, 0x5a, 0xd4, 0x06, 0x8d, 0xdd, 0x30,
		0x36, // note this isn't a real tx
	},
		addrhash: nil,
		class:    btcscript.ScriptStrange,
	},
	{script: []byte{btcscript.OP_HASH160, btcscript.OP_DATA_20,
		0x63, 0xbc, 0xc5, 0x65, 0xf9, 0xe6, 0x8e, 0xe0,
		0x18, 0x9d, 0xd5, 0xcc, 0x67, 0xf1, 0xb0, 0xe5,
		0xf0, 0x2f, 0x45, 0xcb,
		btcscript.OP_EQUAL,
	},
		addrhash: []byte{0x63, 0xbc, 0xc5, 0x65, 0xf9, 0xe6, 0x8e, 0xe0,
			0x18, 0x9d, 0xd5, 0xcc, 0x67, 0xf1, 0xb0, 0xe5,
			0xf0, 0x2f, 0x45, 0xcb},
		class: btcscript.ScriptPayToScriptHash,
	},
	{script: []byte{btcscript.OP_DATA_36,
		0x02, 0x40, 0x05, 0xc9, 0x45, 0xd8, 0x6a, 0xc6,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xb0, 0x1f, 0xb0, 0x42, 0x58, 0x34, 0x5a, 0xbe,
		0xa7, 0xa8, 0x45, 0xbd,
		// note this isn't a real tx
	},
		addrhash: nil,
		class:    btcscript.ScriptStrange,
	},
	{script: []byte{},
		shouldFail: btcscript.StackErrUnknownAddress,
	},
}

func TestAddresses(t *testing.T) {
	for i, s := range addressTests {
		class, addrhash, err := btcscript.ScriptToAddrHash(s.script)
		if s.shouldFail != nil {
			if err != s.shouldFail {
				t.Errorf("Address test %v failed is err [%v] should be [%v]", i, err, s.shouldFail)
			}
		} else {
			if err != nil {
				t.Errorf("Address test %v failed err %v", i, err)
			} else {
				if !bytes.Equal(s.addrhash, addrhash) {
					t.Errorf("Address test %v mismatch is [%v] want [%v]", i, addrhash, s.addrhash)
				}
				if s.class != class {
					t.Errorf("Address test %v class mismatch is [%v] want [%v]", i, class, s.class)
				}
			}
		}

	}
}

type stringifyTest struct {
	name       string
	scripttype btcscript.ScriptType
	stringed   string
}

var stringifyTests = []stringifyTest{
	stringifyTest{
		name:       "unknown",
		scripttype: btcscript.ScriptUnknown,
		stringed:   "Unknown",
	},
	stringifyTest{
		name:       "addr",
		scripttype: btcscript.ScriptAddr,
		stringed:   "Addr",
	},
	stringifyTest{
		name:       "pubkey",
		scripttype: btcscript.ScriptPubKey,
		stringed:   "Pubkey",
	},
	stringifyTest{
		name:       "strange",
		scripttype: btcscript.ScriptStrange,
		stringed:   "Strange",
	},
	stringifyTest{
		name:       "generation",
		scripttype: btcscript.ScriptGeneration,
		stringed:   "Generation",
	},
	stringifyTest{
		name:       "not in enum",
		scripttype: btcscript.ScriptType(255),
		stringed:   "Invalid",
	},
}

func TestStringify(t *testing.T) {
	for _, test := range stringifyTests {
		typeString := test.scripttype.String()
		if typeString != test.stringed {
			t.Errorf("%s: got \"%s\" expected \"%s\"", test.name,
				typeString, test.stringed)
		}
	}
}

type multiSigTest struct {
	script     []byte
	reqSigs    int
	addrhashes [][]byte
	shouldFail error
	class      btcscript.ScriptType
}

var multiSigTests = []multiSigTest{
	{script: []byte{},
		class:      btcscript.ScriptUnknown,
		shouldFail: btcscript.StackErrUnknownAddress,
	},
	{script: []byte{btcscript.OP_DATA_65,
		0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
		0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc,
		0x1e, 0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48,
		0x2e, 0xca, 0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a,
		0x5c, 0xb2, 0xe0, 0xea, 0xdd, 0xfb, 0x84, 0xcc,
		0xf9, 0x74, 0x44, 0x64, 0xf8, 0x2e, 0x16, 0x0b,
		0xfa, 0x9b, 0x8b, 0x64, 0xf9, 0xd4, 0xc0, 0x3f,
		0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56, 0xb4, 0x12,
		0xa3, btcscript.OP_CHECKSIG},
		class:      btcscript.ScriptUnknown,
		shouldFail: btcscript.StackErrUnknownAddress,
	},
	{script: []byte{
		btcscript.OP_1, btcscript.OP_DATA_65,
		0x04, 0xcc, 0x71, 0xeb, 0x30, 0xd6, 0x53, 0xc0,
		0xc3, 0x16, 0x39, 0x90, 0xc4, 0x7b, 0x97, 0x6f,
		0x3f, 0xb3, 0xf3, 0x7c, 0xcc, 0xdc, 0xbe, 0xdb,
		0x16, 0x9a, 0x1d, 0xfe, 0xf5, 0x8b, 0xbf, 0xbf,
		0xaf, 0xf7, 0xd8, 0xa4, 0x73, 0xe7, 0xe2, 0xe6,
		0xd3, 0x17, 0xb8, 0x7b, 0xaf, 0xe8, 0xbd, 0xe9,
		0x7e, 0x3c, 0xf8, 0xf0, 0x65, 0xde, 0xc0, 0x22,
		0xb5, 0x1d, 0x11, 0xfc, 0xdd, 0x0d, 0x34, 0x8a,
		0xc4, btcscript.OP_DATA_65,
		0x04, 0x61, 0xcb, 0xdc, 0xc5, 0x40, 0x9f, 0xb4,
		0xb4, 0xd4, 0x2b, 0x51, 0xd3, 0x33, 0x81, 0x35,
		0x4d, 0x80, 0xe5, 0x50, 0x07, 0x8c, 0xb5, 0x32,
		0xa3, 0x4b, 0xfa, 0x2f, 0xcf, 0xde, 0xb7, 0xd7,
		0x65, 0x19, 0xae, 0xcc, 0x62, 0x77, 0x0f, 0x5b,
		0x0e, 0x4e, 0xf8, 0x55, 0x19, 0x46, 0xd8, 0xa5,
		0x40, 0x91, 0x1a, 0xbe, 0x3e, 0x78, 0x54, 0xa2,
		0x6f, 0x39, 0xf5, 0x8b, 0x25, 0xc1, 0x53, 0x42,
		0xaf, btcscript.OP_2, btcscript.OP_CHECKMULTISIG},
		class:   btcscript.ScriptMultiSig,
		reqSigs: 1,
		addrhashes: [][]byte{
			[]byte{
				0x66, 0x0d, 0x4e, 0xf3, 0xa7, 0x43, 0xe3, 0xe6,
				0x96, 0xad, 0x99, 0x03, 0x64, 0xe5, 0x55, 0xc2,
				0x71, 0xad, 0x50, 0x4b,
			},
			[]byte{
				0x64, 0x1a, 0xd5, 0x05, 0x1e, 0xdd, 0x97, 0x02,
				0x9a, 0x00, 0x3f, 0xe9, 0xef, 0xb2, 0x93, 0x59,
				0xfc, 0xee, 0x40, 0x9d,
			},
		},
	},

	// from real tx 60a20bd93aa49ab4b28d514ec10b06e1829ce6818ec06cd3aabd013ebcdc4bb1, vout 0
	{script: []byte{
		btcscript.OP_1, btcscript.OP_DATA_65,
		0x04, 0xcc, 0x71, 0xeb, 0x30, 0xd6, 0x53, 0xc0,
		0xc3, 0x16, 0x39, 0x90, 0xc4, 0x7b, 0x97, 0x6f,
		0x3f, 0xb3, 0xf3, 0x7c, 0xcc, 0xdc, 0xbe, 0xdb,
		0x16, 0x9a, 0x1d, 0xfe, 0xf5, 0x8b, 0xbf, 0xbf,
		0xaf, 0xf7, 0xd8, 0xa4, 0x73, 0xe7, 0xe2, 0xe6,
		0xd3, 0x17, 0xb8, 0x7b, 0xaf, 0xe8, 0xbd, 0xe9,
		0x7e, 0x3c, 0xf8, 0xf0, 0x65, 0xde, 0xc0, 0x22,
		0xb5, 0x1d, 0x11, 0xfc, 0xdd, 0x0d, 0x34, 0x8a,
		0xc4, btcscript.OP_DATA_65,
		0x04, 0x61, 0xcb, 0xdc, 0xc5, 0x40, 0x9f, 0xb4,
		0xb4, 0xd4, 0x2b, 0x51, 0xd3, 0x33, 0x81, 0x35,
		0x4d, 0x80, 0xe5, 0x50, 0x07, 0x8c, 0xb5, 0x32,
		0xa3, 0x4b, 0xfa, 0x2f, 0xcf, 0xde, 0xb7, 0xd7,
		0x65, 0x19, 0xae, 0xcc, 0x62, 0x77, 0x0f, 0x5b,
		0x0e, 0x4e, 0xf8, 0x55, 0x19, 0x46, 0xd8, 0xa5,
		0x40, 0x91, 0x1a, 0xbe, 0x3e, 0x78, 0x54, 0xa2,
		0x6f, 0x39, 0xf5, 0x8b, 0x25, 0xc1, 0x53, 0x42,
		0xaf, btcscript.OP_2, btcscript.OP_CHECKMULTISIG},
		class:   btcscript.ScriptMultiSig,
		reqSigs: 1,
		addrhashes: [][]byte{
			[]byte{
				0x66, 0x0d, 0x4e, 0xf3, 0xa7, 0x43, 0xe3, 0xe6,
				0x96, 0xad, 0x99, 0x03, 0x64, 0xe5, 0x55, 0xc2,
				0x71, 0xad, 0x50, 0x4b,
			},
			[]byte{
				0x64, 0x1a, 0xd5, 0x05, 0x1e, 0xdd, 0x97, 0x02,
				0x9a, 0x00, 0x3f, 0xe9, 0xef, 0xb2, 0x93, 0x59,
				0xfc, 0xee, 0x40, 0x9d,
			},
		},
	},

	// from real tx: 691dd277dc0e90a462a3d652a1171686de49cf19067cd33c7df0392833fb986a, vout 0
	{script: []byte{
		btcscript.OP_1, btcscript.OP_DATA_65,
		0x1c, 0x22, 0x00, 0x00, 0x73, 0x53, 0x45, 0x58,
		0x57, 0x69, 0x6b, 0x69, 0x6c, 0x65, 0x61, 0x6b,
		0x73, 0x20, 0x43, 0x61, 0x62, 0x6c, 0x65, 0x67,
		0x61, 0x74, 0x65, 0x20, 0x42, 0x61, 0x63, 0x6b,
		0x75, 0x70, 0x0a, 0x0a, 0x63, 0x61, 0x62, 0x6c,
		0x65, 0x67, 0x61, 0x74, 0x65, 0x2d, 0x32, 0x30,
		0x31, 0x30, 0x31, 0x32, 0x30, 0x34, 0x31, 0x38,
		0x31, 0x31, 0x2e, 0x37, 0x7a, 0x0a, 0x0a, 0x44,
		0x6f, btcscript.OP_DATA_65,
		0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x20, 0x74,
		0x68, 0x65, 0x20, 0x66, 0x6f, 0x6c, 0x6c, 0x6f,
		0x77, 0x69, 0x6e, 0x67, 0x20, 0x74, 0x72, 0x61,
		0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
		0x73, 0x20, 0x77, 0x69, 0x74, 0x68, 0x20, 0x53,
		0x61, 0x74, 0x6f, 0x73, 0x68, 0x69, 0x20, 0x4e,
		0x61, 0x6b, 0x61, 0x6d, 0x6f, 0x74, 0x6f, 0x27,
		0x73, 0x20, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f,
		0x61, btcscript.OP_DATA_65,
		0x64, 0x20, 0x74, 0x6f, 0x6f, 0x6c, 0x20, 0x77,
		0x68, 0x69, 0x63, 0x68, 0x0a, 0x63, 0x61, 0x6e,
		0x20, 0x62, 0x65, 0x20, 0x66, 0x6f, 0x75, 0x6e,
		0x64, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x72, 0x61,
		0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
		0x20, 0x36, 0x63, 0x35, 0x33, 0x63, 0x64, 0x39,
		0x38, 0x37, 0x31, 0x31, 0x39, 0x65, 0x66, 0x37,
		0x39, 0x37, 0x64, 0x35, 0x61, 0x64, 0x63, 0x63,
		0x64, btcscript.OP_3, btcscript.OP_CHECKMULTISIG},
		class:   btcscript.ScriptMultiSig,
		reqSigs: 1,
		addrhashes: [][]byte{
			[]byte{
				0x7c, 0xe9, 0x23, 0x01, 0x00, 0x75, 0x98, 0xc8,
				0xb9, 0x31, 0x5b, 0x49, 0x71, 0x21, 0xa0, 0x18,
				0x20, 0x90, 0xa3, 0x0c,
			},
			[]byte{
				0x76, 0xf6, 0x95, 0x05, 0xb3, 0x9f, 0xba, 0x56,
				0x6d, 0xe8, 0x40, 0x89, 0xda, 0xd8, 0xbd, 0x35,
				0x01, 0xec, 0x47, 0xed,
			},
			[]byte{
				0x7b, 0xfb, 0xfe, 0x5d, 0x66, 0x80, 0x22, 0x2e,
				0xe6, 0x38, 0xf9, 0x41, 0xdc, 0xfa, 0x61, 0x6d,
				0xd4, 0x5f, 0x11, 0xee,
			},
		},
	},

	// from real tx: 691dd277dc0e90a462a3d652a1171686de49cf19067cd33c7df0392833fb986a, vout 44
	{script: []byte{
		btcscript.OP_1, btcscript.OP_DATA_65,
		0x34, 0x63, 0x33, 0x65, 0x63, 0x32, 0x35, 0x39,
		0x63, 0x37, 0x34, 0x64, 0x61, 0x63, 0x65, 0x36,
		0x66, 0x64, 0x30, 0x38, 0x38, 0x62, 0x34, 0x34,
		0x63, 0x65, 0x66, 0x38, 0x63, 0x0a, 0x63, 0x36,
		0x36, 0x62, 0x63, 0x31, 0x39, 0x39, 0x36, 0x63,
		0x38, 0x62, 0x39, 0x34, 0x61, 0x33, 0x38, 0x31,
		0x31, 0x62, 0x33, 0x36, 0x35, 0x36, 0x31, 0x38,
		0x66, 0x65, 0x31, 0x65, 0x39, 0x62, 0x31, 0x62,
		0x35, btcscript.OP_DATA_65,
		0x36, 0x63, 0x61, 0x63, 0x63, 0x65, 0x39, 0x39,
		0x33, 0x61, 0x33, 0x39, 0x38, 0x38, 0x61, 0x34,
		0x36, 0x39, 0x66, 0x63, 0x63, 0x36, 0x64, 0x36,
		0x64, 0x61, 0x62, 0x66, 0x64, 0x0a, 0x32, 0x36,
		0x36, 0x33, 0x63, 0x66, 0x61, 0x39, 0x63, 0x66,
		0x34, 0x63, 0x30, 0x33, 0x63, 0x36, 0x30, 0x39,
		0x63, 0x35, 0x39, 0x33, 0x63, 0x33, 0x65, 0x39,
		0x31, 0x66, 0x65, 0x64, 0x65, 0x37, 0x30, 0x32,
		0x39, btcscript.OP_DATA_33,
		0x31, 0x32, 0x33, 0x64, 0x64, 0x34, 0x32, 0x64,
		0x32, 0x35, 0x36, 0x33, 0x39, 0x64, 0x33, 0x38,
		0x61, 0x36, 0x63, 0x66, 0x35, 0x30, 0x61, 0x62,
		0x34, 0x63, 0x64, 0x34, 0x34, 0x0a, 0x00, 0x00,
		0x00, btcscript.OP_3, btcscript.OP_CHECKMULTISIG},
		class:   btcscript.ScriptMultiSig,
		reqSigs: 1,
		addrhashes: [][]byte{
			[]byte{
				0xfb, 0xf0, 0x08, 0x0b, 0xc5, 0xf9, 0xd7, 0x2a,
				0x9e, 0x64, 0x6f, 0x16, 0x46, 0x46, 0x1c, 0x43,
				0x19, 0xc3, 0xb6, 0xd4,
			},
			[]byte{
				0xc6, 0x00, 0xe7, 0x69, 0xc3, 0xae, 0x20, 0xd4,
				0xa0, 0x50, 0x08, 0xd1, 0xe3, 0xad, 0x06, 0x33,
				0xf2, 0x7b, 0x77, 0xa2,
			},
			[]byte{
				0xad, 0x34, 0x62, 0xcb, 0xa3, 0x5b, 0xee, 0x04,
				0xef, 0xd4, 0x20, 0x8c, 0xcd, 0x7f, 0x41, 0xf4,
				0xc8, 0x55, 0xf2, 0x73,
			},
		},
	},
}

func TestMultiSigs(t *testing.T) {
	for i, s := range multiSigTests {
		class, reqSigs, addrhashes, err := btcscript.ScriptToAddrHashes(s.script)
		if s.shouldFail != nil {
			if err != s.shouldFail {
				t.Errorf("MultiSig test %v failed is err [%v] should be [%v]", i, err, s.shouldFail)
			}
		} else {
			if err != nil {
				t.Errorf("MultiSig test %v failed err %v", i, err)
			} else {
				if s.class != class {
					t.Errorf("MultiSig test %v class mismatch is [%v] want [%v]", i, class, s.class)
				}
				if len(addrhashes) != len(s.addrhashes) {
					t.Errorf("MultiSig test %v num addrhashes expected is [%d] want [%d]", i, len(addrhashes), len(s.addrhashes))
				}
				if reqSigs != s.reqSigs {
					t.Errorf("MultiSig test %v reqSigs expected is [%d] want [%d]", i, reqSigs, s.reqSigs)
				}
				for j := 0; j < len(addrhashes); j++ {
					if !bytes.Equal(s.addrhashes[j], addrhashes[j]) {
						t.Errorf("MultiSig test %v addrhash %d mismatch is [%v] want [%v]", i, j, addrhashes[j], s.addrhashes[j])
					}
				}
			}
		}

	}
}

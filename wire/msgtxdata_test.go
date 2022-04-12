package wire

// TODO decide on licensing

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// TestTxDataWire tests the MsgTxData for various data sizes
func TestTxDataWire(t *testing.T) {
	tests := []struct {
		dataLen uint32
	}{
		// Latest protocol version with no transactions.
		{0},
		{20},
		{5000},
		{50000},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		var tx = *multiTx
		var commitment = generaterateRandomCommitment(1, 1, test.dataLen, 0, 96)
		tx.PosCommitment = commitment
		var data = make([]byte, test.dataLen)
		rand.Read(data)

		var txWithData = NewMsgTxData(&tx, data)

		// Encode the message to wire format.
		var buf bytes.Buffer
		err := txWithData.BtcEncode(&buf, 0, WitnessEncoding)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}

		// Decode the message from wire format.
		var msg MsgTxData
		rbuf := bytes.NewReader(buf.Bytes())
		err = msg.BtcDecode(rbuf, 0, WitnessEncoding)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, txWithData) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(&msg), spew.Sdump(txWithData))
			continue
		}
	}
}

func TestTxDataWireTooLargeData(t *testing.T) {
	noTx := NewMsgTx(1)
	noTx.Version = 1
	noTxEncoded := []byte{
		0x01, 0x00, 0x00, 0x00, // Version
		0x00,                   // Varint for number of input transactions
		0x00,                   // Varint for number of output transactions
		0x00, 0x00, 0x00, 0x00, // Lock time
		0x00, // No commitment
		0xfd, 0x51, 0xc3,
	}

	var msg MsgTxData
	rbuf := bytes.NewReader(noTxEncoded)
	err := msg.BtcDecode(rbuf, 0, WitnessEncoding)

	var expectedError = &MessageError{}

	if reflect.TypeOf(err) != reflect.TypeOf(expectedError) {
		t.Errorf("MsgTxData decode, wrong error type, expected MessageError")
	}
}

package wire

// TODO-Babylon: decide on licensing

import (
	"fmt"
	"io"
)

type MsgTxData struct {
	Tx   MsgTx
	Data []byte
}

func (msg *MsgTxData) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	err := msg.Tx.BtcDecode(r, pver, enc)

	if err != nil {
		fmt.Println("Err decoding tx")
		return err
	}

	// TODO-Babylon: decide if we want to validate datasize against transaciton commitent here
	msg.Data, err = ReadVarBytes(r, pver, MaxPosDataSize, "txdata data")

	return err
}

func (msg *MsgTxData) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	size := len(msg.Data)
	if size > MaxPosDataSize {
		str := fmt.Sprintf("txdata size too large for message "+
			"[size %v, max %v]", size, MaxPosDataSize)
		return messageError("MsgTxData.BtcEncode", str)
	}

	err := msg.Tx.BtcEncode(w, pver, enc)

	if err != nil {
		return err
	}

	// TODO-Babylon: decide if we want to validate datasize against transaciton commitent here
	return WriteVarBytes(w, pver, msg.Data)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgTxData) Command() string {
	return CmdTxData
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgTxData) MaxPayloadLength(pver uint32) uint32 {
	return msg.Tx.MaxPayloadLength(pver) + uint32(VarIntSerializeSize(MaxPosDataSize)) +
		MaxPosDataSize
}

// NewMsgTxData returns a new bitcoin block message that conforms to the
// Message interface
func NewMsgTxData(tx *MsgTx, data []byte) *MsgTxData {
	return &MsgTxData{
		Tx:   *tx,
		Data: data,
	}
}

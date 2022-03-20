package wire

// All constant values related to commitments and data stored and relayed by
// babylon nodes
const (
	// Maximal size of tag stored in commitment
	TagSize = 32

	// TODO define what should be final data size
	MaxPosDataSize = 50000

	// Maximum size of Proof of stake chain signature
	MaxPosSigSize = 128

	CurrentCommitmentVersion = 0
)

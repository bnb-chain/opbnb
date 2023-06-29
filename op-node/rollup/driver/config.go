package driver

type Config struct {
	// VerifierConfDepth is the distance to keep from the L1 head when reading L1 data for L2 derivation.
	VerifierConfDepth uint64 `json:"verifier_conf_depth"`

	// SequencerConfDepth is the distance to keep from the L1 head as origin when sequencing new L2 blocks.
	// If this distance is too large, the sequencer may:
	// - not adopt a L1 origin within the allowed time (rollup.Config.MaxSequencerDrift)
	// - not adopt a L1 origin that can be included on L1 within the allowed range (rollup.Config.SeqWindowSize)
	// and thus fail to produce a block with anything more than deposits.
	SequencerConfDepth uint64 `json:"sequencer_conf_depth"`

	// SequencerEnabled is true when the driver should sequence new blocks.
	SequencerEnabled bool `json:"sequencer_enabled"`

	// SequencerStopped is false when the driver should sequence new blocks.
	SequencerStopped bool `json:"sequencer_stopped"`

	// SequencerMaxSafeLag is the maximum number of L2 blocks for restricting the distance between L2 safe and unsafe.
	// Disabled if 0.
	SequencerMaxSafeLag uint64 `json:"sequencer_max_safe_lag"`

	// SequencerCoordinatorEnabled is true when the driver should request permission from op-coordinator before
	// building new blocks.
	SequencerCoordinatorEnabled bool `json:"sequencer_coordinator_enabled"`

	// SequencerCoordinatorAddr is address of the Coordinator JSON-RPC endpoint.
	SequencerCoordinatorAddr string `json:"sequencer_coordinator_addr"`

	// SequencerCoordinatorSequencerId is used by Coordinator to distinguish different sequencers. It must be unique
	// and same as the name of the sequencer node configured in the Coordinator service.
	SequencerCoordinatorSequencerId string `json:"sequencer_coordinator_sequencer_id"`
}

package consensus

type AggregateAndProof struct {
	Index          uint64       `json:"aggregator_index"`
	Aggregate      *Attestation `json:"aggregate"`
	SelectionProof [96]byte     `json:"selection_proof" ssz-size:"96"`
}

type Checkpoint struct {
	Epoch uint64 `json:"epoch"`
	Root  Root   `json:"root" ssz-size:"32"`
}

type AttestationData struct {
	Slot            uint64      `json:"slot"`
	Index           uint64      `json:"index"`
	BeaconBlockHash [32]byte    `json:"beacon_block_root" ssz-size:"32"`
	Source          *Checkpoint `json:"source"`
	Target          *Checkpoint `json:"target"`
}

type Attestation struct {
	AggregationBits []byte           `json:"aggregation_bits" ssz:"bitlist" ssz-max:"2048"`
	Data            *AttestationData `json:"data"`
	Signature       Signature        `json:"signature" ssz-size:"96"`
}

type DepositData struct {
	Pubkey                [48]byte  `json:"pubkey" ssz-size:"48"`
	WithdrawalCredentials [32]byte  `json:"withdrawal_credentials" ssz-size:"32"`
	Amount                uint64    `json:"amount"`
	Signature             Signature `json:"signature" ssz-size:"96"`
	Root                  [32]byte  `ssz:"-"`
}

type Deposit struct {
	Proof [33][32]byte `ssz-size:"33,32"`
	Data  *DepositData
}

type DepositMessage struct {
	Pubkey                [48]byte `json:"pubkey" ssz-size:"48"`
	WithdrawalCredentials [32]byte `json:"withdrawal_credentials" ssz-size:"32"`
	Amount                uint64   `json:"amount"`
}

type IndexedAttestation struct {
	AttestationIndices []uint64         `json:"attesting_indices" ssz-max:"2048"`
	Data               *AttestationData `json:"data"`
	Signature          Signature        `json:"signature" ssz-size:"96"`
}

type PendingAttestation struct {
	AggregationBits []byte           `json:"aggregation_bits" ssz:"bitlist" ssz-max:"2048"`
	Data            *AttestationData `json:"data"`
	InclusionDelay  uint64           `json:"inclusion_delay"`
	ProposerIndex   uint64           `json:"proposer_index"`
}

type Fork struct {
	PreviousVersion [4]byte `json:"previous_version" ssz-size:"4"`
	CurrentVersion  [4]byte `json:"current_version" ssz-size:"4"`
	Epoch           uint64  `json:"epoch"`
}

type Validator struct {
	Pubkey                     [48]byte `json:"pubkey" ssz-size:"48"`
	WithdrawalCredentials      [32]byte `json:"withdrawal_credentials" ssz-size:"32"`
	EffectiveBalance           uint64   `json:"effective_balance"`
	Slashed                    bool     `json:"slashed"`
	ActivationEligibilityEpoch uint64   `json:"activation_eligibility_epoch"`
	ActivationEpoch            uint64   `json:"activation_epoch"`
	ExitEpoch                  uint64   `json:"exit_epoch"`
	WithdrawableEpoch          uint64   `json:"withdrawable_epoch"`
}

type VoluntaryExit struct {
	Epoch          uint64 `json:"epoch"`
	ValidatorIndex uint64 `json:"validator_index"`
}

type SignedVoluntaryExit struct {
	Exit      *VoluntaryExit `json:"message"`
	Signature Signature      `json:"signature" ssz-size:"96"`
}

type Eth1Data struct {
	DepositRoot  Root     `json:"deposit_root" ssz-size:"32"`
	DepositCount uint64   `json:"deposit_count"`
	BlockHash    [32]byte `json:"block_hash" ssz-size:"32"`
}

type SigningRoot struct {
	ObjectRoot Root   `json:"object_root" ssz-size:"32"`
	Domain     []byte `json:"domain" ssz-size:"8"`
}

type ProposerSlashing struct {
	Header1 *SignedBeaconBlockHeader `json:"signed_header_1"`
	Header2 *SignedBeaconBlockHeader `json:"signed_header_2"`
}

type AttesterSlashing struct {
	Attestation1 *IndexedAttestation `json:"attestation_1"`
	Attestation2 *IndexedAttestation `json:"attestation_2"`
}

type BeaconBlock struct {
	Slot          uint64                 `json:"slot"`
	ProposerIndex uint64                 `json:"proposer_index"`
	ParentRoot    Root                   `json:"parent_root" ssz-size:"32"`
	StateRoot     Root                   `json:"state_root" ssz-size:"32"`
	Body          *BeaconBlockBodyAltair `json:"body"`
}

type SignedBeaconBlock struct {
	Block     *BeaconBlock `json:"message"`
	Signature Signature    `json:"signature" ssz-size:"96"`
}

type Transfer struct {
	Sender    uint64    `json:"sender"`
	Recipient uint64    `json:"recipient"`
	Amount    uint64    `json:"amount"`
	Fee       uint64    `json:"fee"`
	Slot      uint64    `json:"slot"`
	Pubkey    [48]byte  `json:"pubkey" ssz-size:"48"`
	Signature Signature `json:"signature" ssz-size:"96"`
}

type BeaconBlockBody struct {
	RandaoReveal      Signature              `json:"randao_reveal" ssz-size:"96"`
	Eth1Data          *Eth1Data              `json:"eth1_data"`
	Graffiti          [32]byte               `json:"graffiti" ssz-size:"32"`
	ProposerSlashings []*ProposerSlashing    `json:"proposer_slashings" ssz-max:"16"`
	AttesterSlashings []*AttesterSlashing    `json:"attester_slashings" ssz-max:"2"`
	Attestations      []*Attestation         `json:"attestations" ssz-max:"128"`
	Deposits          []*Deposit             `json:"deposits" ssz-max:"16"`
	VoluntaryExits    []*SignedVoluntaryExit `json:"voluntary_exits" ssz-max:"16"`
}

type SignedBeaconBlockHeader struct {
	Header    *BeaconBlockHeader `json:"message"`
	Signature Signature          `json:"signature" ssz-size:"96"`
}

type BeaconBlockHeader struct {
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposer_index"`
	ParentRoot    Root   `json:"parent_root" ssz-size:"32"`
	StateRoot     Root   `json:"state_root" ssz-size:"32"`
	BodyRoot      Root   `json:"body_root" ssz-size:"32"`
}

type ForkData struct {
	CurrentVersion        [4]byte `json:"current_version" ssz-size:"4"`
	GenesisValidatorsRoot Root    `json:"genesis_validators_root" ssz-size:"32"`
}

type SigningData struct {
	ObjectRoot Root     `json:"object_root" ssz-size:"32"`
	Domain     [32]byte `json:"domain" ssz-size:"32"`
}

// Altair fork

type SignedBeaconBlockAltair struct {
	Block     *BeaconBlockAltair `json:"message"`
	Signature Signature          `json:"signature" ssz-size:"96"`
}

type BeaconBlockAltair struct {
	Slot          uint64                 `json:"slot"`
	ProposerIndex uint64                 `json:"proposer_index"`
	ParentRoot    Root                   `json:"parent_root" ssz-size:"32"`
	StateRoot     Root                   `json:"state_root" ssz-size:"32"`
	Body          *BeaconBlockBodyAltair `json:"body"`
}

type BeaconBlockBodyAltair struct {
	RandaoReveal      Signature              `json:"randao_reveal" ssz-size:"96"`
	Eth1Data          *Eth1Data              `json:"eth1_data"`
	Graffiti          [32]byte               `json:"graffiti" ssz-size:"32"`
	ProposerSlashings []*ProposerSlashing    `json:"proposer_slashings" ssz-max:"16"`
	AttesterSlashings []*AttesterSlashing    `json:"attester_slashings" ssz-max:"2"`
	Attestations      []*Attestation         `json:"attestations" ssz-max:"128"`
	Deposits          []*Deposit             `json:"deposits" ssz-max:"16"`
	VoluntaryExits    []*SignedVoluntaryExit `json:"voluntary_exits" ssz-max:"16"`
	SyncAggregate     *SyncAggregate         `json:"sync_aggregate"`
}

type SyncAggregate struct {
	SyncCommiteeBits      [64]byte  `json:"sync_committee_bits" ssz-size:"64"`
	SyncCommiteeSignature Signature `json:"sync_committee_signature" ssz-size:"96"`
}

type SyncCommittee struct {
	PubKeys         [512][48]byte `json:"pubkeys" ssz-size:"512,48"`
	AggregatePubKey [48]byte      `json:"aggregate_pubkey" ssz-size:"48"`
}

// bellatrix

type BeaconBlockBodyBellatrix struct {
	RandaoReveal      Signature              `json:"randao_reveal" ssz-size:"96"`
	Eth1Data          *Eth1Data              `json:"eth1_data"`
	Graffiti          [32]byte               `json:"graffiti" ssz-size:"32"`
	ProposerSlashings []*ProposerSlashing    `json:"proposer_slashings" ssz-max:"16"`
	AttesterSlashings []*AttesterSlashing    `json:"attester_slashings" ssz-max:"2"`
	Attestations      []*Attestation         `json:"attestations" ssz-max:"128"`
	Deposits          []*Deposit             `json:"deposits" ssz-max:"16"`
	VoluntaryExits    []*SignedVoluntaryExit `json:"voluntary_exits" ssz-max:"16"`
	SyncAggregate     *SyncAggregate         `json:"sync_aggregate"`
	ExecutionPayload  *ExecutionPayload      `json:"execution_payload"`
}

type ExecutionPayload struct {
	ParentHash    [32]byte  `ssz-size:"32" json:"parent_hash"`
	FeeRecipient  [20]byte  `ssz-size:"20" json:"fee_recipient"`
	StateRoot     [32]byte  `ssz-size:"32" json:"state_root"`
	ReceiptsRoot  [32]byte  `ssz-size:"32" json:"receipts_root"`
	LogsBloom     [256]byte `ssz-size:"256" json:"logs_bloom"`
	PrevRandao    [32]byte  `ssz-size:"32" json:"prev_randao"`
	BlockNumber   uint64    `json:"block_number"`
	GasLimit      uint64    `json:"gas_limit"`
	GasUsed       uint64    `json:"gas_used"`
	Timestamp     uint64    `json:"timestamp"`
	ExtraData     []byte    `ssz-max:"32" json:"extra_data"`
	BaseFeePerGas [32]byte  `ssz-size:"32" json:"base_fee_per_gas"`
	BlockHash     [32]byte  `ssz-size:"32" json:"block_hash"`
	Transactions  [][]byte  `ssz-max:"1048576,1073741824" ssz-size:"?,?" json:"transactions"`
}

type ExecutionPayloadHeader struct {
	ParentHash       [32]byte  `json:"parent_hash" ssz-size:"32"`
	FeeRecipient     [20]byte  `json:"fee_recipient" ssz-size:"20"`
	StateRoot        [32]byte  `json:"state_root" ssz-size:"32"`
	ReceiptsRoot     [32]byte  `json:"receipts_root" ssz-size:"32"`
	LogsBloom        [256]byte `json:"logs_bloom" ssz-size:"256"`
	PrevRandao       [32]byte  `json:"prev_randao" ssz-size:"32"`
	BlockNumber      uint64    `json:"block_number"`
	GasLimit         uint64    `json:"gas_limit"`
	GasUsed          uint64    `json:"gas_used"`
	Timestamp        uint64    `json:"timestamp"`
	ExtraData        []byte    `json:"extra_data" ssz-max:"32"`
	BaseFeePerGas    [32]byte  `json:"base_fee_per_gas" ssz-size:"32"`
	BlockHash        [32]byte  `json:"block_hash" ssz-size:"32"`
	TransactionsRoot [32]byte  `json:"transactions_root" ssz-size:"32"`
}

type SyncAggregatorSelectionData struct {
	Slot              uint64 `json:"slot"`
	SubCommitteeIndex uint64 `json:"subcommittee_index"`
}

// SyncCommitteeContribution is the Ethereum 2 sync committee contribution structure.
type SyncCommitteeContribution struct {
	Slot              uint64    `json:"slot"`
	BeaconBlockRoot   Root      `json:"beacon_block_root" ssz-size:"32"`
	SubcommitteeIndex uint64    `json:"subcommittee_index"`
	AggregationBits   []byte    `json:"aggregation_bits" ssz-size:"16"` // bitvector
	Signature         Signature `json:"signature" ssz-size:"96"`
}

type ContributionAndProof struct {
	AggregatorIndex uint64                     `json:"aggregator_index"`
	Contribution    *SyncCommitteeContribution `json:"contribution"`
	SelectionProof  Signature                  `json:"selection_proof" ssz-size:"96"`
}

type SignedContributionAndProof struct {
	Message   *ContributionAndProof `json:"message"`
	Signature Signature             `json:"signature" ssz-size:"96"`
}

type SyncCommitteeMessage struct {
	Slot           uint64    `json:"slot"`
	BlockRoot      Root      `json:"beacon_block_root" ssz-size:"32"`
	ValidatorIndex uint64    `json:"validator_index"`
	Signature      Signature `json:"signature" ssz-size:"96"`
}

type SignedAggregateAndProof struct {
	Message   *AggregateAndProof `json:"message"`
	Signature Signature          `json:"signature" ssz-size:"96"`
}

type Eth1Block struct {
	Timestamp    uint64 `json:"timestamp"`
	DepositRoot  Root   `json:"deposit_root" ssz-size:"32"`
	DepositCount uint64 `json:"deposit_count"`
}

type PowBlock struct {
	BlockHash       [32]byte `json:"block_hash" ssz-size:"32"`
	ParentHash      [32]byte `json:"parent_hash" ssz-size:"32"`
	TotalDifficulty [32]byte `json:"total_difficulty" ssz-size:"32"`
}

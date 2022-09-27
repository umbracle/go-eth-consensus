package http

type NodeEndpoint struct {
	c *Client
}

func (c *Client) Node() *NodeEndpoint {
	return &NodeEndpoint{c: c}
}

type Identity struct {
	PeerID             string            `json:"peer_id"`
	ENR                string            `json:"enr"`
	P2PAddresses       []string          `json:"p2p_addresses"`
	DiscoveryAddresses []string          `json:"discovery_addresses"`
	Metadata           *IdentityMetadata `json:"metadata"`
}

type IdentityMetadata struct {
	SeqNumber uint64 `json:"seq_number"`
	AttNets   string `json:"attnets"`
	SyncNets  string `json:"syncnets"`
}

// Identity returns the node network identity
func (n *NodeEndpoint) Identity() (*Identity, error) {
	var out *Identity
	err := n.c.Get("/eth/v1/node/identity", &out)
	return out, err
}

type Peer struct {
	PeerID              string `json:"peer_id"`
	Enr                 string `json:"enr"`
	LastSeendP2PAddress string `json:"last_seen_p2p_address"`
	State               string `json:"state"`
	Direction           string `json:"direction"`
}

func (n *NodeEndpoint) Peers() ([]*Peer, error) {
	var peers []*Peer
	err := n.c.Get("/eth/v1/node/peers", &peers)
	return peers, err
}

func (n *NodeEndpoint) GetPeer(peerID string) (*Peer, error) {
	var peer *Peer
	err := n.c.Get("/eth/v1/node/peers/"+peerID, &peer)
	return peer, err
}

type PeerCount struct {
	Disconnected  uint64 `json:"disconnected"`
	Connecting    uint64 `json:"connecting"`
	Connected     uint64 `json:"connected"`
	Disconnecting uint64 `json:"disconnecting"`
}

func (n *NodeEndpoint) PeerCount() (*PeerCount, error) {
	var peerCount *PeerCount
	err := n.c.Get("/eth/v1/node/peer_count", &peerCount)
	return peerCount, err
}

func (n *NodeEndpoint) Version() (string, error) {
	var out struct {
		Version string `json:"version"`
	}
	err := n.c.Get("/eth/v1/node/version", &out)
	return out.Version, err
}

type Syncing struct {
	HeadSlot     uint64 `json:"head_slot"`
	SyncDistance string `json:"sync_distance"`
	IsSyncing    bool   `json:"is_syncing"`
	IsOptimistic bool   `json:"is_optimistic"`
}

func (n *NodeEndpoint) Syncing() (*Syncing, error) {
	var out Syncing
	err := n.c.Get("/eth/v1/node/syncing", &out)
	return &out, err
}

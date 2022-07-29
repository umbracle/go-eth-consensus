package http

import consensus "github.com/umbracle/go-eth-consensus"

type ConfigEndpoint struct {
	c *Client
}

func (c *Client) Config() *ConfigEndpoint {
	return &ConfigEndpoint{c: c}
}

func (c *ConfigEndpoint) Spec() (*consensus.Spec, error) {
	var spec *consensus.Spec
	err := c.c.Get("/eth/v1/config/spec", &spec)
	return spec, err
}

type DepositContract struct {
	ChainID uint64 `json:"chain_id"`
	Address string `json:"address"`
}

func (c *ConfigEndpoint) DepositContract() (*DepositContract, error) {
	var depositContract *DepositContract
	err := c.c.Get("/eth/v1/config/deposit_contract", &depositContract)
	return depositContract, err
}

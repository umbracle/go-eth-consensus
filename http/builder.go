package http

type BuilderEndpoint struct {
	c *Client
}

func (c *Client) Builder() *BuilderEndpoint {
	return &BuilderEndpoint{c: c}
}

func (b *BuilderEndpoint) RegisterValidator() error {
	return nil
}

func (b *BuilderEndpoint) GetExecutionPayload() error {
	return nil
}

func (b *BuilderEndpoint) SubmitBlindedBlock() error {
	return nil
}

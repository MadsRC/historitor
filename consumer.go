package historitor

import (
	"bytes"
	"encoding/gob"
)

type Consumer struct {
	name string
}

// NewConsumer creates a new Consumer with the provided options.
func NewConsumer(options ...ConsumerOption) Consumer {
	opts := defaultConsumerOptions
	for _, opt := range globalConsumerOptions {
		opt.apply(&opts)
	}
	for _, opt := range options {
		opt.apply(&opts)
	}
	return Consumer{
		name: opts.Name,
	}
}

// GetName returns the name of the Consumer.
func (c *Consumer) GetName() string {
	return c.name
}

// externalConsumer is used to represent a [Consumer] that can easily be encoded and decoded using the gob package.
type externalConsumer struct {
	Name string
}

func (c Consumer) MarshalBinary() ([]byte, error) {
	ec := externalConsumer{
		Name: c.name,
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(ec)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Consumer) UnmarshalBinary(data []byte) error {
	var ec externalConsumer
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&ec)
	if err != nil {
		return err
	}
	c.name = ec.Name
	return nil
}

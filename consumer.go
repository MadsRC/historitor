package historitor

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

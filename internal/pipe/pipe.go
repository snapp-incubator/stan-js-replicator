package pipe

// Pipe interface pipes data from input to output.
// it has streaming to jetstream or nats to jetstream implementation.
// please note that we always use queue subscription so you need to specify
// the group id.
type Pipe interface {
	Pipe(topic, group string)
}

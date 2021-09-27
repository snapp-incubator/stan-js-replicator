package pipe

// Pipe interface pipes data from input to output.
// it has streaming to jetstream or nats to jetstream implementation.
type Pipe interface {
	Pipe(string)
}

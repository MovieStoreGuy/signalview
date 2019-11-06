package reporter

import "io"

// Formater is an interface that can be used to generate an admin output
// of the endpoint that was hit
type Formater interface {
	Writer(dev io.Writer) error
}

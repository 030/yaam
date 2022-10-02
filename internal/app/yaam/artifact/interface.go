package artifact

// Preserver is the interface that wraps the basic Preserve method.
//
// Preserve downloads an artifact from an external repository and writes it to
// disk.
type Preserver interface {
	Preserve(urlStrings ...string) error
}

// Publisher is the interface that wraps the basic Publish method.
//
// Publish writes an artifact to disk.
type Publisher interface {
	Publish() error
}

// Reader is the interface that wraps the basic Read method.
//
// Reader reads an artifact from disk.
type Reader interface {
	Read() error
}

// Unifier is the interface that wraps the basic Unify method.
//
// Unify groups multiple repositories.
type Unifier interface {
	Unify(name string) error
}

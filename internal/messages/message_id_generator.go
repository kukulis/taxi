package messages

type MessageIdGenerator interface {
	Generate() string
}

// TODO fixed generator
// TODO random generator

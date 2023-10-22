package state

type Options struct {
	MaxSessions uint
	MaxPlayer   uint
}

func DefaultOptions() Options {
	return Options{
		MaxSessions: 10,
		MaxPlayer:   10,
	}
}

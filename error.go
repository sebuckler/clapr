package clapr

type ErrHelp struct {
	Help string
}

func (e *ErrHelp) Error() string {
	return e.Help
}

type errTerm struct {
	index int
}

func (e *errTerm) Error() string {
	return "arguments terminated"
}

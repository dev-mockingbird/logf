package logf

type nothing struct{}

func Nothing() Logger {
	return &nothing{}
}

func (n *nothing) Prefix(string) Logger {
	return n
}

func (n *nothing) Logf(Level, string, ...any) {
	return
}

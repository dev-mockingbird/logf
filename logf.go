package logf

type Level uint8

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

const (
	CallerDepth = 2
)

var levelString map[Level]string

func init() {
	levelString = map[Level]string{
		Trace: "TRACE",
		Debug: "DEBUG",
		Info:  "INFO",
		Warn:  "WARN",
		Error: "ERROR",
		Fatal: "FATAL",
	}
}

type Prefixable interface {
	Prefix(string) Logger
}

type Logger interface {
	Prefixable
	Logfer
}

type Logfer interface {
	Logf(level Level, format string, v ...interface{})
}

type Logf func(level Level, format string, v ...interface{})

func (logf Logf) Logf(level Level, format string, v ...interface{}) {
	logf(level, format, v...)
}

type logger struct {
	config *Config
}

func New(opts ...Option) Logger {
	return &logger{config: getConfig(opts...)}
}

func (l logger) Logf(level Level, format string, v ...any) {
	if level < l.config.LogLevel {
		return
	}
	record := CollectRecord(level, l.config.CallerDepth, format, v...)
	l.config.Printer.Print(l.config.Prefix, record)
}

func (l *logger) Prefix(prefix string) Logger {
	return &logger{
		config: l.config.WithPrefix(prefix),
	}
}

func LevelString(level Level) string {
	if s, ok := levelString[level]; ok {
		return s
	}
	return "UNKOWN"
}

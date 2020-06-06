package log

type Format string

func (f Format) value() string {
	return string(f)
}

const (
	ConsoleFormat Format = "console"
	JsonFormat    Format = "json"
)

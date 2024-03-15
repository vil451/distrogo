package config

type Flags struct {
	LogFile *string
}

var (
	AppLogFile string
)

func NewFlags() *Flags {
	return &Flags{
		LogFile: strPtr(AppLogFile),
	}
}

func strPtr(s string) *string {
	return &s
}

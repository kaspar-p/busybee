package environment

type Mode int64

const (
	PRODUCTION Mode = iota
	DEVELOPMENT
	TESTING
)

func (m Mode) String() string {
	switch m {
	case PRODUCTION:
		return "production"
	case DEVELOPMENT:
		return "development"
	case TESTING:
		return "testing"
	}

	return ""
}

func (m Mode) ConfigString() string {
	switch m {
	case PRODUCTION:
		return "PROD"
	case DEVELOPMENT:
		return "DEV"
	case TESTING:
		return "TEST"
	}

	return ""
}

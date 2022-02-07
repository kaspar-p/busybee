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

func (m Mode) ConfigFile() string {
	switch m {
	case PRODUCTION:
		return "env.prod"
	case DEVELOPMENT:
		return "env.dev"
	case TESTING:
		return "env.test"
	}

	return ""
}

package rules

type Severity string

const (
	SeverityError Severity = "error"
)

type File struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	ID         string   `yaml:"id"`
	Name       string   `yaml:"name"`
	Severity   Severity `yaml:"severity"`
	Include    []string `yaml:"include"`
	Exclude    []string `yaml:"exclude"`
	Match      Match    `yaml:"match"`
	Message    string   `yaml:"message"`
	Expected   string   `yaml:"expected"`
	Suggestion string   `yaml:"suggestion"`
}

type Match struct {
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern"`
}

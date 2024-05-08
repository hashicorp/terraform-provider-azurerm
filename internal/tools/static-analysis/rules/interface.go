package rules

type Rule interface {
	Run() []error
	Name() string
	Description() string
}

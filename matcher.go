package quamina

type matcher interface {
	addPattern(x X, pat string) error
	matchesForFields(fields []Field) ([]X, error)
	deletePatterns(x X) error
	IsNameUsed(label []byte) bool
	Paths() map[string]bool
}

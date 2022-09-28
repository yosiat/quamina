package quamina

import (
	"os"
	"testing"
)

var (
	topMatches []X
	topFields  []Field
)

type tracker struct {
	names map[string]bool
}

func (t tracker) IsNameUsed(label []byte) bool {
	_, ok := t.names[string(label)]
	return ok
}

const PatternContext = `{ "context": { "user_id": [9034], "friends_count": [158] } }`

func Benchmark_JxFlattener_ContextFields(b *testing.B) {
	var localFields []Field

	event, err := os.ReadFile("./status.json")
	if err != nil {
		b.Fatal(err)
	}

	paths := newPathsIndex()

	paths.add("context\nuser_id")
	paths.add("context\nfriends_count")

	flattener := newJxFlattener(paths)
	t := tracker{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fields, err := flattener.Flatten(event, t)
		if err != nil {
			b.Fatal(err)
		}

		localFields = fields
	}

	topFields = localFields
}

func Benchmark_JsonFlattener_ContextFields(b *testing.B) {
	var localFields []Field

	event, err := os.ReadFile("./status.json")
	if err != nil {
		b.Fatal(err)
	}

	flattener := newJSONFlattener()

	t := tracker{names: make(map[string]bool)}
	t.names["user_id"] = true
	t.names["friends_count"] = true

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fields, err := flattener.Flatten(event, t)
		if err != nil {
			b.Fatal(err)
		}

		localFields = fields
	}

	topFields = localFields
}

func Benchmark_JsonFlattner_Evaluate_ContextFields(b *testing.B) {
	q, err := New()

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternContext)
}

func Benchmark_JxFlattner_Evaluate_ContextFields(b *testing.B) {
	paths := newPathsIndex()

	paths.add("context\nuser_id")
	paths.add("context\nfriends_count")

	flattener := newJxFlattener(paths)

	q, err := New(WithFlattener(flattener))

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternContext)
}

func RunBenchmarkEvaluate(b *testing.B, q *Quamina, pattern string) {
	var localMatches []X

	err := q.AddPattern(1, pattern)
	if err != nil {
		b.Fatalf("Failed adding pattern: %+v", err)
	}

	event, err := os.ReadFile("./status.json")
	if err != nil {
		b.Fatal(err)
	}

	matches, err := q.MatchesForEvent(event)
	if err != nil {
		b.Fatalf("failed matching: %s", err)
	}

	if len(matches) != 1 {
		b.Fatalf("in-correct matching: %+v", matches)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		matches, err := q.MatchesForEvent(event)
		if err != nil {
			b.Fatalf("failed matching: %s", err)
		}

		if len(matches) != 1 {
			b.Fatalf("in-correct matching: %+v", matches)
		}

		localMatches = matches
	}

	topMatches = localMatches
}

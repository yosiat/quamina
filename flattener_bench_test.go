package quamina

import (
	"os"
	"strings"
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
const PatternMiddleNestedField = `{ "payload": { "user": { "id_str": ["903487807"] } } }`
const PatternLastField = `{ "payload": { "lang_value": ["ja"] } }`

func Benchmark_JxFlattener_ContextFields(b *testing.B) {
	RunBenchmarkWithJxFlattener(b, "context\nuser_id", "context\nfriends_count")
}

func Benchmark_JsonFlattener_ContextFields(b *testing.B) {
	RunBehcmarkWithJSONFlattener(b, "context", "user_id", "friends_count")
}

func Benchmark_JxFlattener_MiddleNestedField(b *testing.B) {
	RunBenchmarkWithJxFlattener(b, "payload\nuser\nid_str")
}

func Benchmark_JsonFlattener_MiddleNestedField(b *testing.B) {
	RunBehcmarkWithJSONFlattener(b, "payload", "user", "id_str")
}

func Benchmark_JxFlattener_LastField(b *testing.B) {
	RunBenchmarkWithJxFlattener(b, "payload\nlang_value")
}

func Benchmark_JsonFlattener_LastField(b *testing.B) {
	RunBehcmarkWithJSONFlattener(b, "payload", "lang_value")
}

func RunBenchmarkWithJxFlattener(b *testing.B, fields ...string) {
	b.Helper()
	var localFields []Field

	event, err := os.ReadFile("./status.json")
	if err != nil {
		b.Fatal(err)
	}

	paths := newPathsIndex()

	for _, field := range fields {
		paths.add(field)
	}

	flattener := newJxFlattener(paths)
	t := tracker{}

	results, err := flattener.Flatten(event, t)
	if err != nil {
		b.Fatal(err)
	}
	PrintFields(b, results)

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

func RunBehcmarkWithJSONFlattener(b *testing.B, fields ...string) {
	b.Helper()
	var localFields []Field

	event, err := os.ReadFile("./status.json")
	if err != nil {
		b.Fatal(err)
	}

	flattener := newJSONFlattener()

	t := tracker{names: make(map[string]bool)}
	for _, field := range fields {
		t.names[field] = true
	}
	results, err := flattener.Flatten(event, t)
	if err != nil {
		b.Fatal(err)
	}
	PrintFields(b, results)

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

func Benchmark_JsonFlattner_Evaluate_MiddleNestedField(b *testing.B) {
	q, err := New()

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternMiddleNestedField)
}

func Benchmark_JxFlattner_Evaluate_MiddleNestedField(b *testing.B) {
	paths := newPathsIndex()

	paths.add("payload\nuser\nid_str")

	flattener := newJxFlattener(paths)

	q, err := New(WithFlattener(flattener))

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternMiddleNestedField)
}

func Benchmark_JsonFlattner_Evaluate_LastField(b *testing.B) {
	q, err := New()

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternLastField)
}

func Benchmark_JxFlattner_Evaluate_LastField(b *testing.B) {
	paths := newPathsIndex()

	paths.add("payload\nlang_value")

	flattener := newJxFlattener(paths)

	q, err := New(WithFlattener(flattener))

	if err != nil {
		b.Fatal(err)
	}

	RunBenchmarkEvaluate(b, q, PatternLastField)
}

func RunBenchmarkEvaluate(b *testing.B, q *Quamina, pattern string) {
	b.Helper()
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

func PrintFields(b *testing.B, fields []Field) {
	b.Helper()

	b.Logf("> Fields\n")

	for _, field := range fields {
		b.Logf("Path [%s] Val [%s] ArrayTrail [%+v]\n", strings.ReplaceAll(string(field.Path), "\n", "->"), field.Val, field.ArrayTrail)
	}
	b.Logf("\n")
}

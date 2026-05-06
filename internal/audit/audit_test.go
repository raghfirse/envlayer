package audit_test

import (
	"strings"
	"testing"

	"github.com/user/envlayer/internal/audit"
)

func TestRecord_AddedKey(t *testing.T) {
	l := &audit.Log{}
	before := map[string]string{}
	after := map[string]string{"APP_ENV": "production"}
	l.Record(before, after, ".env.production")

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	e := l.Entries[0]
	if e.Action != "added" || e.Key != "APP_ENV" || e.NewValue != "production" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecord_RemovedKey(t *testing.T) {
	l := &audit.Log{}
	before := map[string]string{"OLD_KEY": "value"}
	after := map[string]string{}
	l.Record(before, after, ".env")

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	e := l.Entries[0]
	if e.Action != "removed" || e.Key != "OLD_KEY" || e.OldValue != "value" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecord_ChangedKey(t *testing.T) {
	l := &audit.Log{}
	before := map[string]string{"DB_URL": "localhost"}
	after := map[string]string{"DB_URL": "prod-db.internal"}
	l.Record(before, after, ".env.production")

	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	e := l.Entries[0]
	if e.Action != "changed" || e.OldValue != "localhost" || e.NewValue != "prod-db.internal" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestRecord_NoChanges(t *testing.T) {
	l := &audit.Log{}
	env := map[string]string{"KEY": "val"}
	l.Record(env, env, ".env")
	if len(l.Entries) != 0 {
		t.Errorf("expected no entries for identical maps, got %d", len(l.Entries))
	}
}

func TestRecord_MultipleChanges(t *testing.T) {
	l := &audit.Log{}
	before := map[string]string{"A": "1", "B": "2"}
	after := map[string]string{"A": "99", "C": "3"}
	l.Record(before, after, ".env.local")

	if len(l.Entries) != 3 {
		t.Fatalf("expected 3 entries (changed A, removed B, added C), got %d", len(l.Entries))
	}
}

func TestPrint_Output(t *testing.T) {
	l := &audit.Log{}
	before := map[string]string{"X": "old"}
	after := map[string]string{"X": "new", "Y": "added"}
	l.Record(before, after, "testfile")

	var sb strings.Builder
	l.Print(&sb)
	out := sb.String()

	if !strings.Contains(out, "changed") && !strings.Contains(out, "~") {
		t.Errorf("expected changed marker in output, got: %s", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected added marker in output, got: %s", out)
	}
	if !strings.Contains(out, "testfile") {
		t.Errorf("expected source in output, got: %s", out)
	}
}

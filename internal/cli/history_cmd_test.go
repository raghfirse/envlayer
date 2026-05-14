package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envlayer/internal/cli"
)

func makeHistoryDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envhistory-cli-*")
	if err != nil {
		t.Fatalf("makeHistoryDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func writeHistoryEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeHistoryEnv: %v", err)
	}
	return p
}

func TestRunHistoryRecord_CreatesEntry(t *testing.T) {
	envDir := makeHistoryDir(t)
	histDir := makeHistoryDir(t)
	p := writeHistoryEnv(t, envDir, ".env", "APP=hello\nPORT=9000\n")

	var buf bytes.Buffer
	if err := cli.RunHistoryRecord(histDir, "v1", []string{p}, &buf); err != nil {
		t.Fatalf("RunHistoryRecord: %v", err)
	}
	if !strings.Contains(buf.String(), "recorded entry") {
		t.Errorf("unexpected output: %q", buf.String())
	}
}

func TestRunHistoryRecord_NoFiles_ReturnsError(t *testing.T) {
	histDir := makeHistoryDir(t)
	var buf bytes.Buffer
	if err := cli.RunHistoryRecord(histDir, "v1", nil, &buf); err == nil {
		t.Error("expected error for no files")
	}
}

func TestRunHistoryRecord_EmptyLabel_ReturnsError(t *testing.T) {
	envDir := makeHistoryDir(t)
	histDir := makeHistoryDir(t)
	p := writeHistoryEnv(t, envDir, ".env", "K=V\n")
	var buf bytes.Buffer
	if err := cli.RunHistoryRecord(histDir, "", []string{p}, &buf); err == nil {
		t.Error("expected error for empty label")
	}
}

func TestRunHistoryList_ShowsEntries(t *testing.T) {
	envDir := makeHistoryDir(t)
	histDir := makeHistoryDir(t)
	p := writeHistoryEnv(t, envDir, ".env", "X=1\n")

	_ = func() { cli.RunHistoryRecord(histDir, "snap1", []string{p}, &bytes.Buffer{}) }
	cli.RunHistoryRecord(histDir, "snap1", []string{p}, &bytes.Buffer{})

	var buf bytes.Buffer
	if err := cli.RunHistoryList(histDir, &buf); err != nil {
		t.Fatalf("RunHistoryList: %v", err)
	}
	if !strings.Contains(buf.String(), "snap1") {
		t.Errorf("expected snap1 in output, got: %q", buf.String())
	}
}

func TestRunHistoryList_EmptyDir_PrintsMessage(t *testing.T) {
	histDir := makeHistoryDir(t)
	var buf bytes.Buffer
	if err := cli.RunHistoryList(histDir, &buf); err != nil {
		t.Fatalf("RunHistoryList: %v", err)
	}
	if !strings.Contains(buf.String(), "no history") {
		t.Errorf("expected no-history message, got: %q", buf.String())
	}
}

func TestRunHistoryShow_PrintsVars(t *testing.T) {
	envDir := makeHistoryDir(t)
	histDir := makeHistoryDir(t)
	p := writeHistoryEnv(t, envDir, ".env", "COLOR=blue\n")

	var recBuf bytes.Buffer
	cli.RunHistoryRecord(histDir, "show-test", []string{p}, &recBuf)

	// extract the ID from the recorded output
	parts := strings.Fields(recBuf.String())
	if len(parts) < 3 {
		t.Fatalf("unexpected record output: %q", recBuf.String())
	}
	id := parts[2]

	var buf bytes.Buffer
	if err := cli.RunHistoryShow(histDir, id, &buf); err != nil {
		t.Fatalf("RunHistoryShow: %v", err)
	}
	if !strings.Contains(buf.String(), "COLOR") {
		t.Errorf("expected COLOR in output, got: %q", buf.String())
	}
}

func TestRunHistoryDelete_RemovesEntry(t *testing.T) {
	envDir := makeHistoryDir(t)
	histDir := makeHistoryDir(t)
	p := writeHistoryEnv(t, envDir, ".env", "D=1\n")

	var recBuf bytes.Buffer
	cli.RunHistoryRecord(histDir, "del-test", []string{p}, &recBuf)
	parts := strings.Fields(recBuf.String())
	id := parts[2]

	var buf bytes.Buffer
	if err := cli.RunHistoryDelete(histDir, id, &buf); err != nil {
		t.Fatalf("RunHistoryDelete: %v", err)
	}
	if !strings.Contains(buf.String(), "deleted") {
		t.Errorf("expected deleted in output, got: %q", buf.String())
	}
}

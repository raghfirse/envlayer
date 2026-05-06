// Package watcher provides lightweight polling-based file watching for
// .env files used by envlayer.
//
// It compares MD5 checksums of file contents between poll intervals to
// detect modifications reliably, even when filesystem timestamps have
// low resolution. Three change kinds are reported: created, modified,
// and removed.
//
// Usage:
//
//	done := make(chan struct{})
//	defer close(done)
//	events := watcher.Watch(paths, 500*time.Millisecond, done)
//	for ev := range events {
//		fmt.Println(ev.Kind, ev.Path)
//	}
package watcher

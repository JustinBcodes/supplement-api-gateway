package config

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch watches a file path for changes and invokes onChange with a freshly loaded config.
func Watch(path string, onChange func(*GatewayConfig)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		return err
	}

	debounce := time.NewTimer(0)
	if !debounce.Stop() {
		<-debounce.C
	}

	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				if !debounce.Stop() {
					select {
					case <-debounce.C:
					default:
					}
				}
				debounce.Reset(200 * time.Millisecond)
			}
		case <-debounce.C:
			cfg, err := Load(path)
			if err != nil {
				log.Printf("config reload error: %v", err)
				continue
			}
			onChange(cfg)
		case err := <-watcher.Errors:
			log.Printf("watcher error: %v", err)
		}
	}
}



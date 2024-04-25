package storage

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
)

// Watch watches storage directory on the specified path and calls onChange when storage changes.
func (s *Storage) Watch(storageChange chan<- string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create new watcher: %w", err)
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Println("Error closing watcher")
		}
	}(watcher)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					log.Println("Watcher event was not ok")
					return
				}
				msg := fmt.Sprintf("Event: %s\n", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					msg += fmt.Sprintf("Modified file: %s\n", event.Name)
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					msg += fmt.Sprintf("Created file: %s\n", event.Name)
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					msg += fmt.Sprintf("Deleted file: %s\n", event.Name)
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					msg += fmt.Sprintf("Renamed file: %s\n", event.Name)
				}

				storageChange <- msg
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Println("Watcher error event was not ok")
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(s.Path)
	if err != nil {
		log.Printf("Failed to add path %s to watcher", s.Path)
	}

	select {}
}

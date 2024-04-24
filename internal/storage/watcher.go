package storage

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

// Watch watches storage directory on the specified path and calls onChange when storage changes.
func (s *Storage) Watch(storageChange chan<- string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR CREATING NEW WATCHER")
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			fmt.Println("ERROR WATCHER CLOSE")
		}
	}(watcher)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					fmt.Println("EVENT NOT OK")
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

				//onChange()
				storageChange <- msg
			case err, ok := <-watcher.Errors:
				if !ok {
					fmt.Println("WATCHER ERRORS")
					return
				}
				fmt.Println("ERROR:", err)
			}
		}
	}()

	err = watcher.Add(s.Path)
	if err != nil {
		fmt.Println("ERROR WATCHER ADD")
	}

	select {}
}

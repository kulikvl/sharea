package storage

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type Storage struct {
	Path     string
	Capacity int64 // available storage space in bytes
}

type FileInfo struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	DownloadLink string    `json:"downloadLink"`
	ModTime      time.Time `json:"modTime"`
}

func (s *Storage) GetFilesInfo() ([]FileInfo, error) {
	files, err := os.ReadDir(s.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", s.Path, err)
	}

	var fileInfos []FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get info of file %s: %w", file.Name(), err)
		}

		fileInfos = append(fileInfos, FileInfo{
			Name:         file.Name(),
			Size:         fileInfo.Size(),
			DownloadLink: fmt.Sprintf("/api/download/%s", file.Name()),
			ModTime:      fileInfo.ModTime(),
		})
	}

	// sort file info entries by file modification time
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime.After(fileInfos[j].ModTime)
	})

	return fileInfos, nil
}

func (s *Storage) CalculateAvailableSpace() (int64, error) {
	var occupiedSpace int64 = 0

	files, err := os.ReadDir(s.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to read directory %s: %w", s.Path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			return 0, fmt.Errorf("failed to get info of file %s: %w", file.Name(), err)
		}

		occupiedSpace += fileInfo.Size()
	}

	return s.Capacity - occupiedSpace, nil
}

func (s *Storage) FileExists(filename string) (bool, error) {
	fileInfos, err := s.GetFilesInfo()
	if err != nil {
		return false, fmt.Errorf("failed to get files info: %w", err)
	}

	for _, f := range fileInfos {
		if f.Name == filename {
			return true, nil
		}
	}

	return false, nil
}

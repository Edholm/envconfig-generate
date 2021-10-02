package list

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path"
)

type listEntry struct {
	Dir     string   `json:"Dir"`
	GoFiles []string `json:"GoFiles"`
}

func AllGoFiles() ([]string, error) {
	golist := exec.Command("go", "list", "-json", "./...")

	stdoutPipe, err := golist.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to setup output pipe: %w", err)
	}
	defer stdoutPipe.Close()

	err = golist.Start()
	if err != nil {
		return nil, fmt.Errorf("failed start go list: %w", err)
	}

	allFiles := make([]string, 0, 50)
	decoder := json.NewDecoder(stdoutPipe)
	for {
		var entry listEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				return allFiles, nil
			}
			return nil, fmt.Errorf("failed to decode entry: %w", err)
		}

		if len(entry.GoFiles) == 0 {
			continue
		}

		files := make([]string, 0, len(entry.GoFiles))
		for _, file := range entry.GoFiles {
			files = append(files, path.Join(entry.Dir, file))
		}
		allFiles = append(allFiles, files...)
	}
}

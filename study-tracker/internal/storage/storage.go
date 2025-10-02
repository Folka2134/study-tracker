package storage

import (
	"encoding/csv"
	"os"
	"time"
)

const csvFile = "sessions.csv"

func SaveSession(task string, elapsed time.Duration) error {
	file, err := os.OpenFile(csvFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Check if the file is new to write the header
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		if err := writer.Write([]string{"date", "task", "elapsed"}); err != nil {
			return err
		}
	}

	row := []string{
		time.Now().Format(time.RFC3339),
		task,
		elapsed.Round(time.Second).String(),
	}

	return writer.Write(row)
}

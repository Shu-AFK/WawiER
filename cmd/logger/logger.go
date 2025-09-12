package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var Log *log.Logger

func InitLogger(mode string, filepath string) error {
	var writers []io.Writer

	switch mode {
	case "none":
		writers = append(writers, io.Discard)

	case "console":
		writers = append(writers, os.Stdout)

	case "file":
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, file)

	case "both":
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		writers = append(writers, os.Stdout, file)

	default:
		// fallback: console only
		writers = append(writers, os.Stdout)
	}

	multiWriter := io.MultiWriter(writers...)
	Log = log.New(multiWriter, "[WawiER] ", log.Ldate|log.Ltime)

	return nil
}

package logs

import (
	"fmt"
	"os"
)

/*
This function will check if the log file exists and open it in append mode.
If the log file doesn't exist, it will create a new one.
*/
func OpenLogFile() (*os.File) {
	logFile, err := os.OpenFile("./logs/errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v", err)
		return nil
	}
	return logFile
}

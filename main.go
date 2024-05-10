package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

type Config struct {
    Interval int `json:"interval"`
    Times []Time `json:"times"`
}

type Time struct {
    Time string `json:"time"`
    File string `json:"file"`
}

func main() {
    //-----------------------------------------------------------
    // @todo: make a function to create this path according to the OS
    config_path := fmt.Sprintf("/home/%s/.config/walter.json", os.Getenv("USER"))

    file, err := os.ReadFile(config_path)
    if err != nil {
        log.Println("Failed to read config file: ")
        log.Println("", err)
        os.Exit(0)
    }

    data := Config{}
    if err := json.Unmarshal([]byte(file), &data); err != nil {
        log.Println("Failed to parse config: ")
        log.Println("", err)
        os.Exit(0)
    }

    //-----------------------------------------------------------
    checkChan := make(chan struct{})
    stopChan := make(chan struct{})
    sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt) // Notify for interrupt signal
    //-----------------------------------------------------------
    currentItem := -1
    //-----------------------------------------------------------
    // Start a goroutine to run the code every 3 seconds
    go func() {

        for {
            select {
            case <-stopChan:
                fmt.Println("Exiting...")
                checkChan <- struct{}{} // Signal main goroutine to exit
                return
            default:
                checkTheTimes(data.Times, &currentItem)
                time.Sleep(time.Duration(data.Interval) * time.Second)
            }
        }
    }()

	// Wait for interrupt signal or exit signal
	select {
	case <-sigChan:
		fmt.Println("\nReceived interrupt signal")
		stopChan <- struct{}{}
	case <-checkChan:
		fmt.Println("Exited.")
	}
}

func checkTheTimes(configTimes []Time, currentItem *int) {
    //maxItems := len(configTimes)

    for index, entry := range configTimes {
        // Get current date and time in local timezone
        now := time.Now()

        // Parse the specified time "09:00"
        targetTime, err := time.Parse("15:04", entry.Time)
        if err != nil {
            // @todo: log the error, do we exit the process since the time is invalid or just keep the process running?
            log.Println("Error parsing time:", err)
            return
        }

        // Create a new time instance with the same date as now and the specified time
        entryDateTime := time.Date(
            now.Year(), 
            now.Month(), 
            now.Day(), 
            targetTime.Hour(), 
            targetTime.Minute(), 
            0, 
            0, 
            now.Location(),
        )

        timePassed := now.After(entryDateTime)

        if index > *currentItem {

            if (entryDateTime == now) || timePassed {
                *currentItem = index
                osFilePath := fmt.Sprintf("file://%s", entry.File)
                cmd1 := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri", osFilePath)
                cmd1.Run()
                cmd2 := exec.Command("gsettings", "set", "org.gnome.desktop.background", "picture-uri-dark", osFilePath)
                cmd2.Run()
            }
        }
    }

}

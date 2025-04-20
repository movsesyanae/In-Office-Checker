package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Metadata struct {
	office_wifi_name string
	path             string
}

func (metadata *Metadata) get_wifi_script_path() string {
	return filepath.Join(metadata.path, "scripts", "get_wifi.sh")
}

func (metadata *Metadata) get_wifi_name() string {
	cmd := exec.Command("sh", metadata.get_wifi_script_path())
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return strings.TrimSpace(string(stdout))
}

func (metadata *Metadata) load() {
	path, non_standard_path := os.LookupEnv("IN_OFFICE_PREFIX")
	if non_standard_path {
		metadata.path = path
	} else {
		metadata.path = "~/.my_tools/in_office/"
	}

	var file_path string = filepath.Join(metadata.path, "metadata.txt")
	body, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	// read file into lines list
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) == 0 {
		log.Fatalf("metadata file missing line of office wifi name")
	}
	// first line should be the wifi name
	metadata.office_wifi_name = lines[0]
}

func (metadata *Metadata) create_empty() {
	path, non_standard_path := os.LookupEnv("IN_OFFICE_PREFIX")
	if non_standard_path {
		metadata.path = path
	} else {
		metadata.path = "~/.my_tools/in_office/"
	}
	metadata.office_wifi_name = ""

	// save get wifi name command into a script
	metadata.create_get_wifi_script()

}

func (metadata *Metadata) create_get_wifi_script() {
	wifi_name_command :=
		`#!/bin/bash
ipconfig getsummary "$(networksetup -listallhardwareports | awk '/Wi-Fi|AirPort/{getline; print $NF}')" | grep '  SSID : ' | awk -F ': ' '{print $2}'
`
	script_dir := filepath.Join(metadata.path, "scripts")
	// Create directory if not present
	if _, err := os.Stat(script_dir); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(script_dir, 0755)
		} else {
			fmt.Printf("Error checking directory: %v\n", err)
		}
	}
	script_file_path := filepath.Join(script_dir, "get_wifi.sh")
	file, err := os.Create(script_file_path)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(wifi_name_command)
	file.Chmod(0770)
	file.Close()

}

func (metadata *Metadata) get_csv_location() string {
	return filepath.Join(metadata.path, "days_in_office.csv")
}
func (metadata *Metadata) save() {
	var file_path string = filepath.Join(metadata.path, "metadata.txt")
	f, err := os.Create(file_path)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	_, err = f.WriteString(metadata.office_wifi_name + "\n")
}

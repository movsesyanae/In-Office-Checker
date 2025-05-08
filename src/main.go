package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Uses will be
	// ./in_office <year> <month> <day>
	// 		saves the date as in office
	// ./in_office check
	// 		checks if in office today (via wifi name)
	// 		and updates data if necessary
	// ./in_office days since YYYY-MM-DD
	// 		prints how many days I have been in offic
	// 		since the specified date
	// ./in_office set office_wifi_name <office wifi name>
	// 		sets wifi name to compare (saves in metadata)
	// ./in_office init clean
	// 		initializes empy data to disk with empty df
	// // ./in_office init
	// 		initializes empy data to disk without
	// 		overwriting file for df
	args := os.Args
	switch args[1] {
	case "init":
		in_office_data := new(In_Office)
		if len(args) > 2 {
			// clean case
			if args[2] != "clean" {
				log.Fatal("Incorrect call of init\n Must follow ./in_office init clean?")
			}
			in_office_data.create_empty(true)
		} else {
			in_office_data.create_empty(false)
		}
		in_office_data.save()
	case "set":
		if len(args) < 4 {
			log.Fatal("Incorrect use of set\nMust be ./in_office set <var_name> <var_value>\nAvailable variables are: office_wifi_name")
		}
		metadata := new(Metadata)
		metadata.create_empty()
		switch args[2] {
		case "office_wifi_name":
			metadata.office_wifi_name = strings.Join(args[3:], " ")
		default:
			log.Fatal("Variable does not exist")
		}
		metadata.save()
	case "check":
		in_office_data := new(In_Office)
		in_office_data.load()
		in_office_data.check_today()
		in_office_data.save()
	case "days":
		if (len(args) != 6 && len(args) != 5) || args[2] != "since" {
			log.Fatal("Incorrect usage of days since\nShould be ./in_office days since <month> <day> <year?>")
		}
		month, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("Incorrect usage of days since\nShould be ./in_office days since <month> <day> <year?>")
		}
		day, err := strconv.Atoi(args[4])
		if err != nil {
			log.Fatal("Incorrect usage of days since\nShould be ./in_office days since <month> <day> <year?>")
		}
		var year int = time.Now().Year()
		if len(args) == 6 {
			year, err = strconv.Atoi(args[5])
			if err != nil {
				log.Fatal("Incorrect usage of days since\nShould be ./in_office days since <month> <day> <year?>")
			}
		}
		in_office_data := new(In_Office)
		in_office_data.load()
		in_office_data.days_since(year, month, day)
	case "help":
		fmt.Println("Available commands are:")
		fmt.Println("\tinit")
		fmt.Println("\tset")
		fmt.Println("\tcheck")
		fmt.Println("\tdays since")
		fmt.Println("\thelp")

	default:
		if len(args) != 4 {
			log.Fatal("Incorrect usage")
		}
		year, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("Incorrect usage\nShould be ./in_office <year> <month> <day>")
		}
		month, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatal("Incorrect usage\nShould be ./in_office <year> <month> <day>")
		}
		day, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatal("Incorrect usage\nShould be ./in_office <year> <month> <day>")
		}
		in_office_data := new(In_Office)
		in_office_data.load()
		in_office_data.log_day(year, month, day)
		in_office_data.save()
	}
}

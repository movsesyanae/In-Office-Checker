package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

type In_Office struct {
	metadata *Metadata
	df       dataframe.DataFrame
}

func create_empty_df() dataframe.DataFrame {
	years := series.New(
		[]int{},
		series.Int,
		"year",
	)
	months := series.New(
		[]int{},
		series.Int,
		"month",
	)
	days := series.New(
		[]int{},
		series.Int,
		"day",
	)
	df := dataframe.New(years, months, days)
	return df
}

func (in_office_data *In_Office) add_date(year int, month int, day int) {
	new_year := series.New(
		[]int{
			year,
		},
		series.Int,
		"year",
	)
	new_month := series.New(
		[]int{
			month,
		},
		series.Int,
		"month",
	)
	new_day := series.New(
		[]int{
			day,
		},
		series.Int,
		"day",
	)
	new_row := dataframe.New(new_year, new_month, new_day)

	in_office_data.df = in_office_data.df.Concat(new_row)
}

func read_from_csv(file string) dataframe.DataFrame {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	df := dataframe.ReadCSV(f)
	f.Close()
	return df
}

func (in_office_data *In_Office) is_in_office() bool {
	current_wifi_name := in_office_data.metadata.get_wifi_name()

	check := current_wifi_name == in_office_data.metadata.office_wifi_name
	if check {
		return true
	} else {
		return false
	}
}

func (in_office_data *In_Office) load() {
	in_office_data.metadata = new(Metadata)
	in_office_data.metadata.load()
	in_office_data.df = read_from_csv(in_office_data.metadata.get_csv_location())
	if in_office_data.df.Nrow() == 0 {
		in_office_data.df = create_empty_df()
	}
}

func (in_office_data *In_Office) log_day(year int, month int, day int) bool {

	filtered := in_office_data.df.FilterAggregation(
		dataframe.And,
		dataframe.F{0, "year", series.Eq, year},
		dataframe.F{0, "month", series.Eq, month},
		dataframe.F{0, "day", series.Eq, day},
	)
	if filtered.Nrow() == 0 {
		fmt.Printf("\tMarking %v/%v/%v as in office\n", month, day, year)
		in_office_data.add_date(year, month, day)
		return true
	} else {
		fmt.Printf("\t%v/%v/%v was already marked as in office\n", month, day, year)
	}
	return false
}

func (in_office_data *In_Office) check_today() bool {
	today := time.Now()
	fmt.Printf("Checking if in office at %v via automatic check\n", today)
	if in_office_data.is_in_office() {
		return in_office_data.log_day(today.Year(), int(today.Month()), today.Day())
	}
	return false
}

func (in_office_data *In_Office) save() {
	f, err := os.Create(in_office_data.metadata.get_csv_location())
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}

	in_office_data.df.WriteCSV(f)
	f.Close()
	in_office_data.metadata.save()
}

func (in_office_data *In_Office) create_empty(reset_df bool) {
	in_office_data.metadata = new(Metadata)
	in_office_data.metadata.create_empty()
	if reset_df {
		in_office_data.df = create_empty_df()
	} else {
		in_office_data.df = read_from_csv(in_office_data.metadata.get_csv_location())
	}
}

func (in_office_data *In_Office) days_since(year int, month int, day int) {
	filt_0 := in_office_data.df.Filter(
		dataframe.F{0, "year", series.Greater, year})
	filt_1 := in_office_data.df.FilterAggregation(
		dataframe.And,
		dataframe.F{0, "year", series.Eq, year},
		dataframe.F{0, "month", series.Greater, month},
	)
	filt_2 := in_office_data.df.FilterAggregation(
		dataframe.And,
		dataframe.F{0, "year", series.Eq, year},
		dataframe.F{0, "month", series.Eq, month},
		dataframe.F{0, "day", series.GreaterEq, day},
	)
	filtered := filt_0.Concat(filt_1).Concat(filt_2)
	fmt.Printf("There have been %v days in office since %v-%v-%v\n", filtered.Nrow(), year, month, day)
}

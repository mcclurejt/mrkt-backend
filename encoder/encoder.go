package encoder

import (
	"encoding/csv"
	"log"
	"os"
)

func CSV(filename string, headers []string, values [][]string) {
	file, err := os.Create("data/time_series/" + filename + ".csv")
	if err != nil {
		log.Fatal("Could not create file", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)
	if err := w.Write(headers); err != nil {
		//write failed do something
		log.Fatal("Could not encode headers")
	}
	if err := w.WriteAll(values); err != nil {
		//write failed do something
		log.Fatal("Could not encode values")
	}
}

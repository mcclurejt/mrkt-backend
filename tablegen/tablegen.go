package tablegen

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mcclurejt/mrkt-backend/api/iexcloud"
	"github.com/mcclurejt/mrkt-backend/config"
	"github.com/tealeg/xlsx/v3"
	"github.com/urfave/cli/v2"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Action(c *cli.Context) error {
	ctx := context.Background()
	conf := config.New() //env
	content, err := ioutil.ReadFile("symbols.config")
	check(err)
	symbols := strings.Split(strings.TrimSpace(string(content)), "\n")
	fmt.Println(symbols)
	client := iexcloud.NewIEXCloudClient(conf.Api.IEXCloudAPIKey)
	summaries := make([]*iexcloud.Company, len(symbols))
	var summary *iexcloud.Company
	for i, symbol := range symbols {
		if symbol != "" {
			summary, err = client.Company.Get(ctx, symbol)
			if err != nil {
				log.Printf("\n'%s' is an invalid symbol", symbol)
			}
			summaries[i] = summary
		}
	}
	currentPrices := make([]float64, len(symbols))
	oneWeekPercentChange := make([]float64, len(symbols))
	var ohlcvs []iexcloud.OHLCV
	for i, symbol := range symbols {
		ohlcvs, err = client.Chart.Get(ctx, symbol, iexcloud.ChartRange5d, &iexcloud.ChartOptions{ChartCloseOnly: true})
		if err != nil {
			log.Printf("\nFailing on %s", symbol)
			continue
		}
		currentPrice := ohlcvs[len(ohlcvs)-1].Close
		currentPrices[i] = currentPrice
		oneWeekPercentChange[i] = 100 * (currentPrice - ohlcvs[0].Close) / ohlcvs[0].Close
	}

	wb := xlsx.NewFile()
	sh, err := wb.AddSheet("Main")
	check(err)
	header := sh.AddRow()
	headers := []string{"Ticker", "Company Name", "Description", "Sector", "1WPercentChange", "Current Price"}
	for i := range headers {
		cell := header.AddCell()
		cell.SetString(headers[i])
	}

	for i, symbol := range symbols {
		row := sh.AddRow()

		symbolCell := row.AddCell()
		symbolStyle := xlsx.NewStyle()
		symbolStyle.Alignment.WrapText = true
		symbolCell.SetString(symbol)
		symbolCell.SetStyle(symbolStyle)

		cNameCell := row.AddCell()
		cNameStyle := xlsx.NewStyle()
		cNameStyle.Alignment.WrapText = true
		cNameCell.SetString(summaries[i].CompanyName)
		cNameCell.SetStyle(cNameStyle)

		descriptionCell := row.AddCell()
		descriptionStyle := xlsx.NewStyle()
		descriptionStyle.Alignment.WrapText = true
		descriptionCell.SetString(summaries[i].Description)
		descriptionCell.SetStyle(descriptionStyle)

		sectorCell := row.AddCell()
		sectorStyle := xlsx.NewStyle()
		sectorStyle.Alignment.WrapText = true
		sectorCell.SetString(summaries[i].Sector)
		sectorCell.SetStyle(sectorStyle)

		oneWeekPercentCell := row.AddCell()
		oneWeekPercentCell.SetFloat(oneWeekPercentChange[i])

		currentPriceCell := row.AddCell()
		currentPriceCell.SetFloat(currentPrices[i])
	}

	err = sh.ForEachRow(rowVisitor)
	check(err)
	f, err := os.Create("output.xlsx")
	check(err)
	err = sh.File.Write(f)
	check(err)
	return nil
}

func rowVisitor(r *xlsx.Row) error {
	fmt.Println("")
	r.ForEachCell(cellVisitor)
	return nil
}

func cellVisitor(c *xlsx.Cell) error {
	value, err := c.FormattedValue()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf(" %s ", value)
	}
	return err
}

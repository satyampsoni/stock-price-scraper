package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// Stock struct to hold stock data
type Stock struct {
    Company string
    Price   string
    Change  string
}

// fetchStockData fetches and parses stock data from Yahoo Finance
func fetchStockData(ticker string) (*Stock, error) {
    url := fmt.Sprintf("https://finance.yahoo.com/quote/%s/", ticker)
    response, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("error fetching data for %s: %v", ticker, err)
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        return nil, fmt.Errorf("error fetching data for %s: %d", ticker, response.StatusCode)
    }

    doc, err := goquery.NewDocumentFromReader(response.Body)
    if err != nil {
        return nil, fmt.Errorf("error parsing data for %s: %v", ticker, err)
    }

    company := doc.Find("section.container.yf-3a2v0c.paddingRight h1.yf-3a2v0c").Text()
    price := doc.Find("fin-streamer[data-field=regularMarketPrice]").AttrOr("data-value", "N/A")
    change := doc.Find("fin-streamer[data-field=regularMarketChangePercent]").AttrOr("data-value", "N/A")

    company = strings.TrimSpace(company)
    price = strings.TrimSpace(price)
    change = strings.TrimSpace(change)

    return &Stock{
        Company: company,
        Price:   price,
        Change:  change,
    }, nil
}

func main() {
    tickers := []string{
        "MSFT", "IBM", "GE", "UNP", "COST", "MCD", "V", "WMT",
        "DIS", "MMM", "INTC", "AXP", "AAPL", "BA", "CSCO", "GS",
        "JPM", "CRM", "VZ",
    }

    var stocks []*Stock

    for _, ticker := range tickers {
        fmt.Printf("Visiting: https://finance.yahoo.com/quote/%s/\n", ticker)
        stock, err := fetchStockData(ticker)
        if err != nil {
            log.Printf("Error: %v\n", err)
            continue
        }
        stocks = append(stocks, stock)
        fmt.Printf("Company: %s, Price: %s, Change: %s\n", stock.Company, stock.Price, stock.Change)
    }

    // Write to CSV
    file, err := os.Create("stocks.csv")
    if err != nil {
        log.Fatalf("Failed to create CSV file: %v", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write CSV headers
    writer.Write([]string{"Company", "Price", "Change"})
    // Write stock data
    for _, stock := range stocks {
        writer.Write([]string{stock.Company, stock.Price, stock.Change})
    }
}

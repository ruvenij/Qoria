package data

import (
	"Qoria/internal/model"
	"bufio"
	"encoding/csv"
	"github.com/shopspring/decimal"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func LoadCsvData() ([]*model.Transaction, error) {
	allTransactions := make([]*model.Transaction, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	splitFiles, err := filepath.Glob("././data/split/chunk_*.csv")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(splitFiles); i++ {
		filePath := splitFiles[i]
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			transactions, err := loadCsvFile(file)
			if err != nil {
				log.Println("Error occurred while loading transactions from file, file path : ",
					filePath, ", error : ", err)
				return
			}

			mu.Lock()
			allTransactions = append(allTransactions, transactions...)
			mu.Unlock()
		}(filePath)
	}

	wg.Wait()
	return allTransactions, nil
}

func loadCsvFile(filePath string) ([]*model.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []*model.Transaction{}, err
	}

	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	transactions := make([]*model.Transaction, 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error occurred while reading csv, error : ", err)
			continue
		}

		transaction := parseRecord(record)
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func parseRecord(record []string) *model.Transaction {
	price, _ := decimal.NewFromString(record[8])
	quantity, _ := strconv.ParseInt(record[9], 10, 64)
	total, _ := decimal.NewFromString(record[10])
	stock, _ := strconv.ParseInt(record[11], 10, 64)
	txDate, _ := time.Parse("2006-01-02", record[1])
	addedDate, _ := time.Parse("2006-01-02", record[12])

	return &model.Transaction{
		TransactionId:   record[0],
		TransactionDate: txDate,
		UserId:          record[2],
		Country:         record[3],
		Region:          record[4],
		ProductId:       record[5],
		ProductName:     record[6],
		Category:        record[7],
		Price:           price,
		Quantity:        quantity,
		TotalPrice:      total,
		StockQuantity:   stock,
		AddedDate:       addedDate,
	}
}

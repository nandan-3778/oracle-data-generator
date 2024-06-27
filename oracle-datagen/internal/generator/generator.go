/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

var mu sync.Mutex

type Generator interface {
	Table() string
	CSVHeaders() string
	CSVColumnMapping() string
	FakeRecord() (string, int)
}

func Generate(runID string, n, batchSize int, generator Generator) (string, error, float64) {
	table := generator.Table()
	headers := generator.CSVHeaders()
	columns := generator.CSVColumnMapping()

	dataGenStart := time.Now()
	wg := errgroup.Group{}

	var processed atomic.Int64
	var totalSizeMBperRun atomic.Value
	totalSizeMBperRun.Store(0.0)

	for i := range n {
		wg.Go(func() error {
			totalSizeMBforOneBatch, err := generateCSVDataFile(runID, table, headers, i, batchSize, generator.FakeRecord)
			if err != nil {
				return err
			} else {
				processed.Add(int64(batchSize))
				mu.Lock()
				currentTotal := totalSizeMBperRun.Load().(float64)
				totalSizeMBperRun.Store(currentTotal + totalSizeMBforOneBatch)
				mu.Unlock()
			}
			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return "", err, 0
	}

	controlFileName := filepath.Join(".", table, fmt.Sprintf("control-%s.ctl", runID))

	err = createControlFile(runID, table, columns, controlFileName, n)
	if err != nil {
		return "", err, 0
	}

	dataGenStartElapsed := time.Since(dataGenStart)
	log.Printf("generated %d csv files for table: %s runID: %s controlFileName: %s time elapsed: %s", processed.Load(), table, runID, controlFileName, dataGenStartElapsed)
	return controlFileName, nil, totalSizeMBperRun.Load().(float64)
}

func createControlFile(runID, table, columns, controlFileName string, n int) error {
	err := os.MkdirAll(table, 0755)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create directory %s: %s", table, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	controlFile, err := os.Create(controlFileName)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create and open file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	_, err = fmt.Fprintln(controlFile, "LOAD DATA")
	if err != nil {
		errMsg := fmt.Sprintf("failed to write LOAD DATA to file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	for i := range n {
		csvFileName := filepath.Join(".", table, runID, fmt.Sprintf("batchNumber-%d.csv", i))
		_, err = fmt.Fprintf(controlFile, "\nINFILE '%s'", csvFileName)
		if err != nil {
			errMsg := fmt.Sprintf("failed to write INFILE to file %s: %s", controlFileName, err.Error())
			log.Print(errMsg)
			return fmt.Errorf(errMsg)
		}
	}

	_, err = fmt.Fprintln(controlFile, "\nAPPEND")
	if err != nil {
		errMsg := fmt.Sprintf("failed to write APPEND to file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	_, err = fmt.Fprintf(controlFile, "INTO TABLE %s\n", table)
	if err != nil {
		errMsg := fmt.Sprintf("failed to write INTO TABLE to file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	_, err = fmt.Fprintln(controlFile, "FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '\"'")
	if err != nil {
		errMsg := fmt.Sprintf("failed to write FIELDS TERMINATED BY ',' to file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}

	_, err = fmt.Fprintln(controlFile, columns)
	if err != nil {
		errMsg := fmt.Sprintf("failed to write fields to file %s: %s", controlFileName, err.Error())
		log.Print(errMsg)
		return fmt.Errorf(errMsg)
	}
	return err
}

func generateCSVDataFile(runID, table, headers string, batchNumber, batchSize int, lineFn func() (string, int)) (float64, error) {
	producerWg := errgroup.Group{}

	err := os.MkdirAll(filepath.Join(table, runID), 0755)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create directory %s: %s", table, err.Error())
		log.Print(errMsg)
		return 0, fmt.Errorf(errMsg)
	}

	csvFileName := filepath.Join(".", table, runID, fmt.Sprintf("batchNumber-%d.csv", batchNumber))

	csvFile, err := os.Create(csvFileName)
	if err != nil {
		errMsg := fmt.Sprintf("failed to create and open file %s: %s", csvFileName, err.Error())
		log.Print(errMsg)
		return 0, fmt.Errorf(errMsg)
	}

	_, err = fmt.Fprintln(csvFile, headers)
	if err != nil {
		errMsg := fmt.Sprintf("failed to write record to file %s: %s", csvFileName, err.Error())
		log.Print(errMsg)
		return 0, fmt.Errorf(errMsg)
	}
	var mu sync.Mutex
	var totalSize atomic.Int64
	for range batchSize {
		producerWg.Go(func() error {
			line, lineSize := lineFn()
			//lineSize is the size of the 1 record in a batch , to get the total size of data per batch , we need to add size of all records in the batch
			mu.Lock() // Lock the mutex before modifying totalSize
			totalSize.Add(int64(lineSize))
			mu.Unlock() // Unlock the mutex after modifying
			_, err = fmt.Fprintln(csvFile, line)
			if err != nil {
				errMsg := fmt.Sprintf("failed to write record to file %s: %s", csvFileName, err.Error())
				log.Print(errMsg)
				return fmt.Errorf(errMsg)
			}
			return nil
		})
	}

	producerWg.Wait()

	csvFile.Sync()
	totalSizeMBforBatch := float64(totalSize.Load()) / (1024 * 1024) // Convert total size from bytes to MB
	// log.Printf("Total data size for batch %d: %.2f MB", batchNumber, totalSizeMBforBatch) // Log the total size of data in MB

	return totalSizeMBforBatch, nil
}

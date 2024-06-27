package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator"
	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/generator/tables"
	"git.betsol.com/zmanda/zmandapro/automations/data-generators/oracle-datagen/internal/loaders/sqlldr"
	"github.com/brianvoe/gofakeit"
	"github.com/godoes/go-figure"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type RootOptions struct {
	Cycles    int
	Batches   int
	BatchSize int
	Table     string
	KeepData  bool
}

var ropts RootOptions
var sqlldr_max_retries = 5

var rootCmd = &cobra.Command{
	Use:   "tdkodo [flags]",
	Short: "Populate an Oracle Database with test data to test backup and recovery scripts.",
	Run: func(cmd *cobra.Command, args []string) {
		Run(ropts.Batches, ropts.BatchSize, ropts.Cycles, ropts.Table)
	},
}

func init() {
	rootCmd.Flags().IntVarP(&ropts.Cycles, "cycles", "c", 20, "Number of cycles to run")
	rootCmd.Flags().IntVarP(&ropts.Batches, "batches", "b", 1000, "Number of batches to run")
	rootCmd.Flags().IntVarP(&ropts.BatchSize, "batchSize", "s", 1000, "Number of records per batch")
	rootCmd.Flags().StringVarP(&ropts.Table, "table", "t", "calls", "Table to populate")
	rootCmd.Flags().BoolVarP(&ropts.KeepData, "keepData", "k", false, "Keep data after run")
}

func setupLogger() {
	logFile, err := os.OpenFile("data_insertion.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func Execute() {
	setupLogger()
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %s", err)
	}
}

type controlFile struct {
	table           string
	runID           string
	controlFileName string
}

func Run(n, batchSize, cycles int, table string) {
	wg := errgroup.Group{}

	controlFiles := make(chan controlFile, cycles)
	loadedControlFiles := make(chan controlFile, cycles)

	start := time.Now()

	cleanupWg := errgroup.Group{}
	if !ropts.KeepData {
		cleanupWg.Go(func() error {
			// CLEANER: will delete data which is loaded into the database
			for controlFile := range loadedControlFiles {
				deleteFile := filepath.Join(".", controlFile.table, controlFile.runID)
				err := os.RemoveAll(deleteFile)
				if err != nil {
					log.Printf("failed to delete directory for runID: %s : %s", controlFile.runID, err.Error())

				}

			}
			return nil
		})
	}

	wg.Go(func() error {
		// LOADER: will load any given control file, not dependent on the table name
		// totalDataLoaded := 0.0
		count := 0
		for controlFile := range controlFiles {
			err := sqlldr.Load(controlFile.runID, controlFile.controlFileName)
			if err != nil {
				log.Printf("failed to run sqlldr for runID: %s controlFileName: %s : %s", controlFile.runID, controlFile.controlFileName, err.Error())
				count += 1
				if count == sqlldr_max_retries {
					log.Fatal("Terminating program due to sqlldr failure.")
				}

			}
			loadedControlFiles <- controlFile
		}

		return nil
	})

	g := getGenerator(table)
	var totalSizeofDataInMB float64
	myFigure := figure.NewFigure("ORACLE DATA LOADER", "standard", true)

	log.Println(myFigure)
	log.Println("Starting oracle Data population............")
	currentTime := time.Now()
	log.Println("************************************************************")
	log.Printf("Date and time: %s", currentTime.Format("2006-01-02 15:04:05"))
	log.Println("************************************************************")
	for i := range cycles {
		runID := gofakeit.UUID()
		log.Printf("Data Loading: Cycle %d\n", i+1)
		controlFileName, err, totalSizeMBperRun := generator.Generate(runID, n, batchSize, g)
		if err != nil {
			log.Printf("failed to generate data for runID: %s : %s", runID, err.Error())
			return
		}

		controlFiles <- controlFile{table: g.Table(), runID: runID, controlFileName: controlFileName}
		totalSizeofDataInMB += totalSizeMBperRun
		log.Println("************************************************************")
		log.Printf("Amount of data loaded till now: %.2f MB", totalSizeofDataInMB)
		log.Println("************************************************************")
	}
	log.Println("Data loading complete")
	log.Printf("Total Amount of data loaded: %.2f MB", totalSizeofDataInMB)

	close(controlFiles)

	err := wg.Wait()
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	close(loadedControlFiles)

	err = cleanupWg.Wait()
	if err != nil {
		log.Printf("error: %s", err.Error())
	}

	elapsed := time.Since(start)
	log.Default().Printf("total time taken: %s", elapsed)
}

func getLoadedDataSize(logFile string) (float64, error) {

	file, err := os.Open(logFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`Read\s+buffer\s+bytes:\s+(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if len(match) == 2 {
			bytesRead, _ := strconv.ParseInt(match[1], 10, 64)
			mbLoaded := float64(bytesRead) / (1024 * 1024)
			return mbLoaded, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("could not find 'Bytes Read' entry in the log file")
}

func getGenerator(table string) generator.Generator {
	switch table {
	case "calls":
		return tables.Call{}
	case "users":
		return tables.User{}
	case "shipments":
		return tables.Shipments{}
	case "bookings":
		return tables.Bookings{}
	case "inventory":
		return tables.Inventory{}
	case "productDescription":
		return tables.ProductDescription{}
	case "payments":
		return tables.Payment{}
	case "orders":
		return tables.Orders{}
	case "products":
		return tables.Products{}
	case "customers":
		return tables.Customer{}
	case "suppliers":
		return tables.Suppliers{}
	case "admin":
		return tables.Admin{}
	default:
		return tables.Call{}
	}
}

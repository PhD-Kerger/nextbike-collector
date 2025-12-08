package compactor

import (
	"encoding/json"
	"fmt"
	"nextbike-scraper/internal/logging"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/parquet-go/parquet-go/compress/zstd"

	"nextbike-scraper/internal/types"
)

type JsonFile struct {
	Path      string
	Timestamp uint32
}

func collectJsonFiles(dir string) []JsonFile {
	re := regexp.MustCompile(`(\d+)\.json$`)
	var files []JsonFile
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && re.MatchString(path) {
			matches := re.FindStringSubmatch(path)
			if len(matches) > 1 {
				var ts uint32
				fmt.Sscanf(matches[1], "%d", &ts)
				files = append(files, JsonFile{Path: path, Timestamp: ts})
			}
		}
		return nil
	})
	return files
}

func convertFolder(inputDir, outputDir string) error {
	files := collectJsonFiles(inputDir)
	logging.Logger.Info("[Compactor] Found files to convert", "count", len(files), "input_dir", inputDir)

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("[Compactor] failed to create output directory: %w", err)
		}
	}

	// Sort files by timestamp to ensure sorted output
	sort.Slice(files, func(i, j int) bool {
		return files[i].Timestamp < files[j].Timestamp
	})

	// Create parquet files with on-disk page buffers
	countryFile, err := os.Create(filepath.Join(outputDir, "Countries.parquet"))
	if err != nil {
		return fmt.Errorf("[Compactor] failed to create Countries.parquet file: %w", err)
	}
	defer countryFile.Close()

	cityFile, err := os.Create(filepath.Join(outputDir, "Cities.parquet"))
	if err != nil {
		return fmt.Errorf("[Compactor] failed to create Cities.parquet file: %w", err)
	}
	defer cityFile.Close()

	placeFile, err := os.Create(filepath.Join(outputDir, "Places.parquet"))
	if err != nil {
		return fmt.Errorf("[Compactor] failed to create Places.parquet file: %w", err)
	}
	defer placeFile.Close()

	bikeFile, err := os.Create(filepath.Join(outputDir, "Bikes.parquet"))
	if err != nil {
		return fmt.Errorf("[Compactor] failed to create Bikes.parquet file: %w", err)
	}
	defer bikeFile.Close()

	// Create parquet writers with on-disk page buffers to reduce memory usage
	countryWriter := parquet.NewGenericWriter[types.FlatCountry](countryFile,
		parquet.Compression(&zstd.Codec{}),
	)
	defer func() {
		if err := countryWriter.Close(); err != nil {
			logging.Logger.Error("[Compactor] Failed to close country writer", "err", err)
		}
	}()

	cityWriter := parquet.NewGenericWriter[types.FlatCity](cityFile,
		parquet.Compression(&zstd.Codec{}),
	)
	defer func() {
		if err := cityWriter.Close(); err != nil {
			logging.Logger.Error("[Compactor] Failed to close city writer", "err", err)
		}
	}()

	placeWriter := parquet.NewGenericWriter[types.FlatPlace](placeFile,
		parquet.Compression(&zstd.Codec{}),
	)
	defer func() {
		if err := placeWriter.Close(); err != nil {
			logging.Logger.Error("[Compactor] Failed to close place writer", "err", err)
		}
	}()

	bikeWriter := parquet.NewGenericWriter[types.FlatBike](bikeFile,
		parquet.Compression(&zstd.Codec{}),
	)
	defer func() {
		if err := bikeWriter.Close(); err != nil {
			logging.Logger.Error("[Compactor] Failed to close bike writer", "err", err)
		}
	}()

	// Process files in small batches to reduce memory usage
	const batchSize = 40 // Process 120 files at a time
	for i := 0; i < len(files); i += batchSize {
		end := i + batchSize
		if end > len(files) {
			end = len(files)
		}

		batch := files[i:end]
		if err := processBatch(batch, countryWriter, cityWriter, placeWriter, bikeWriter); err != nil {
			logging.Logger.Error("[Compactor] Failed to process batch", "batch_start", i, "batch_end", end, "err", err)
			return err
		}

		logging.Logger.Info("[Compactor] Processed batch", "batch_start", i, "batch_end", end, "total", len(files))
	}

	// delete the input directory after processing
	if err := os.RemoveAll(inputDir); err != nil {
		logging.Logger.Error("[Compactor] Failed to delete input directory after processing", "input_dir", inputDir, "err", err)
		return fmt.Errorf("[Compactor] failed to delete input directory: %w", err)
	}

	logging.Logger.Info("[Compactor] Compaction completed", "input_dir", inputDir, "output_dir", outputDir)
	return nil
}

// processBatch processes a batch of JSON files and writes the data to parquet writers
func processBatch(files []JsonFile, countryWriter *parquet.GenericWriter[types.FlatCountry], cityWriter *parquet.GenericWriter[types.FlatCity], placeWriter *parquet.GenericWriter[types.FlatPlace], bikeWriter *parquet.GenericWriter[types.FlatBike]) error {
	// Create small buffers for batch processing (much smaller than before)
	var countries []types.FlatCountry
	var cities []types.FlatCity
	var places []types.FlatPlace
	var bikes []types.FlatBike

	for _, file := range files {
		logging.Logger.Info("[Compactor] Processing file", "path", file.Path, "timestamp", file.Timestamp)

		f, err := os.Open(file.Path)
		if err != nil {
			logging.Logger.Error("[Compactor] Failed to open JSON file during compaction", "path", file.Path, "err", err)
			continue
		}

		// Read the JSON file and unmarshal it into the Root struct
		var root types.Root
		if err := json.NewDecoder(f).Decode(&root); err != nil {
			logging.Logger.Error("[Compactor] Failed to decode JSON file during compaction", "path", file.Path, "err", err)
			f.Close()
			continue
		}
		f.Close()

		// Process and collect data from this file
		for _, country := range root.Countries {
			flatCountry := country.ToFlat(file.Timestamp)
			countries = append(countries, flatCountry)

			for _, city := range country.Cities {
				flatCity := city.ToFlat(file.Timestamp, country.Name)
				cities = append(cities, flatCity)

				for _, place := range city.Places {
					flatPlace := place.ToFlat(file.Timestamp, city.UID)
					places = append(places, flatPlace)

					for _, bike := range place.BikeList {
						flatBike := bike.ToFlat(file.Timestamp, place.UID)
						bikes = append(bikes, flatBike)
					}
				}
			}
		}
	}

	// Sort the data to maintain order
	sort.Slice(countries, func(i, j int) bool {
		if countries[i].Name != countries[j].Name {
			return countries[i].Name < countries[j].Name
		}
		return countries[i].Timestamp < countries[j].Timestamp
	})

	sort.Slice(cities, func(i, j int) bool {
		if cities[i].UID != cities[j].UID {
			return cities[i].UID < cities[j].UID
		}
		return cities[i].Timestamp < cities[j].Timestamp
	})

	sort.Slice(places, func(i, j int) bool {
		if places[i].UID != places[j].UID {
			return places[i].UID < places[j].UID
		}
		return places[i].Timestamp < places[j].Timestamp
	})

	sort.Slice(bikes, func(i, j int) bool {
		if bikes[i].Number != bikes[j].Number {
			return bikes[i].Number < bikes[j].Number
		}
		return bikes[i].Timestamp < bikes[j].Timestamp
	})

	// Write the sorted data to parquet files
	if len(countries) > 0 {
		if _, err := countryWriter.Write(countries); err != nil {
			return fmt.Errorf("[Compactor] failed to write countries: %w", err)
		}
	}

	if len(cities) > 0 {
		if _, err := cityWriter.Write(cities); err != nil {
			return fmt.Errorf("[Compactor] failed to write cities: %w", err)
		}
	}

	if len(places) > 0 {
		if _, err := placeWriter.Write(places); err != nil {
			return fmt.Errorf("[Compactor] failed to write places: %w", err)
		}
	}

	if len(bikes) > 0 {
		if _, err := bikeWriter.Write(bikes); err != nil {
			return fmt.Errorf("[Compactor] failed to write bikes: %w", err)
		}
	}

	logging.Logger.Info("[Compactor] Batch written to parquet files", "countries", len(countries), "cities", len(cities), "places", len(places), "bikes", len(bikes))
	return nil
}

// target refers to the scraping target.
func Compact(inputTargetDirPath, outputTargetDirPath string) {
	logging.Logger.Info("[Compactor] Starting data compaction process")
	logging.Logger.Info("[Compactor] Using output directory", "path", outputTargetDirPath)
	logging.Logger.Info("[Compactor] Reading directories from input directory", "input_dir", inputTargetDirPath)

	var input_output_folders []struct {
		input  string
		output string
	}

	// Read the entries in the input target directory
	day_input_folders, err := os.ReadDir(inputTargetDirPath)
	if err != nil {
		logging.Logger.Error("[Compactor] Failed to read input target directory", "inputTargetDirPath", inputTargetDirPath, "err", err)
		return
	}

	for _, day_input_folder := range day_input_folders {
		if day_input_folder.IsDir() {
			// check if the date of the folder (yyy-mm-dd) is before today
			if day_input_folder.Name() >= time.Now().Format("2006-01-02") {
				logging.Logger.Info("[Compactor] Skipping folder", "folder", day_input_folder.Name(), "reason", "as it is not before today")
				continue
			}
			if day_input_folder.IsDir() {
				input_output_folders = append(input_output_folders, struct {
					input  string
					output string
				}{
					input:  filepath.Join(inputTargetDirPath, day_input_folder.Name()),
					output: filepath.Join(outputTargetDirPath, day_input_folder.Name()),
				})
			}

		}
	}

	logging.Logger.Info("[Compactor] Found input/output folders", "count", len(input_output_folders))

	sem := make(chan struct{}, 1) // threadpool with 2 workers
	errCh := make(chan error, len(input_output_folders))

	for _, ioPair := range input_output_folders {
		sem <- struct{}{}
		go func(ioPair struct{ input, output string }) {
			defer func() { <-sem }()
			if err := convertFolder(ioPair.input, ioPair.output); err != nil {
				errCh <- fmt.Errorf("[Compactor] failed to convert folder %s: %w", ioPair.input, err)
			} else {
				errCh <- nil
			}
		}(ioPair)
	}

	// Wait for all goroutines to finish
	for i := 0; i < len(input_output_folders); i++ {
		if err := <-errCh; err != nil {
			logging.Logger.Error("[Compactor] Error during folder conversion", "err", err)
		}
	}
}

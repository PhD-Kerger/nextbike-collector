package main

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"

	"nextbike-scraper/internal/compactor"
	"nextbike-scraper/internal/config"
	"nextbike-scraper/internal/logging"
	"nextbike-scraper/internal/scraping"

	"github.com/fsnotify/fsnotify"
	"github.com/hirochachacha/go-smb2"
	"github.com/robfig/cron/v3"
)

func main() {
	var configFilePath string
	if len(os.Args) < 2 {
		configFilePath = "./values.yaml"
	} else {
		configFilePath = os.Args[1]
	}

	if err := logging.SetupLogger("./log"); err != nil {
		os.Exit(1)
	}
	logger := logging.Logger

	currentConfig, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Error("[Scraper] Error loading config", "err", err)
		os.Exit(1)
	}
	logger.Info("[Scraper] Loaded config", "targets", len(currentConfig.Targets))

	cronJobs := SetupScrapingCronJobs(currentConfig)
	logging.Logger.Info("[Scraper] Scraping cron jobs set up", "jobs_count", len(currentConfig.Targets))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Error("[Scraper] Error creating file watcher", "err", err)
		os.Exit(1)
	}
	defer watcher.Close()

	err = watcher.Add(configFilePath)
	if err != nil {
		logger.Error("[Scraper] Error adding file to watcher", "file", configFilePath, "err", err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logger.Info("[Scraper] Config file changed, reloading", "file", event.Name)
					newConfig, err := config.LoadConfig(configFilePath)
					if err != nil {
						logger.Error("[Scraper] Error reloading config", "err", err)
						continue
					}
					currentConfig = newConfig
					logger.Info("[Scraper] Reloaded config", "targets", len(currentConfig.Targets))
					cronJobs.Stop()
					cronJobs = SetupScrapingCronJobs(currentConfig)
					logger.Info("[Scraper] Scraping cron jobs updated", "jobs_count", len(currentConfig.Targets))
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("[Scraper] Watcher error", "err", err)
			}
		}
	}()
	// Block forever to keep the main goroutine running
	select {}
}

func moveParquetFilesToSambaShare(config *config.ScraperConfig) {

	logging.Logger.Info("[Samba Move] Starting to move Parquet files to Samba share", "smb_address", config.SMBAddress, "smb_username", config.SMBUsername)
	conn, err := net.Dial("tcp", config.SMBAddress)
	if err != nil {
		logging.Logger.Error("[Samba Move] Error connecting to Samba server", "smb_address", config.SMBAddress, "err", err)
		return
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     config.SMBUsername,
			Password: config.SMBPassword,
			Domain:   config.SMBWorkgroup,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		logging.Logger.Error("[Samba Move] Error dialing Samba server", "smb_address", config.SMBAddress, "err", err)
		return
	}
	defer s.Logoff()
	logging.Logger.Info("[Samba Move] Connected to Samba share", "smb_address", config.SMBAddress, "smb_username", config.SMBUsername)

	fs, err := s.Mount(config.SMBPathToMount)
	if err != nil {
		logging.Logger.Error("[Samba Move] Error mounting Samba share", "smb_path_to_mount", config.SMBPathToMount, "err", err)
		return
	}
	logging.Logger.Info("[Samba Move] Mounted Samba share", "smb_path_to_mount", config.SMBPathToMount)
	// move all files from output-parquet dir and sub dirs to the samba share (./docker-nextbike-scraper)
	outputPath := config.OutputRootPathParquet
	err = filepath.Walk(outputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Logger.Error("[Samba Move] Error walking through output path", "path", path, "err", err)
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Open the source file
		srcFile, err := os.Open(path)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error opening source file", "src", path, "err", err)
			return err
		}
		defer srcFile.Close()

		// Compute the relative path from outputPath to the current file
		relPath, err := filepath.Rel(outputPath, path)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error computing relative path", "base", outputPath, "target", path, "err", err)
			return err
		}
		destPath := filepath.Join(config.SMBPathInsideMount, relPath)

		// Ensure the destination directory exists on the SMB share
		destDir := filepath.Dir(destPath)
		err = fs.MkdirAll(destDir, 0755)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error creating directory on samba share", "dir", destDir, "err", err)
			return err
		}
		dstFile, err := fs.Create(destPath)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error creating file on samba share", "dest", destPath, "err", err)
			return err
		}
		defer dstFile.Close()

		// Copy contents
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error copying file to samba share", "src", path, "dest", destPath, "err", err)
			return err
		}

		// Optionally remove the source file after copying
		err = os.Remove(path)
		if err != nil {
			logging.Logger.Error("[Samba Move] Error removing source file after moving", "src", path, "err", err)
			return err
		}

		logging.Logger.Info("[Samba Move] Moved file to samba share", "src", path, "dest", destPath)
		return nil
	})

	if err != nil {
		logging.Logger.Error("[Samba Move] Error walking through output path", "outputPath", outputPath, "err", err)
		return
	}
	logging.Logger.Info("[Samba Move] Finished moving Parquet files to Samba share", "smb_address", config.SMBAddress, "smb_username", config.SMBUsername)

}

func SetupScrapingCronJobs(config *config.ScraperConfig) *cron.Cron {
	cronJobs := cron.New()
	for _, target := range config.Targets {
		ctx := context.Background()
		targetURL := config.URL
		// validate that no traling slash is at the end of the URL
		if len(targetURL) > 0 && targetURL[len(targetURL)-1] == '/' {
			targetURL = targetURL[:len(targetURL)-1]
			logging.Logger.Warn("[Scraper] Removed trailing slash from URL", "url", targetURL)
		}
		if target.Filter != nil && *target.Filter != "" {
			targetURL = targetURL + *target.Filter
		}
		logging.Logger.Info("[Scraper] Scheduling scraping job", "target", target.Name, "url", targetURL, "cron_expression", target.ScrapeCronExpression)
		ctx = context.WithValue(ctx, "targetName", target.Name)
		ctx = context.WithValue(ctx, "targetURL", targetURL)
		ctx = context.WithValue(ctx, "outputPath", config.OutputRootPathJSON+"/"+target.Name)
		_, err := cronJobs.AddFunc(target.ScrapeCronExpression, func() {
			scraping.ScrapeNextbike(&ctx)
		})
		if err != nil {
			logging.Logger.Error("[Scraper] Failed to schedule scraping job", "target", target.Name, "err", err)
			continue
		}

		// add compactor job
		logging.Logger.Info("[Scraper] Scheduling compactor job", "target", target.Name, "cron_expression", target.CompactCronExpression)
		_, err = cronJobs.AddFunc(target.CompactCronExpression, func() {
			compactor.Compact(config.OutputRootPathJSON+"/"+target.Name, config.OutputRootPathParquet+"/"+target.Name)
		})
		if err != nil {
			logging.Logger.Error("[Scraper] Failed to schedule compactor job", "target", target.Name, "err", err)
			continue
		}
	}
	// Add cron job to move parquet files to samba share
	_, err := cronJobs.AddFunc(config.MoveParquetFilesToSambaShareCronExpression, func() {
		moveParquetFilesToSambaShare(config)
	})
	if err != nil {
		logging.Logger.Error("[Scraping Cron] Failed to schedule move parquet files job", "err", err)
	}

	cronJobs.Start()
	return cronJobs
}

// Function for compacting historical data from a specific target directory

// func main() {
// 	if err := logging.SetupLogger("./log"); err != nil {
// 		os.Exit(1)
// 	}
// 	logger := logging.Logger

// 	inputTargetDirPath := "/home/daniel/Desktop/Promotion/nextbike-collector/input"
// 	outputTargetDirPath := "/home/daniel/Desktop/Promotion/nextbike-collector/compacted/"
// 	logger.Info("Starting compaction process", "input_dir", inputTargetDirPath, "output_dir", outputTargetDirPath)
// 	compactor.Compact(inputTargetDirPath, outputTargetDirPath)
// 	logger.Info("Compaction process completed")
// }

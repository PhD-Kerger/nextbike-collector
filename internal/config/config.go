package config

import (
	"nextbike-scraper/internal/logging"
	"os"

	"gopkg.in/yaml.v3"
)

// URL: https://maps.nextbike.net/maps/nextbike-live.json
// output_root_path: /home/finn/Desktop/git/nextbike-scraper/output
// targets:
//   - name: Germany
//     filter: ?country=Germany
//     cron_expression: '0 0 * * *'  # Every day at midnight
//   - name: Global
//     filter: ''
//     cron_expression: '0 0 * * *'  # Every day at midnight

type ScraperConfig struct {
	URL                                        string   `yaml:"url"`
	OutputRootPathJSON                         string   `yaml:"output_root_path_json"`
	OutputRootPathParquet                      string   `yaml:"output_root_path_parquet"`
	Targets                                    []Target `yaml:"targets"`
	MoveParquetFilesToSambaShareCronExpression string   `yaml:"move_parquet_files_to_samba_share_cron_expression"`
	SMBUsername                                string   `yaml:"SMB_username"`
	SMBPassword                                string   `yaml:"SMB_password"`
	SMBWorkgroup                               string   `yaml:"SMB_workgroup"`
	SMBPathToMount                             string   `yaml:"SMB_path_to_mount"`
	SMBPathInsideMount                         string   `yaml:"SMB_path_inside_mount"` // Path inside the Samba share where the Parquet files will be stored
	SMBAddress                                 string   `yaml:"SMB_address"`
}

type Target struct {
	Name                  string  `yaml:"name"`
	Filter                *string `yaml:"filter,omitempty"` // Pointer to allow nil value
	ScrapeCronExpression  string  `yaml:"scrape_cron_expression"`
	CompactCronExpression string  `yaml:"compact_cron_expression"` // Pointer to allow nil value
}

// LoadConfig reads a YAML configuration file from the given path and unmarshals it into a ScraperConfig struct.
// It returns a pointer to the ScraperConfig and an error if reading or unmarshalling fails.
func LoadConfig(path string) (*ScraperConfig, error) {
	configContent, err := os.ReadFile(path)
	if err != nil {
		logging.Logger.Error("[Scraper] Failed to read config file", "path", path, "err", err)
		return nil, err
	}

	var config ScraperConfig
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		logging.Logger.Error("[Scraper] Failed to unmarshal config file", "path", path, "err", err)
		return nil, err
	}

	logging.Logger.Info("[Scraper] Config loaded", "path", path, "targets_loaded", len(config.Targets))
	return &config, nil
}

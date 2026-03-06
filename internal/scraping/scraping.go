package scraping

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"nextbike-scraper/internal/logging"
	"nextbike-scraper/internal/types"
	"os"
	"path/filepath"
	"time"
)

func ScrapeNextbike(ctx_ptr *context.Context) {
	ctx := *ctx_ptr
	logging.Logger.Info("["+ctx.Value("targetName").(string)+"] Starting scraping job", "target", ctx.Value("targetName"), "url", ctx.Value("targetURL"))
	req, err := http.NewRequestWithContext(ctx, "GET", ctx.Value("targetURL").(string), nil)
	if err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to create HTTP request", "error", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to execute HTTP request", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Received non-OK HTTP status", "status", resp.Status, "target", ctx.Value("targetName"))
		return
	}
	// Try to unmarshal the response body into the Root struct
	var root types.Root
	if err := json.NewDecoder(resp.Body).Decode(&root); err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to decode response body", "error", err)
		return
	}

	// Save the scraped data to a file
	// path is outputRootPath/targetName/yyyy-mm-dd/timestamp.json
	outputPath := ctx.Value("outputPath").(string) + "/" + time.Now().Format("2006-01-02") + "/" + fmt.Sprintf("%d.json", time.Now().Unix())
	// Ensure the output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to create output directory", "path", filepath.Dir(outputPath), "error", err)
		return
	}
	// Write the scraped data to the output file
	file, err := os.Create(outputPath)
	if err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to create output file", "path", outputPath, "error", err)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(root); err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to write to output file", "path", outputPath, "error", err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		logging.Logger.Error("["+ctx.Value("targetName").(string)+"] Failed to stat output file", "path", outputPath, "error", err)
		return
	}
	logging.Logger.Info("["+ctx.Value("targetName").(string)+"] Successfully saved scraped data", "path", outputPath, "size_bytes", fileInfo.Size())
}

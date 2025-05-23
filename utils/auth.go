package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// drawProgressBar prints a progress bar with ETA
func drawProgressBar(filename string, current, total int64, start time.Time) {
	percent := float64(current) / float64(total)
	width := 30 // characters wide
	filled := int(percent * float64(width))
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat(" ", width-filled) + "]"

	elapsed := time.Since(start).Seconds()
	speed := float64(current) / elapsed // bytes per second
	remaining := float64(total-current) / speed
	eta := time.Duration(remaining) * time.Second

	fmt.Printf("\r📤 Uploading %-20s %s %6.2f%% ⏳ ETA: %s", filename, bar, percent*100, eta.Truncate(time.Second))
}

// progressReader wraps an io.Reader to show progress with ETA
type progressReader struct {
	reader    io.Reader
	total     int64
	current   int64
	filename  string
	startTime time.Time
}

func (pr *progressReader) Read(p []byte) (int, error) {
	if pr.startTime.IsZero() {
		pr.startTime = time.Now()
	}
	n, err := pr.reader.Read(p)
	pr.current += int64(n)
	drawProgressBar(pr.filename, pr.current, pr.total, pr.startTime)
	return n, err
}

func configPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("❌ Unable to get user config directory: %v", err)
	}
	appConfigDir := filepath.Join(configDir, "gdrivecli")
	os.MkdirAll(appConfigDir, os.ModePerm)
	return filepath.Join(appConfigDir, "token.json")
}

func Authorize() error {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return fmt.Errorf("unable to parse credentials: %v", err)
	}

	tokenFile := configPath()
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}

	fmt.Println("✅ Authorization complete.")
	return nil
}

func UploadFile(filePath string) error {
	ctx := context.Background()
	tokenFile := configPath()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("unable to read credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return fmt.Errorf("unable to parse credentials: %v", err)
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return fmt.Errorf("unable to load saved token: %v", err)
	}

	client := config.Client(ctx, tok)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	stat, _ := f.Stat()
	progress := &progressReader{
		reader:   f,
		total:    stat.Size(),
		filename: filepath.Base(filePath),
	}

	file := &drive.File{Name: filepath.Base(filePath)}
	res, err := srv.Files.Create(file).
		Media(progress, googleapi.ChunkSize(10*1024*1024)).
		SupportsAllDrives(true).
		Fields("id").
		Do()

	if err != nil {
		return fmt.Errorf("error uploading file: %v", err)
	}

	fmt.Printf("\r📤 Uploading %-20s [██████████████████████████████] 100.00%% ✅ Completed\n", filepath.Base(filePath))
	fmt.Printf("✅ File uploaded successfully. File ID: %s\n", res.Id)
	return nil
}

// Helpers

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("💾 Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("🔑 Open the following URL in the browser and authorize:\n%v\n", authURL)

	// Try to open browser automatically
	exec.Command("xdg-open", authURL).Start()                                // Linux
	exec.Command("open", authURL).Start()                                    // macOS
	exec.Command("rundll32", "url.dll,FileProtocolHandler", authURL).Start() // Windows

	fmt.Print("Enter the authorization code: ")
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"gopkg.in/yaml.v2"
)

var (
	conf Conf
)

type (
	Conf struct {
		Bucket   string `yaml:"bucket"`
		Object   string `yaml:"object"`
		Filename string `yaml:"filename"`
	}
	update struct {
		current, total int64
	}
)

func init() {
	var configPath string
	if os.Getenv("ENV") != "" {
		configPath = "/root/"
	}
	configPath += "files/config.yaml"

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	driveService := InitDriveClient()
	gcsObject := GetGCSObject(conf.Bucket, conf.Object)

	ctx := context.Background()
	attr, err := gcsObject.Attrs(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	gcsSize := attr.Size
	r, err := gcsObject.NewReader(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	driveFile := &drive.File{
		Name: conf.Filename,
	}
	media := driveService.Files.Create(driveFile).Media(r, googleapi.ChunkSize(10000))

	updaterChan := make(chan update, 1000)
	go updater(updaterChan)
	media.ProgressUpdater(func(current, total int64) {
		updaterChan <- update{
			current: current,
			total:   gcsSize,
		}
	})
	_, err = media.Do()
	fmt.Println(err)
	fmt.Println("done")
}

func updater(ch <-chan update) {
	for u := range ch {
		currentS := size(u.current)
		totalS := size(u.total)
		fmt.Printf("Current %s Total %s (%.2f%%)\n", currentS, totalS, float64(u.current)/float64(u.total)*100)
	}
}

func size(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

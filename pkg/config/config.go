package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	PublicUrl       string `json:"public_url"`
}

func configPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".termup.json"), nil
}

func Exists() (bool, error) {
	path, err := configPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.PublicUrl == "" {
		cfg.PublicUrl = "https://your-bucket.s3.amazonaws.com/"
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func PromptForConfig() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Access Key ID: ")
	accessKeyID, _ := reader.ReadString('\n')

	fmt.Print("Enter Secret Access Key: ")
	secretAccessKey, _ := reader.ReadString('\n')

	fmt.Print("Enter Bucket Name: ")
	bucket, _ := reader.ReadString('\n')

	fmt.Print("Enter Endpoint: ")
	endpoint, _ := reader.ReadString('\n')

	fmt.Print("Enter Public URL (default: https://your-bucket.s3.amazonaws.com/): ")
	PublicUrl, _ := reader.ReadString('\n')
	PublicUrl = strings.TrimSpace(PublicUrl)
	if PublicUrl == "" {
		PublicUrl = "https://your-bucket.s3.amazonaws.com/"
	}

	return &Config{
		AccessKeyID:     strings.TrimSpace(accessKeyID),
		SecretAccessKey: strings.TrimSpace(secretAccessKey),
		Bucket:          strings.TrimSpace(bucket),
		Endpoint:        strings.TrimSpace(endpoint),
		PublicUrl:       PublicUrl,
	}, nil
}

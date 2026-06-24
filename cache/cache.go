package cache

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/daegalus/xxh3"
)

func LoadImageHashes() (map[string]string, error) {
	cacheDir := "cache"
	filePath := filepath.Join(cacheDir, "images.json")

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the file if it does not exist
		if err := os.WriteFile(filePath, []byte("{}"), 0644); err != nil {
			return nil, err
		}
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a map
	var imageHashes map[string]string
	if err := json.Unmarshal(data, &imageHashes); err != nil {
		return nil, err
	}

	return imageHashes, nil
}

func SaveImageHashes(imageHashes map[string]string) error {
	cacheDir := "cache"
	filePath := filepath.Join(cacheDir, "images.json")

	// Marshal the map into JSON
	data, err := json.MarshalIndent(imageHashes, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	return nil
}

func CalculateImageHashes() (map[string]string, error) {
	imageHashes := make(map[string]string)

	// Walk the images directory
	err := filepath.Walk("data/content/images", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate the hash of the image
		hash, err := HashFile(path)
		if err != nil {
			return err
		}

		// Add the hash to the map
		imageHashes[path] = hash

		return nil
	})
	if err != nil {
		return nil, err
	}

	return imageHashes, nil
}

func HashFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	hash := xxh3.Hash128(data).Bytes()
	return hex.EncodeToString(hash[:]), nil
}

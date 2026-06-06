package cloudinary

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var AllowedMIMEs = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

func Upload(ctx context.Context, data []byte, folder string) (string, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	log.Printf("[cloudinary] cloud=%q key=%q secret_len=%d", cloudName, apiKey, len(apiSecret))

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Printf("[cloudinary] NewFromParams error: %v", err)
		return "", fmt.Errorf("cloudinary config error: %w", err)
	}

	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{
		Folder:       folder,
		ResourceType: "auto",
	})
	if err != nil {
		log.Printf("[cloudinary] upload error: %v", err)
		return "", fmt.Errorf("cloudinary upload error: %w", err)
	}

	return resp.SecureURL, nil
}

package cloudinary

import (
	"bytes"
	"context"
	"fmt"

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
	cld, err := cloudinary.New()
	if err != nil {
		return "", fmt.Errorf("cloudinary config error: %w", err)
	}

	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{
		Folder:       folder,
		ResourceType: "auto",
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload error: %w", err)
	}

	return resp.SecureURL, nil
}

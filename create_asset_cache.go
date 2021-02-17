package pack

import (
	"context"
	"fmt"
	"github.com/buildpacks/pack/internal/blob"
	"github.com/buildpacks/pack/internal/dist"
	"github.com/google/go-containerregistry/pkg/name"
)

type CreateAssetCacheOptions struct {
	ImageName        string
	Assets           []dist.Asset
}

type AssetMetadata map[string]dist.Asset

func (c *Client) CreateAssetCache(ctx context.Context, opts CreateAssetCacheOptions) error {
	validOpts, err := validateConfig(opts)
	if err != nil {
		return err
	}

	// TODO -Dan- add support for remote image creation here
	img, err := c.imageFactory.NewImage(validOpts.ImageName, true)
	if err != nil {
		return fmt.Errorf("unable to create asset cache image: %q", err)
	}

	// TODO -Dan- parallelize these downloads using a threadpool.
	assetMap, err := c.downloadAssets(opts.Assets)
	if err != nil {
		return err
	}

	assetCacheImage := dist.NewAssetCacheImage(img, assetMap, opts.Assets)
	return assetCacheImage.Save()
}

func (c *Client) downloadAssets(assets []dist.Asset) (map[string]blob.Blob, error) {
	result := make(map[string]blob.Blob)
	for _, asset := range assets {
		b, err := c.downloader.Download(context.Background(), asset.URI, blob.RawDownload, blob.ValidateDownload(asset.Sha256))
		if err != nil {
			return map[string]blob.Blob{}, err
		}
		result[asset.Sha256] = b
	}
	return result, nil
}

func validateConfig(cfg CreateAssetCacheOptions) (CreateAssetCacheOptions, error) {
	tag, err := name.NewTag(cfg.ImageName, name.WeakValidation)
	if err != nil {
		return CreateAssetCacheOptions{}, fmt.Errorf("invalid asset cache image name: %q", err)
	}
	return CreateAssetCacheOptions{
		ImageName: tag.String(),
	}, nil
}

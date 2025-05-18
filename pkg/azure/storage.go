// pkg/azure/storage.go
package azure

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type BlobStorage struct {
	containerURL azblob.ContainerURL
}

func NewBlobStorage(accountName, accountKey, containerName string) (*BlobStorage, error) {
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %w", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	return &BlobStorage{
		containerURL: containerURL,
	}, nil
}

func (bs *BlobStorage) CreateContainer(ctx context.Context) error {
	_, err := bs.containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok {
			if serr.ServiceCode() == azblob.ServiceCodeContainerAlreadyExists {
				return nil
			}
		}
		return fmt.Errorf("failed to create container: %w", err)
	}
	return nil
}

func (bs *BlobStorage) UploadBlob(ctx context.Context, blobName string, data io.ReadSeeker, contentType string) (string, error) {
	blobURL := bs.containerURL.NewBlockBlobURL(blobName)

	options := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: contentType,
		},
	}

	_, err := azblob.UploadFileToBlockBlob(ctx, data, blobURL, options)
	if err != nil {
		return "", fmt.Errorf("failed to upload blob: %w", err)
	}

	sasURL, err := bs.GetBlobSasURL(ctx, blobName, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return sasURL, nil
}

func (bs *BlobStorage) GetBlobSasURL(ctx context.Context, blobName string, expiration time.Duration) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (bs *BlobStorage) DownloadBlob(ctx context.Context, blobName string, writer io.Writer) error {
	blobURL := bs.containerURL.NewBlockBlobURL(blobName)

	response, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		return fmt.Errorf("failed to download blob: %w", err)
	}

	reader := response.Body(azblob.RetryReaderOptions{})
	defer reader.Close()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("failed to read blob: %w", err)
	}

	return nil
}

func (bs *BlobStorage) DeleteBlob(ctx context.Context, blobName string) error {
	blobURL := bs.containerURL.NewBlockBlobURL(blobName)

	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}

	return nil
}
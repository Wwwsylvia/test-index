package main

import (
	"context"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/registry/remote"
)

func generateManifest(ctx context.Context, repo *remote.Repository, tag string) (v1.Descriptor, error) {
	blob := []byte(`test`)
	blobDesc := content.NewDescriptorFromBytes("test/blob", blob)
	opts := oras.PackManifestOptions{
		Layers: []v1.Descriptor{blobDesc},
	}
	desc, err := oras.PackManifest(ctx, repo, oras.PackManifestVersion1_1_RC4, "test/artifact", opts)
	if err != nil {
		return v1.Descriptor{}, err
	}
	if tag == "" {
		return desc, nil
	}
	return oras.Tag(ctx, repo, desc.Digest.String(), tag)
}

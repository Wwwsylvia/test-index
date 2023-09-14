package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content"
	"oras.land/oras-go/v2/errdef"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/credentials"
)

func main() {
	ctx := context.Background()

	// set up repository
	repoName := "myregistry.com/myrepo" // modify this
	repo, err := remote.NewRepository(repoName)
	if err != nil {
		panic(err)
	}
	credsStore, err := credentials.NewStoreFromDocker(credentials.StoreOptions{DetectDefaultNativeStore: true})
	if err != nil {
		panic(err)
	}
	authClient := auth.DefaultClient
	authClient.Credential = credentials.Credential(credsStore)
	repo.Client = authClient

	// push manifest
	m1, err := generateManifest(ctx, repo, "") // manifest with no tag
	if err != nil {
		panic(err)
	}
	m2, err := generateManifest(ctx, repo, "linux") // manifest with tag "v1"
	if err != nil {
		panic(err)
	}
	subject, err := generateManifest(ctx, repo, "subject") // manifest with tag "subject"
	if err != nil {
		panic(err)
	}

	// push index
	manifests := []v1.Descriptor{m1, m2}
	index, err := generateIndex(ctx, repo, manifests, &subject, "index") // index with tag "index"
	if err != nil {
		panic(err)
	}

	fmt.Println("pushed index digest:", index.Digest)
}

func generateManifest(ctx context.Context, repo *remote.Repository, tag string) (v1.Descriptor, error) {
	blob := []byte(`test`)
	blobDesc := content.NewDescriptorFromBytes("test/blob", blob)
	opts := oras.PackManifestOptions{
		Layers: []v1.Descriptor{blobDesc},
	}
	if err := repo.Blobs().Push(ctx, blobDesc, bytes.NewReader(blob)); err != nil && !errors.Is(err, errdef.ErrAlreadyExists) {
		return v1.Descriptor{}, err
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

func generateIndex(ctx context.Context, repo *remote.Repository, manifests []v1.Descriptor, subject *v1.Descriptor, tag string) (v1.Descriptor, error) {
	opts := oras.PackIndexOptions{
		Subject: subject,
	}
	desc, err := oras.PackIndex(ctx, repo, manifests, opts)
	if err != nil {
		return v1.Descriptor{}, err
	}
	if tag == "" {
		return desc, nil
	}
	return oras.Tag(ctx, repo, desc.Digest.String(), tag)
}

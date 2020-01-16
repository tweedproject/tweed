package stencil

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/apex/log"
	"github.com/apex/log/handlers/discard"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/types"
	imgspecv1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/openSUSE/umoci"
	"github.com/openSUSE/umoci/oci/cas/dir"
	"github.com/openSUSE/umoci/oci/casext"

	"github.com/google/uuid"
)

func (r *registry) loadStencilBundle(stencil string) (string, error) {
	// create a deterministic guid based on the stencil reference
	id := uuid.NewSHA1(uuid.Nil, []byte(stencil))
	bundlePath := path.Join(r.stencilsDir, id.String())

	// src reveres to the stencil image in the local docker registry
	ctx := context.Background()
	srcCtx := types.SystemContext{
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	srcRef, err := docker.ParseReference(fmt.Sprintf("//%s/%s", localRegistry, stencil))
	if err != nil {
		return "", fmt.Errorf("failed to parse source stencil image reference: %s", err)
	}

	// dst reveres to a temporary destination to which the image will be copied (in oci layout)
	dstDir, err := ioutil.TempDir("", "oci-image")
	if err != nil {
		return "", fmt.Errorf("failed to create tmp dir for stencil image: %s", err)
	}
	defer os.RemoveAll(dstDir)

	dstRef, err := layout.ParseReference(fmt.Sprintf("/%s:tmp", dstDir))
	if err != nil {
		return "", fmt.Errorf("failed to parse destination stencil image reference: %s", err)
	}

	// our local docker registry does not use ssl so use an insecure policy
	policy := &signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}}
	policyContext, err := signature.NewPolicyContext(policy)
	if err != nil {
		return "", fmt.Errorf("failed to create policy context: %s", err)
	}
	defer policyContext.Destroy()

	// now do the actual copy
	_, err = copy.Image(ctx, policyContext, dstRef, srcRef, &copy.Options{
		RemoveSignatures:      true,
		ReportWriter:          ioutil.Discard,
		SourceCtx:             &srcCtx,
		DestinationCtx:        &types.SystemContext{},
		ForceManifestMIMEType: imgspecv1.MediaTypeImageManifest,
		ImageListSelection:    copy.CopySystemImage,
	})
	if err != nil {
		return "", fmt.Errorf("failed to copy stencil oci image: %s", err)
	}

	// configure umoci (which is used to extract an oci image into a runc bundle)
	var meta umoci.Meta
	meta.Version = umoci.MetaVersion
	meta.MapOptions.KeepDirlinks = true

	// open the oci image
	engine, err := dir.Open(dstDir)
	if err != nil {
		return "", fmt.Errorf("failed to open stencil oci image: %s", err)
	}

	engineExt := casext.NewEngine(engine)
	defer engine.Close()

	// discard the image layer processing logs
	log.SetHandler(discard.New())
	// unpack the oci image to the stencil dir
	err = umoci.Unpack(engineExt, "tmp", bundlePath, meta.MapOptions)
	if err != nil {
		return "", fmt.Errorf("failed to unpack stencil rootfs: %s", err)
	}

	return bundlePath, nil
}

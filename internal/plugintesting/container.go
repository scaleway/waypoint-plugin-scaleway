package plugintesting

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var testDockerIMG = "docker.io/library/nginx:alpine"

// type used for error returned by docker registry
type errorRegistryMessage struct {
	Error string
}

// ContainerNamespace returns a new container namespace for the current test
// Add the namespace to the list of resource to be cleaned up after test
func (tt *TestTools) ContainerNamespace() *container.Namespace {
	api := container.NewAPI(tt.scwClient)
	namespace, err := api.CreateNamespace(&container.CreateNamespaceRequest{
		Name:        tt.genName("container"),
		Description: scw.StringPtr("container namespace creating by waypoint tests"),
	})
	if err != nil {
		tt.t.Fatal(fmt.Errorf("failed to create container namespace: %w", err))
	}
	tt.cleanupFunctions = append(tt.cleanupFunctions, func() error {
		namespace, err := api.DeleteNamespace(&container.DeleteNamespaceRequest{
			Region:      namespace.Region,
			NamespaceID: namespace.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete container namespace: %w", err)
		}

		_, err = api.WaitForNamespace(&container.WaitForNamespaceRequest{
			NamespaceID: namespace.ID,
			Region:      namespace.Region,
		})
		if err != nil && !is404Error(err) {
			return fmt.Errorf("failed to wait for container namespace to be deleted: %w", err)
		}

		_, err = registry.NewAPI(tt.scwClient).DeleteNamespace(&registry.DeleteNamespaceRequest{
			Region:      namespace.Region,
			NamespaceID: namespace.RegistryNamespaceID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete registry namespace: %w", err)
		}

		return err
	})
	return namespace
}

// UploadTestImage upload a test image to target namespace
// Skipped when not recording cassettes
func (tt *TestTools) UploadTestImage(namespace *container.Namespace) (string, error) {
	// Complete link to the image
	scwTag := namespace.RegistryEndpoint + "/nginx:test"

	if !UpdateCassettes {
		return scwTag, nil
	}

	api := container.NewAPI(tt.scwClient)
	namespace, err := api.WaitForNamespace(&container.WaitForNamespaceRequest{
		NamespaceID: namespace.ID,
		Region:      namespace.Region,
	})
	if err != nil {
		return "", fmt.Errorf("error waiting namespace: %v", err)
	}

	var errorMessage errorRegistryMessage

	accessKey, _ := tt.scwClient.GetAccessKey()
	secretKey, _ := tt.scwClient.GetSecretKey()
	authConfig := types.AuthConfig{
		ServerAddress: namespace.RegistryEndpoint,
		Username:      accessKey,
		Password:      secretKey,
	}

	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv, docker.WithAPIVersionNegotiation())
	if err != nil {
		return "", fmt.Errorf("failed to init docker client: %w", err)
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", fmt.Errorf("could not marshal auth config: %v", err)
	}

	ctx := context.Background()
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	tt.t.Log("Pulling test image")

	out, err := dockerClient.ImagePull(ctx, testDockerIMG, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("could not pull image: %v", err)
	}
	defer out.Close()
	buffIOReader := bufio.NewReader(out)
	for {
		streamBytes, errPull := buffIOReader.ReadBytes('\n')
		if errPull == io.EOF {
			break
		}
		err = json.Unmarshal(streamBytes, &errorMessage)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal: %v", err)
		}

		if errorMessage.Error != "" {
			return "", fmt.Errorf(errorMessage.Error)
		}
	}

	imageTag := testDockerIMG

	err = dockerClient.ImageTag(ctx, imageTag, scwTag)
	if err != nil {
		return "", fmt.Errorf("could not tag image: %v", err)
	}

	tt.t.Log("Pushing test image")

	pusher, err := dockerClient.ImagePush(ctx, scwTag, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return "", fmt.Errorf("could not push image: %v", err)
	}
	defer pusher.Close()

	buffIOReader = bufio.NewReader(pusher)
	for {
		streamBytes, errPush := buffIOReader.ReadBytes('\n')
		if errPush == io.EOF {
			break
		}
		err = json.Unmarshal(streamBytes, &errorMessage)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal: %v", err)
		}

		if errorMessage.Error != "" {
			return "", fmt.Errorf(errorMessage.Error)
		}
	}

	_, err = api.WaitForNamespace(&container.WaitForNamespaceRequest{
		NamespaceID: namespace.ID,
		Region:      namespace.Region,
	})
	if err != nil {
		return "", fmt.Errorf("error waiting namespace: %v", err)
	}
	return scwTag, nil
}

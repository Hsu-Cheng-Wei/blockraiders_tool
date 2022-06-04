package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"path/filepath"
	"strings"
)

type HostDocker struct {
	Cli *client.Client
}

func NewHost() (*HostDocker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		return nil, err
	}

	return &HostDocker{
		Cli: cli,
	}, nil
}

func (host *HostDocker) Build(path string, tag string, args map[string]*string) error {
	ctx := context.Background()

	dir := filepath.Dir(path)
	file := filepath.Base(path)

	fmt.Printf("Build Project %s \r\n", dir)
	fmt.Printf("Tag: %s \r\n", tag)
	tar, err := archive.TarWithOptions(dir, &archive.TarOptions{})

	if err != nil {
		return err
	}
	reps, err := host.Cli.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Dockerfile: file,
		Tags:       []string{tag},
		NoCache:    true,
		Remove:     true,
		BuildArgs:  args,
	})

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reps.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	_ = reps.Body.Close()

	return nil
}

func (host *HostDocker) Push(tag string, password string) error {
	ctx := context.Background()
	auth := types.AuthConfig{
		Username: "AWS",
		Password: password,
	}

	encoding, err := json.Marshal(auth)
	if err != nil {
		return err
	}

	authStr := base64.URLEncoding.EncodeToString(encoding)

	reader, err := host.Cli.ImagePush(ctx, tag, types.ImagePushOptions{
		RegistryAuth: authStr,
	})

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}

func (host *HostDocker) Login(region string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return "", err
	}
	svc := ecr.New(sess, aws.NewConfig().WithMaxRetries(10).WithRegion(region))

	resp, err := svc.GetAuthorizationToken(&ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return "", err
	}

	decoding, _ := base64.StdEncoding.DecodeString(*resp.AuthorizationData[0].AuthorizationToken)

	return strings.Split(string(decoding), ":")[1], nil
}

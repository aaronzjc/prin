package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

type Target struct {
	Os   string
	Arch string
	Out  string
}

const (
	AppVersion = "1.4"
)

func main() {
	var action string
	if len(os.Args) >= 2 {
		action = os.Args[1]
	}

	ctx := context.Background()
	switch action {
	case "backend":
		wrapStep(ctx, "backend", buildBackend)
	case "frontend":
		wrapStep(ctx, "frontend", buildFrontend)
	case "image":
		wrapStep(ctx, "backend", buildBackend)
		wrapStep(ctx, "frontend", buildFrontend)
		wrapStep(ctx, "image", buildAndPushImage)
	case "deploy":
		wrapStep(ctx, "backend", buildBackend)
		wrapStep(ctx, "frontend", buildFrontend)
		wrapStep(ctx, "image", buildAndPushImage)
		wrapStep(ctx, "deploy", deploy)
	default:
		fmt.Println("usage: go run ./scripts/dagger.go [backend|frontend|image|deploy]")
	}
}

func wrapStep(ctx context.Context, name string, fn func(context.Context) error) {
	if err := fn(ctx); err != nil {
		fmt.Printf("[x]run %s failed, err = %v\n", name, err)
	}
	fmt.Printf("[o]run %s done\n", name)
}

func buildFrontend(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Workdir()
	npm := client.Container().From("node:14-alpine")
	npm = npm.WithMountedDirectory("/src", src).WithWorkdir("/src/web")

	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"rm", "-rf", "node_modules"},
	})
	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "install", "--sass_binary_site=https://npm.taobao.org/mirrors/node-sass/"},
	})
	npm = npm.Exec(dagger.ContainerExecOpts{
		Args: []string{"npm", "run", "build"},
	})
	build, err := npm.Stdout().Contents(ctx)
	if err != nil {
		return err
	}
	output := npm.Directory("/src/public")
	if _, err := output.Export(ctx, "dagger/frontend"); err != nil {
		return err
	}
	fmt.Println("npm stdout", build)
	return nil
}

func buildBackend(ctx context.Context) error {
	fmt.Println("start build with dagger")
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	// 获取本地项目路径
	src := client.Host().Workdir()

	targets := []Target{
		{
			Os:   "linux",
			Arch: "amd64",
			Out:  "prin",
		},
		{
			Os:   "windows",
			Arch: "amd64",
			Out:  "prin.exe",
		},
	}
	outputs := client.Directory()
	golang := client.Container().From("golang:1.19-alpine3.15")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")
	golang = golang.WithEnvVariable("GO111MODULE", "on")
	golang = golang.WithEnvVariable("GOPROXY", "https://goproxy.cn,direct")
	golang = golang.WithEnvVariable("CGO_ENABLED", "0")
	for _, target := range targets {
		build := golang.WithEnvVariable("GOOS", target.Os)
		build = build.WithEnvVariable("GOARCH", target.Arch)
		path := fmt.Sprintf("dagger/backend/%s/", target.Os)
		build = build.Exec(dagger.ContainerExecOpts{
			Args: []string{"go", "build", "-o", path + target.Out, "cmd/main.go"},
		})
		outputs = outputs.WithDirectory(path, build.Directory(path))
	}

	if _, err := outputs.Export(ctx, "."); err != nil {
		return err
	}

	return nil
}

func buildAndPushImage(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Workdir()
	docker := client.Container()
	docker = docker.Build(src, dagger.ContainerBuildOpts{Dockerfile: "./scripts/Dockerfile"})
	resp, err := docker.Publish(ctx, "aaronzjc/prin:"+AppVersion)
	if err != nil {
		return err
	}
	fmt.Println(resp)

	return nil
}

func deploy(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return err
	}
	defer client.Close()

	kubeconfig := client.Host().Workdir().File("./scripts/kubeconf.yaml")
	deployment := client.Host().Workdir().File("./scripts/k8s/Deployment.yaml")
	kubectl := client.Container().From("bitnami/kubectl")
	kubectl = kubectl.WithMountedFile("/.kube/config", kubeconfig)
	kubectl = kubectl.WithMountedFile("/tmp/deployment.yaml", deployment)
	kubectl = kubectl.Exec(dagger.ContainerExecOpts{
		Args: []string{"apply", "-f", "/tmp/deployment.yaml", "-n", "k3s-apps"},
	})
	logs, err := kubectl.Stdout().Contents(ctx)
	if err != nil {
		return err
	}
	fmt.Println(logs)
	return nil
}

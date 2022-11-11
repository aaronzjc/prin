package main

import (
	"context"
	"fmt"
	"os"
	"prin/pkg/flow"

	"dagger.io/dagger"
)

const (
	AppVersion = "1.9"
)

func main() {
	var action string
	if len(os.Args) >= 2 {
		action = os.Args[1]
	}

	fw := flow.NewFlow(context.Background())
	switch action {
	case "backend":
		fw.Step("backend", buildBackend)
	case "frontend":
		fw.Step("frontend", buildFrontend)
	case "image":
		fw.Step("backend", buildBackend)
		fw.Step("frontend", buildFrontend)
		fw.Step("image", buildAndPushImage)
	case "deploy":
		fw.Step("backend", buildBackend)
		fw.Step("frontend", buildFrontend)
		fw.Step("image", buildAndPushImage)
		fw.Step("deploy", deploy)
	default:
		fmt.Println("usage: go run ./scripts/dagger.go [backend|frontend|image|deploy]")
		return
	}
	fw.Run()
}

func buildFrontend(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	src := client.Host().Directory("./web", dagger.HostDirectoryOpts{
		Exclude: []string{"node_modules"},
	})
	npm := client.Container().From("node:14-alpine")
	npm = npm.WithMountedDirectory("/src/web", src).WithWorkdir("/src/web")
	npm = npm.WithEnvVariable("VERSION", AppVersion)
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
	if _, err := npm.Directory("/src/public").Export(ctx, "dagger/frontend"); err != nil {
		return err
	}
	fmt.Println("npm stdout", build)
	return nil
}

func buildBackend(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		return err
	}
	defer client.Close()

	envs := map[string]string{
		"GO111MODULE": "on",
		"GOPROXY":     "https://goproxy.cn,direct",
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}

	// 获取本地项目路径
	src := client.Host().Workdir()
	golang := client.Container().From("golang:1.19-alpine3.15")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")
	for k, v := range envs {
		golang = golang.WithEnvVariable(k, v)
	}
	path := "dagger/backend/"
	golang = golang.Exec(dagger.ContainerExecOpts{
		Args: []string{"go", "build", "-o", path + "prin", "cmd/main.go"},
	})

	if _, err := golang.Directory(path).Export(ctx, path); err != nil {
		return err
	}

	return nil
}

func buildAndPushImage(ctx context.Context) error {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
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
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
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

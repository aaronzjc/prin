package main

import (
	"context"
	"fmt"
	"os"
	"prin/pkg/flow"
	"strings"

	"dagger.io/dagger"
)

const (
	Version = "1.26"
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
	npm = npm.WithEnvVariable("VERSION", Version)
	npm = npm.WithExec([]string{"npm", "config", "set", "registry", "https://registry.npmmirror.com"})
	npm = npm.WithExec([]string{"npm", "install", "--sass_binary_site=https://npm.taobao.org/mirrors/node-sass/"})
	npm = npm.WithExec([]string{"npm", "run", "build"})
	build, err := npm.Stdout(ctx)
	if err != nil {
		return err
	}

	dst := "dagger/frontend"
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if _, err := npm.Directory("/src/public").Export(ctx, dst); err != nil {
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
	src := client.Host().Directory(".")
	golang := client.Container().From("golang:1.19-alpine3.15")
	golang = golang.WithMountedDirectory("/src", src).WithWorkdir("/src")
	for k, v := range envs {
		golang = golang.WithEnvVariable(k, v)
	}
	path := "dagger/backend/"
	golang = golang.WithExec([]string{"go", "build", "-o", path + "prin", "cmd/main.go"})

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

	src := client.Host().Directory(".")
	docker := client.Container()
	docker = docker.Build(src, dagger.ContainerBuildOpts{Dockerfile: "./scripts/Dockerfile"})
	resp, err := docker.Publish(ctx, "aaronzjc/prin:"+Version)
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

	// 处理版本
	oldTag, newTag := "latest", Version
	file := "./scripts/k8s.yaml"
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out := strings.ReplaceAll(string(data), oldTag, newTag)
	os.WriteFile(file, []byte(out), 0666)
	defer os.WriteFile(file, data, 0666)

	kubectl := client.Container().From("bitnami/kubectl")
	kubeconfig := client.Host().Directory(".").File("./scripts/kubeconf.yaml")
	kubectl = kubectl.WithMountedFile("/.kube/config", kubeconfig)
	deployment := client.Host().Directory(".").File(file)
	kubectl = kubectl.WithMountedFile("/tmp/deployment.yaml", deployment)

	kubectl = kubectl.WithExec([]string{"apply", "-f", "/tmp/deployment.yaml", "-n", "k3s-apps"})
	logs, err := kubectl.Stdout(ctx)
	if err != nil {
		return err
	}
	fmt.Println(logs)
	return nil
}

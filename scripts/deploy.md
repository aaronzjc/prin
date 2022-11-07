## dagger

本项目采用`dagger`进行构建部署。

### 构建前端

```shell
$ go run ./scripts/dagger.go frontend
```

### 构建后端

```shell
$ go run ./scripts/dagger.go backend
```

### 构建镜像

```shell
$ go run ./scripts/dagger.go image
```

### 部署至k8s

```shell
$ go run ./scripts/dagger.go deploy
```
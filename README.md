# envconfig-generate

Parse [`envconfig`](https://github.com/sethvargo/go-envconfig) struct tags and list them all in one convenient place

## How to use

### Whole project

```shell
envconfig-generate ./... > CONFIG.md
```

### Individual files

```shell
envconfig-generate $(find ~/my-project -name '*config*.go') > CONFIG.md
```

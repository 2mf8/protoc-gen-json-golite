# protoc-gen-json-golite

基于 `protoc-gen-go` 修改而来

<img width="100%" src="https://starify.komoridevs.icu/api/starify?owner=2mf8&repo=protoc-gen-json-golite" alt="starify" />

## ![alt text](0019E536.png)星星曲线
[![Stargazers over time](https://starchart.cc/2mf8/protoc-gen-json-golite.svg?variant=adaptive)](https://starchart.cc/2mf8/protoc-gen-json-golite)

## proto 文件 快速生成 json Golang 结构体

# 安装

```
go install github.com/2mf8/protoc-gen-json-golite@latest
```

# 使用
与 `proto` 文件同目录

```
protoc --json-golite_out=. *.proto
```

## 注： 暂不支持 `oneof`

示例请查看 [example](https://github.com/2mf8/protoc-gen-json-golite/blob/main/example/onebot/group.json.go)

package(default_visibility = ["//visibility:public"])

load("@rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "almostio",
    srcs = [
        "fixedsizewriter.go",
        "multiwritecloser.go",
        "marshal.go",
        "overlay.go",
    ],
    importpath = "github.com/lanseg/golang-commons/almostio",
    deps = [

    ],
)

go_test(
    name = "almostio_test",
    size = "small",
    srcs = [
        "overlay_test.go",
    ],
    embed = [
        ":almostio",
    ],
)


go_test(
    name = "almostio_benchmark",
    size = "large",
    srcs = [
        "overlay_benchmark.go",
    ],
    embed = [
        ":almostio",
    ],
)

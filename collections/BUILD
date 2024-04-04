load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "collections",
    srcs = [
        "collections.go",
        "iterator.go",
        "streams.go",
        "tuples.go",
    ],
    importpath = "github.com/lanseg/golang-commons/collections",
    deps = [
    ],
)

go_test(
    name = "collections_test",
    size = "small",
    srcs = [
        "collections_test.go",
        "iterator_test.go",
        "streams_test.go",
        "tuples_test.go",
    ],
    embed = [
        ":collections",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "collections",
    srcs = [
        "tuples.go",
        "iterator.go",
    ],
    importpath = "github.com/lanseg/golang-commons/collections",
    deps = [
    ],
)

go_test(
    name = "collections_test",
    size = "small",
    srcs = [
        "tuples_test.go",
        "iterator_test.go",
    ],
    embed = [
        ":collections",
    ],
)

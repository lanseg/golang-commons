load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "optional",
    srcs = [
        "optional.go",
    ],
    importpath = "github.com/lanseg/golang-commons/optional",
    deps = [
    ],
)

go_test(
    name = "optional_test",
    size = "small",
    srcs = [
        "optional_test.go",
    ],
    embed = [
        ":optional",
    ],
)

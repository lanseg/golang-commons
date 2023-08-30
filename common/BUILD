load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "common",
    srcs = [
        "common.go",
    ],
    importpath = "github.com/lanseg/golang-commons/common",
    deps = [
    ],
)

go_test(
    name = "common_test",
    size = "small",
    srcs = [
        "common_test.go",
    ],
    embed = [
        ":common",
    ],
)

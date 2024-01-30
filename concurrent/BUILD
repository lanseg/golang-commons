load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "concurrent",
    srcs = [
        "waiter.go",
    ],
    importpath = "github.com/lanseg/golang-commons/concurrent",
    deps = [
        "//optional"
    ],
)

go_test(
    name = "concurrent_test",
    size = "small",
    srcs = [
        "waiter_test.go",
    ],
    embed = [
        ":concurrent",
    ],
)

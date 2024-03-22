load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "concurrent",
    srcs = [
        "executor.go",
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
        "executor_test.go",
        "waiter_test.go",
    ],
    embed = [
        ":concurrent",
    ],
)

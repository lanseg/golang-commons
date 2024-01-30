load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "common",
    srcs = [
        "common.go",
        "config.go",
        "encoding.go",
        "logging.go",
        "strings.go",
        "uuid.go",
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
        "config_test.go",
        "encoding_test.go",        
        "logging_test.go",
        "strings_test.go",
        "uuid_test.go",
    ],
    data = glob(["testdata/*"]),     
    embed = [
        ":common",
    ],
)

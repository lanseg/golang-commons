load("@rules_go//go:def.bzl", "go_library", "go_test")

package(default_visibility = ["//visibility:public"])

go_library(
    name = "almosthtml",
    srcs = [
        "tokenizer.go",
        "html.go"
    ],
    importpath = "github.com/lanseg/golang-commons/almosthtml",
    deps = [
        "//optional",
        "//collections",
    ],
)

go_test(
    name = "almosthtml_test",
    size = "small",
    srcs = [
        "html_test.go",
        "tokenizer_test.go",
    ],
    data = glob(["testdata/*"]),
    embed = [
        ":almosthtml",
    ],
)

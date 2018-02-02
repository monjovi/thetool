package cmd

import (
	"strings"
	"text/template"

	"github.com/solo-io/thetool/pkg/downloader"
	"github.com/solo-io/thetool/pkg/feature"
)

const (
	buildContent = `package(default_visibility = ["//visibility:public"])

load(
    "@envoy//bazel:envoy_build_system.bzl",
    "envoy_cc_binary",
    "envoy_cc_library",
    "envoy_cc_test",
)

envoy_cc_binary(
    name = "envoy",
    repository = "@envoy",
    deps = [{{range .}}
		"@{{.Name}}//:filter_lib",{{end}}
		"@envoy//source/exe:envoy_main_entry_lib",
    ],
)
`

	workspaceContent = `load('@bazel_tools//tools/build_defs/repo:git.bzl', 'git_repository')
{{range .}}
local_repository(
    name = "{{.Name}}",
    path = "{{. | path}}",
)

{{end}}
bind(
    name = "boringssl_crypto",
    actual = "//external:ssl",
)

ENVOY_SHA = "29989a38c017d3be5aa3c735a797fcf58b754fe5"
http_archive(
    name = "envoy",
    strip_prefix = "envoy-" + ENVOY_SHA,
    url = "https://github.com/envoyproxy/envoy/archive/" + ENVOY_SHA + ".zip",
)

load("@envoy//bazel:repositories.bzl", "envoy_dependencies")
envoy_dependencies()
load("@envoy//bazel:cc_configure.bzl", "cc_configure")
cc_configure()
load("@envoy_api//bazel:repositories.bzl", "api_dependencies")
api_dependencies()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@com_lyft_protoc_gen_validate//bazel:go_proto_library.bzl", "go_proto_repositories")
go_proto_repositories(shared=0)
go_rules_dependencies()
go_register_toolchains()
load("@io_bazel_rules_go//proto:def.bzl", "proto_register_toolchains")
proto_register_toolchains()
`

	dockerfile = `FROM ubuntu:16.04

ADD WORKSPACE /etc/envoy.WORKSPACE
ADD envoy /usr/local/bin/envoy

CMD /usr/local/bin/envoy -c /etc/envoy.yaml --service-cluster $CLUSTER --service-node $NODE`

	helmValuesContent = `#custom values for Glue Helm chart
#gateway
gw:
  type: NodePort
  replicaCount: 1
  port: 80
  image: "{{ .EnvoyImage }}"
  imageTag: {{ .EnvoyTag }}
  serviceCluster: cluster
  serviceNode: node
`
)

var (
	buildTemplate      *template.Template
	workspaceTemplate  *template.Template
	helmValuesTemplate *template.Template
)

func init() {
	buildTemplate = template.Must(template.New("build").Parse(buildContent))
	funcMap := template.FuncMap{
		"path": path,
	}
	workspaceTemplate = template.Must(template.New("workspace").
		Funcs(funcMap).Parse(workspaceContent))
	helmValuesTemplate = template.Must(template.New("helm").Parse(helmValuesContent))
}

func path(f feature.Feature) string {
	if strings.HasSuffix(f.Repository, ".git") {
		return "external/" + downloader.RepoDir(f.Repository) + "/envoy"
	}

	if isGitHubHTTP(f.Repository) {
		return "external/" + f.Name + "-" + f.Version + "/envoy"
	}

	return "external/" + downloader.RepoDir(f.Repository) + "/envoy"
}

func isGitHubHTTP(url string) bool {
	return strings.HasPrefix(url, "https://github.com/")
}

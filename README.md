# dsm

Data structure manipulation. Used in e.g. deployment pipelines.



```
$ dsm -h
Data structure manipulation. Used in e.g. deployment pipelines.

Usage:
  dsm [flags]
  dsm [command]

Available Commands:
  completion  Generate shell completions.
  help        Help about any command
  search      Search for values within YAML or JSON data structures.
  update      Update values within YAML or JSON data structures.
  verify      Verify the consistency of values within YAML or JSON data structures.
  version     Print version information of this command line tool.

Flags:
  -h, --help   help for dsm

Use "dsm [command] --help" for more information about a command.
```



```
$ dsm search -h
Search for values within YAML or JSON data structures. Consider the following HelmRelease CR
defining a Docker image tag in its spec

    apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
    kind: "HelmRelease"
    metadata:
      name: "apiserver"
    spec:
      values:
        image:
          tag: "8469445410f8a74d72af0cf430ed8dd44fb6b8fa"

The following example shows how to extract the image tag from the YAML file.

    $ dsm search -r HelmRelease -n apiserver -k spec.values.image.tag
    8469445410f8a74d72af0cf430ed8dd44fb6b8fa

Usage:
  dsm search [flags]

Flags:
  -h, --help              help for search
  -k, --key string        JSON path key to work with.
  -n, --name string       Metadata name of the resources to work with.
  -r, --resource string   Resource kind to work with.
  -s, --source string     Source directory to traverse. (default ".")
```



```
$ dsm update -h
Update values within YAML or JSON data structures. Consider the following HelmRelease CR
defining a Docker image tag in its spec

    apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
    kind: "HelmRelease"
    metadata:
      name: "apiserver"
    spec:
      values:
        image:
          tag: "8469445410f8a74d72af0cf430ed8dd44fb6b8fa"

The following example shows how to modify the image tag of the YAML file.

    dsm update -r HelmRelease -n apiserver -k spec.values.image.tag -v <new-sha>

Usage:
  dsm update [flags]

Flags:
  -h, --help              help for update
  -k, --key string        JSON path key to work with.
  -n, --name string       Metadata name of the resources to work with.
  -r, --resource string   Resource kind to work with.
  -s, --source string     Source directory to traverse. (default ".")
  -v, --value string      JSON path value to work with.
```



```
$ dsm verify -h
Verify the consistency of values within YAML or JSON data structures.
Consider multiple of the following HelmRelease CRs defining a Docker image
tag in its spec.

    apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
    kind: "HelmRelease"
    metadata:
      name: "apiserver"
    spec:
      values:
        image:
          tag: "8469445410f8a74d72af0cf430ed8dd44fb6b8fa"

The following example shows how to validate the image tag of all the YAML
files. In case the validation succeeds the command exits with an exit code 0
and without output. If the validation fails the command exits with an exit
code 1 and with the output of a stack trace.

    $ dsm verify -r HelmRelease -n apiserver -k spec.values.image.tag
    program panic

    {
        "kind": "invalidValueError",
        "stck": [
            "/Users/xh3b4sd/project/xh3b4sd/dsm/cmd/verify/runner.go:87",
            "/Users/xh3b4sd/project/xh3b4sd/dsm/cmd/verify/runner.go:30",
            "/Users/xh3b4sd/project/xh3b4sd/dsm/main.go:47"
        ],
        "type": "*tracer.Error"
    }

Usage:
  dsm verify [flags]

Flags:
  -h, --help              help for verify
  -k, --key string        JSON path key to work with.
  -n, --name string       Metadata name of the resources to work with.
  -r, --resource string   Resource kind to work with.
  -s, --source string     Source directory to traverse. (default ".")
```

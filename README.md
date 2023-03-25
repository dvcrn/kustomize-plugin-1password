# 1Password CLI Kustomize Plugin

[Kustomize Plugin](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/go_plugins/) to generate secrets directly from 1Password values

**Note:** This is **not** for using the 1Password API or secret automation. This plugin aims to pull data out of the locally installed 1Password and generates secrets with those.

## Installation

### Requirements

- Make sure you have the `op` CLI installed: https://1password.com/downloads/command-line/
- Locally installed Golang

Make sure you can actually use `op` by trying a command like: `op vault list`.

### Build + Install

Because of the [skew problem](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/go_plugins/#the-skew-problem), you are _very likely_ required to build both `kustomize` as well as this plugin with the same Golang versions.

You're welcome to try without doing this, but you may run into issues.

1. Install `kustomize` through Golang:

```
go install sigs.k8s.io/kustomize/kustomize/v5@latest
```

2. Prepare the plugin dir if you haven't yet

```bash
mkdir -p ~/.config/kustomize/plugin/sh.d.kustomize/v1/opclisecret/
```

3. Build this plugin

```
go build -buildmode plugin \
        -o ~/.config/kustomize/plugin/sh.d.kustomize/v1/opclisecret/OpCLISecret.so .
```

4. Put the plugin into `~/.config/kustomize/plugin/sh.d.kustomize/v1/opclisecret/OpCLISecret.so` (the command under 3. is already doing that for you)

## Usage

Add a new manifest with `kind: OpCLISecret`

```yaml
# example: netflixSecret.yaml

apiVersion: sh.d.kustomize/v1
kind: OPCLISecret
metadata:
  name: myopsecret
  namespace: default # default is default
type: Opaque # opaque is default
values:
  - key: mySecretKey
    opPath: /Kustomize/Netflix/username
# options:
#    disableNameSuffixHash: true
```

Specify under `values`:

- `key`: What key do you want to use within the secret?
- `opPath`: 1Password path to the field you wish to use for this secret in the form `/vault/item/field`. `Vault` and `Item` can be names or IDs. I recommend using IDs

Add the generator to `kustomization.yaml`:

```yaml
generators:
  - netflixSecret.yaml
```

Run kustomize with the `--enable-alpha-plugins` flag: `kustomize build --enable-alpha-plugins`

The above example will generate the following secret:

```yaml
apiVersion: v1
data:
  mySecretKey: <value of "username" of item 'Netflix' within 1Password>
kind: Secret
metadata:
  name: myopsecret-24b5hmbhk5
  namespace: default
type: Opaque
```

## Troubleshooting

**How can I get rid of the "attempting to load plugin" message? I can't pipe into `kubectl` anymore!**

Use sed:

```bash
kustomize build --enable-alpha-plugins | sed 1d | kubectl apply -f -
```

**I'm getting "plugin was built with a different version of package x**

This usually means that the dependencies of this plugin and `kustomize` are no longer in sync.

Go to the `go.mod` file of the kustomize repo (https://github.com/kubernetes-sigs/kustomize/blob/master/kustomize/go.mod), copy the contents and put them into `go.mod` of this repository. Don't delete the `github.com/dvcrn/go-op-cli` key, you need that.

Oh and please make a Pull Request to this repo with the new dependencies :)

Make sure `which kustomize` is actually the version that you installed with Go (same version that built this plugin), and not something installed through brew for example.

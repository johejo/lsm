# Language Server Manager (LSM)

[![ci](https://github.com/johejo/lsm/workflows/ci/badge.svg)](https://github.com/johejo/lsm/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/johejo/lsm/branch/master/graph/badge.svg)](https://codecov.io/gh/johejo/lsm)

## Motivation

Language Server is difficult to set up and is tightly coupled with each editor plugin.<br>
https://github.com/neovim/nvim-lsp/issues/200#issuecomment-618807759

There are also lots of scripts to support different platforms and Language Servers, which is becoming difficult to maintain.<br>
https://github.com/mattn/vim-lsp-settings/tree/master/installer

LSM aims to be a simple and cross platform installation manager for Language Server.

## Notice
Windows support is still incomplete.

## Language Server Destination

### macOS, Linux

```
$HOME/.local/share/lsm/servers
```

If you defined `$XDG_DATA_HOME`

```
$XDG_DATA_HOME/lsm/servers
```

### Windows

```
%LOCALAPPDATA%\lsm\servers
```

## Install

go get
```
go get github.com/johejo/lsm
```

pre-built binary
```
# COMMING SOON
```

## Usage

For example, when gopls is the target

install
```
lsm install gopls
```

uninstall
```
lsm uninstall gopls
```

list (Installation status of Language Servers)

```
lsm list
```

## Supported Language Servers

- [bash-language-server](https://github.com/bash-lsp/bash-language-server)
- [cmake-language-server](https://github.com/regen100/cmake-language-server)
- [dockerfile-language-server-nodejs](https://github.com/rcjsuen/dockerfile-language-server-nodejs)
- [eclipse.jdt.ls](https://github.com/eclipse/eclipse.jdt.ls)
- [efm-langserver](https://github.com/mattn/efm-langserver)
- [eslint-server(vscode-eslint)](https://github.com/microsoft/vscode-eslint)
- [fortran-language-server](https://github.com/hansec/fortran-language-server)
- [gopls](https://github.com/golang/tools/tree/master/gopls)
- [graphql-language-service-cli](https://github.com/graphql/graphiql)
- [kotlin-language-server](https://github.com/fwcd/kotlin-language-server)
- [lemminx](https://github.com/eclipse/lemminx) ([vscode-xml](https://github.com/redhat-developer/vscode-xml))
- [metals](https://scalameta.org/metals/)
- [purescript-language-server](https://github.com/nwolverson/purescript-language-server)
- [python-language-server](https://github.com/palantir/python-language-server)
- [reason-language-server](https://github.com/jaredly/reason-language-server)
- [rust-analyzer](https://rust-analyzer.github.io/)
- [sqls](https://github.com/lighttiger2505/sqls)
- [svelte-language-server](https://github.com/sveltejs/language-tools/tree/master/packages/language-server)
- [terraform-ls](https://github.com/hashicorp/terraform-ls)
- [terraform-lsp](https://github.com/juliosueiras/terraform-lsp)
- [typescript-language-server](https://github.com/theia-ide/typescript-language-server)
- [vim-language-server](https://github.com/iamcco/vim-language-server)
- [vls](https://github.com/vuejs/vetur/tree/master/server)
- [vscode-css-languagesever](https://github.com/vscode-langservers/vscode-css-languageserver)
- [vscode-html-languagesever](https://github.com/vscode-langservers/vscode-html-languageserver)
- [vscode-json-languagesever](https://github.com/vscode-langservers/vscode-json-languageserver)
- [yaml-language-sever](https://github.com/redhat-developer/yaml-language-server/)

We plan to add more Language Servers.<br>
If you want to add a new Language Server, please create a PR.<br>

## Thanks

Language Server Manager inherits the concept of [vim-lsp-settings](https://github.com/mattn/vim-lsp-settings).

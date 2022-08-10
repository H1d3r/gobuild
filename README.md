# gobuild

[![Test](https://github.com/caixw/gobuild/workflows/Test/badge.svg)](https://github.com/caixw/gobuild/actions?query=workflow%3ATest)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/caixw/gobuild)](https://pkg.go.dev/github.com/caixw/gobuild)
![Go version](https://img.shields.io/github/go-mod/go-version/caixw/gobuild)
![License](https://img.shields.io/github/license/caixw/gobuild)

gobuild 是一个简单的 Go 代码热编译工具。
会实时监控指定目录下的文件变化(重命名，删除，创建，添加)，并编译和运行程序。

## 命令行语法:

主要包含了 watch 和 init 两个子命令。

### init

```shell
gobuild init github.com/owner/mod
```

### watch

```shell
gobuild watch [options] [dependents]

options:
 -h    显示当前帮助信息；
 -v    显示 gobuild 和 Go 程序的版本信息；
 -r    是否搜索子目录，默认为 true；
 -i    是否显示被标记为 IGNORE 的日志内容，默认为 false，即不显示；
 -o    执行编译后的可执行文件名；
 -x    传递给编译程序的参数；
 -ext  需要监视的扩展名，默认值为"go"，区分大小写，会去掉每个扩展名的首尾空格。
       若需要监视所有类型文件，请使用 *，传递空值代表不监视任何文件；
 -main 指定需要编译的文件，默认为""。

dependents:
 指定其它依赖的目录，只能出现在命令的尾部。
```

#### 常见用法:

```shell
# 监视当前目录下的文件，若发生变化，则触发 go build -main="*.go"
gobuild

# 监视当前目录和 ~/Go/src/github.com/issue9/term 目录下的文件，
# 若发生变化，则触发 go build -main="main.go"
gobuild -main=main.go ~/Go/src/github.com/issue9/term
```

## 支持平台

平台支持依赖 [colors](https://github.com/issue9/term) 与 [fsnotify](https://github.com/fsnotify/fsnotify) 两个模块，
目前支持以下平台：windows, linux, macOS, BSD。

## 安装

macOS 和 linux 用户可以直接使用 brew 进行安装：

```shell
brew tap caixw/brew
brew install caixw/brew/gobuild
```

常用平台可以从 <https://github.com/caixw/gobuild/releases> 下载，并将二进制文件放入 `PATH` 即可。

如果不存在你当前平台的二进制，可以自己编译：

```shell
git clone https://github.com/caixw/gobuild.git
cd gobuild
./build.sh
```

## 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。

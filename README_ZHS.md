# sysx

<a href="https://github.com/krau/sysx/blob/main/README.md"> English </a> / **简体中文**

将命令作为 Systemd 服务运行

## 概述

sysx 是一个用 Go 编写的简单的命令行工具，简化了在 Linux 上将命令作为 systemd 服务运行的过程。它会自动生成任何你想要在后台运行的命令的 systemd 服务文件，并使用 systemctl 管理服务。

## 特性

- 自动生成 systemd 服务文件
- 自动重载 systemd 守护进程、启动服务并设为开机启动
- 通过 -n 标志自定义服务命名
- 极为简单的使用方法：只需在命令前加上 sysx

## 要求

- 支持 systemd 的 Linux 发行版
- 需要 root 权限（sudo）以写入系统目录和管理服务

## 安装

### 从预编译二进制文件安装

1. 在 [releases page](https://github.com/krau/sysx/releases) 下载适合你的系统的最新版本
2. 解压
3. 将 sysx 二进制文件移动到你的 PATH 中的一个目录（例如 /usr/local/bin）

```shell
tar -xzf sysx-*-linux-*.tar.gz
sudo mv sysx /usr/local/bin
```

### 通过 Go 安装

```shell
go install github.com/krau/sysx@latest
```

## 使用

只需在需要后台运行的命令前加上 sysx 即可将其作为 systemd 服务运行。例如，要运行一个 Python HTTP 服务器：

```shell
sudo sysx python -m http.server
```

这行命令会：

- 在 /etc/systemd/system/ 中生成一个 systemd 服务文件（默认情况下以命令的第一个单词命名）
- 重载 systemd 守护进程
- 将服务设为开机启动
- 立即启动服务

### 自定义服务名

使用 `-n` 标志可以手动指定服务名。例如：

```shell
sudo sysx -n mycustomservice python -m http.server
```

这会创建一个名为 "mycustomservice.service" 的服务文件，而不是默认的名字。

// SPDX-License-Identifier: MIT

// Package watch 监视文件变化并编译
package watch

import (
	"context"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/caixw/gobuild/internal/local"
	"github.com/caixw/gobuild/log"
)

// Watch 执行热编译服务
//
// 如果初始化参数有误，则反错误信息，如果是编译过程中出错，将直接将错误内容输出到 logs。
func Watch(ctx context.Context, logs chan<- *log.Log, opt *Options) error {
	if err := opt.sanitize(); err != nil {
		return err
	}

	b := opt.newBuilder(logs)

	env, err := local.GoVersion()
	if err != nil {
		return err
	}
	b.logf(log.Info, "当前环境参数如下：%s", env)

	b.logf(log.Info, "给程序传递了以下参数：%s", b.appArgs) // 输出提示信息

	switch { // 提示扩展名
	case len(b.exts) == 0: // 允许不监视任意文件，但输出一信息来警告
		b.logf(log.Warn, "将 ext 设置为空值，意味着不监视任何文件的改变！")
	case len(b.exts) > 0:
		b.logf(log.Info, "系统将监视以下类型的文件：%s", b.exts)
	}

	b.logf(log.Info, "输出文件为：%s", b.appName) // 提示 appName

	go b.build() // 第一次主动编译程序，后续的才是监视变化。

	return b.watch(ctx, opt.paths)
}

// 开始监视 paths 中指定的目录或文件
func (b *builder) watch(ctx context.Context, paths []string) error {
	watcher, err := b.initWatcher(paths)
	if err != nil {
		return err
	}
	defer watcher.Close()

	var buildTime time.Time
	for {
		select {
		case <-ctx.Done():
			b.kill()
			b.logf(log.Info, "用户取消")
			return nil
		case event := <-watcher.Events:
			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				b.logf(log.Ignore, "watcher.Events:忽略 %s 事件", event.String())
				continue
			}

			if b.isIgnore(event.Name) { // 不需要监视的扩展名
				b.logf(log.Ignore, "watcher.Events:忽略不被监视的文件：%s", event.String())
				continue
			}

			if time.Since(buildTime) <= b.watcherFreq {
				b.logf(log.Ignore, "watcher.Events:忽略短期内频繁修改的文件：%s", event.Name)
				continue
			}

			buildTime = time.Now()

			if event.Name == "go.mod" && b.autoTidy {
				b.logf(log.Info, "watcher.Events:%s 事件触发了 go mod tidy", event.String())
				go b.tidy()
			} else {
				b.logf(log.Info, "watcher.Events:%s 事件触发了编译", event.String())
				go b.build()
			}
		case err := <-watcher.Errors:
			b.logf(log.Warn, "watcher.Errors：%s", err.Error())
			return nil
		} // end select
	}
}

func (b *builder) initWatcher(paths []string) (*fsnotify.Watcher, error) {
	b.logf(log.Info, "初始化监视器...")

	// 初始化监视器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	paths = b.filterPaths(paths)

	b.logf(log.Ignore, "以下路径或是文件将被监视：")
	for _, path := range paths {
		b.log(log.Ignore, path) // 路径不翻译
	}

	for _, path := range paths {
		if err := watcher.Add(path); err != nil {
			watcher.Close()
			return nil, err
		}
	}

	return watcher, nil
}

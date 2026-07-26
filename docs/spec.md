# SPEC FOR SS
 
ss(so simple) 规约定义了该项目始终应当遵循的设计原则. ss 没有一个实际的含义, 不过主要的核心是围绕着让一切变得简单来设计的. ss 只是一个 modules 用于 import 到其他项目当中使用, ss 也许后续会带一定的 cli bin 功能, 取决于实际的必要性.

每个规约代表一个主题紧接着是对应的内容, 足够清晰的情况下没有多余的东西需要解释和记录.

> 包

```go
///! 包: 命名直观即可, 通常属于一个独立的 domain, mod 内唯一.

// package.config - 配置相关
package config
// package.cmd - 指令相关
package cmd
// package.err - 错误相关
package err
// package.server - 服务器处理相关
package server
// package.log - 日志相关
package log
// package.db - 数据库相关
package db
// package.service - 服务相关
package service
```

> 接口

```go
///! 接口: [<Prefix>]<Action>[<Object>]([args...]) [returns] {}
///! 通常情况下 Prefix 一般省略, 以下主要集中在 Action 和 Object 的规约.
///! 借口规约只对对应的语义处理需要遵循, 自定义 Action 完全可以, 如果一个 Actioin 能够具备较好的语义, 需要记录到这里.

// Action.Init - 三方对象或者包内对象需要创建一次并完成实际的初始化.
func Init_()
// Action.New - 包内对象需要多次独立创建, 不必在创建过程中完全初始化.
func New_()
// Action.Build - 可复用的对象构建过程
func Build_()
// Action.Make - 默认方式构建, 只包含必要设置的对象字段
func Make_()
// Action.Set - 设置一个对象的属性, 不包含复杂逻辑
func Set_()
// Action.Add - 增加一个对象, 包含复杂逻辑
func Add_()
// Action.Create - 创建一个对象, 包含复杂逻辑
func Create_()
// Action.Reset - 重置对象的字段, 无删除含义
func Reset_()
// Action.Register - 注册一个实例, 一般用于三方接管的注册或者map
func Register_()
// Action.Append - 追加一个实例, 一般用于 slice
func Append_()
// Action.Mount - 挂载一个对象部分的关联性操作, 包含复杂逻辑
func Mount_()
// Action.Get - 获取一个对象的属性
func Get_()
// Action.Use - 使用一个对象, 自动创建并初始化成可用状态
func Use_()
// Action.Default - GetDefault 的压缩 or 默认值返回器
func Default_()
// Action.Has - 对象有指定的属性字段
func Has_()
// Action.Check - 检查逻辑
func Check_()
// Action.Del - 移除
func Del_()
// Prefix.Must - 必须, 包含内部 panic
func Must_()
// Prefix.G - 全局状态
func G_()
// Prefix.L - 本地状态
func L_()
// Prefix.To - 转化
func To_()
// Prefix.Is - 断言
func Is_()
```

> 测试

```go
// Action.Test - 测试, golang 中规定
func Test_()
```

> 领域

**err**

```go
///! err 规约包含: 
///! 1. errno, code, msg 等字段规范
///! 2. 封装原理和err追溯链
///! 3. 实践
```

**config**

```go
///! config 规约包含:
///! 1. flag, env, config key 规范
///! 2. 优先级定义和遵循原理
///! 3. spec配置查找位置
///! 4. 实践
```

**server**

```go
///! server 规约包含:
///! 1. rest, rpc 等启动式服务规范
///! 2. runner 设计
///! 3. 实践
```

**server.middleware**

```go
///! server.middleware 规约包含:
///! 1. rest, rpc 等服务中间件编写规范
///! 2. 相关中间件逻辑记录
///! 3. 实践
```

**cmd**

```go
///! cmd 规约包含:
///! 1. root entry 入口指令
///! 2. 指令挂载原理
///! 3. 实践
```

**log**

```go
///! log 规约包含:
///! 1. 日志封装原理
///! 2. 实践

/// 1. 日志封装原理
/// Init(logCtx).UseXXX() (xxx.Logger, error)
/// 核心: 保留原始日志库的使用, 统一部分配置, 实际创建初始化依赖 UseXXX() 方法

/// 2. 实践
/// 如下为一个使用示例
const (
    lv_info int iota = 1
    lv_warn
    lv_error
    lv_panic
    lv_fatal
)

type MutateLogger struct {
    lcontext *LogContext
}

func (ml *MutateLogger) UseZap() *zap.Logger {
    logger, _ := zap.NewExample()
    defer logger.Sync()
    return logger
}

type LogContext struct {
    level int
}

type LogOption func(*LogContext)

func WithLevel(lv int) LogOption {
    return func(ctx *LogContext) {
        ctx.level = lv
    }
}

func Init(opts ...LogOption) MutateLogger {
    mLogger := MutateLogger{
        lcontext: new(LogContext)
    }
    for _,opt := range opts {
        opt(mLogger.lcontext)
    }
    return mLogger
}

// main.go
func main() {
    logger := Init(WithLevel(lv_warn)).UseZap()
    defer logger.Sync()
    logger.Info("this is example for zap logger usage.")
}

```

**db**

```go
///! db 规约包含:
///!
```

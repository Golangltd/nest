# nest公共类库使用说明

## 简介
1. 本项目含有多个子模块：tools下是与业务无关的基本模块；其他文件夹下是和业务相关在各项目之间通用的代码;
2. leaf文件夹里面是对leaf框架的补充代码；
3. 在对本项目进行设计的时候，请尽量不要添加太多和业务耦合较深的代码，如无必要勿增实体；
4. 由于使用go mod，仅支持使用最新的golang版本进行构建(>=1.11)；
5. 必要的时候对具体模块进行再封装；
6. 一般情况下，**禁止对模块的导出接口进行删减改名**，这会导致库不再向下兼容；
7. 如果你能手动把其他使用该库的项目都改过来，也可以无视第6条；


## 使用
1. 如果不需要更新本项目代码，不必clone本项目到本地；
2. 由于本项目未存放在github之上，属于私有repo，因此你必须手动设置对应gitlab的地址，在bash中输入

```bash
git config --global url."yourgit:".insteadOf "https://lol.com/"
```

**请将上文中的git地址替换成本地的git仓库**，作为开发人员，请确保你有该项目的访问权限（只读权限即可）
3. 如果新建项目需要依赖本项目，请在`go.mod`里面`require lol.com/server/nest.git latest`，然后使用`go mod download`即可将其下载到本地；
4. 如果`go.mod`中的依赖版本与本地不符合，可以使用`go mod tidy`进行更新，其他`go mod`命令请自行学习；
5. NOTE: 如果本地go语言环境不满足需求，可以手动git clone本项目到`$GOPATH/lol.com/server`文件夹下，但是以后就需要手动git pull来进行更新了；

## 项目结构介绍

```
├── cache             # 公用cache，包含rpc/闪告/金币锁
├── ginutils          # gin的一些工具函数，供短链接API使用
├── heartbeat         # 通用心跳模块
├── leaf              # leaf的一些扩展函数，包括自定义的protobuf消息解析器
├── log               # 通用的log封装
├── proto             # 通用的protobuf消息
├── mg                # mongo中数据统计相关的公用接口
├── tools
│   ├── collection    # 集合类的扩展方法
│   ├── database      # 快速初始化数据库连接
│   ├── deepcopy      # 通用深拷贝（使用反射）
│   ├── fs            # 文件系统/配置解析
│   ├── ip            # ip地址库
│   ├── jsonutils     # json工具库
│   ├── mem           # 常用的内存缓存类
│   ├── num           # 基础数字类型工具函数
│   ├── sample        # 随机抽样函数
│   └── tz            # 时间函数
└── user              # 用户/流水相关公用接口
```


# Nav

一个适合个人使用的导航管理系统，包含：

1. 前台导航首页
2. 后台管理面板
3. SQLite 数据存储
4. Docker 一键部署能力

项目的目标很直接：把常用站点、工具、平台入口统一收拢，既能在前台快速查找和打开，也能在后台方便地维护菜单、应用、用户和站点设置。

## 项目特点

### 1. 前台导航页

1. 展示所有已启用的菜单和应用
2. 支持一级菜单和二级菜单分组展示
3. 支持按名称、描述、分组进行搜索
4. 支持快速目录导航和返回顶部
5. 支持登录用户在前台直接进入后台管理

### 2. 后台管理

1. 仪表盘总览
2. 应用管理
3. 菜单管理
4. 用户管理
5. 站点设置

### 3. 数据存储

1. 当前项目仅使用 `SQLite`
2. 数据文件路径为 `data/nav.db`
3. 适合个人使用、轻量部署和简单备份

## 项目结构

```text
nav/
├─ backend/                 # Go 后端
│  ├─ config/               # 配置加载
│  ├─ controllers/          # 接口控制器
│  ├─ database/             # 数据库初始化
│  ├─ middlewares/          # 中间件
│  ├─ models/               # 数据模型
│  ├─ routes/               # 路由
│  ├─ utils/                # 工具函数
│  ├─ data/                 # SQLite 数据文件目录
│  ├─ .env                  # 后端环境配置
│  ├─ Dockerfile            # 后端镜像构建文件
│  └─ main.go               # 后端入口
├─ frontend/                # Vue 前端
│  ├─ src/
│  ├─ .env                  # 前端环境变量
│  ├─ Dockerfile            # 前端镜像构建文件
│  └─ nginx.conf            # 前端 Nginx 配置
├─ docker-compose.yml       # Docker Compose 启动文件
├─ go.mod                   # Go 模块文件
└─ README.md
```

## 运行方式

## 一、本地开发运行

### 1. 后端启动

在项目根目录执行：

```bash
go run ./backend
```

后端默认监听：

```text
http://localhost:8080
```

### 2. 前端启动

在 `frontend` 目录执行：

```bash
npm install
npm run serve
```

前端默认访问：

```text
http://localhost:8081
```

具体端口以 Vue CLI 实际输出为准。

## 二、数据库说明

项目当前使用 SQLite。

后端环境变量文件：

[backend/.env]

当前数据库配置：

```env
DB_PATH=./data/nav.db
JWT_SECRET=你的密钥
```

备份时只需要备份这个文件。

## 三、默认账号

首次启动时如果数据库为空，系统会自动创建默认管理员账号：

```text
用户名：admin
密码：admin
```

首次登录后建议立即修改密码。


## Docker Compose 部署

在项目根目录执行：

```bash
docker compose up -d --build
```

启动后：

1. 前端访问地址：`http://localhost`
2. 后端接口地址：`http://localhost:8080`

## 适用场景

这个项目当前更适合：

1. 个人导航站
2. 单机部署
3. 轻量后台管理
4. 小规模内部自用

如果以后要做多实例部署、共享数据库、复杂权限体系，再考虑把 SQLite 升级为 MySQL 或 PostgreSQL 会更稳。

## 建议

1. 首次部署完成后先登录后台检查菜单和应用数据
2. 及时修改默认管理员密码
3. 定期备份 `data/nav.db`
4. 如果只做个人使用，当前这套 Docker + SQLite 已经足够轻量稳定

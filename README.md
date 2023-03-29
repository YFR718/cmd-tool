# cmd-tool

## 基本使用


```bash
# 下载cmd-tool
git clone git@github.com:YFR718/cmd-tool.git
# 修改配置文件
cd cmd-tool
vim config.yaml
# 安装cmd-tool
go install github.com/YFR718/cmd-tool@latest
# 部署网盘服务器
cd server/cloud-disk
go run main.go

# 使用cmd-tool
tool file
>ls 
>cd 
>mkdir xxx
>push xxx
>pull xxx
```



## 主要功能

- 基本工具
  - SQL 工具
    - 获取表的结构体
    - 表生成虚拟数据
  - string
    - json转结构体



- 网盘cli
  - 目录查看
  - 文件上传
  - 文件下载
  - 数据同步
- 压测cli
  - tcp 压测
  - http 压测
  - mysql压测
- 计算机性能测试
  - CPU
  - 内存
  - 磁盘
  - GPU
  - 网络
- 机器人
  - chatgpt 机器人




## 时间预计

基本工具
- 随时更新
- 网盘cli
  - 前置任务：网盘服务器搭建（2月下 7day）
- 压测cli
  - 压测服务端开发（2月下）

# distributed-computer-monitor
多主机资源监测（DCM）

## 项目介绍

DCM 一个简单的主机资源监测上报项目。项目分为服务端和采集端。

服务端接收采集端上报的数据，存储到执行目录同级 db 文件夹下 data.db n文件中。同级目录的 web 文件夹为前端页面。服务启动后访问 localhost:10020 就可以看到前端页面。

采集端每隔一分钟获取主机的 cpu、内存等信息并上报给服务端。

## 项目打包

项目采用 docker 进行打包处理。

### linux

```
cd ./server

docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang go build -v
```

### windows

在 linux 上交叉编译 windiows 平台对应的包。

```
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e GOOS=windows -e GOARCH=amd64 golang go build -v
```

## 部署

自行打包或者下载 release 文件。

### 服务端

在运行包之前，需要创建相关文件夹和文件。

* 配置文件：保留 config 文件夹就可以。或者自行创建相关文件夹和文件。

* 数据库：创建 db 文件夹就可以。如果配置文件的 databse-sqlite-path 字段修改了，请创建对应文件夹。

运行：

```
nohup ./xxx &
// 会出现下面的提示，直接回车
nohup: ignoring input and appending output to 'nohup.out'
```

或者用 screen：

```
screen -S DCM
./xxx
// Ctrl-a d 退出
```

### 采集端

在运行包之前，需要创建相关文件夹和文件。

* 配置文件：config.json
  ```
  {
    "hostname": "", // 采集端主机名称，作为主机标识，比如：阿里ECS
    "server": {
      "address": "" // 上报地址。比如：http://server.dcm.com:10020
    }
  }
  ```

运行：同服务端。


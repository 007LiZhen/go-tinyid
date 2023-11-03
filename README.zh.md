# go-tinyid

#### 介绍
本项目主要利用go语言(go1.20)实现了一种id生成器，并提供了http和grpc两种访问方式。项目中采用的生成算法主要基于数据库号段算法实现。关于这个算法可以参考
[美团Left](https://tech.meituan.com/2017/04/21/mt-leaf.html)。

#### 项目结构
    main.go    - 程序入口，项目初始化，并实现了平滑停服
    router     - 路由
    controller - api接口
    model      - 数据模型定义
    dao        - 数据表操作
    logic      - 逻辑操作
        grpcserver - grpc服务器
        idsequence - 实现了数据号段生成算法
    conf        - 数据库配置信息
    common      - 公共库
        config  - viper配置
        dto     - 请求响应/返回值结构体
        merrors - 错误码、错误信息定义
        mysql   - 数据库连接池
        xgrpc   - grpc server的proto定义

#### 使用说明
1. 项目采用go1.20编写，采用go mod进行包管理
2. 编译运行 go build && ./go-tinyid
3. 项目提供http和grpc两种访问方式，可自行选择

#### 核心流程说明
##### 1. 定义id生成器结构体
```
   type IdSequence struct {
      idListLength int64           // 号段长度，可根据业务qps自行设置
      biz          string          // 业务类型
      ids          chan int64      // 生成的id list, chan通道
      stopMonitor  chan bool       // 停止标志channel类型
   }
```

##### 2. id生成器共有Monitor，GetOne, Close三个对外暴露的方法。
       Monitor方法主要实现对id list的监控，当检测到id list为空时，会调用add方法，向id list中添加idListLength个新id，在添加新id过程中，
    会使用mysql 乐观锁，以防止其他进程也在更新获取到的最新id;
       GetOne方法主要会从id list里面获取一个新的id;
       Close方法主要是关闭channel，停止写入新的id;
##### 3. 数据表结构
```
create table if not exists test.sequence
(
    id          bigint unsigned auto_increment primary key,
    biz         varchar(128) default ''                not null comment '业务类型',
    value       bigint       default 0                 not null comment 'id值',
    version     bigint       default 0                 not null comment '乐观锁',
    is_del      tinyint      default 0                 not null comment '是否软删标志',
    create_time timestamp    default CURRENT_TIMESTAMP not null comment '创建时间',
    update_time timestamp    default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP,
    unique (version)
) charset = utf8mb4;
```

#### 参与贡献
欢迎大家积极提issue, 共建golang版本的tinyid

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request
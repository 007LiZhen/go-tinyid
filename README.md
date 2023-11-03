# go-tinyid

#### Introduction
[中文介绍](https://github.com/007LiZhen/go-tinyid/blob/master/README.md)

This project mainly utilizes the Go language (go1.20) to implement an ID generator, and provides two access methods: HTTP and GRPC.
The generation algorithm used in the project is mainly based on the database number segment algorithm. You can refer to this algorithm for reference
[MeiTuan Left](https://tech.meituan.com/2017/04/21/mt-leaf.html)。

#### Project Structure
    main.go    - Program entry, project initialization, and smooth server shutdown achieved
    router     - HTTP routing
    controller - api controller
    model      - data model definition
    dao        - database operations
    logic      - logic operation
        grpcserver - grpc server
        idsequence - implemented data number segment generation algorithm
    conf        - the configures of database
    common      - common lib
        config  - viper config
        dto     - request response/return data structure
        merrors - error number and error message
        mysql   - the connection pool of database
        xgrpc   - the proto defination of grpc server

#### Instructions For Use
1. The project is written using go1.20 and package management is carried out using go mod
2. build and run: go build && ./go-tinyid
3. The project provides two access methods: HTTP and GRPC, which can be chosen by oneself

#### Core Process Description
##### 1. Define the ID generator structure
```
   type IdSequence struct {
      idListLength int64           // 号段长度，可根据业务qps自行设置
      biz          string          // 业务类型
      ids          chan int64      // 生成的id list, chan通道
      stopMonitor  chan bool       // 停止标志channel类型
   }
```

##### 2. The ID generator has three exposed methods: Monitor(), GetOne(), and Close().
       The Monitor method mainly monitors the ID list. When it is detected that the ID list is empty, 
    the add method will be called to add a new ID of idListLength to the ID list. During the process 
    of adding a new ID, MySQL optimistic locks will be used to prevent other processes from updating 
    the latest IDs obtained;
       The GetOne method mainly obtains a new ID from the ID list;
       The Close method mainly closes the channel and stops writing new IDs.

##### 3. Data Table Structure
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

#### Participation Contribution
Welcome everyone to actively raise issues and jointly build the Golang version of tinyid

1.  Fork this project
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create new Pull Request

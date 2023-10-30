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
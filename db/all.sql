create database beyond character set utf8mb4 collate utf8mb4_general_ci;


CREATE TABLE `article` (
                           `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                           `title` varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
                           `content` text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
                           `cover` varchar(255) NOT NULL DEFAULT '' COMMENT '封面',
                           `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
                           `author_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '作者ID',
                           `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:待审核 1:审核不通过 2:可见 3:用户删除',
                           `comment_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论数',
                           `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
                           `collect_num` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
                           `view_num` int(11) NOT NULL DEFAULT '0' COMMENT '浏览数',
                           `share_num` int(11) NOT NULL DEFAULT '0' COMMENT '分享数',
                           `tag_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '标签ID',
                           `publish_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发布时间',
                           `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                           `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                           PRIMARY KEY (`id`),
                           KEY `ix_author_id` (`author_id`),
                           KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='文章表';


insert into article(title, content, author_id, like_num, publish_time) values ('文章标题1', '文章内容1', 1, 1, '2023-11-25 17:01:01');
insert into article(title, content, author_id, like_num, publish_time) values ('文章标题2', '文章内容2', 1, 10, '2023-11-25 15:01:01');

CREATE TABLE `follow` (
                          `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                          `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
                          `followed_user_id` bigint(20) UNSIGNED NOT NULL COMMENT '被关注用户ID',
                          `follow_status` tinyint(1) UNSIGNED NOT NULL DEFAULT '1' COMMENT '关注状态：1-关注，2-取消关注',
                          `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `uk_user_id_followed_user_id` (`user_id`,`followed_user_id`),
                          KEY `ix_followed_user_id` (`followed_user_id`),
                          KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '关注表';

CREATE TABLE `follow_count` (
                                `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                `user_id` bigint(20) UNSIGNED NOT NULL COMMENT '用户ID',
                                `follow_count` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '关注数',
                                `fans_count` int(10) UNSIGNED NOT NULL DEFAULT '0' COMMENT '粉丝数',
                                `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                                PRIMARY KEY (`id`),
                                UNIQUE KEY `uk_user_id` (`user_id`),
                                KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT '关注计数表';


CREATE TABLE `like_record` (
                               `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                               `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
                               `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '点赞对象id',
                               `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户ID',
                               `like_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型 0:点赞 1:点踩',
                               `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                               PRIMARY KEY (`id`),
                               KEY `ix_update_time` (`update_time`),
                               UNIQUE KEY `uk_biz_obj_uid` (`biz_id`,`obj_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='点赞记录表';

CREATE TABLE `like_count` (
                              `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                              `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
                              `obj_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '点赞对象id',
                              `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
                              `dislike_num` int(11) NOT NULL DEFAULT '0' COMMENT '点踩数',
                              `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                              `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                              PRIMARY KEY (`id`),
                              KEY `ix_update_time` (`update_time`),
                              UNIQUE KEY `uk_biz_obj` (`biz_id`,`obj_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='点赞计数表';


CREATE TABLE `reply_count` (
                               `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                               `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
                               `target_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '评论目标id',
                               `reply_num` int(11) NOT NULL DEFAULT '0' COMMENT '评论总数',
                               `reply_root_num` int(11) NOT NULL DEFAULT '0' COMMENT '根评论总数',
                               `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                               `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                               PRIMARY KEY (`id`),
                               KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='评论计数';

CREATE TABLE `reply` (
                         `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                         `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
                         `target_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '评论目标id',
                         `reply_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '评论用户ID',
                         `be_reply_user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '被回复用户ID',
                         `parent_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '父评论ID',
                         `content` text COLLATE utf8_unicode_ci NOT NULL COMMENT '内容',
                         `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态 0:正常 1:删除',
                         `like_num` int(11) NOT NULL DEFAULT '0' COMMENT '点赞数',
                         `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                         `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                         PRIMARY KEY (`id`),
                         KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='评论表';


CREATE TABLE `tag` (
                       `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                       `tag_name` varchar(32) NOT NULL DEFAULT '' COMMENT '标签名',
                       `tag_desc` varchar(128) NOT NULL DEFAULT '' COMMENT '标签描述',
                       `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                       `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                       PRIMARY KEY (`id`),
                       KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='标签表';

CREATE TABLE `tag_resource` (
                                `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
                                `biz_id` varchar(64) NOT NULL DEFAULT '' COMMENT '业务ID',
                                `target_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '内容id',
                                `tag_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '标签id',
                                `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT '用户ID',
                                `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                                `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
                                PRIMARY KEY (`id`),
                                KEY `ix_update_time` (`update_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='标签资源表';


CREATE TABLE `user`
(
    `id`          bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `username`    varchar(32)  NOT NULL DEFAULT '' COMMENT '用户名',
    `avatar`      varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
    `mobile`      varchar(128) NOT NULL DEFAULT '' COMMENT '手机号',
    `create_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
    PRIMARY KEY (`id`),
    KEY           `ix_update_time` (`update_time`),
    UNIQUE KEY `uk_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';

insert into user(username, avatar, mobile)
values ('张三', 'https://beyond-blog.oss-cn-beijing.aliyuncs.com/avatar/2021/01/01/1609488000.jpg', '13800138000');
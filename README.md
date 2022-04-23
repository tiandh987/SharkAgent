# SharkAgent

2022/04/23
- apiserver 用户相关接口
1. 创建用户
   
2. 删除用户
3. 批量删除用户
   
4. 修改用户密码
5. 修改用户属性
   
6. 查询用户信息
7. 查询用户列表

用户数据存储在数据库，所以需要安装 mariadb；
安装 mariadb 后需要 创建数据库，创建 user 表, 插入一条 admin 记录。
sql文件放置在 configs/iam.sql
```text
+--------------+---------------------+------+-----+---------------------+-------------------------------+
| Field        | Type                | Null | Key | Default             | Extra                         |
+--------------+---------------------+------+-----+---------------------+-------------------------------+
| id           | bigint(20) unsigned | NO   | PRI | NULL                | auto_increment                |   用户id
| instanceID   | varchar(32)         | YES  | UNI | NULL                |                               |   
| name         | varchar(45)         | NO   | UNI | NULL                |                               |   用户名
| status       | int(1)              | YES  |     | 1                   |                               |   用户状态
| nickname     | varchar(30)         | NO   |     | NULL                |                               |   昵称
| password     | varchar(255)        | NO   |     | NULL                |                               |   密码
| email        | varchar(256)        | NO   |     | NULL                |                               |   邮箱
| phone        | varchar(20)         | YES  |     | NULL                |                               |   手机号
| isAdmin      | tinyint(1) unsigned | NO   |     | 0                   |                               |   是否为 admin
| extendShadow | longtext            | YES  |     | NULL                |                               |   
| loginedAt    | timestamp           | YES  |     | NULL                |                               |   登录时间
| createdAt    | timestamp           | NO   |     | current_timestamp() |                               |   创建时间
| updatedAt    | timestamp           | NO   |     | current_timestamp() | on update current_timestamp() |   更新时间
+--------------+---------------------+------+-----+---------------------+-------------------------------+
```

apiserver 对外提供接口，通过使用 gin 框架提供 HTTP 服务
为了安全需要使用 x509 证书对通信进行加密和认证。

安装目录：/opt/iam
配置目录：/etc/iam
日志目录：/var/log/iam
数据目录：/data/iam
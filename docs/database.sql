CREATE TABLE pic_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nickname VARCHAR(255) NULL COMMENT '昵称',
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    user_type tinyint(1) unsigned not null default 0 comment '用户类型: 0,普通用户; 1,超级管理员',
    status  tinyint(1) unsigned not null default 0 comment '是否激活: 0,未激活; 1,已激活',
    deleted_at TIMESTAMP NULL COMMENT '软删除标记',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4 
COMMENT='图床用户表';

CREATE TABLE pic_repositories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL COMMENT '仓库所属用户',
    repo_name VARCHAR(255) NOT NULL COMMENT '仓库名称',
    repo_url VARCHAR(255) NOT NULL COMMENT '仓库链接',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX `idx_user_repo` (`user_id`, `repo_url`)
)
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4 
COMMENT='存储仓库表';


CREATE TABLE pic_files (
    id INT AUTO_INCREMENT PRIMARY KEY,
    repo_id INT NOT NULL COMMENT '仓库ID',
    user_id INT NOT NULL COMMENT '仓库和用户',
    filename VARCHAR(200) NOT NULL COMMENT '文件名称',
    url VARCHAR(255) NOT NULL COMMENT '文件路径',
    hash_value VARCHAR(200) NULL COMMENT '文件散列值',
    raw_filename VARCHAR(250) NULL COMMENT '文件上传时的原始名称',
    filesize INT UNSIGNED NULL COMMENT '文件大小，单位bit	',
    width INT UNSIGNED NULL comment '图片宽度',
    height INT UNSIGNED NULL comment '图片高度',
    mime varchar(50) NULL COMMENT '文件类型',
    filetype tinyint(1) UNSIGNED not null default 0 comment '文件类型: 0,未知;1,图片;2,视频;3,音频;4,文本;5,其他',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    index `idx_user_id` (`user_id`),
    index `idx_repo_id` (`repo_id`),
    unique index `url` (`url`)
)
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4 
COMMENT='图床文件表';

CREATE TABLE pic_backup_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    repo_id INT NOT NULL COMMENT '备份仓库ID',
    backup_path VARCHAR(255) NOT NULL,
    backup_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4 
COMMENT='图床数据备份记录表';

CREATE TABLE `pic_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` INT NOT NULL default 0 comment '用户ID，默认0为系统配置',
  `type` varchar(24) DEFAULT NULL COMMENT '类型',
  `name` varchar(32) NOT NULL COMMENT '名称',
  `value` varchar(500) COMMENT '值',
  `remark` varchar(50) COMMENT '此值标记',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  unique idx_user_key (user_id, type, name)
) 
ENGINE=InnoDB 
DEFAULT CHARSET=utf8mb4
COMMENT='配置表';


INSERT INTO pic_config
(`user_id`, `type`, `name`, `value`, `remark`)
values 
(0, 'pichub', 'github_token', 'xxxxxx', `github授权token`),
(0, 'pichub', 'cdn_domain', 'https://pic.mysticalpower.uk', `cdn加速域名`)
;


CREATE TABLE pic_scheduled_tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,                 -- 任务名称
    description TEXT,                           -- 任务描述
    task_type ENUM('bash_script', 'db_backup', 'url_access', 'dir_backup') NOT NULL, -- 任务类型
    cron_expression VARCHAR(255) NOT NULL,      -- CRON 表达式
    timezone VARCHAR(100) DEFAULT 'UTC',        -- 时区信息
    command TEXT,                               -- Bash 脚本路径或命令
    db_connection_info TEXT,                    -- 数据库连接信息
    url VARCHAR(2048),                          -- 访问的 URL
    dir_path TEXT,                              -- 备份目录的路径
    status ENUM('pending', 'running', 'success', 'failed') DEFAULT 'pending',  -- 任务状态
    last_run_at DATETIME DEFAULT NULL,          -- 上次运行时间
    next_run_at DATETIME DEFAULT NULL,          -- 下次运行时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- 更新时间
    is_active BOOLEAN DEFAULT TRUE              -- 任务是否激活
);

CREATE TABLE pic_task_run_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id BIGINT UNSIGNED NOT NULL,           -- 关联到 scheduled_tasks 表
    run_at DATETIME NOT NULL,                   -- 任务执行时间
    status ENUM('success', 'failed') NOT NULL,  -- 执行状态
    result TEXT,                                -- 执行结果或输出
    error_message TEXT,                         -- 错误信息（如果有）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 日志记录创建时间
    FOREIGN KEY (task_id) REFERENCES scheduled_tasks(id) ON DELETE CASCADE
);

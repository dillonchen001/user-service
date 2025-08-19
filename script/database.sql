-- 创建数据库
CREATE DATABASE IF NOT EXISTS auth_db;
USE auth_db;

-- 用户表
CREATE TABLE users (
  id bigint AUTO_INCREMENT PRIMARY KEY comment '自增id',
  user_id bigint not null default 0 UNIQUE comment '用户id',
  name VARCHAR(100) not null default '' comment '昵称',
  email VARCHAR(255) not null default '' UNIQUE comment '邮箱',
  phone VARCHAR(20) not null default '' UNIQUE comment '手机号',
  avatar VARCHAR(255) not null default '' comment '图像',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP comment '创建时间',
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci comment '用户表';

-- 认证表 (存储不同登录方式的关联)
CREATE TABLE auth_providers (
  id bigint AUTO_INCREMENT PRIMARY KEY comment '自增id',
  user_id bigint not null default 0 comment '用户id',
  provider_type ENUM('phone', 'facebook', 'apple', 'google', 'snapchat') NOT NULL  comment '类型',
  provider_id VARCHAR(255) NOT NULL default '' comment '登陆id',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP comment '创建时间',
  UNIQUE KEY unique_provider (provider_type, provider_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci comment '授权登陆用户表';

-- 验证码表
CREATE TABLE verification_codes (
  id bigint AUTO_INCREMENT PRIMARY KEY comment '自增id',
  phone VARCHAR(20) not null default '' comment '手机号',
  code VARCHAR(10) NOT NULL default '' comment '验证码',
  expires_at TIMESTAMP NOT NULL default CURRENT_TIMESTAMP comment '失效时间',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP comment '生成时间',
  index phone(phone)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci comment '验证码表';

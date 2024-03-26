-- 01, table for storing policy rules
CREATE TABLE IF NOT EXISTS auth_rules (
  id INTEGER,                               -- ID
  ptype VARCHAR(32) NOT NULL DEFAULT '',    -- 策略类型，例如 p 表示权限检查，g 表示角色继承
  v0 VARCHAR(255) NOT NULL DEFAULT '',      -- 主体，例如用户或角色
  v1 VARCHAR(255) NOT NULL DEFAULT '',      -- 资源
  v2 VARCHAR(255) NOT NULL DEFAULT '',      -- 操作
  v3 VARCHAR(255) NOT NULL DEFAULT '',      -- 条件
  v4 VARCHAR(255) NOT NULL DEFAULT '',      -- 其他参数
  v5 VARCHAR(255) NOT NULL DEFAULT '',      -- 其他参数
  PRIMARY KEY (id)
);
-- -- 允许用户 "admin" 访问所有资源
-- INSERT INTO auth_rules (ptype, v0, v1, v2) VALUES ('p', 'admin', '*', '*');
-- -- 允许角色 "manager" 访问所有资源
-- INSERT INTO auth_rules (ptype, v0, v1, v2) VALUES ('g', 'manager', '*', '*');
-- -- 允许角色 "user" 访问 "read" 操作
-- INSERT INTO auth_rules (ptype, v0, v1, v2) VALUES ('p', 'user', '*', 'read');

-- 02, 用户表
CREATE TABLE IF NOT EXISTS auth_users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username varchar(64) UNIQUE NOT NULL,
  password varchar(128) NOT NULL,
  email varchar(128) UNIQUE NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 03, 角色表
CREATE TABLE IF NOT EXISTS auth_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name varchar(64) UNIQUE NOT NULL,
  description varchar(128),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 04, 用户角色关系表
CREATE TABLE IF NOT EXISTS auth_user_roles (
  user_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  PRIMARY KEY (user_id, role_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (role_id) REFERENCES roles(id)
);

-- -- -- 05, 权限表
-- -- CREATE TABLE IF NOT EXISTS permissions (
-- --   id INTEGER PRIMARY KEY AUTOINCREMENT,
-- --   name varchar(128) UNIQUE NOT NULL,
-- --   description varchar(128),
-- --   created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
-- --   updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- -- );

-- -- -- 角色权限关系表
-- -- CREATE TABLE IF NOT EXISTS role_permissions (
-- --   role_id INTEGER NOT NULL,
-- --   permission_id INTEGER NOT NULL,
-- --   PRIMARY KEY (role_id, permission_id),
-- --   FOREIGN KEY (role_id) REFERENCES roles(id),
-- --   FOREIGN KEY (permission_id) REFERENCES permissions(id)
-- -- );

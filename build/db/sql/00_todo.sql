CREATE DATABASE IF NOT EXISTS todo;
CREATE TABLE IF NOT EXISTS todo.todo (
  id        BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title     VARCHAR(128) NOT NULL
);
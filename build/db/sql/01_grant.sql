CREATE USER IF NOT EXISTS 'app'@'%' IDENTIFIED BY 'app';
GRANT SELECT,INSERT,UPDATE,DELETE ON todo.todo TO 'app'@'%';
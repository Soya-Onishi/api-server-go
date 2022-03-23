INSERT INTO todo.todo_list (title) VALUES ('prepare hot water');
INSERT INTO todo.todo_list (title) VALUES ('wait for three minutes');
INSERT INTO todo.todo_list (title) VALUES ('eat ramen');

INSERT INTO auth.users(username, passwd, session_hash) VALUES (
  'Taro',
  SHA2('Taro', 256),
  NULL
);
INSERT INTO auth.users(username, passwd, session_hash) VALUES (
  'Hanako',
  SHA2('Hanako', 256),
  NULL
);
INSERT INTO auth.users(username, passwd, session_hash) VALUES (
  'Ryota',
  SHA2('Ryota', 256),
  NULL
);

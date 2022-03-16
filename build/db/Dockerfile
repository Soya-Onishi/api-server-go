FROM mysql:8.0-debian

ENV MYSQL_DATABASE todo

COPY custom.cnf /etc/mysql/conf.d/
COPY sql /etc/docker-entrypoint-initdb.d
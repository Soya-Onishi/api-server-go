FROM mysql:8.0-debian

ENV MYSQL_DATABASE todo

COPY ./custom.cnf /etc/mysql/conf.d/
COPY ./sql /docker-entrypoint-initdb.d

EXPOSE 3306:3306
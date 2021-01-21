FROM postgres:latest

WORKDIR /usr/app
COPY init.sql /usr/app/init.sql
RUN psql -U postgres -f init.sql
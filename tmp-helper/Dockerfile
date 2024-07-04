FROM postgres:latest

COPY dbhandler/script-db /docker-entrypoint-initdb.d/

RUN chmod +x /docker-entrypoint-initdb.d/script-db
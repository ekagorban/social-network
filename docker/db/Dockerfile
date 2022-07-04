FROM mysql:8.0

COPY ./* /docker-entrypoint-initdb.d/
RUN chmod -R 777 /docker-entrypoint-initdb.d/
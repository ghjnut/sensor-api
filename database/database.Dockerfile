FROM postgres:latest
ENV POSTGRES_PASSWORD=pgpassword
ENV POSTGRES_USER=pguser
EXPOSE 5432
COPY ./database/create_table.sql /docker-entrypoint-initdb.d/create_table.sql

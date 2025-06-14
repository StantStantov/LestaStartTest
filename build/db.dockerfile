FROM postgres:latest AS base

COPY ./build/sql/schema.sql /docker-entrypoint-initdb.d/001-schema.sql

FROM base AS build
RUN echo "exit 0" > /docker-entrypoint-initdb.d/100-exit_before_boot.sh
ENV PGDATA=/pgdata
RUN --mount=type=secret,id=ps_user,env=POSTGRES_USER \
    --mount=type=secret,id=ps_password,env=POSTGRES_PASSWORD \
    --mount=type=secret,id=ps_db,env=POSTGRES_DB \
    docker-entrypoint.sh postgres

FROM postgres:latest AS runner
ENV PGDATA=/pgdata
COPY --chown=postgres:postgres --from=build /pgdata /pgdata

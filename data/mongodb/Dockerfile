FROM mongo:latest
LABEL maintainer="Tatiana @mtm Mironova"

COPY ./init/ /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/mongo-init.sh

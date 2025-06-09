FROM postgres:15

# Set environment variables for PostgreSQL
ENV POSTGRES_DB=socialnetwork
ENV POSTGRES_USER=admin
ENV POSTGRES_PASSWORD=adminpassword

# Expose PostgreSQL port
EXPOSE 5432

# Optional: Add initialization scripts
# COPY init.sql /docker-entrypoint-initdb.d/

# The postgres image automatically runs initdb and starts the server
# Any .sql files in /docker-entrypoint-initdb.d/ will be executed on first run

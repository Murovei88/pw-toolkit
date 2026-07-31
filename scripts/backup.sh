#!/bin/bash
# Бэкап MariaDB
docker exec $(docker ps -q -f name=toolkit_mariadb) mysqldump -u root -p$MARIADB_ROOT_PASSWORD pw_toolkit > /opt/backup/mariadb_$(date +%Y%m%d).sql
# Бэкап MinIO (можно использовать mc)
echo "Backup completed"

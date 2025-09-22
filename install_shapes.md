
# Команда для загрузки шейп-файлов на linux(ubuntu, WSL)

```bash
  shp2pgsql .data/admin_4.cpg public.district_shapes | PGPASSWORD={DATABASE PASSWORD} psql -h {DATABASE_HOST} -p {DATABASE_PORT} -d {DATABASE_NAME} -U {POSTGRES_USER}
```

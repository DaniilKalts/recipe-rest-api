<h1 align="center">
  <img width="200" src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/b410dc17-55f9-440e-a00a-12e0e27e8b33" align="center" />
  Recipe-Rest-API
   <img width="230" src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/128b2020-05e4-43ab-abae-8669cd4ba643" align="center" />
</h1>

## ⚙️ Tech Stack:

- Go
- Docker
- Postgres
- pgAdmin

## ❗❗ Prerequisites:

- Download <a href="https://www.docker.com/products/docker-desktop/">Docker Desktop</a> to run the project.
- Download <a href="https://www.postman.com/downloads/">Postman</a> to play around with API.

## 🧾 Instruction:

1. Clone the repository.

```
git clone https://github.com/DaniilKalts/recipe-rest-api
```

2. Create .env (environment variables) file and fill it with your values.

```
DB_USER=db_user
DB_PASSWORD=db_password
DB_NAME=db_name
DB_PORT=db_port
DB_HOST=db_host

PGADMIN_DEFAULT_EMAIL=pgdmin_default_email
PGADMIN_DEFAULT_PASSWORD=pgadmin_deffault_password
```

3. Run docker compose up to create and start containers.

```
docker-compose up
```

4. Open localhost:5050 and paste values from the .env file to log in in pgAdmin.
   <img src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/633d91dd-dbbe-40f6-ac4c-0437eac464be" />

5. Create a new server in pgAdmin.
   <img src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/63b7e0d3-0a21-4117-9c06-ff1757de3e5c" />

6. Go to CMD (Command Prompt) and execute these queries to intialize and insert values into Recipes Table.

```
docker exec -it postgres-db sh
psql -U db_user -d db_name -f /docker-entrypoint-initdb.d/create_recipe_table.sql
psql -U db_user -d db_name -f /docker-entrypoint-initdb.d/insert_recipe_table.sql
```

7. Run the following query (highlighted with red pencil) in pgAdmin to make sure, the table is created and populated.
   <img src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/5337585b-0cff-4fc9-b58d-a3639888ad8c" />

8. Open Postman (or any other API platform) and send CRUD requests.
   <img src="https://github.com/DaniilKalts/recipe-rest-api/assets/109500182/201bfa36-3280-4159-a07d-2d3d05cf97e8" />

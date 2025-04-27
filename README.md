# Lab 9

La explicación de la arquitectura esta dividida en 2 videos:

1. [PARTE 1](https://youtu.be/HUII7wCScWo)
2. [PARTE 2]()


# Como correr

Para generar data:

```bash
cd dataGen/
go run main.go
```

Para subir los datos

```bash
cd dataUpload/
python main.go
```

Se necesitará tener un archivo .env con la estructura 

```
DB_URI=<URI>
DB_NAME=<DB_NAME>
```
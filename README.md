# Lab 9

La explicación de la arquitectura esta dividida en 2 videos:

1. [PARTE 1](https://youtu.be/HUII7wCScWo)
2. [PARTE 2](https://www.youtube.com/watch?v=Q4UG_1kyCfI)

## Dashboard
![image](https://github.com/user-attachments/assets/9f70ab2c-6022-4dc7-8d2e-bd72d6394499)
![image](https://github.com/user-attachments/assets/2bc3cfe2-41c8-4068-8352-214b5926b4eb)
![image](https://github.com/user-attachments/assets/9eb17e79-2963-4325-b415-8673bd371ebf)
![image](https://github.com/user-attachments/assets/bfc776b9-3c07-4733-ba73-fbee7d88a84b)
![image](https://github.com/user-attachments/assets/9cb3de17-3f57-473d-9ecd-394887c6d7b7)


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

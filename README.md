# Trips Service (Servicio de Viajes y Medios)

## Descripción
El servicio de viajes gestiona la creación y organización de viajes, permitiendo asociar fotos y videos directamente a cada viaje. Ofrece control granular sobre la privacidad y los metadatos de cada medio almacenado.

## Características

- Creación y gestión integral de viajes.
- Almacenamiento y administración de fotos y videos asociados directamente a viajes.
- Configuración individual de privacidad (público/privado) para cada medio.
- Almacenamiento de metadatos detallados (ubicación GPS, fecha, etiquetas).
- Integración con MinIO para almacenamiento eficiente de imágenes y videos.

## Tecnologías Utilizadas

- **Lenguaje:** Go
- **Base de Datos:** PostgreSQL
- **Almacenamiento:** MinIO
- **Cache:** Redis **NOT YET**
- **Orquestación:** Docker

## Instalación

Clona el repositorio:

```bash
git clone <repo-url>
cd trips-service
```

Configura las variables de entorno en un archivo `.env`.

Construye y ejecuta el servicio con Docker:

```bash
docker-compose up --build -d
```

## Endpoints

### Viajes

| Método | Ruta                    | Descripción                          |
|--------|-------------------------|--------------------------------------|
| POST   | `/trips`                | Crea un nuevo viaje                  |
| GET    | `/trips`                | Lista los viajes propios del usuario |
| GET    | `/trips/:trip_id`       | Obtiene detalles específicos de un viaje |
| PUT    | `/trips/:trip_id`       | Actualiza información del viaje      |
| DELETE | `/trips/:trip_id`       | Elimina un viaje                     |

### Medios

| Método | Ruta                                      | Descripción                             |
|--------|-------------------------------------------|-----------------------------------------|
| POST   | `/trips/:trip_id/media`                   | Añade una imagen/video al viaje         |
| GET    | `/trips/:trip_id/media`                   | Obtiene medios asociados a un viaje     |
| PUT    | `/trips/:trip_id/media/:media_id`         | Actualiza privacidad y metadatos de un medio |
| DELETE | `/trips/:trip_id/media/:media_id`         | Elimina un medio del viaje              |

## Seguridad

- **Autenticación:** Implementada mediante JWT.
- **Control de acceso:** Basado en permisos.
- **Privacidad:** Respeta privacidad granular por cada medio.


# CAS to OAuth2

[![es](https://img.shields.io/badge/lang-es-yellow.svg)](https://github.com/vti-uchile/cas-to-oauth2/blob/develop/README.es.md)
[![en](https://img.shields.io/badge/lang-en-red.svg)](https://github.com/vti-uchile/cas-to-oauth2/blob/develop/README.md)

## Descripción

"CAS a OAuth2" es un proyecto que implementa una interfaz de autenticación diseñada para traducir entre el protocolo CAS (Central Authentication Service v2.0) y OAuth 2.0 con OpenID. Este proyecto facilita la integración de sistemas que utilizan CAS con un servidor OAuth2, proporcionando una solución robusta y versátil para la gestión de autenticaciones y autorizaciones en entornos complejos.

## Tecnologías utilizadas

- **Go**: Lenguaje principal utilizado para el desarrollo del backend.
- **Gin**: Framework web de Go para manejar peticiones HTTP.
- **MongoDB**: Base de datos usada para almacenar la información relevante.

## Requisitos

- Instalar Go v1.21.1.
- Instalar y configurar MongoDB.

## Instalación y configuración

1. Clona el repositorio:

```bash
git clone https://github.com/vti-uchile/cas-to-oauth2.git
cd cas-to-oauth2
```

2. Instala las dependencias:

```bash
go mod tidy
```

3. Configura las variables de entorno:

- Copia o renombra el archivo `.env.local` a `.env`.
- Completa el archivo `.env` con los datos apropiados (conexión a MongoDB, URL del servicio OAuth2, nombre de cookie, etc.).

## Ejecución del proyecto

Una vez completada la configuración, ejecuta el proyecto con:

```bash
go run cmd/main.go
```

## Uso

El proyecto se utiliza como un servidor CAS estándar. Para una verificación rápida de que todo funciona correctamente, puedes ejecutar las pruebas unitarias disponibles en el directorio `tests`.

## Configuraciones adicionales

> Si la variable `USE_APM` en el archivo `.env` está establecida en `true`, también debes configurar las siguientes variables: `ELASTIC_APM_SERVICE_NAME`, `ELASTIC_APM_SERVER_URL`, `ELASTIC_APM_SECRET_TOKEN` y `ELASTIC_APM_ENVIRONMENT`.

#!/bin/bash
# Cárcel de los Gemelos - Deployment Script (Production)
echo "Iniciando despliegue de La Cárcel de los Gemelos..."
docker-compose up --build -d
echo "¡Despliegue completado! El Panóptico (Next.js) está en el puerto 3000, y el Motor Core (Go) en el 8080."

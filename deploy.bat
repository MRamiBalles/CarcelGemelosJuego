@echo off
echo Iniciando despliegue de La Carcel de los Gemelos en Windows...
docker-compose up --build -d
echo Despliegue completado.
echo Backend (Go Engine): http://localhost:8080
echo Frontend (Panoptico Next.js): http://localhost:3000
pause

# La Cárcel de los Gemelos (Interactive Reality Game)

**Creador y Propiedad Intelectual:** Manuel Ramírez Ballesteros  
**Contacto Comercial:** ramiballes96@gmail.com  
**Versión del Motor:** V2.1 (Authoritative Engine Pivot)

---

## 👁️ Visión del Producto
"La Cárcel de los Gemelos" no es solo un videojuego, es una experiencia híbrida pionera en la categoría de **Interactive Reality Games (IRG)**. Diseñado específicamente para la economía de los creadores de contenido (Streaming/Twitch), el proyecto entrelaza dinámicas de supervivencia psicológica hardcore (aislamiento, inanición, dilemas morales) con la intervención directa y monetizada de la audiencia en tiempo real. 

El servidor actúa como "Gran Hermano", un ente omnisciente impulsado por IA que juzga las acciones de los prisioneros, gestiona los votos de la audiencia (Sushi vs. Tortura) y mantiene un registro inmutable de cada traición.

## 📁 Documentación Oficial y Comercial
El proyecto está estructurado no solo como un hito técnico, sino como un producto comercializable y escalable:

- 📊 [**Plan de Negocio y Financiación**](docs/business_plan.md) (Monetización, Revenue Share, Inversión Seed).
- ⚖️ [**Marco Legal y Compliance**](docs/legal.md) (GDPR, IA EU Act, EULA, Limitación de Responsabilidad).
- 📜 [**La Constitución (Core Design)**](docs/constitution.md) (Filosofía de Game Design y el Dilema del Prisionero).
- 🛠️ [**Especificaciones Técnicas**](docs/spec.md) (Arquitectura y Mecánicas).

## ⚙️ Arquitectura Tecnológica
Construido bajo el paradigma de **Spec-Driven Development (SDD)** y **Clean Architecture**.

- **`/server` (Backend Autoridad - Go):** El corazón del proyecto. Un servidor autoritativo concurrente en Golang que usa **Event Sourcing** (`VAR Replay`) para persistir cada interacción en memoria y en disco (`SQLite`). Previene trucos y asegura que la lógica del negocio ($$$ y votos) jamás resida en el cliente.
- **`/client` (Panel de Control - Next.js):** El "Panóptico" para la administración y visualización. Conectado vía WebSockets (`ws://`) y API REST al núcleo en Go, permite monitorizar en vivo la cordura, lanzar audios de tortura u organizar encuestas de Twitch.
- **`Los Gemelos` (Capa de IA Cognitiva):** Un sistema LLM agnóstico (OpenAI/Anthropic) que audita transcripciones del juego, vigila alianzas, gestiona el "Módulo de Aislamiento" y sanciona el romper las reglas basadas en un estricto *System Prompt* de rol.

---
© Todos los derechos reservados. Manuel Ramírez Ballesteros.

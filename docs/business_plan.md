# Plan de Negocio y Financiación: "La Cárcel de los Gemelos"
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 1.0 - Confidencial

## 1. Resumen Ejecutivo
"La Cárcel de los Gemelos" es una experiencia interactiva asimétrica que fusiona el formato de *Reality Show* tradicional (tipo Gran Hermano) con un ecosistema de videojuegos multijugador masivo. El modelo de negocio se basa en la monetización directa de la audiencia a través de plataformas de streaming (Twitch/YouTube), donde los espectadores no solo consumen contenido, sino que financian activamente castigos, recompensas y dinámicas del juego en tiempo real. 

El producto busca posicionarse como el pionero en la categoría de "Juegos de Realidad Interactiva" (IRG), capitalizando la economía de la atención y el ecosistema de creadores de contenido.

## 2. Modelo de Asignación de Valor (Monetización)

El núcleo de ingresos ("The House Wins") está diseñado para ser altamente escalable y basarse en microtransacciones impulsadas por la emoción y el salseo (drama) en directo:

1. **Economía de la Audiencia (Twitch Integrations):**
   - **Bits/Donaciones Directas:** La audiencia paga dinero real (convertido a través de la API de Twitch) para desencadenar eventos en el juego: *Audio Torture*, cortes de luz, o envíos de comida premium (Sushi) a sus prisioneros favoritos.
   - **Suscripciones Prime/Tier:** Los suscriptores de pago del canal obtienen multiplicadores de voto en las encuestas de "Audience Polling" que deciden el destino diario de los prisioneros.
2. **El Bote (Prize Pool) Dinámico:**
   - Una parte del revenue generado por la audiencia se destina al gran premio final que se disputan los dúos el Día 21. Esto incentiva a los espectadores a donar más para engordar el bote de su creador favorito.
3. **Patrocinios In-Game (Product Placement Branded):**
   - Los suministros diarios (cajas de comida, agua, ropa de prisión) son espacios publicitarios dinámicos. Marcas reales pueden patrocinar los "Loot Events" o los desafíos del Patio.

## 3. Arquitectura Comercial y Operativa
Para operar este sistema a nivel comercial, la arquitectura técnica (Authoritative Server en Go + Next.js Frontend) garantiza:
- **Tolerancia a Fraude (Anti-Cheat):** Al ser un servidor autoritativo, es imposible que terceros manipulen las votaciones o los balances de dinero. Toda transacción financiera está protegida y encriptada, crucial para la confianza de los inversores.
- **Coste de Infraestructura (FinOps):** Implementación de *Cachés* y bases de datos ligeras (SQLite centralizado con write-through) que reducen el coste de computación en la nube (AWS/GCP) a un mínimo operativo, maximizando márgenes.

## 4. Estrategia de Financiación (Funding)

Para escalar el proyecto de un entorno de prueba a una producción masiva con grandes creadores de contenido, se persiguen las siguientes vías de capitalización:

### 4.1. Financiación Semilla (Seed-Stage)
*   **Business Angels (Sector Gaming/Media):** Búsqueda de un ticket inicial de 100k€ - 250k€ destinado estrictamente a escalar los costes de servidores concurrentes y licencias de APIs comerciales.
*   **Subvenciones de Innovación:** Aplicación a fondos europeos (NextGenerationEU) o ENISA, enmarcando el proyecto dentro del I+D+i en tecnologías de entretenimiento interactivo y gamificación social.

### 4.2. Acuerdos de "Revenue Share" (B2B2C)
En lugar de gastos directos en marketing de adquisición de usuarios (CAC), el modelo penetra en el mercado mediante la cesión de licencias temporales a grandes *Streamers* o Agencias de Representación (e.g., Vizz Agency). 
- **Estructura:** El servidor y el producto se ceden gratis a cambio de un porcentaje (20-30%) del volumen total de bits/donaciones generadas por la audiencia durante los 21 días de duración del evento.

## 5. Viabilidad y Escalabilidad Futura (Roadmap)
Una vez validada la Temporada 1, el motor del juego ("Cognition Engine" de Los Gemelos) es un SaaS (Software as a Service) empaquetable. Se podrán vender instancias privadas de "La Cárcel" a comunidades más pequeñas o a ligas de eSports como herramientas de retención de audiencia.

---
*Este documento constituye la hoja de ruta comercial de la propiedad intelectual de M. Ramírez Ballesteros. Prohibida su copia o distribución sin autorización expresa.*

# Marco Legal y Cumplimiento Normativo
**Autor:** Manuel Ramírez Ballesteros
**Contacto:** ramiballes96@gmail.com
**Versión:** 1.0 - Documento Oficial Confidencial

## 1. Introducción
El presente documento describe el marco legal estructural, los Términos de Servicio (ToS) base y las políticas de cumplimiento normativo vinculadas a la explotación comercial del software interactivo "La Cárcel de los Gemelos". 

Dado que la naturaleza del software implica la interacción masiva de usuarios en tiempo real, el procesamiento de datos biométricos simulados, transacciones económicas e interacciones sociales con Inteligencia Artificial, este marco es de estricto cumplimiento para operar en la Unión Europea e Internacionalmente.

## 2. Cumplimiento del RGPD (Reglamento General de Protección de Datos - UE)

El motor del juego ("Authoritative Server") y su frontend asumen la figura de **Responsable del Tratamiento** de cara a los jugadores (Streamers/Participantes), y **Encargado del Tratamiento** de cara a la audiencia integrada a través de APIs de terceros (Twitch/YouTube).

### 2.1. Minimización y Retención de Datos
- **Jugadores (Prisoners):** El almacenamiento en `SQLite/PostgreSQL` de estados psicológicos (`Sanity`, `Dignity`) y perfiles de personalidad está disociado del PII (Información Personal Identificable) real tras el cierre del ciclo de 21 días.
- **Audiencia (Voters):** No se almacenarán IPs, correos electrónicos ni tokens persistentes de la audiencia. Los votos del módulo "Audience Polling" usan tokens efímeros hasheados u OAuth anónimo de la plataforma de origen solo con fines de contabilización estadística de la encuesta, borrándose tras el evento `POLL_RESOLVED`.

### 2.2. Derecho al Olvido (Art. 17 RGPD)
Los participantes tienen el derecho de solicitar el borrado de sus perfiles in-game. El sistema ejecutará una ofuscación de la base de datos (Ej: `UPDATE prisoners SET name = 'Redacted', archetype = 'UNKNOWN' WHERE prisoner_id = ?`) para mantener la integridad relacional del `EventLog` (VAR) sin vulnerar la privacidad del usuario.

## 3. Inteligencia Artificial y Normativa (EU AI Act)
El componente autónomo "Los Gemelos" supervisa y castiga las acciones de los jugadores.
- **Transparencia (Art. 52):** Los jugadores serán explícitamente informados (mediante un EULA de aceptación obligatoria) de que están interactuando bajo un sistema algorítmico y de IA generativa.
- **Supervisión Humana (Human-in-the-Loop):** Para mitigar riesgos de clasificación algorítmica abusiva en la "Celda de Aislamiento" o censura de comunicación, el panel de administración (`AdminPanel.tsx`) conserva un mecanismo de anulación (Override) de las decisiones de la IA en todo momento.

## 4. Condiciones de Servicio (ToS) y Política de Conducta (EULA)

Para eximir a la propiedad intelectual y a sus creadores (M. Ramírez Ballesteros) de responsabilidad legal civil:
- **Cláusula de Daño Psicológico Simulado:** El participante reconoce que la pérdida de cordura ("Sanity Drain"), los castigos o traiciones ("Snitching") son mecánicas ficticias de rol de supervivencia. Renuncian expresamente a acciones legales por "angustia emocional".
- **Sistema Financiero y Donaciones:** Las donaciones de la audiencia para afectar al juego se considerarán contribuciones irrevocables de entretenimiento voluntario ("tips"). No garantizan la supervivencia del destinatario y no están sujetas a la ley de juego de azar ni reembolsos (Chargebacks policy strict restriction).

## 5. Propiedad Intelectual
Todos los derechos de explotación, el código base ("Cognition Engine"), esquemas de bases de datos, y los protocolos de sincronización de eventos ("VAR EventLog") son propiedad exclusiva de **Manuel Ramírez Ballesteros**. Cualquier acuerdo "B2B2C" (Revenue Share con agencias) representa una cesión de uso limitado (Licencia de Software), sin transferencia de titularidad de la IP de "La Cárcel de los Gemelos".

---
*Este documento está diseñado como base conceptual para revisión por un equipo jurídico especializado antes del lanzamiento comercial. Queda reservado cualquier otro derecho no especificado.*

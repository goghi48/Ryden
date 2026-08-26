# Diagram rendering demo

This page shows the same Ryden deployment view in formats commonly used in GitHub repositories.

## Mermaid

GitHub renders the fenced `mermaid` block directly in Markdown.

```mermaid
flowchart LR
    User[Browser] -->|HTTPS| Nginx[React + nginx]
    Nginx -->|REST API| API[Go modular monolith]
    Nginx -->|SSE updates| API
    API -->|pgx / SQL| DB[(PostgreSQL)]
    API --> Metrics[Prometheus metrics]

    subgraph Ryden backend
        API
        Auth[Auth]
        Meeting[Meetings]
        Polls[Polls and votes]
        Prep[Preparation]
        API --- Auth
        API --- Meeting
        API --- Polls
        API --- Prep
    end
```

## PlantUML

GitHub displays a `.puml` file as source code rather than rendering it natively. The generated SVG below is committed next to the [PlantUML source](ryden-architecture.puml), so it is visible in Markdown without a browser extension or external rendering at view time.

![Ryden architecture rendered from PlantUML](ryden-architecture.svg)

## Files

- [Mermaid embedded in Markdown](README.md#mermaid)
- [PlantUML source](ryden-architecture.puml)
- [PlantUML-generated SVG](ryden-architecture.svg)

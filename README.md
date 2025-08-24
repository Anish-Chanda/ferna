# Ferna

> **STATUS**: This project is in active development and may introduce breaking changes until the v1.0.0 release. Use the project for evaluation or development. Production stability is not guaranteed yet. Expect data migrations, API changes, and breaking refactors as the project evolves.

A practical app to help you keep your houseplants healthy.

Ferna is an open source plant care app that you can run yourself and grow with over time. It starts with the basics of tracking plants, setting reminders that can be snoozed or adjusted, and keeping notes and photos, but the long-term goal is much bigger. The project is designed around an API-first backend, making it easy to connect to other tools, build new clients, or integrate with sensors. Planned features include machine learning services for plant identification and diagnostics, so you can snap a photo to recognize a plant and get care suggestions that fit. Ferna is built to be more than a checklist app, it is meant to become a complete toolkit for plant care that stays private, extensible, and community-driven.

## Key features

- Personalized reminders for watering, fertilizing, and other care tasks
- Photos, notes, and a timeline so you can track each plant's history
- Early diagnostics and tooling (work in progress)
- Modular architecture designed to support future additions: machine learning, sensors, integrations
- API-friendly design to allow third-party integrations and custom clients

## Project status

Active development. The immediate goal is a polished MVP that reliably covers tracking and reminders. After that, expect incremental improvements to diagnostics and portability.

## High level roadmap  

- **Core**: Stable plant tracking, snoozable reminders, notes, photos, and a reliable mobile interface  
- **Near term**: Machine learning–based plant identification, public API endpoints (auth via API keys), data export/import tools.
- **Mid term**: Disease and health diagnostics powered by machine learning, and other features.
- **Long term**: Sensor integrations, home automation hooks, extensible APIs for community add-ons, and internationalization  

## Contributing

Contributions are welcome. If you'd like to help:

- Open an issue describing the problem or feature
- Submit small pull requests with focused changes and tests where applicable
- Changes to plant care information (schedules, intervals, recommended care notes) and other domain data are welcome
- If you plan to work on larger features, open an issue first so we can coordinate

Please add tests for new behavior when practical.

## Security & data

Ferna stores local data (SQLite, with optional PostgreSQL support planned). There are no external telemetry or tracking calls by default.

## License

This project is open source. See the `LICENSE` file for details.

Made by [Ferna Labs](https://fernalabs.com)

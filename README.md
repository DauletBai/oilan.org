# Oilan: AI for Self-Healing

[![Status](https://img.shields.io/badge/status-in_development-orange.svg)](https://github.com/your-username/oilan)
[![Go Version](https://img.shields.io/badge/go-1.24.5-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](/LICENSE)

Oilan is an experimental AI-powered platform designed to facilitate self-healing by helping users uncover the psychosomatic root causes of their conditions through a guided, empathetic dialogue.

## 📖 Philosophy & Mission

Modern medicine excels at treating symptoms, but often overlooks the deep connection between our mind, unresolved emotional traumas, and physical health. Many chronic conditions are not just biochemical imbalances but physical manifestations of a "bug in our mental code"—a maladaptive program recorded in our subconscious during a past traumatic event.

Oilan's mission is to provide a safe, private, and accessible tool that acts as an "AI-powered guide." It doesn't give commands or diagnoses. Instead, it uses a unique, non-directive dialogue methodology to help users navigate their own consciousness, identify the root of their issue, and, through that awareness, activate their own innate capacity for healing.

Our goal is to create a new paradigm in digital therapeutics, making deep inner work accessible and affordable for everyone, everywhere.

## ⚙️ Technology Stack

We believe in building a robust, maintainable, and scalable platform from the ground up. Our technology stack is deliberately minimalist and powerful, relying on proven technologies without unnecessary frameworks.

* **Backend:** Go (Golang) 1.24.5 (using only the standard library)
* **Database:** PostgreSQL
* **Frontend:** Bootstrap 5.3 (for a clean, responsive UI)
* **Deployment:** Docker

## 🚀 Project Status

The project is currently in the MVP (Minimum Viable Product) development stage. The core architecture is based on Clean Architecture principles, ensuring a clear separation of concerns and long-term maintainability.

We are building upon the experience gained from a previously developed, fully functional prototype which proved the viability of the core concepts.

## 🤝 Invitation to Collaborate

Oilan is more than just a software project; it's a mission to fundamentally change how we approach wellbeing. We are actively preparing to apply to leading startup programs such as Microsoft for Startups Founders Hub and Google for Startups Accelerator.

We are open to collaboration with researchers, mental health professionals, AI engineers, and potential strategic partners who share our vision. If you believe in the power of technology to unlock human potential, let's connect.

## 🏁 Getting Started (For Developers)

To get the project running locally using Docker:

1.  Clone the repository:
    ```bash
    git clone [https://github.com/your-username/oilan.git](https://github.com/your-username/oilan.git)
    cd oilan
    ```

2.  Create a `config.yml` in the `configs/` directory based on `config.example.yml`.

3.  Build and run the application using Docker Compose:
    ```bash
    docker-compose up --build
    ```

The application will be available at `http://localhost:8080`.
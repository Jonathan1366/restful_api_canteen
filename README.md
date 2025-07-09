# GoCanteen API - A Microservices-based E-commerce Backend

![Go Version](https://img.shields.io/badge/go-1.22.x-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-green)

A high-performance, scalable backend system for a modern canteen e-commerce application, built with Golang and a microservices architecture. This project is a personal learning endeavor to implement and explore enterprise-level backend patterns and technologies.

---

### 🚧 Project Status: Actively Under Development 🚧

**Please Note:** This is a work-in-progress project. The primary goal is to serve as a learning platform and a portfolio piece to demonstrate expertise in backend development, not as a production-ready application. Core features are being actively developed and the architecture is continuously being refined.

---

## 🏛️ Architecture Overview

This project implements a **Microservices Architecture** to ensure separation of concerns, independent scalability, and resilience. Each service is designed to handle a specific business domain.

Communication between services is planned to use a combination of:
* **RESTful APIs:** For synchronous, external-facing requests from the client.
* **gRPC:** For high-performance, low-latency internal communication between services.

<br>

**Services:**
* `auth_service`: Manages user registration, login, JWT generation, and identity verification.
* `product_service`: Handles all product-related logic, including catalogue, inventory, and details.
* `payment_service`: Responsible for integrating with third-party payment gateways and managing transaction lifecycle.
* `chat_service`: A planned real-time chat feature between users and canteen owners, likely using WebSockets or gRPC streams.

---

## 🛠️ Tech Stack & Tools

This project leverages a modern, cloud-native tech stack chosen for performance, scalability, and developer experience.

| Category                  | Technology / Service                                                                                                   |
| ------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| **Backend** | **Golang**, **Fiber** (Web Framework), **gRPC** (Inter-service Communication)                                            |
| **Database** | **PostgreSQL** (hosted on **Supabase**), **pgx** (Go Driver)                                                             |
| **DevOps & Infrastructure** | **Docker**, **Railway** (Deployment & Hosting), **Cloudflare** (DNS, Caching, Security)                                |
| **Cloud Services (Planned)** | **AWS / GCP** (e.g., S3 for object storage, Pub/Sub for asynchronous tasks)                                              |
| **API Testing & Docs** | **Postman**, **Swagger/OpenAPI** (Planned)                                                                             |

---

## ✨ Features

### ✅ Implemented Features

* **Authentication Service:**
    * User Registration with password hashing.
    * User Login with JWT-based authentication.
* **Product Service:**
    * `GET /api/products`: Retrieve a list of all available products.
    * `GET /api/products/:id`: Retrieve details for a single product.

### 🔄 Roadmap & Planned Features

* **gRPC Integration:** Refactor internal API calls to use gRPC for improved performance.
* **Complete Product CRUD:** Implement `POST`, `PUT`, `DELETE` endpoints for product management by admins.
* **Payment Gateway Integration:** Integrate with a payment provider like Midtrans or Stripe.
* **Real-time Chat:** Build the `chat_service` using WebSockets.
* **Role-Based Access Control (RBAC):** Differentiate permissions between `user` and `admin` roles.
* **Containerization:** Full Docker support for consistent development and deployment environments.
* **CI/CD Pipeline:** Automate testing and deployment using GitHub Actions.

---

## 🚀 Getting Started

### Prerequisites
* Go 1.22 or higher
* PostgreSQL / Supabase account
* Docker (Optional, for future use)

### Installation & Running
1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/](https://github.com/)[your-github-username]/restful_api_canteen.git
    cd restful_api_canteen
    ```

2.  **Navigate to a service directory:**
    ```bash
    cd auth_service
    ```

3.  **Set up environment variables:**
    Create a `.env` file in the service's root directory and populate it with the required variables (e.g., `DATABASE_URL`, `JWT_SECRET`). Refer to a `.env.example` file if provided.

4.  **Install dependencies and run the service:**
    ```bash
    go mod tidy
    go run main.go
    ```
    The service will start, typically on a port like `:3000`. Repeat for other services as needed.

---

## 📚 API Endpoints & Documentation

A Postman collection is being prepared. For now, here are some key endpoints you can test:

**Register a new user:**
```bash
curl -X POST http://localhost:PORT/api/auth/register \
-H "Content-Type: application/json" \
-d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword123"
}'
```

**Login and get JWT:**
```bash
curl -X POST http://localhost:PORT/api/auth/login \
-H "Content-Type: application/json" \
-d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
}'
```

---

## 👤 Author

* **Jonathan1366**
* GitHub: [@Jonathan1366](https://github.com/Jonathan1366)

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

# BookMySeat

**BookMySeat** is a microservices-based movie ticket booking system designed to handle the entire process of booking movie tickets. It leverages modern technologies and best practices for scalability, security, and maintainability.

## Features

- **Microservices Architecture**: The system is composed of several microservices, including User Management, Booking, Event Management, and Payment services.
- **RabbitMQ for Communication**: Inter-service communication is handled via RabbitMQ, ensuring reliable message delivery and decoupling of services.
- **JWT Authentication**: Secure user authentication is implemented using JSON Web Tokens (JWT).
- **Payment Integration**: The system integrates with Stripe as the payment gateway for processing payments securely.
- **Email Notifications**: Users receive email notifications for booking confirmations and other updates using Gmail SMTP.
- **AWS API Gateway Integration**: The API endpoints are exposed via AWS API Gateway, ensuring scalability and secure access.

## Technologies Used

- **Programming Languages**: Go (Golang), SQL
- **Microservices Framework**: Gin
- **Message Broker**: RabbitMQ
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Payment Gateway**: Stripe
- **Email Service**: Gmail SMTP
- **Cloud Services**: AWS API Gateway

## Getting Started

### Prerequisites

- Go (Golang) installed
- Docker and Docker Compose installed
- PostgreSQL database set up
- RabbitMQ server running
- Stripe account for payment processing
- Gmail account for email notifications

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/abidaziz98762/BookMySeat.git
    cd BookMySeat
    ```

2. Set up environment variables:

   Create a `.env` file in the root directory and add the following variables:

    ```env
    DATABASE_URL=postgres://username:password@localhost:5432/bookmyseat
    RABBITMQ_URL=amqp://guest:guest@localhost:5672/
    JWT_SECRET=your_jwt_secret_key
    STRIPE_SECRET_KEY=your_stripe_secret_key
    SMTP_USERNAME=your_gmail_username
    SMTP_PASSWORD=your_gmail_password
    ```


    ```

3. Set up AWS API Gateway:

   - Deploy your API endpoints to AWS API Gateway for secure and scalable access.

## Usage

1. **User Registration**: Users can register and authenticate via the User Management service.
2. **Booking Tickets**: Users can browse events, select seats, and book tickets through the Booking service.
3. **Payment Processing**: Payments are securely processed through Stripe.
4. **Email Notifications**: Users receive booking confirmation and updates via email.

## Project Structure

```plaintext
BookMySeat/
├── user-service/
├── booking-service/
├── event-service/
├── payment-service/
└── README.md

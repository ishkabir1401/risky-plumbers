# risky-plumbers
This project is a Go-based service. Below are the steps to run the service, install dependencies, and execute any tests.

## Prerequisites

Ensure you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

## Getting Started

### 1. Clone the Repository

Start by cloning this repository to your local machine:

```bash
git clone https://github.com/ishkabir1401/risky-plumbers.git
cd risky-plumbers
```
### 2. Install Dependencies

To install the required Go dependencies for the project, run the following command:

```bash
go get .
```

### 3. Run the Service

To run the service, use the following command:

```bash
go run main.go
```

### 4. Run Test

This command will run all the tests in your project. Specifically, if you want to run the tests in the routes folder, you can run:


```bash
go test ./routes
```
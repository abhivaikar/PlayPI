
# PlayPI

PlayPI (pronounced as Play-P-I, similar to API) is an open-source, simple, and local API playground that allows software engineers to test and experiment with various types of APIs. It is designed for hands-on learning and testing without requiring an internet connection or complex setup. 

With PlayPI, you can practice API testing across multiple technologies and protocols, including RESTful, gRPC, GraphQL, and WebSocket (and more to come). You can also use this playground if you are conducting a hands-on API testing workshop or bootcamp.

----------

## Why PlayPI?

PlayPI stands out as a versatile, multi-protocol API playground:

-   **Multiple API protocols**: Includes RESTful, gRPC, GraphQL, and WebSocket APIs (and more to come).
-   **Realistic use cases**: Each API implements meaningful functionalities such as inventory management, task management, user registration, and live chat.
-   **Offline testing**: No internet connection is required; all APIs run locally.
-   **Ease of use**: Simple CLI and Docker-based installation options make it beginner-friendly.

## Quick Start Using CLI

### Download the Binary

1.  Go to the [Releases](https://github.com/abhivaikar/playpi/releases) page of this repository.
2.  Download the binary for your platform (only macOS & Linux supported currently).
3.  Make the binary executable (if required):
    `chmod +x playpi` 
    
### Run the Playground
Use the following command to start the desired service:
`./playpi start [api-type]` 

Replace `[api-type]` with one of the following:
-   `restful-inventory-manager`
-   `graphql-inventory-manager`
-   `grpc-inventory-manager`
-   `websocket-inventory-manager`
-   `restful-task-manager`
-   `grpc-user-registration`
-   `websocket-live-chat`

### Example:
`./playpi start restful-inventory-manager`

## Docker Installation and Usage
If you are a docker fan and prefer not downloading the binary, you can run the playground using a docker image too!

### Pull the Docker Image
`docker pull abhijeetvaikar/playpi:latest` 

### Run the Playground
Run a specific API service:
`docker run -p <port>:<port> abhijeetvaikar/playpi start [api-type]` 

Replace `<port>` and `[api-type]` as needed.

### Example:
-   Start RESTful Inventory Manager:
    `docker run -p 8080:8080 abhijeetvaikar/playpi start restful-inventory-manager`
- Start RESTful Task Management API:
`docker run -p 8085:8085 abhijeetvaikar/playpi start restful-task-manager`
- Start GraphQL inventory Management API:
`docker run -p 8081:8081 abhijeetvaikar/playpi start graphql-inventory-manager`
- Start gRPC inventory management API:
`docker run -p 8082:8082 abhijeetvaikar/playpi start grpc-inventory-manager`
- Start gRPC user registration API:
`docker run -p 8084:8084 abhijeetvaikar/playpi start grpc-user-registration`
- Start websocket inventory manager API:
`docker run -p 8083:8083 abhijeetvaikar/playpi start websocket-inventory-manager`
- Start websocket live chat API:
`docker run -p 8086:8086 abhijeetvaikar/playpi start websocket-live-chat`

## Accessing a service once started
Which client you want to use to access the service on the playground is entirely upto you. But here are some suggestions.
-   **RESTful API**: Use Postman or curl or via your favourite programming language at `http://localhost:<port>`.
-   **gRPC API**: Test with grpcurl or Postman or via your favourite programming language on `localhost:<port>`.
-   **WebSocket API**: Connect using a WebSocket client like Postman or WebSocket King at `ws://localhost:<port>`.
- **GraphQL API**: Connect using a GraphQL client like Postman, GraphiQL, Insomnia or your favourite programming language at `http://localhost:8081/graphql`

## APIs and Use Cases

### RESTful API - Inventory Management
CRUD operations to manage an inventory of items. Example operations include adding, updating, deleting, and retrieving items.

#### Create item
HTTP Method: `POST`
URL: `/items`
Payload:
```json
{
"name": "Laptop",
"description": "A high-performance laptop",
"price": 1500.99,
"quantity": 10
}
```

**Validation and business rules**
- Name:
  - Must be between 3 and 50 characters.
  - Error: "name must be between 3 and 50 characters"
- Description:
  - Cannot exceed 200 characters.
  - Error: "description cannot exceed 200 characters"
- Price:
  - Must be a positive number not exceeding 10,000.
  - Error: "price must be a positive number not exceeding 10,000"
- Quantity:
  - Must be at least 0.
  - Error: "quantity must be at least 0"

#### Update item
HTTP Method: `PUT`
URL: `items/{id}`
Payload:
```json
{
"name": "Laptop",
"description": "A high-performance laptop",
"price": 1500.99,
"quantity": 1
}
```

**Validation and business rules**
- ID:
  - Must correspond to an existing item.
  - Error: "item not found"
- Name:
  - Must be between 3 and 50 characters.
  - Error: "name must be between 3 and 50 characters"
- Description:
  - Cannot exceed 200 characters.
  - Error: "description cannot exceed 200 characters"
- Price:
  - Must be a positive number not exceeding 10,000.
  - Error: "price must be a positive number not exceeding 10,000"
- Quantity:
  - Must be at least 0.
  - Error: "quantity must be at least 0"

#### Delete item
HTTP Method: `DELETE`
URL: `items/{id}`
Payload: No payload

**Validation and business rules**
- ID:
  - Must correspond to an existing item.
  - Error: "item not found"

#### Get all items
HTTP Method: `GET`
URL: `/items`

#### Update item - specific field
HTTP Method: `PATCH`
URL: `items/{id}`
Payload:
```json
{
"quantity" : 1
}
```
**Validation and business rules**
- ID:
  - Must correspond to an existing item.
  - Error: "item not found"
- Name:
  - If provided, must be between 3 and 50 characters.
  - Error: "name must be between 3 and 50 characters"
- Description:
  - If provided, cannot exceed 200 characters.
  - Error: "description cannot exceed 200 characters"
- Price:
  - If provided, must be a positive number not exceeding 10,000.
  - Error: "price must be a positive number not exceeding 10,000"
- Quantity:
  - If provided, must be at least 0.
  - Error: "quantity must be at least 0"

### RESTful API - Task Management
Manage tasks with fields such as title, description, due date, and status. Tasks can be marked as overdue based on their due date.

#### Create a new task
HTTP Method: `POST`
URL: `/tasks`
Payload:
```json
{
"title": "Write documentation",
"description": "Document all APIs for the PlayPI project",
"due_date": "2025-01-15",
"priority": "high"
}
```
**Validation and business rules**
- Title:
  - Must be between 3 and 100 characters.
  - Error: "title must be between 3 and 100 characters"
- Description:
  - Cannot exceed 500 characters.
  - Error: "description cannot exceed 500 characters"
- Due Date:
  - Must follow the format YYYY-MM-DD.
  - Cannot be in the past.
  - Error: "due date must follow the format YYYY-MM-DD"
  - Error: "due date cannot be in the past"
- Priority:
  - Must be one of: low, medium, high.
  - Error: "priority must be one of: low, medium, high"


#### Update a task
HTTP Method: `PUT`
URL: `/tasks/{id}`
Payload:
```json
{
"title": "Finalize documentation",
"description": "Update and finalize all API docs",
"due_date": "2025-01-20",
"priority": "medium",
"status": "in progress"
}
```

**Validation and business rules**
- ID:
  - Must correspond to an existing task.
  - Error: "task not found"
- Title:
  - Must be between 3 and 100 characters.
  - Error: "title must be between 3 and 100 characters"
- Description:
  - Cannot exceed 500 characters.
  - Error: "description cannot exceed 500 characters"
- Due Date:
  - Must follow the format YYYY-MM-DD.
  - Cannot be in the past.
  - Error: "due date must follow the format YYYY-MM-DD"
  - Error: "due date cannot be in the past"
- Priority:
  - Must be one of: low, medium, high.
  - Error: "priority must be one of: low, medium, high"

#### Mark task as complete
HTTP Method: `PUT`
URL: `/tasks/{id}/complete`
Payload: N/A

**Validation and business rules**
- ID:
  - Must correspond to an existing task.
  - Error: "task not found"
- Status:
  - If the task is already marked as completed, return an error.
  - Error: "task is already marked as completed"

#### Get all tasks
HTTP Method: `GET`
URL: `/tasks`

**Validation and business rules**
- No specific validation rules.
- If no tasks are created, return a message: "no tasks created"

#### Get a task
HTTP Method: `GET`
URL: `/tasks/{id}`

**Validation and business rules**
-  ID:
  - Must correspond to an existing task.
  - Error: "task not found"
- If no tasks are created, return a message: "no tasks created"

#### Delete a task
HTTP Method: `DELETE`
URL: `/tasks/{id}`

**Validation and business rules**
- ID:
  - Must correspond to an existing task.
  - Error: "task not found"

### gRPC API - Inventory Management
- Full CRUD support for managing inventory.
- Proto file for generating client can be found in `services/grpc/inventory_management/pb/inventory.proto`

#### AddItem
Payload
```json
{
"description": "reprehenderit ex et anim",
"name": "Ut mollit",
"price": 10,
"quantity": 12
}
```
**Validation and business rules**
- Name:
  - Must be between 3 and 50 characters.
- Description:
  - Cannot exceed 200 characters.
- Price:
  - Must be a positive number not exceeding 10,000.
- Quantity:
  - Must be at least 0.
- Error:
  - If any validation rule is violated, an appropriate error message will be returned.

#### GetItem
Payload
```json
{
"id": 1
}
```
**Validation and business rules**
- ID:
  - Must correspond to an existing item.
- Error:
  - If the item with the specified ID does not exist, an error message "item not found" will be returned.

#### ListItems
No payload

#### UpdateItem
Payload
```json
{
"description": "dolor et nostrud reprehenderit",
"id": 1,
"name": "ex ut velit fugiat eiusmod",
"price": 6668482.157076433,
"quantity": 1800787081
}
```
**Validation and business rules**
- ID:
  - Must correspond to an existing item.
- Name:
  - If provided, must be between 3 and 50 characters.
- Description:
  - If provided, cannot exceed 200 characters.
- Price:
  - If provided, must be a positive number not exceeding 10,000.
- Quantity:
  - If provided, must be at least 0.
- Error:
  - If the item with the specified ID does not exist, an error message "item not found" will be returned.
  - If any validation rule is violated, an appropriate error message will be returned.

#### DeleteItem
Payload
```json
{
"id": 1
}
```

**Validation and business rules**
- ID:
  - Must correspond to an existing item.
- Error:
  - If the item with the specified ID does not exist, an error message "item not found" will be returned.
  - If the item has a quantity greater than 0, an error message "cannot delete an item with stock remaining" will be returned.

### gRPC API - User Registration and Sign-In
- Register a new user, sign in with a username and password, view profiles, and update/delete account details.
- Proto file for generating client can be found in `services/grpc/user_registration/pb/user.proto`

##### RegisterUser
Payload:
```json
{
"user": {
"address": "ut",
"email": "consectetur@gmail.com",
"full_name": "test",
"password": "test",
"phone": "2323431",
"username": "testuser"
 }
}
```

**Validations and business rules**
- Username:
  - Must be between 3 and 50 characters.
  - Must be unique (no duplicate usernames allowed).
- Password:
  - Must be at least 8 characters long.
- Full Name:
  - Cannot be empty.
- Email:
  - Must be a valid email format.
- Phone:
  - Must be numeric and between 10-15 digits.
- Address:
  - Cannot exceed 100 characters.

#### SignIn
Payload:
```json
{
"password": "test",
"username": "testuser"
}
```
**Validation and business rules**
- Username and Password:
  - Both fields are required.
  - Must match an existing user's credentials.

#### GetProfile
Payload:
```json
{
"token": "token_from_signin_response"
}
```

**Validation and business rules**
- Token:
  - Must be valid and correspond to an existing session.

#### UpdateProfile
Payload:
```json
{
"token": "token_from_signin",
"user": {
"address": "sint",
"email": "xyz@gmail.com",
"full_name": "cillum sed",
"password": "awdad",
"phone": "43141",
"username": "updated_username"
 }
}
```
**Validation and business rules**
- Token:
  - Must be valid and correspond to an existing session.
- Username:
  - If provided, must be between 3 and 50 characters.
  - If changed, must be unique.
- Password:
  - If provided, must be at least 8 characters long.
- Full Name:
  - If provided, cannot be empty.
- Email:
  - If provided, must be a valid email format.
- Phone:
  - If provided, must be numeric and between 10-15 digits.
- Address:
  - If provided, cannot exceed 100 characters.

#### DeleteAccount
Payload:
```json
{
"token": "token_from_signin"
}
```

**Validation and business rules**
- Token:
  - Must be valid and correspond to an existing session.

### GraphQL API - Inventory management
Query and mutate inventory data with a flexible schema.

#### Get all items
```
{
  items {
    id
    name
    description
    price
    quantity
  }
}
```
#### Get item by item id
```
{
item(id: 1) {
name
description
quantity
}
}
```

**Validation and business rules**
- invalid item id will return error "item not found"

#### Add item
```
mutation {
addItem(name: "Headphones", description: "Noise-canceling headphones", price: 299.99, quantity: 10) {
id
name
description
 }
}
```

**Validation and business rules**
- **Name**:
  - Must be between 3 and 50 characters.
  - Must be unique (no duplicate names allowed).
- **Description**:
  - Cannot exceed 200 characters.
- **Price**:
  - Must be a positive number not exceeding 10,000.
- **Quantity**:
  - Must be at least 1.

#### Update item
```
mutation {
updateItem(id: 1, name: "Updated Laptop", quantity: 20) {
id
name
description
price
quantity
 }
}
```

**Validation and business rules**
- **Name**:
  - Must be between 3 and 50 characters.
- **Description**:
  - Cannot exceed 200 characters.
- **Price**:
  - Must be a positive number not exceeding 10,000.
- **Quantity**:
  - Cannot be negative.
- **Item id**:
  - Has to be a valid item id

#### Delete item
```
mutation {
deleteItem(id: 2)
}
```

**Validation and business rules**
- **Item Not Found**:
  - If the item with the specified ID does not exist, an error message "item not found" will be returned.
- **Cannot Delete Item with Stock Remaining**:
  - If the item has a quantity greater than 0, an error message "cannot delete an item with stock remaining" will be returned.

### WebSocket API - Live Chat
Multi-user chat functionality with join/leave notifications, public and private messaging.

#### Connect to the chat server as a user
- Assigns a random username to the client the moment a user connects to the websocket server.
- Sends a system welcome message to the user immediately post connection.
- Broadcasts a system message to all existing connected users that a user has joined. Example: "Cherry1737899928256010000 has joined the chat"

**Validation**
- Max 5 users allowed.

#### Broadcast a message to all chat users
- Sends a message to all connected users except the sender.

Request payload
```
{
    "type": "chat",
    "username": "john_doe",
    "message": "Hello, everyone!"
}
```

**Validation**
- Message type must be "chat".
- Message or username cannot be empty.
- Sender does not receive their own message.

#### Send a private message to a specific user
Request payload
```
{
    "type": "private",
    "username": "john_doe",
    "message": "Hi Jane!",
    "to": "jane_doe"
}
```

**Validation**
- Message type must be "private".
- recipient must exist.
- Message cannot be empty.
- Sender does not receive their own message.

#### Leave chat
You just need to disconnect from the websocket connection and all other clients connected to the websocket will get a message that you have left.
Example: "username has left the chat"

## Contributing to PlayPI

We welcome contributions to make PlayPI even better! Whether it's fixing a bug, suggesting a new feature, or improving the documentation, your input is valuable.

*The most important contributions needed will be to add new and real-world use cases (beyond inventory, task manager, chat etc) and/or protocols (beyond restful, gRPC, websocket etc).*

### How to Contribute
1.  Fork the repository
2.  Clone your fork locally
3.  Create a new branch
4.  Make your changes
    -   Add new features, fix bugs, refactor existing code, add tests or improve documentation.
5.  Test your changes locally to ensure they work as expected.
	- Currently the project does not have automated tests for the services. We will be adding them soon and create an automated build and test Github actions pipeline.
6. Commit and push    
7. Submit a pull request
    -   Go to your forked repository on GitHub.
    -   Click the "Compare & pull request" button.
    -   Describe your changes and submit the pull request.

### Building and running from source locally
In order to build the source locally, please run below commands:
```
# this compiles the binary
go build -o playpi main.go

# test running the binary
./playpi start 
```

### Creating binaries for distribution
Use Go's cross-compilation feature to build binaries for multiple platforms:
```
GOOS=windows GOARCH=amd64 go build -o playpi-windows.exe main.go
GOOS=darwin GOARCH=amd64 go build -o playpi-mac main.go
GOOS=linux GOARCH=amd64 go build -o playpi-linux main.go
```
## Reporting Issues and Feedback

We value your feedback to improve PlayPI. If you encounter any issues or have suggestions, here's how you can raise them:

1.  **Search for existing issues**:
    
    -   Check the [Issues](https://github.com/abhivaikar/playpi/issues) tab on GitHub to see if your issue has already been reported.
2.  **Create a new issue**:
    -   If your issue or feedback is new, click on the "New Issue" button in the [Issues](https://github.com/abhivaikar/playpi/issues) tab.
    -   Provide a clear and detailed description of the issue or feedback:
        -   What were you trying to do?
        -   What happened instead?
        -   Steps to reproduce the issue (if applicable).
        -   Any relevant logs or screenshots.
3.  **Feature requests**:
    -   Label your issue as a "Feature Request" and provide details about the functionality you'd like to see.

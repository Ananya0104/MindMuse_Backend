# MindMuse Backend

## Setup Guide: Running Locally

### Prerequisites
- Go 1.23 or later installed: https://go.dev/dl/
- AWS credentials configured (for DynamoDB access)
- (Optional) [direnv](https://direnv.net/) or similar tool to auto-load `.env` variables

### 1. Clone the Repository
```sh
git clone https://github.com/Ananya0104/MindMuse_Backend
cd MindMuse_Backend
```

### 2. Set Up Environment Variables
Create a `.env` file in the project root with the following content:

```
HUGGINGFACE_API_KEY=your_huggingface_api_key_here
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_REGION=ap-south-1
PORT=8080
```

> Replace the values with your actual credentials and API keys.

### 3. Install Go Dependencies
```sh
go mod tidy
```

### 4. Run the Server Locally
```sh
go run main.go
```

The server will start on `http://localhost:8080` by default.

### 5. API Endpoints
- Health check: `GET /health`
- Auth, journal, emergency, survey, and chat endpoints are available under `/api/`

### 6. Running Tests
```sh
go test ./...
```

## Notes
- The backend uses AWS DynamoDB for data storage. Make sure your AWS credentials have the necessary permissions.
- The Hugging Face API key is required for chat/AI features.
- For local development, you may use [AWS DynamoDB Local](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html) if you don't want to connect to a real AWS account.

---

For any issues or contributions, please open an issue or pull request on the repository.

## Open Source Libraries, Frameworks, and Tools Used

| Name                                   | Version      | License         | Role in Build                | Source Link                                              |
|----------------------------------------|--------------|-----------------|------------------------------|---------------------------------------------------------|
| Go (Golang)                           | 1.23.0+      | BSD-3-Clause    | Main programming language     | https://go.dev/                                          |
| github.com/gin-gonic/gin              | 1.10.1       | MIT             | Web framework (API server)    | https://github.com/gin-gonic/gin                         |
| github.com/aws/aws-lambda-go          | 1.48.0       | Apache-2.0      | AWS Lambda integration       | https://github.com/aws/aws-lambda-go                     |
| github.com/aws/aws-sdk-go-v2          | 1.36.3+      | Apache-2.0      | AWS SDK for DynamoDB, etc.   | https://github.com/aws/aws-sdk-go-v2                     |
| github.com/awslabs/aws-lambda-go-api-proxy | 0.16.2  | Apache-2.0      | Lambda <-> Gin adapter       | https://github.com/awslabs/aws-lambda-go-api-proxy       |
| github.com/dgrijalva/jwt-go           | 3.2.0        | MIT             | JWT token handling           | https://github.com/dgrijalva/jwt-go                      |
| github.com/google/uuid                 | 1.6.0        | BSD-3-Clause    | UUID generation              | https://github.com/google/uuid                           |
| github.com/stretchr/testify            | 1.10.0       | MIT             | Testing utilities            | https://github.com/stretchr/testify                      |
| golang.org/x/crypto                   | 0.38.0       | BSD-3-Clause    | Password hashing, crypto     | https://pkg.go.dev/golang.org/x/crypto                   |
| google.golang.org/api                 | 0.236.0      | Apache-2.0      | Google API integration       | https://pkg.go.dev/google.golang.org/api                 |

> **Note:** Some dependencies are used transitively (indirectly) via the above libraries. See `go.mod` for the full list.

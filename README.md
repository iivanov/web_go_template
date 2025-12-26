# Database Setup with GORM v2

This project uses GORM v2 as the ORM with PostgreSQL.

## Quick Start

1. Start PostgreSQL using Docker Compose:
```bash
make docker-up
```

2. Set environment variables (optional, defaults are provided):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=gonewproject
export DB_USER=gonewproject
export DB_PASSWORD=gonewproject
export DB_SSLMODE=disable
```

3. Run the application:
```bash
make run
```

## Repository Pattern

The project implements the repository pattern for data access:

### User Repository Interface
```go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id uint) (*entities.User, error)
    GetByUID(ctx context.Context, uid string) (*entities.User, error)
    GetAll(ctx context.Context) ([]*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uint) error
}
```

### Usage Example
```go
// The repository is injected via dependency injection (fx)
type UserService struct {
    repository repositories.UserRepository
}

func (s *UserService) GetUser(ctx context.Context, uid string) (*dto.UserResponse, error) {
    user, err := s.repository.GetByUID(ctx, uid)
    if err != nil {
        return nil, err
    }
    
    return &dto.UserResponse{
        UID:   user.UID,
        Name:  user.Name,
        Email: user.Email,
    }, nil
}
```

## Database Migrations

GORM automatically handles migrations when the application starts. The `User` entity will be migrated automatically:

```go
fx.Invoke(func(db *gorm.DB) error {
    return db.AutoMigrate(&entities.User{})
})
```

## Makefile Commands

- `make docker-up` - Start PostgreSQL container
- `make docker-down` - Stop PostgreSQL container
- `make docker-logs` - View PostgreSQL logs
- `make build` - Build the application with GOEXPERIMENT=jsonv2
- `make run` - Build and run the application

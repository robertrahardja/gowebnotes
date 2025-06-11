# DB Call Flow

```mermaid
flowchart TD

    App -->|store:| NewPostgresStorage
    App -->|config:| Config
    App --> Mount[App.mount]
    App --> Run[App.run]
    Run --> Mount
    Mount --> Router

    Config --> DbConfig[dbConfig struct]
    Config -->|from env| Port[Port Addr]
    DbConfig -->|from env| Addr
    DbConfig -->|from env| maxOpenConns
    DbConfig -->|from env| maxIdleConns
    DbConfig -->|from env| maxIdleTimes

    Addr --> db.new
    maxOpenConns --> db.new
    maxIdleConns --> db.new
    maxIdleTimes --> db.new

    db --> db.new
    db --> NewPostgresStorage
    NewPostgresStorage --> StoreFunctions


    subgraph StoreFunctions
      GetByID
      Create
    end

    StoreFunctions --> PostStore[PostStore struct]
    StoreFunctions --> UserStore[UserStore struct]

    PostStore --> PostModel[Post struct]
    UserStore --> UserModel[User struct]

    Router --> health["/health"]
    health --> healthcheckhandler

    health--> post["/posts"]
    post --> createPostHandler
    createPostHandler --> Create

    post --> postID["/{postID}"]
    postID --> getPostHandler
    getPostHandler --> GetByID

    PostStore --> DB[(PostgreSQL)]
    UserStore --> DB

    PostModel -.-> PostCreate[PostStore.Create]
    UserModel -.-> UserCreate[UserStore.Create]

    PostCreate --> DB
    UserCreate --> DB


    classDef Main fill:#ee1313
    class App,getPostHandler,createPostHandler Main

    classDef DB fill:#2f4fee

    classDef Store fill:#a721f0
    class GetByID,Create,PostStore,UserStore Store
```

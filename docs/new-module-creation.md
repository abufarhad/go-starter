#### When creating a new module, follow these steps. Let's consider the creation of the "user" module.
1. In the root directory, create a folder named "user" (rename as needed).
2. Within the "user" folder, create three subfolders: "controller," "service," and "repository."
- Controller: Define routes and controller methods here.
- Service: Implement business logic in this folder.
- Repository: Handle database logic and queries in this folder.
3. In the domain folder declare the db entity and inside "domain/dto" jot down the api request and response
4. In the "cmd/serve.go" file, specify dependency injections for the module (User module) in the "ApisToServe" function.
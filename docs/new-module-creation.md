#### When anyone is going to create a new module, follow these steps:

1. Create a folder in the root directory named as - user (you can rename it to whatever suits best).
2. Inside the user folder, create three new folders - controller, service, and repository.
   - Controller: Write down all routes and controller methods here.
   - Service: Write down business logic here.
   - Repository: Write the database logics/queries here.
3. Write the dependency injections for the module (User module) inside the "cmd/serve.go" file (ApisToServe function).


Please refer to the example we have in our codebase for the "system_check" module.


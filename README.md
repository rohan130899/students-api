Create a 'students' api (GET, POST, PUT, DELETE, GET by ID)
- Should be accessible only after authentication (JWT)
- Returns appropriate http status codes
- Create an Authentication API
- Invoke the API with user-id & password
- Returns JWT token
- The 'Student' Model should have the following attributes
  CreatedBy
  CreatedOn
  UpdatedBy
  UpdatedOn
- Important: the user-id should be passed from http transport layer to the db layer
- Persist the configuration information in a config file
- Log in a log file

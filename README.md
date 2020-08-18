# iWant
an api + db to record "wants" (tasks, meetings, reviews, etc) by team members slack

**required environment variables:**
- SLACK_TOKEN (this is your slack apps "bot user OAuth access token")
- DB_USERNAME
- DB_PASSWORD
- DB_PROTOCOL
- DB_HOST
- DB_PORT
- DB_NAME
- APP_ADMIN_USERS (this should be a comma separated list of slack usernames)

**required slack scopes (add these to your apps "bot token scopes"):**
- chat:write
- commands
- im:write
- users:read

**slash command urls**
- http://\<your host\>/get-wants
- http://\<your host\>/delete-want
- http://\<your host\>/update-want
- https://\<your host\>/create-want

# rust-server-map-deleter
**What**: a simple lambda function that automatically deletes map files for the game Rust

**Why**: Rust has an extremely buggy gameserver that sometimes corrupts maps between server restarts. By deleting the map and allowing the server to regenerate it on startup, the game server will start successfully instead of intermittently failing due to a corrupt map.

# Requirements
- Go 1.x (Navigate to https://go.dev to install the binaries for your OS)
- OSX or Linux (or WSL2 if running on windows)

# Deployment
It's recommended to deploy this to AWS as it's completely free for this use case.
1. Sign up for an AWS account
2. Create a Lambda function, make sure you select `go 1.x` as the runtime
3. Change the handler to `bootstrap`
4. Make sure the input is set to Event Bridge with the following settings (make sure to change the schedule to any time just before your game server is scheduled to restart - if you aren't familiar with cron expressions, check out https://crontab.guru):

![system-design](https://github.com/bakatz/rust-server-map-deleter/assets/1575240/3ddaff01-e89e-4094-8a2b-0371dd8f7396)

5. Add the following Environment Variables to the Lambda function you just created:
```
SFTP_HOST_PORT: "hostipaddress:port"
SFTP_USERNAME: "yourusername"
SFTP_PASSWORD: "yourpassword"
```
If you don't know the port, it's probably `21`.

6. On your local machine, run ./build.sh which will then output a lambda-handler.zip file
7. Back in AWS lambda, upload the zip file
8. To test and make sure everything is working, use the Test menu in the AWS Lambda Console to send a test event to the lambda function. It should report back "success." Otherwise, wait until the scheduled time that you configured as a cron expression and the function will automatically execute.

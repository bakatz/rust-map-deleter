# rust-server-map-deleter
**What**: a simple lambda function that automatically deletes map files for the game Rust (on the server side, only helpful for server owners)

**Why**: Rust has an extremely delicate gameserver that sometimes corrupts maps between server restarts depending on how the server is restarted. This can especially happen when using plugins like SmoothRestarter or relying on your host's restart functionality which may not be implemented properly. By deleting the map and allowing the server to regenerate it on startup, the game server will start successfully instead of intermittently failing due to a corrupt map. Rust map generation is deterministic anyway, so as long as you don't change the seed the same map will be generated every time/there's no risk of data loss.

**How do you know if this is a useful tool for your server**: If you've gotten an error like `LZ4 block is corrupted, or invalid length has been given.` when starting your server, deleting the .map file will typically fix it. If you don't want to do that manually every day or your game server host doesn't have a modern control panel with scheduled tasks, use the approach in this repository to automatically clean up your maps.

# Requirements (only necessary if you want to build from source, otherwise just skip to the deployment instructions)
- Go 1.x (Navigate to https://go.dev to install the binaries for your OS)
- OSX or Linux (or WSL2 if running on windows)
- `zip` should be installed with `sudo apt install zip` or similar

# Deployment instructions
It's recommended to deploy this to AWS as it's completely free for this use case.
1. Sign up for an AWS account
2. Create a Lambda function, make sure you select `Custom runtime on Amazon Linux 2` as the runtime
3. Change the handler to `bootstrap`
4. Make sure the input is set to Event Bridge with the following settings (make sure to change the schedule to any time just **_before_** your game server is scheduled to restart - if you aren't familiar with cron expressions, check out https://crontab.guru):

![system-design](https://github.com/bakatz/rust-server-map-deleter/assets/1575240/3ddaff01-e89e-4094-8a2b-0371dd8f7396)

5. Add the following Environment Variables to the Lambda function you just created:
```
SFTP_HOST_PORT: "hostipaddress:port"
SFTP_USERNAME: "yourusername"
SFTP_PASSWORD: "yourpassword"
GAME_SERVER_BASE_PATH: "/some/absolute/path/to/rust"
DISCORD_WEBHOOK_URL: "https://discord.com/api/webhooks/someidgoeshere/somesecretgoeshere" (this is optional, feel free to leave it out it if you don't want a webhook sent or you don't have a Discord for your game server)
```
If you don't know the port, it's probably `21`. For game server base path, use an absolute path (starting with a slash) and make sure the directory is the one with the .map file in it.

6. Go to the latest releases page: https://github.com/bakatz/rust-server-map-deleter/releases and download the lambda-handler.zip file. Alternatively, on your local machine, run ./build.sh which will then output a lambda-handler.zip file.
7. Back in AWS lambda, upload the zip file from the above step under the "Code" menu
8. To test and make sure everything is working, use the Test menu in the AWS Lambda Console to send a test event to the lambda function. It should report back "success." You can also just wait until the scheduled time that you configured as a cron expression and the function will automatically execute.

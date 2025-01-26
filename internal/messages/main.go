package messages

var Greeting = `Hello, you are using Telegram bot for Todoist integration!

This bot can periodically notify about your planned tasks marked with specific labels:
- tg-notify-7d - notifies about task every 7 days;
- tg-notify-14d - notifies about task every 14 days;
- tg-notify-30d - notifies about task every 30 days;
- tg-notify-1M - notifies about task every month.

Every notification is calculated according to your 'due to' settings in Todoist task.`

var GettingStarted = `Caution: Fine-granted authorization process is in development! Developer access token provides full access to your account. Do it on your own risk! 

Go to your Profile -> Settings -> Integrations -> Developer and copy access token from there.`

var AddingTokenSuccessed = `Token added successfully!`

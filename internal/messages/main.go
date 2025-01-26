package messages

var Greeting = `Hello, you are using Telegram bot for Todoist integration!

This bot can periodically notify about your planned tasks marked with specific labels:
- tg-notify-7d - notifies about task every 7 days;
- tg-notify-14d - notifies about task every 14 days;
- tg-notify-30d - notifies about task every 30 days;
- tg-notify-1M - notifies about task every month.

Every notification is calculated according to your 'due to' settings in Todoist task.`

var GettingStarted = `Open the button below to authorize the bot and get access to your Todoist account.

Bot will get permissions to read all your tasks and projects.
If you don't have a Todoist account yet, you can sign up here: https://todoist.com/signup.`

var AuthorizationSuccessful = `Authorization successful!`

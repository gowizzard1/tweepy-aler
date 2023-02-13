## Tweepy-alert

Tweepy-alert is a simple program that uses the Twitter API to send an email whenever a specified keyword is mentioned on Twitter.

### Requirements
- Go 1.16 or later
- A Twitter API key
- An SMTP email server
### Usage

```php
go build
./egotter -smtpHost <SMTP_HOST> -smtpPort <SMTP_PORT> -smtpUser <SMTP_USER> -smtpPass <SMTP_PASS> -notificationEmail <NOTIFICATION_EMAIL> -monitoredKeywords <MONITORED_KEYWORDS>
```
or with environment variables set environment variables

```bash
export SMTP_HOST=<SMTP_HOST>
export SMTP_PORT=<SMTP_PORT>
export SMTP_USER=<SMTP_USER>
export SMTP_PASS=<SMTP_PASS>
export NOTIFICATION_EMAIL=<NOTIFICATION_EMAIL>
export MONITORED_KEYWORDS=<MONITORED_KEYWORDS>

go build

./egotter
```
### Options
```vbnet
-monitoredKeywords string
Comma-separated list of keywords to monitor on Twitter (default "egotter")
-notificationEmail string
Email address to send notifications to
-smtpHost string
SMTP server host (default "localhost")
-smtpPass string
SMTP server password
-smtpPort string
SMTP server port (default "25")
-smtpUser string
SMTP server username
```

### Running with Docker
You can also run Tweepy-alert using Docker and Docker Compose. To do so, simply build the Docker image and run it using Docker Compose, as described in the Dockerfile and docker-compose.yml files.

### License
Tweepy-alert is licensed under the MIT License. See LICENSE for details.




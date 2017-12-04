# slackstatus
Go app to set your slack status automatically depending of your connected wifi SSID or outgoing IP

# Installation
`slackstatus` is a simple binary to set your slack status. The best way to use is to run it as a cron command.

* You can download the binary corresponding to your architecture in the [release page](https://github.com/jvermillard/slackstatus/releases).
* if you have Go installed just run `go get github.com/jvermillard/slackstatus`

# Usage

For using `slackstatus` you need a slack API access token.
You can get it from: https://api.slack.com/custom-integrations/legacy-tokens

You also need a CSV file listing your known WiFi network names or the known outgoing IP and the corresponding status and emoji to associate.

For examples:
```
vrmvrm,remote,:house_with_garden:
daofficewifi,,
coprwifi,office,:office:
85.13.160.132,Office,:office:
```
`slackstatus` will look for the WiFi SSID names, then look for the exit IP (using http://ip-api.com).

Finally if nothing match, `slackstatus` will put you in status `traveling` with a flag corresponding to the detected country.

# Run it

`./slackstatus -token ****-*****-*****-*****-***** -file location.txt`

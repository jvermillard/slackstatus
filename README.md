# slackstatus
App written in go to set your slack status automatically depending on the SSID of the Wifi you are connected to, or your outgoing IP

# Installation
`slackstatus` is a simple binary to set your slack status. The best way to use it is to run it as a cron command.

* You can download the binary corresponding to your architecture in the [release page](https://github.com/jvermillard/slackstatus/releases).
* if you have Go installed just run `go get github.com/jvermillard/slackstatus`

# Usage

To user `slackstatus` you need a slack API access token.
You can get it from: https://api.slack.com/custom-integrations/legacy-tokens

You also need a CSV file listing your known WiFi network names or the known outgoing IP and the corresponding status and emoji you want to associate to them.

For examples:
```
homewifi,remote,:house_with_garden:
daofficewifi,,
corpwifi,office,:office:
85.13.160.132,Office,:office:
```
`slackstatus` will look for the WiFi SSID names, then look for the outgoing IP (using http://ip-api.com).

Finally if nothing matches, `slackstatus` will put you in status `traveling` with a flag corresponding to the detected country.

# Run it

`./slackstatus -token ****-*****-*****-*****-***** -file location.txt`

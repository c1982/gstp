# Gmail Subject Tracker for Prometheus (GSTP)

This tool follows your gmail mails and adds prometheus metrics when they match the filter rules.

## Motivation

Notifications of some back-end services come only by e-mail (for example: accounting system backup or notification about an important issue within the office). I wanted to monitor these notifications in the way I used to.

## Obtaining OAuth2 Credentials

* Go to the Google Developers API Console: https://console.developers.google.com/apis/credentials
* Click to "Create Credentials" button
* Select "Create OAuth client ID"
* Check "Other" 
* Enter your client name and click create button
* You can see in the "OAuth 2.0 Client IDs" your client. Download it.
* Copy to in gstp binary directory

## Installation

* Download

## First Run

* In order for the tool to work, you must first allow your gmail account.
* After putting credential.json file in the same directory as gstp. Issue the ./gstp command through the terminal. It will give you a URL address.
* Open this URL from your browser and continue with your gmail account.
* After giving the necessary permissions, it will give you an authorisation code.
* After entering this code in the relevant field, your token.js file will be created.

## Configuration

```yaml
#You can use filter your Gmail search results. More info: https://support.google.com/mail/answer/7190?hl=en
query: in:inbox is:unread newer_than:1d

#Your mail box. If you use over organization you must change this
userid: me

#Google Cloud Credentials file. You can get from https://console.developers.google.com/apis/credentials
credential: ./credential.json

#This file create automatically.
token: ./token.json

#Default prometheus export path. Default: http://localhost:8080/metrics
webpath: /metrics

#Default prometheus exporter http port.
port: :8080

#Interval of email check. Default: 5m (5 minutes) Syntax: 1m, 1h, 1d, 30m
check_interval: 5m

#Filter rules of email subjects.
filters:
  - filter:
    #E-mail topics are matched according to this regex.
    #Regex Syntax: https://github.com/google/re2/wiki/Syntax
    subject_regex: ^\[SOCRadar Incident\].+
    #Prometheus metric label
    label: socradar_incidend
  - filter:
    subject_regex: .+MySQL Daily Full Backup.+
    label: mysql_backup
  - filter:
    subject_regex: .+CEO Here.+
    label: ceo_in_the_building
```

## Example prometheus Metrics

```sh
# HELP gstp_counter subject count of filtered email message
# TYPE gstp_counter counter
gstp_counter{subject="socradar_incidend"} 5
gstp_counter{subject="mysql_backup"} 18
gstp_counter{subject="ceo_in_the_building"} 0
```
## Running tests

> go test -v ./...

## Using Docker

Also you can run easily with docker this tool.

* git clone git@github.com:c1982/gstp.git
* cd gstp
* follow the build and run sections below

##### build docker

```sh
GOOS=linux
go build -ldflags "-s -w" -o gstp
docker build -t gstp:latest .
```

##### run docker

```sh
docker run -d -rm -p 8080:8080 gstp
```

## Daemon

for debian:

```sh
mkdir -p /opt/gstp
cp gstp config.yaml /opt/gstp
chmod +x /opt/gstp/gstp
cp gstp.services /etc/systemd/system/gstp.service
systemctl enable gstp
systemctl start gstp
systemctl status gstp
```

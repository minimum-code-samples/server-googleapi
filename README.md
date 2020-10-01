# Google API Access to Google Sheets with Go

## Running the demo

> go run cmd/webserver/main.go

## Description

This demo runs a server that prompts the user to sign in. The first sign-in will request for a access+refresh token. This refresh token will be stored on the server as "secrets/google_admin_token.json".

If this file is present, subsequent sign-ins will request for access token (without the refresh token). This file is checked during start-up of the server.

Two endpoints are present to test access to a spreadsheet via the refresh token and the access token. Both of these endpoints are served using a single function.

The access token is stored in session whereas the one with the refresh token is stored as a file.

The client secrets JSON file may be specified either in the configuration file "config/web.yaml" or via the run flag "gac". A different configuration file may also be set using the "conf" flag.

## What's the point?

This demo shows that a refresh token can be used to access a protected (i.e non-public) Google Sheets document without requiring active user input.

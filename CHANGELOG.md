# CHANGELOG

## Version 1.3.2
- Fix the Pipelined Authentication Read Buffer Carryover bug by clearing the read buffer upon recieving an invalid auth method.

## Version 1.3.1
- Fix a bug where the `authRetriesBeforeSuccess` setting output in logs was accidentally set to the value of immediateTlsWindow.
- Fix JSON format in settings.json by removing comments.
- Use native uuid library instead of importing google's

## Version 1.3.0
- Reject the first few authentication attempts to encourage bots to provide more credentials. The number of attempts before auth is accepted is controlled via `authRetriesBeforeSuccess`

## Version 1.2.1
- Added `settings` field to the honeypot logs. Contains the setting values used during the connection.

## Version 1,2,0
- Added settings file which allows for partial rollouts of settings values
- Upgraded to golang 1.27.0

## Version 1.1.0
- Added SMTP_HONEYPOT_IMMEDIATE_TLS_WINDOW. This envvar controls the time the honeypot waits before sending the SMTP banner. Higher values make it more likely for slow immediate-tls connections to succeed, but may cause plaintext connections to time out. See https://github.com/Radicalius/smtp_honeypot/issues/3

## Version 1.0.0
- Added version string and incorportated it into honeypot logs
# CHANGELOG

## Version 1.1.0
- Added SMTP_HONEYPOT_IMMEDIATE_TLS_WINDOW. This envvar controls the time the honeypot waits before sending the SMTP banner. Higher values make it more likely for slow immediate-tls connections to succeed, but may cause plaintext connections to time out. See https://github.com/Radicalius/smtp_honeypot/issues/3

## Version 1.0.0
- Added version string and incorportated it into honeypot logs
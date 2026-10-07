import base64

from smtp_test_runner import run_smtp_test


def test_big_body():
    run_smtp_test(
        command=r'''swaks --server localhost --port 2525 \
  --to recipient@test.com --from recipient@test.com \
  --h-Date: "Thu, 01 Jan 2020 00:00:00 +0000" \
  --h-Subject: "test" \
  --h-Message-Id: "<fixed@test>" \
  --ehlo fixed.host \
  --body @big_body.data''',
        expected_output='''<-  220 mx01.example.net ESMTP
 -> EHLO fixed.host
<-  250-mx01.example.net
<-  250-PIPELINING
<-  250-SIZE 10240000
<-  250-VRFY
<-  250-ETRN
<-  250-STARTTLS
<-  250-AUTH LOGIN PLAIN
<-  250 8BITMIME
 -> MAIL FROM:<recipient@test.com>
<-  250 2.1.0 OK
 -> RCPT TO:<recipient@test.com>
<-  250 2.1.5 OK
 -> DATA
<-  354 3.0.0 Start mail input
 -> Date: Thu, 01 Jan 2020 00:00:00 +0000
 -> To: recipient@test.com
 -> From: recipient@test.com
 -> Subject: test
 -> Message-Id: <fixed@test>
 -> X-Mailer: swaks v20240103.0 jetmore.org/john/code/swaks/
 -> 
 -> ''' + "a" * 10240 + '''
 -> 
 -> 
 -> .
<-  250 2.0.0 OK
 -> QUIT
<-  221 2.0.0 Bye''',
        expected_log={
            "hostname": "fixed.host",
            "transactions": [
                {
                    "status": 2,
                    "from": ["recipient@test.com"],
                    "to": ["recipient@test.com"],
                    "data": base64.b64encode(
                        b"Date: Thu, 01 Jan 2020 00:00:00 +0000"
                        b"To: recipient@test.com"
                        b"From: recipient@test.com"
                        b"Subject: test"
                        b"Message-Id: <fixed@test>"
                        b"X-Mailer: swaks v20240103.0 jetmore.org/john/code/swaks/"
                        + b"a" * 10240
                    ).decode(),
                }
            ],
        },
    )
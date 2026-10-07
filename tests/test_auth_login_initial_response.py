from smtp_test_runner import run_smtp_test


def test_auth_login_initial_response():
    run_smtp_test(
        command=r'''echo -e "ehlo localhost\nauth login bm9hdXRo\ndGVzdHBhc3M=\nquit\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
250-mx01.example.net
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME
334 UGFzc3dvcmQ6
235 2.7.0 Authentication successful
221 2.0.0 Bye''',
        expected_log={
            "hostname": "localhost",
            "transactions": None,
            "authentication": [
                {
                    "type": "LOGIN",
                    "authorizationId": "",
                    "username": "noauth",
                    "password": "testpass",
                }
            ],
        },
    )
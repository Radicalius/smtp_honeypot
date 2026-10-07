from smtp_test_runner import run_smtp_test


def test_vrfy():
    run_smtp_test(
        command=r'''echo -e "EHLO test\r\nVRFY test@test.com\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
250-mx01.example.net
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME
250 <test@test.com>
221 2.0.0 Bye''',
        expected_log={
            "hostname": "test",
            "transactions": None,
            "authentication": None,
            "verifiedAddrs": ["test@test.com"],
            "tls": False,
            "tlsInfo": None,
            "extended": True,
            "etrn": False,
            "etrnNode": "",
        },
    )
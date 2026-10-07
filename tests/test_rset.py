from smtp_test_runner import run_smtp_test


def test_rset():
    run_smtp_test(
        command=r'''echo -e "EHLO test\r\nDATA\r\na\r\n.\r\nRSET\r\nDATA\r\nb\r\n.\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
250-mx01.example.net
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME
354 3.0.0 Start mail input
250 2.0.0 OK
250 2.0.0 OK
354 3.0.0 Start mail input
250 2.0.0 OK
221 2.0.0 Bye''',
        expected_log={
            "hostname": "test",
            "transactions": [
                {"status": 1, "from": None, "to": None, "data": "YQ=="},
                {"status": 2, "from": None, "to": None, "data": "Yg=="},
            ],
            "authentication": None,
            "verifiedAddrs": None,
            "tls": False,
            "tlsInfo": None,
            "extended": True,
            "etrn": False,
            "etrnNode": "",
        },
    )
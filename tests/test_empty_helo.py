from smtp_test_runner import run_smtp_test


def test_empty_helo():
    run_smtp_test(
        command=r'''printf "HELO\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
250 example.com
221 2.0.0 Bye''',
        expected_log={
            "hostname": "",
            "transactions": None,
            "authentication": None,
            "verifiedAddrs": None,
            "tls": False,
            "tlsInfo": None,
            "extended": False,
            "etrn": False,
            "etrnNode": "",
        },
    )
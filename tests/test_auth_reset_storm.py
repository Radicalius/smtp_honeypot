from smtp_test_runner import run_smtp_test


def test_auth_reset_storm():
    run_smtp_test(
        command=r'''echo -e "EHLO test\r\nAUTH NTLM\r\nTlRMTVNTUAABAAAAB4IIoAAAAAAAAAAAAAAAAAAAAAA=\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
250-mx01.example.net
250-PIPELINING
250-SIZE 10240000
250-VRFY
250-ETRN
250-STARTTLS
250-AUTH LOGIN PLAIN
250 8BITMIME
504 5.5.4 Unrecognized authentication type''',
        expected_log={
            "hostname": "test",
            "transactions": None,
            "authentication": None,
            "verifiedAddrs": None,
            "tls": False,
            "tlsInfo": None,
            "extended": True
        },
    )
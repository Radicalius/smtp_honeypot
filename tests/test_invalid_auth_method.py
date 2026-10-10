from smtp_test_runner import run_smtp_test


def test_invalid_auth_method():
    run_smtp_test(
        command=r'''printf "AUTH INVALID\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
504 5.5.4 Unrecognized authentication type''',
        expected_log={},
    )
from smtp_test_runner import run_smtp_test


def test_invalid_command():
    run_smtp_test(
        command=r'''printf "MGLNDD test\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
500 5.5.1 Command unrecognized
221 2.0.0 Bye''',
        expected_log={},
    )
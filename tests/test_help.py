from smtp_test_runner import run_smtp_test


def test_help():
    run_smtp_test(
        command=r'''printf "HELP\r\nHELP test\r\nhelp\r\nhelp test\r\nQUIT\r\n" | nc localhost 2525''',
        expected_output='''220 mx01.example.net ESMTP
214-This server supports the following commands:
214-HELO EHLO MAIL RCPT DATA RSET VRFY EXPN HELP NOOP QUIT STARTTLS AUTH
214 End of HELP info
214 Help entry for test
214-This server supports the following commands:
214-HELO EHLO MAIL RCPT DATA RSET VRFY EXPN HELP NOOP QUIT STARTTLS AUTH
214 End of HELP info
214 Help entry for test
221 2.0.0 Bye''',
        expected_log={},
    )
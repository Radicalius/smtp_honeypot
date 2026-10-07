import difflib
import json
import os
from pathlib import Path
import subprocess
import time


TEST_DIR = Path(__file__).resolve().parent
PROJECT_DIR = TEST_DIR.parent


def _log_directory():
    log_root = Path(os.environ.get("SMTP_HONEYPOT_LOG_ROOT", PROJECT_DIR / "data"))
    if not log_root.is_absolute():
        log_root = PROJECT_DIR / log_root
    return log_root / "transactions"


def _log_file_sizes(log_directory):
    return {
        path: path.stat().st_size
        for path in log_directory.glob("*.jsonl")
    }


def _new_log_records(log_directory, initial_sizes):
    records = []
    for path in log_directory.glob("*.jsonl"):
        initial_size = initial_sizes.get(path, 0)
        with path.open("rb") as logfile:
            if path.stat().st_size < initial_size:
                initial_size = 0
            logfile.seek(initial_size)
            records.extend(line for line in logfile.read().splitlines() if line)
    return records


def _assert_log_matches(actual, expected):
    if isinstance(expected, dict):
        assert isinstance(actual, dict)
        for key, value in expected.items():
            assert key in actual, f"Expected log key {key!r} is missing"
            _assert_log_matches(actual[key], value)
        return

    if isinstance(expected, list):
        assert isinstance(actual, list)
        assert len(actual) == len(expected)
        for actual_item, expected_item in zip(actual, expected):
            _assert_log_matches(actual_item, expected_item)
        return

    assert actual == expected, f"Expected {expected!r}, got {actual!r}"


def run_smtp_test(command, expected_output, expected_log):
    expected_lines = [
        line for line in expected_output.splitlines() if not line.startswith("===")
    ]

    log_directory = _log_directory()
    initial_sizes = _log_file_sizes(log_directory)

    result = subprocess.run(
        command,
        shell=True,
        cwd=TEST_DIR,
        capture_output=True,
        text=True,
        timeout=60,
    )
    actual_lines = [
        line for line in result.stdout.splitlines() if not line.startswith("===")
    ]

    diff = "\n".join(
        difflib.unified_diff(
            expected_lines,
            actual_lines,
            fromfile="expected",
            tofile="got",
            lineterm="",
        )
    )
    assert not diff, f"SMTP output differed:\n{diff}\n{result.stderr}"

    deadline = time.monotonic() + 10
    records = []
    while time.monotonic() < deadline:
        records = _new_log_records(log_directory, initial_sizes)
        if records:
            break
        time.sleep(0.05)
    assert records, "No transaction log record was written"

    actual_log = json.loads(records[-1])
    _assert_log_matches(actual_log, expected_log)
import argparse
import os
import sys
import difflib

parser = argparse.ArgumentParser(description="Test runner for smtp honeypot")
parser.add_argument(
    "-f", "--filter", type=str, default="", help="Only run matching tests"
)
args = parser.parse_args()

ran = 0
passed = 0
skipped = 0

for file in os.listdir("."):
  if file.endswith(".test"):
    if args.filter != "" and args.filter not in file:
      skipped += 1
      continue

    cont = open(file).read()
    parts = cont.split("\n\n")
    assert len(parts) >= 2
    
    cmd = parts[0]
    res = "\n\n".join(parts[1:])

    out = os.popen(cmd).read()
    ran += 1

    res_lines = [i for i in res.splitlines() if not i.startswith("===")]
    out_lines = [i for i in out.splitlines() if not i.startswith("===")]

    diff = "\n".join(difflib.unified_diff(
      res_lines, out_lines,
      fromfile="expected", tofile="got", lineterm=""
    ))

    if diff != "":
      print(f"Test failed: {file}\n{diff}")
    else:
      passed += 1

print(f"Summary: Tests Passed: {passed}; Tests Run: {ran}; Skipped {skipped}")

if passed < ran:
  sys.exit(1)
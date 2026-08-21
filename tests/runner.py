import argparse
import os
import sys
import difflib
import json

parser = argparse.ArgumentParser(description="Test runner for smtp honeypot")
parser.add_argument(
    "-f", "--filter", type=str, default="", help="Only run matching tests"
)
args = parser.parse_args()

ran = 0
passed = 0
skipped = 0

logfile = os.listdir("../data/transactions")[0]
f = open("../data/transactions/" + logfile, "r")

for file in os.listdir("."):
  if file.endswith(".test"):
    if args.filter != "" and args.filter not in file:
      skipped += 1
      continue

    cont = open(file).read()
    parts = cont.split("\n\n")
    assert len(parts) >= 2
    
    cmd = parts[0]
    res = parts[1]

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
      continue

    json_ = f.read().strip()
    log = json.loads(json_)
    
    failed = False
    if len(parts) > 2:
      explog = json.loads(parts[2])
      for k, v in explog.items():
        if log[k] != v:
          print(f"Test failed: {file}; expected {v} for key {k} but got {log[k]}")
          failed = True

    if failed:
      continue
    
    passed += 1

print(f"Summary: Tests Passed: {passed}; Tests Run: {ran}; Skipped {skipped}")

if passed < ran:
  sys.exit(1)
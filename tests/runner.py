import os
import sys
import difflib

ran = 0
passed = 0

for file in os.listdir("."):
  if file.endswith(".test"):
    cont = open(file).read()
    parts = cont.split("\n\n")
    assert len(parts) >= 2
    
    cmd = parts[0]
    res = "\n\n".join(parts[1:])

    out = os.popen(cmd).read()
    ran += 1

    diff = "\n".join(difflib.unified_diff(
      res.splitlines(), out.splitlines(),
      fromfile="expected", tofile="got", lineterm=""
    ))

    if diff != "":
      print(f"Test failed: {file}\n{diff}")
    else:
      passed += 1

print(f"Summary: Tests Passed: {passed}; Tests Run: {ran}")

if passed < ran:
  sys.exit(1)
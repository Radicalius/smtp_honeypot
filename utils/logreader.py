import os, sys
from datetime import datetime

file = sys.argv[1]
data = open(file, 'rb')

def read_record(data):
  ts = int.from_bytes(data.read(8))
  dir = int.from_bytes(data.read(1))
  length = int.from_bytes(data.read(4))
  data = data.read(length)
  
  dt = datetime.fromtimestamp(ts / 1000.)
  dir_c = "->"
  if dir == 2:
    dir_c = "<-"

  print(f"{dt} {dir_c} {data}")

while len(data.peek(1)) > 0:
  read_record(data)

from os import listdir
from os.path import isfile, join

path = "/home/sftp/uploads"
only_files = [f for f in listdir(path) if isfile(join(path, f))]
report = {}
for file_name in only_files:
    with open(join(path, file_name), 'r') as f:
        data = f.readline()
        user = data.split()[2]
        if user in report:
            report[user] += 1
        else:
            report[user] = 1

for key, value in report.items():
    print(f"User {key} make {value} requests")

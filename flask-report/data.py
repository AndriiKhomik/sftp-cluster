import paramiko
import paramiko.util
from stat import S_ISDIR, S_ISREG
from os.path import join
import os
from pathlib import Path
from dotenv import load_dotenv

load_dotenv()

paramiko.util.log_to_file("paramiko.log")

hosts = os.getenv('FLASK_HOSTS').split(',')
port = int(os.getenv('FLASK_HOST_PORT'))
path = os.getenv('FLASK_SERVER_UPLOADS_PATH')
username = os.getenv('FLASK_SERVER_USERNAME')
password = os.getenv('FLASK_SERVER_PASSWORD')

structured_data = {}


def create_sftp_client(host):
    transport = paramiko.Transport((host, port))
    transport.connect(None, username, password)
    return paramiko.SFTPClient.from_transport(transport), transport


def process_file(sftp, remote_file_path, report):
    try:
        sftp.stat(remote_file_path)  # Check if the file exists
        with sftp.file(remote_file_path, 'r') as f:
            data = f.readline()
            user = data.split()[2]
            date = data.split(',')[0]

            if date not in structured_data:
                structured_data[date] = {}

            if user in structured_data[date]:
                structured_data[date][user] += 1
            else:
                structured_data[date][user] = 1

            if user in report:
                report[user] += 1
            else:
                report[user] = 1
    except FileNotFoundError:
        print(f"File {remote_file_path} does not exist, skipping...")


def process_user_requests(sftp, remote_dir, report):
    for entry in sftp.listdir_attr(remote_dir):
        remote_path = Path(join(remote_dir, entry.filename))

        if S_ISREG(entry.st_mode):
            process_file(sftp, remote_path.as_posix(), report)
        elif S_ISDIR(entry.st_mode):
            process_user_requests(sftp, remote_path.as_posix(), report)


def aggregate_reports():
    report = {}
    for host in hosts:
        sftp, transport = create_sftp_client(host.strip())
        try:
            # Check if the directory exists
            sftp.stat(path)
            print(f"Directory {path} exists on {host}.")
            process_user_requests(sftp, path, report)
        except FileNotFoundError:
            print(f"Directory {path} does not exist on {host}.")
        finally:
            sftp.close()
            transport.close()
    return report


# Aggregate data and generate the report plot
report = aggregate_reports()
# generate_report_plot(report)

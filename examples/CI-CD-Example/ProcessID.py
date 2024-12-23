import subprocess
import sys
import argparse

def get_pid(executable, port):
    try:
        # Run the application with the port argument
        process = subprocess.Popen([executable, port], stdout=subprocess.PIPE, stderr=subprocess.PIPE)

        # Get the PID of the process
        pid = process.pid
        return pid

    except Exception as e:
        print(f'Failed to run the process: {e}')
        return None

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Get the PID of a running process.")
    parser.add_argument("executable", help="Path to the executable.")
    parser.add_argument("port", help="Port number to be used.")

    args = parser.parse_args()

    pid = get_pid(args.executable, args.port)
    if pid is not None:
        print(f'Process ID: {pid}')
    else:
        print("Failed to get the process ID.")

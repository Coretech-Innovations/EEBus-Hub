import subprocess
import json
import sys
import os
import argparse

def run_go_test(go_file, ip_address):
    # Form the ldflags argument with the provided IP address
    ldflags = f'-X "main.BaseIPAddress=http://{ip_address}/api/v1"'

    # Open test.log file for writing logs
    with open("test.log", "w") as log_file:
        # Run the Go file with the ldflags argument
        result = subprocess.run(
            ['go', 'run', go_file, '-ldflags', ldflags], 
            capture_output=True, 
            text=True
        )

        # Capture stdout, stderr, and return status
        stdout = result.stdout
        stderr = result.stderr
        status = result.returncode

        # Write stdout, stderr, and return status to the log file
        log_file.write("STDOUT:\n")
        log_file.write(stdout)
        log_file.write("\nSTDERR:\n")
        log_file.write(stderr)
        log_file.write(f"\nReturn Status: {status}\n")

        # Extract the file name without the extension from the file path
        test_name = os.path.splitext(os.path.basename(go_file))[0]

        # Prepare the JSON data for the new test result
        new_test_result = {
            "test_name": test_name,
            "status": "Passed" if status == 0 else "Failed"
        }

        # Load existing test results if the file exists
        if os.path.exists("test_results.json"):
            with open("test_results.json", "r") as json_file:
                test_results = json.load(json_file)
        else:
            test_results = []

        # Append the new test result
        test_results.append(new_test_result)

        # Write the updated JSON file
        with open("test_results.json", "w") as json_file:
            json.dump(test_results, json_file, indent=4)

        # Print the JSON data (optional)
        print(json.dumps(test_results, indent=4))

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Run Go test and log results.")
    parser.add_argument("go_file", help="Path to the Go file.")
    parser.add_argument("ip_address", help="IP address to be used in the ldflags argument.")
    
    args = parser.parse_args()
    run_go_test(args.go_file, args.ip_address)

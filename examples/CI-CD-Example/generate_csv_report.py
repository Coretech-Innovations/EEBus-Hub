import json
import csv

# Read the JSON file
with open("test_results.json", "r") as json_file:
    test_results = json.load(json_file)

# Define CSV file headers
headers = ["Test Name", "Status"]

# Create a CSV file and write the data
with open("test_results.csv", "w", newline='') as csv_file:
    csv_writer = csv.writer(csv_file)
    csv_writer.writerow(headers)
    
    for result in test_results:
        csv_writer.writerow([result["test_name"], result["status"]])

print("CSV report generated: test_results.csv")

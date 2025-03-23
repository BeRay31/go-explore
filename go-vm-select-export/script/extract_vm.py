import json
import argparse
from collections import defaultdict

def extract_metric_keys_grouped_by_job(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    
    grouped_metric_keys = defaultdict(set)
    for entry in data:
        metric = entry.get("metric", {})
        job_name = metric.get("job")
        if job_name:
            grouped_metric_keys[job_name].update(metric.keys())
    
    return {job: list(keys) for job, keys in grouped_metric_keys.items()}

def main():
    parser = argparse.ArgumentParser(description="Extract Prometheus metric keys grouped by job name.")
    parser.add_argument("file_path", help="Path to the JSON file containing Prometheus metrics.")
    args = parser.parse_args()
    
    metric_keys_by_job = extract_metric_keys_grouped_by_job(args.file_path)
    print(json.dumps(metric_keys_by_job, indent=2))

if __name__ == "__main__":
    main()

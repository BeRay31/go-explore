import json
import argparse
from collections import defaultdict

def extract_metric_keys_grouped_by_job(file_path):
    with open(file_path, 'r') as file:
        data = json.load(file)
    
    grouped_metric_keys = defaultdict(list)
    job_entries = defaultdict(list)
    result = []
    if not isinstance(data, list):
        data_resp = data.get("data", {})
        result = data_resp.get("result", [])
    else:
        result = data
    # Group entries by job
    for entry in result:
        metric = entry.get("metric", {})
        job_name = metric.get("job")
        if job_name:
            job_entries[job_name].append(metric)
    
    # Find common keys for each job (only keys that exist in all entries of the same job)
    for job, metrics_list in job_entries.items():
        common_keys = set(metrics_list[0].keys())
        for metric in metrics_list[1:]:
            common_keys.intersection_update(metric.keys())
        grouped_metric_keys[job] = list(common_keys)
    
    return grouped_metric_keys

def main():
    parser = argparse.ArgumentParser(description="Extract Prometheus metric keys grouped by job name, ensuring only keys that exist in all entries for a specific job.")
    parser.add_argument("file_path", help="Path to the JSON file containing Prometheus metrics.")
    args = parser.parse_args()
    
    metric_keys_by_job = extract_metric_keys_grouped_by_job(args.file_path)
    print(json.dumps(metric_keys_by_job, indent=2))

if __name__ == "__main__":
    main()

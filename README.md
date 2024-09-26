Kubernetes Namespace Cleaner
Kubernetes Namespace Cleaner is a CLI tool designed to identify and remove unused namespaces in your Kubernetes cluster. This tool helps improve resource management and keeps your cluster tidy.
Features

Scans all namespaces in the cluster
Identifies empty namespaces (those without Pods and Services)
Displays a list of remaining resources in each namespace
Deletes empty namespaces after user confirmation
Supports a dry-run mode for safe checking

Installation
bashCopygo get github.com/your-username/kubernetes-namespace-cleaner
Usage
Basic command:
bashCopykubernetes-namespace-cleaner ns-clean
Dry-run mode (check without actual deletion):
bashCopykubernetes-namespace-cleaner ns-clean --dry-run
Options

--dry-run, -d: List namespaces that would be deleted without actually deleting them.

How It Works

The tool scans all namespaces in your Kubernetes cluster.
It checks each namespace for the presence of Pods and Services.
If a namespace has no Pods and Services, it's considered "empty".
For each empty namespace, it shows a list of any remaining resources.
In non-dry-run mode, it prompts for confirmation before deleting each empty namespace.

Caution

This tool considers a namespace "empty" if it has no Pods and Services, even if other resources exist.
Always use the --dry-run flag first to review which namespaces would be deleted.
Be careful not to delete important system namespaces.

Contributing
Bug reports, feature requests, and pull requests are welcome. Please check the contributing guidelines before submitting any contributions.
License
This project is distributed under the MIT License.

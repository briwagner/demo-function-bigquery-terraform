# GCP File Loader to BigQuery

This repo includes:
* Terraform plan to deploy resources on Google Cloud Platform
* GCP bucket to use as trigger for Cloud Function
* BigQuery table and dataset
* Cloud function that responds to file uploads in the bucket and loads data in BigQuery

### Data Loading

There are two example files here to demonstrate loading data into BigQuery. Both contain the same data, but they are in different formats: CSV and JSON. It can be helpful to understand the differences between the two formats when loading data.

**CSV**
* we must instruct BigQuery to skip the first row, if it contains headers
* fields at the end of a row may be skipped

**JSON**
* fields can be skipped, no matter what order they appear in the row
* files tend to be larger in size, as field labels are repeated throughout
* files must be in newline delimited JSON format

What is ETL?

Extract, transform and load. This is a process of transferring large amounts of data from one system into BigQuery. It's important to note that transform occurs before loading, not during the process. If you need to clean up data, or transform some values, that should be done as a separate process before moving to BigQuery. BigQuery operates very efficiently when loading up to millions of lines of data.

It's also possible to stream inserts to BigQuery, which would allow for transforming data in the process. This is slower and likely more expensive. It also involves more restrictive limits on request size.

## Setup

1. Create a file "terraform.tfvars" in the project root
2. Add the following variables:
  * `project-name` = name of GCP project
  * `gcp-region` = GCP region, e.g. "us-central1"
  * `creds-file` = path to GCP credentials file (json format)
3. Change the names and descriptions of the GCP resources
4. `terraform init`
5. Type `terraform plan` to confirm the setup
6. Type `terraform apply` to deploy the resources

Note: this terraform plan assumes we have already created a bucket named `functions-store-bucket` to store the function code. If you need to create that here, it can be added to the terraform plan, similar to the trigger bucket.
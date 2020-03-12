provider "google" {
  project = var.project-name
  region = var.gcp-region
  credentials = file(var.creds-file)
}

variable "project-name" {}
variable "gcp-region" {}
variable "creds-file" {}

locals {
  runtime = "go111"
}

resource "google_storage_bucket" "file_upload_trigger_bucket" {
  name = "file-upload-trigger-bucket"
  location = var.gcp-region
}

data "archive_file" "file-loader-dist" {
  type = "zip"
  source_dir = "./src"
  output_path = "dist/file-loader-function.zip"
}

resource "google_storage_bucket_object" "file-loader-archive" {
  name = "file-loader-archive.zip"
  bucket = "functions-store-bucket"
  source = data.archive_file.file-loader-dist.output_path
}

resource "google_cloudfunctions_function" "file-loader-function" {
  name = "flexible-file-loader-function"
  description = "Load json or csv files and transfer to BigQuery."
  runtime = local.runtime
  available_memory_mb = 128
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource = google_storage_bucket.file_upload_trigger_bucket.name
  }
  entry_point = "LoadData"

  source_archive_bucket = "functions-store-bucket"
  source_archive_object = google_storage_bucket_object.file-loader-archive.name

  environment_variables = {
    PROJECT_ID = var.project-name
    DATASET = google_bigquery_dataset.demo-dataset.dataset_id
    TABLE_ID = google_bigquery_table.marvel_characters.table_id
  }
}

resource "google_bigquery_dataset" "demo-dataset" {
  dataset_id = "demo_loader"
  description = "Demo dataset for Marvel characters"
  location = "US"
}

resource "google_bigquery_table" "marvel_characters" {
  dataset_id = google_bigquery_dataset.demo-dataset.dataset_id
  table_id = "marvel_characters"
}
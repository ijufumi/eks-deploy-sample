import io
import json
import os
import shutil
import tempfile
import zipfile
from datetime import datetime

import boto3
import requests
import yaml

CODEPIPELINE_NAME = os.getenv("CODEPIPELINE_NAME")
BUCKET_NAME = os.getenv("BUCKET_NAME")
ACCESS_TOKEN = os.getenv("ACCESS_TOKEN")
ORGANIZATION_NAME = os.getenv("ORGANIZATION_NAME", "ijufumi")
REPOSITORY_NAME = os.getenv("REPOSITORY_NAME", "eks-deploy-sample")


def get_commit_hash(repository_name: str, ref_name: str) -> str:
    target_url = f"https://api.github.com/repos/{repository_name}/commits/{ref_name}"
    print(f"target_url is {target_url}")
    response = requests.get(target_url, headers={"Authorization": f"token {ACCESS_TOKEN}"})
    print(response)
    if response.status_code != 200:
        raise Exception(f"Get commit hash from {target_url} is failed.")

    return response.json()["sha"]


def download_zip_file(target_url: str):
    print(f"target_url is {target_url}")
    response = requests.get(target_url, headers={"Authorization": f"token {ACCESS_TOKEN}"})
    print(response)
    if response.status_code != 200:
        raise Exception(f"Get zipball from {target_url} is failed.")

    return response.content


def upload_zip_file(repository_name: str, ref_name: str, use_hash: bool = False) -> str:
    target_url = f"https://api.github.com/repos/{repository_name}/zipball/{ref_name}"

    client = boto3.client("s3")
    timestamp = datetime.utcnow().strftime("%Y-%m-%d-%H-%M-%S")
    ref_name2 = ref_name.replace("/", "_")
    file_key = f"source_code/{repository_name}/{timestamp}_{ref_name2}.zip"
    contents = download_zip_file(target_url)
    with tempfile.TemporaryDirectory() as directory:
        with zipfile.ZipFile(io.BytesIO(contents)) as zip_file:
            zip_root_path = zip_file.namelist()[0]
            zip_file.extractall(directory)
        zip_file_path = f"{directory}/{timestamp}"

        root_dir=f"{directory}/{zip_root_path}/app"
        
        with open(f"{root_dir}/tag.txt", "w") as tag_file:
            if use_hash:
                commit_hash = get_commit_hash(repository_name, ref_name)
                tag_file.write(commit_hash[0:7])
            else:
                tag_file.write(ref_name2)

        shutil.make_archive(base_name=zip_file_path, format="zip", root_dir=root_dir)

        client.upload_file(f"{zip_file_path}.zip", BUCKET_NAME, file_key)

    return file_key


def lambda_handler(event, context):
    print(event)
    request_data = json.loads(event.get("body", "{}"))
    # repository_name contains organization
    repository_name = request_data.get("repository")
    branch_name = request_data.get("branch")
    tag_name = request_data.get("tag")
    ref_name = tag_name if tag_name is not None else branch_name

    if repository_name is None:
        if ORGANIZATION_NAME is not None and REPOSITORY_NAME is not None:
            repository_name = f"{ORGANIZATION_NAME}/{REPOSITORY_NAME}"

    if repository_name is None or ref_name is None:
        return {"statusCode": 400, "body": "Missing necessary parameters"}

    print(f"target is {repository_name}/{ref_name}")

    client = boto3.client("codepipeline")
    if ref_name:
        file_key = upload_zip_file(repository_name, ref_name, tag_name is None)

        response = client.get_pipeline(name=CODEPIPELINE_NAME)
        print(response)
        pipeline = response["pipeline"]
        pipeline["stages"][0]["actions"][0]["configuration"]["S3ObjectKey"] = file_key
        response = client.update_pipeline(pipeline=pipeline)
        print(response)

    response = client.start_pipeline_execution(name=CODEPIPELINE_NAME)
    return {"statusCode": 200, "body": "Successfully completed"}

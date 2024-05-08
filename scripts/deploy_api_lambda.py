import subprocess
import os
import shutil
import time
import zipfile

PROJECT_DIR = os.path.abspath(os.path.join(os.path.dirname(os.path.abspath(__file__)), '..'))
PROGRAM_FILE = 'bootstrap'
OUTPUT_ZIP = 'bootstrap.zip'
LAMBDA_NAME = 'go-test'
AWS_REGION = 'us-east-2'

def upload_to_aws_lambda():
    try:
        print('uploading to aws lambda')
        time.sleep(.5)
        command = f'aws lambda update-function-code --function-name {LAMBDA_NAME} --region {AWS_REGION} --zip-file fileb://bootstrap.zip'
        subprocess.run(["powershell", command], capture_output=True, text=True)
        return True
    except Exception as e:
        print(e)
        return False

def zip_project():
    try:
        source_path = f'{PROJECT_DIR}\\{PROGRAM_FILE}'
        zip_path = f'{PROJECT_DIR}\\{OUTPUT_ZIP}'
        print(f'zipping: {zip_path}')
        with zipfile.ZipFile(zip_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
            zipf.write(source_path, arcname='bootstrap')
        return True
    except Exception as e:
        print(e)
        return False

def build_release_project():
    try:
        current_dir = os.getcwd()
        print("Current directory:", current_dir)
        build_type = '$env:GOOS="linux"; $env:GOARCH="arm64"; $env:CGO_ENABLED = "0";'
        build = f'go build -tags lambda.norpc -o .\\{PROGRAM_FILE} .\\cmd\\lambda\\lambda.go;'
        command = f'{build_type} {build}'
        result = subprocess.run(["powershell", command], capture_output=True, text=True)
        time.sleep(.5)
        return True
    except Exception as e:
        print(e)
        return False

def remove_old_artifacts():
    try:
        source_path = f'{PROJECT_DIR}\\{PROGRAM_FILE}'
        if os.path.exists(source_path):
            os.remove(source_path)
            print('deleted release directory')
        else:
            print('no release directory found')
        output_zip_path = f'{PROJECT_DIR}\\{OUTPUT_ZIP}'
        if os.path.exists(output_zip_path):
            os.remove(output_zip_path)
            print('deleted output zip file')
        else:
            print('no output zip file to delete')
        time.sleep(.5)
        return True
    except Exception as e:
        print(e)
        return False

def exit_program():
    print('exiting program...')
    exit()

def change_directory_to(path):
    try:
        os.chdir(path)
        return True
    except Exception as e:
        print(e)
        return False

def main():
    if not remove_old_artifacts():
        exit_program()

    if not change_directory_to(PROJECT_DIR):
        exit_program()

    if not build_release_project():
        exit_program()

    if not zip_project():
        exit_program()

    if not upload_to_aws_lambda():
        exit_program()

    if not remove_old_artifacts():
        exit_program()

if __name__ == "__main__":
    main()
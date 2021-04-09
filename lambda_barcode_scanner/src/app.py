import json
import boto3
import base64
import uuid

bucket_name = 'ocr-zone'
folder_name = 'uploaded'

# Response for success.
def response_success(body: dict):
    return {
        'statusCode': 200,
        'body': json.dumps({
            'data': body
        })
    }

# Response for error.
def response_error(status: int, message: str):
    return {
        'statusCode': status,
        'body': json.dumps({
            'data': [],
            'message': message,
        })
    }


def lambda_handler(event, context):
    """
    Parameters
    ----------
    event: dict, required
        API Gateway Lambda Proxy Input Format

        Event doc: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-input-format

    context: object, required
        Lambda Context runtime methods and attributes

        Context doc: https://docs.aws.amazon.com/lambda/latest/dg/python-context-object.html

    Returns
    ------
    API Gateway Lambda Proxy Output Format: dict

        Return doc: https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
    """
    # Check JSON body.  
    if not 'body' in event or event['body'] == None:
        return response_error(400, 'body not exist.')

    # Parsing request_body json into dictionary.
    try:
        input_data = json.loads(event['body'])
    except json.JSONDecodeError:
        return response_error(400, 'Invalid JSON.')
    
    # Get file base64.
    try:
        img_base64 = input_data['img']
    except KeyError:
        return response_error(400, 'img is required.')
    
    # Init s3.
    s3 = boto3.resource('s3')
    bucket = s3.Bucket(bucket_name)
    filename = str(uuid.uuid4())
    file_path = f'{folder_name}/{filename}.png'

    # Decode image.
    try:
        img = base64.b64decode(img_base64)
    except Exception:
        return response_error(400, 'invalid base64 img')

    # Upload image to s3.
    try:
        temp_path = '/tmp/output'
        with open(temp_path, 'wb') as data:
            data.write(img)
            bucket.upload_file(temp_path, file_path)
    except Exception as err:
        print(err)
        return response_error(400, 'failed uploading image')

    # Scan with recognition.
    try:
        rekognition = boto3.client('rekognition')
        result = rekognition.detect_text(Image={'S3Object': {'Bucket': bucket_name, 'Name': file_path }})
        text_detections = result['TextDetections']
    except Exception as err:
        print(err)
        return response_error(400, 'failed scanning barcode')


    return response_success(text_detections)
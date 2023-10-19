# receipt-processor-challenge
## Build the Docker Image
To build the Docker Image run the following command: 
```console
docker build -t <image-name>:latest .
```
Replace \<image-name\> with the name that you want for the Docker image.

## Run the Docker Container
```console
docker run -p 8080:8080 <your-api-name>:latest
```
After running this command the API should now be accessible at http://localhost:8080.

## Using the API

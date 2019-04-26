# Portainer ECR Credentials Updater

This is a docker container that will update the Elastic Container Registry credentials in Portainer 

## Usage

Build portainer-ecr container  

Configure environment variables (you can use AWS_ROLE_ARN *OR* AWS Keys)
- AWS_ROLE_ARN 
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- AWS_REGION
- PORTAINER_URL
- PORTAINER_USER
- PORTAINER_PASS
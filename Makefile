.PHONY: package deploy build

package:
	sam package \
	    --template-file template.yaml \
	    --output-template-file package.yaml \
	    --s3-bucket ${BUCKET_NAME}

deploy:
	sam deploy \
	    --template-file package.yaml \
	    --stack-name ${STACK_NAME} \
	    --capabilities CAPABILITY_IAM \
		--parameter-overrides Token=${TOKEN} \
		    Username=${USERNAME} \
			Channel=${CHANNEL} 

build: main.go
	GOOS=linux go build -o build/codepipeline_slack

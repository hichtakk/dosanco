secret: ## create dosanco backend db secret
	kubectl create secret generic dosanco-db-secret --from-literal=url=${DOSANCO_DB}

deploy: ## deploy dosanco to kubernetes
	kubectl create -f ./deployment.yaml

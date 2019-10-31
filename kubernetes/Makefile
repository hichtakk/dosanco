secret: ## create dosanco backend db secret
	kubectl create secret generic dosanco-db-secret --from-literal=url=${DOSANCO_DB}

deploy: ## deploy dosanco to kubernetes
	kubectl create -f ./deployment.yaml

rollout:
	kubectl describe deployment dosanco
	kubectl set image deployment/dosanco dosanco=hichtakk/dosanco:${VERSION} --record
	kubectl describe deployment dosanco
	kubectl rollout status deployment/dosanco

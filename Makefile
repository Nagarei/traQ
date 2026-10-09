.PHONY: env
env: ## Generate env file
	#git remote add upstream https://github.com/traPtitech/traQ.git
	git fetch upstream master --tags
	git pull
	echo -n > .env
	echo "TRAQ_VERSION=$$(git describe --tags --abbrev=0)" >> .env
	echo "TRAQ_REVISION=$$(git rev-parse --short HEAD)" >> .env


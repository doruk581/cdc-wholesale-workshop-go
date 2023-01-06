include ./make/config.mk

install:
	@if [ ! -d pact/bin ]; then\
		echo "--- Installing Pact CLI dependencies";\
		curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash;\
    fi

run-consumer:
	@go run consumer/client/cmd/main.go

run-provider:
	@go run provider/cmd/productsvc/main.go

deploy-consumer: install
	@echo "--- ✅ Checking if we can deploy consumer"
	@pact-broker can-i-deploy \
		--pacticipant $(CONSUMER_NAME) \
		--broker-base-url ${PACT_BROKER_PROTO}://$(PACT_BROKER_URL) \
		--latest

deploy-provider: install
	@echo "--- ✅ Checking if we can deploy provider"
	@pact-broker can-i-deploy \
		--pacticipant $(PROVIDER_NAME) \
		--broker-base-url ${PACT_BROKER_PROTO}://$(PACT_BROKER_URL) \
		--latest

publish: install
	@echo "--- 📝 Publishing Pacts"
	go run consumer/client/pact/publish.go
	@echo
	@echo "Pact contract publishing complete!"

unit:
	@echo "--- 🔨Running Unit tests "
	go test github.com/doruk581/cdc-wholesale-workshop-go/consumer/client -run 'TestClientUnit'

consumer: export PACT_TEST := true
consumer: install
	@echo "--- 🔨Running Consumer Pact tests "
	go test -count=1 github.com/doruk581/cdc-wholesale-workshop-go/consumer/client -run 'TestClientPact'

provider: export PACT_TEST := true
provider: install
	@echo "--- 🔨Running Provider Pact tests "
	go test -count=1 -tags=integration github.com/doruk581/cdc-wholesale-workshop-go/provider -run "TestPactProvider"

.PHONY: install deploy-consumer deploy-provider publish unit consumer provider

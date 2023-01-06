SHELL = "/bin/bash"

export PATH := $(PWD)/pact/bin:$(PATH)
export PATH
export PROVIDER_NAME = GoProductService
export CONSUMER_NAME = GoListingService
export PACT_DIR = $(PWD)/pacts
export LOG_DIR = $(PWD)/log
export PACT_BROKER_USERNAME = dXfltyFMgNOFZAxr8io9wJ37iUpY42M
export PACT_BROKER_PASSWORD = O5AIZWxelWbLvqMd8PkAVycBJh2Psyg1
export PACT_BROKER_PROTO = https
export PACT_BROKER_URL = test.pact.dius.com.au
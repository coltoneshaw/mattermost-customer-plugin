# Include custom targets and environment variables here

## Generate mocks.
mocks:
ifneq ($(HAS_SERVER),)
	./bin/mockgen -destination server/bot/mocks/mock_poster.go github.com/coltoneshaw/mattermost-plugin-customers/server/bot Poster
	./bin/mockgen -destination server/sqlstore/mocks/mock_kvapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore KVAPI
	./bin/mockgen -destination server/sqlstore/mocks/mock_storeapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore StoreAPI
	./bin/mockgen -destination server/sqlstore/mocks/mock_configurationapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore ConfigurationAPI
endif

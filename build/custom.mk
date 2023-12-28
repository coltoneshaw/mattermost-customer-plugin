# Include custom targets and environment variables here

## Generate mocks.
mocks:
ifneq ($(HAS_SERVER),)
	mockgen -destination server/bot/mocks/mock_poster.go github.com/coltoneshaw/mattermost-plugin-customers/server/bot Poster
	mockgen -destination server/app/mocks/mock_job_once_scheduler.go github.com/coltoneshaw/mattermost-plugin-customers/server/app JobOnceScheduler
	mockgen -destination server/sqlstore/mocks/mock_kvapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore KVAPI
	mockgen -destination server/sqlstore/mocks/mock_storeapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore StoreAPI
	mockgen -destination server/sqlstore/mocks/mock_configurationapi.go github.com/coltoneshaw/mattermost-plugin-customers/server/sqlstore ConfigurationAPI
endif

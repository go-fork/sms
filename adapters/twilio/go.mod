module github.com/go-fork/sms/adapters/twilio

go 1.18

require (
	github.com/go-resty/resty/v2 v2.7.0
	github.com/spf13/viper v1.15.0
	github.com/go-fork/sms v0.0.0
)

// Use a replace directive for local development
// This should be removed before publishing
replace github.com/go-fork/sms => ../../

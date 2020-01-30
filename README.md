service-framework
=================
[![CircleCI](https://circleci.com/gh/blackbx/service-framework.svg?style=svg&circle-token=ee056b0c6563f64ff0ef888d6d252f0acd140235)](https://circleci.com/gh/blackbx/service-framework)
[![codecov](https://codecov.io/gh/blackbx/service-framework/branch/master/graph/badge.svg?token=7WukXOxzcR)](https://codecov.io/gh/blackbx/service-framework)


service-framework is a package that allows you to quickly and easily build new
services. service-framework follows the principles of [12 Factor Apps][twelve-factor-apps]
and as such allows you to be able to rapidly develop maintainable, and reliable
apps.

### Libraries
service-framework aggregates multiple go libraries, and will inject them into
your service.

service-framework uses the following libraries:

* [Heptio Healthcheck][healthcheck]
* [Redis][redis]
* [Gorrilla Mux][mux]
* [Gorrilla Handlers][handlers]
* [SQLX][sqlx]
* [NewRelic][newrelic]
* [Cobra][cobra]
* [Viper][viper]
* [Uber FX][fx]
* [Uber Zap][zap]


[twelve-factor-apps]: https://12factor.net/
[healthcheck]: https://github.com/heptiolabs/healthcheck
[redis]: https://github.com/go-redis/redis/v7
[mux]: https://github.com/gorilla/mux
[handlers]: https://github.com/gorilla/handlers
[sqlx]: https://github.com/jmoiron/sqlx
[newrelic]: https://github.com/newrelic/go-agent
[cobra]: https://github.com/spf13/cobra
[viper]: https://github.com/spf13/viper
[fx]: https://go.uber.org/fx
[zap]: https://go.uber.org/zap


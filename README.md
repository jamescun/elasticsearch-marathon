Elasticsearch Marathon
======================

Elasticsearch Marathon is a wrapper around the Docker Library's Elasticsearch 2.4 image.

Discovery is done by placing multiple elasticsearch apps within the same grouping in Marathon. For example, all apps (and their associated tasks) started with the prefix `/elasticsearch/mycluster/[app name]` or `/foo/bar/[app name]` can form a cluster. If an app is started outside of a group, Elasticsearch Marathon will only form a cluster with other tasks within the app.


Configuration
-------------

Elasticsearch Marathon is configured with environment variables.

    +----------------------------+----------+---------+--------------------------------+
    | name                       | required | default | description                    |
    +----------------------------+----------+---------+--------------------------------+
    | MARATHON_ADDR              | yes      |         | http://[ip of marathon]:[port] |
    +----------------------------+----------+---------+--------------------------------+

**NOTE:** You must configure the cluster name indentically across all nodes of the cluster otherwise they will not discover each other; do this either in a configuration file, the `--cluster.name` command line flag or `ELASTICSEARCH_CLUSTER_NAME` environment variable.


### Examples

  - [Standalone](examples/standalone/standalone.json)
  - [HA Master](examples/ha/master.json) and [HA Slaves](examples/ha/slaves.json)


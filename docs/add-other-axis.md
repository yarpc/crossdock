[← Integrate Other Repos](integrate-other-repos.md)

# Add Other Test Axis

You can also add additional Axis to the integration test.
Let's update our `docker-compose.yml` file to include an additional Dimension:

```
# ...
- AXIS_BEHAVIOR=dance,run
- AXIS_SPEED=fast,slow
```

Now run Crossdock:

```
$ docker-compose run crossdock

Beginning matrix of tests...

  STATUS  |  CLIENT   |      RESPONSE       | SPEED | BEHAVIOR
+---------+-----------+---------------------+-------+----------+
  PASSED  | client    | ok                  | fast  | dance
  SKIPPED | client    | 404                 | fast  | run
  PASSED  | client    | ok                  | slow  | dance
  SKIPPED | client    | 404                 | slow  | run
  PASSED  | newclient | ok                  | fast  | dance
  PASSED  | newclient | ok                  | fast  | run
  PASSED  | newclient | ok                  | slow  | dance
  PASSED  | newclient | ok                  | slow  | run

```

Awesome - a test was triggered for every combination of "speed" and "behavior".

In this way, we can develop our Test Client's to address multiple dimensions of integration.

[Run During Continuous Integartion →](add-to-ci.md)

job "deployment_auto.nomad" {
  datacenters = ["dc1"]

  group "one" {
    count = 3

    update {
      max_parallel = 3
      auto_promote = true
      canary = 2
    }

    task "one" {
      driver = "raw_exec"

      config {
	command = "/bin/sleep"
	args = ["1000000"]
      }

      resources {
	cpu    = 20
	memory = 20
      }
    }
  }

  group "two" {
    count = 3

    update {
      max_parallel = 2
      auto_promote = true
      canary = 2
      min_healthy_time = "5s"
    }

    task "two" {
      driver = "raw_exec"

      config {
	command = "/bin/sleep"
	args = ["200000"]
      }

      resources {
	cpu    = 20
	memory = 20
      }
    }
  }
}

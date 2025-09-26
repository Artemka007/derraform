resource "docker_image" "nginx" {
  name = "nginx:alpine"
}

resource "docker_image" "redis" {
  name = "redis:alpine"
}

resource "docker_network" "app_network" {
  name = "myapp-network"
  driver = "bridge"
}

resource "docker_container" "web_server" {
  name  = "web-server"
  image = docker_image.nginx.name
  
  ports {
    internal = 80
    external = 8080
  }
  
  networks = [docker_network.app_network.name]
  
  env = {
    "NGINX_HOST" = "localhost"
    "NGINX_PORT" = "80"
  }
  
  volumes {
    host_path      = "./html"
    container_path = "/usr/share/nginx/html"
    read_only      = true
  }
}

resource "docker_container" "cache" {
  name  = "redis-cache"
  image = docker_image.redis.name
  
  networks = [docker_network.app_network.name]
  
  ports {
    internal = 6379
    external = 6379
  }
  
  env = {
    "REDIS_PASSWORD" = "secret123"
  }
  
  healthcheck {
    test     = ["CMD", "redis-cli", "ping"]
    interval = "30s"
    timeout  = "10s"
    retries  = 3
  }
}
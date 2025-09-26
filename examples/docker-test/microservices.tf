resource "docker_network" "backend" {
  name = "backend-network"
}

resource "docker_volume" "db_data" {
  name = "postgres-data"
}

resource "docker_container" "database" {
  name  = "postgres-db"
  image = "postgres:13"
  
  networks = [docker_network.backend.name]
  
  env = {
    "POSTGRES_DB"       = "myapp"
    "POSTGRES_USER"     = "admin"
    "POSTGRES_PASSWORD" = "password123"
  }
  
  volumes {
    volume_name = docker_volume.db_data.name
    container_path = "/var/lib/postgresql/data"
  }
  
  ports {
    internal = 5432
    external = 5432
  }
}

resource "docker_container" "api" {
  name  = "api-server"
  image = "node:18-alpine"
  
  networks = [docker_network.backend.name]
  
  ports {
    internal = 3000
    external = 3000
  }
  
  env = {
    "DATABASE_URL" = "postgresql://admin:password123@postgres-db:5432/myapp"
    "NODE_ENV"     = "production"
  }
  
  depends_on = [docker_container.database]
}
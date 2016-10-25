provider "ipam" {
  address = "127.0.0.1:8000"
  scheme = "http"
}

resource "ipam_pool" "demo_pool" {
  name = "POOL_1"
}

resource "ipam_subnet" "demo_subnet" {
  name = "SUBNET_1"
  pool = "${ipam_pool.demo_pool.id}"
  start = "192.168.1.10"
  end = "192.168.1.20"
}

resource "ipam_reservation" "demo_reservation" {
 name = "RESERVATION_1"
 subnet = "${ipam_subnet.demo_subnet.id}"
}

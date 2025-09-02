resource "opticloud_instance" "instance" {
  name             = "instance"
  service_offering = "Small Instance"
  template         = "CentOS 5.6 64-bit no GUI Simulator"
  zone             = "BR1"
}
